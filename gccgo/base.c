// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <ucontext.h>
#include <sys/syscall.h>
#include <errno.h>
#include <semaphore.h>
#include <pthread.h>

extern void __splitstack_getcontext(void *context[10]);
extern void __splitstack_setcontext(void *context[10]);

extern void simple_cgocall(void (*)(void*), void*);
extern void simple_cgocallback(void (*)(void*), void (*)(void*), void*);

extern void (*cgocallback)(void (*)(void*), void (*)(void*), void*);

extern void runtime_check(void);
extern void runtime_osinit(void);
extern void runtime_schedinit(void);
extern void *__go_go(void (*fn)(void *), void *);
extern void runtime_mstart(void *);
extern void *runtime_m(void);

extern void runtime_entersyscall(void) __asm__("libgo_syscall.syscall.Entersyscall");
extern void runtime_exitsyscall(void) __asm__("libgo_syscall.syscall.Exitsyscall");
extern void runtime_LockOSThread(void) __asm__("libgo_runtime.runtime.LockOSThread");
extern void runtime_UnlockOSThread(void) __asm__("libgo_runtime.runtime.UnlockOSThread");

extern void main_init(void) __asm__ ("__go_init_main");

typedef struct ctxt {
    struct ctxt *next;
    void (*fn)(void*);
    void *arg;
    void (*gfn)(void*);
    void *garg;
    sem_t sem;
} Ctxt;

typedef struct state {
    int in_go;
    ucontext_t gmc, c, g0c;
    void *gm_stack[10];
    void *stack[10];
    void *g0_stack[10];
    void (*fn)(void*);
    void *arg;
    void (*gfn)(void*);
    void *garg;
    struct sigaction handlers[_NSIG];
} S;

static __thread S *s;
static S *s0;

static Ctxt *ctxt_head = NULL;
static Ctxt *ctxt_tail = NULL;
static sem_t ctxt_sem;
static pthread_mutex_t ctxt_mutex = PTHREAD_MUTEX_INITIALIZER;

static volatile int base_init_done = 0;

S *get_s(void) __attribute__((noinline, no_split_stack));

S *get_s(void) {
    return s;
}

static void store_signal_handlers(struct sigaction handlers[_NSIG]) {
    int i;
    for (i = 0; i < _NSIG; i++) {
        sigaction(i, NULL, &handlers[i]);
    }
}

static void restore_signal_handlers(struct sigaction handlers[_NSIG]) {
    int i;
    for (i = 0; i < _NSIG; i++) {
        sigaction(i, &handlers[i], NULL);
    }
}

static void g0_main(void *_) {
    // Initialise state
    s->gfn = NULL;
    s->in_go = 0;

    // Setup return entry point
    __splitstack_getcontext(&s->g0_stack[0]);
    getcontext(&s->g0c);

    // in_go + gfn == return to call gfn from goroutine
    if (s->in_go && s->gfn) return;

    // just in_go == we have come back to run some go code on g0 ...
    if (s->in_go) s->fn(s->arg);

    // reset state
    s->gfn = NULL;
    s->in_go = 0;

    // jump back to whatever non-go code got us here
    __splitstack_setcontext(&s->stack[0]);
    setcontext(&s->c);
}

static void ctxtHandler(void *arg) {
    Ctxt *ctxt = arg;

    runtime_LockOSThread();

    runtime_entersyscall();

    // Create a new state object
    s = calloc(1, sizeof(S));

    // Swap over to the Go signal handlers
    store_signal_handlers(s->handlers);
    restore_signal_handlers(s0->handlers);

    // we are about to go into go
    s->in_go = 1;

    // use s->fn as a "still doing init" flag
    s->fn = (void *)1;

    // Store the current context, this is where we will jump back to
    __splitstack_getcontext(&s->stack[0]);
    getcontext(&s->c);

    // !in_go && s->fn == 1 means we are back from g0_main to setup s
    if (!s->in_go && s->fn == (void *)1) {
        // we are about to go into go
        s->in_go = 1;

        // Copy fn/gfn from ctxt to s
        s->fn   = ctxt->fn;
        s->arg  = ctxt->arg;
        s->gfn  = ctxt->gfn;
        s->garg = ctxt->garg;

        // jump back to g0_main
        __splitstack_setcontext(&s->g0_stack[0]);
        setcontext(&s->g0c);
    }

    runtime_exitsyscall();

    // !in_go && s->fn != 1 means we are done
    if (!s->in_go) {
        // Swap back to our own signal handlers
        restore_signal_handlers(s->handlers);

        // Free state
        free(s);
        s = NULL;

        runtime_UnlockOSThread();

        sem_post(&ctxt->sem);
        return;
    }

    while (1) {
        simple_cgocall(g0_main, NULL);
        // g0_main only returns when we want to run code in a goroutine ...
        s->gfn(s->garg);
    }
}

static Ctxt *recv_ctxt(void) {
    Ctxt *head = NULL;
    while (sem_wait(&ctxt_sem) != 0) {
        if (errno == EINTR) continue;
        fprintf(stderr, "libgopy.ext internal error: sem_wait failed.\n");
        abort();
    }
    pthread_mutex_lock(&ctxt_mutex);
    if (!ctxt_head) {
        fprintf(stderr, "libgopy.ext internal error: ctxt_head NULL.\n");
        abort();
    }
    head = ctxt_head;
    ctxt_head = ctxt_head->next;
    if (ctxt_tail == head) ctxt_tail = ctxt_head;
    pthread_mutex_unlock(&ctxt_mutex);
    return head;
}

static void send_ctxt(Ctxt *ctxt) {
    pthread_mutex_lock(&ctxt_mutex);
    ctxt->next = NULL;
    if (ctxt_tail) {
        ctxt_tail->next = ctxt;
        ctxt_tail = ctxt;
    } else {
        ctxt_head = ctxt;
        ctxt_tail = ctxt;
    }
    pthread_mutex_unlock(&ctxt_mutex);
    sem_post(&ctxt_sem);
}

// This function runs in a goroutine and runs ctxt requests in new goroutines
static void ctxtDispatcher(void *arg __attribute__((unused))) {
    Ctxt *ctxt;
    sem_init(&ctxt_sem, 0, 0);
    while (1) {
        runtime_entersyscall();
        ctxt = recv_ctxt();
        runtime_exitsyscall();
        __go_go(ctxtHandler, ctxt);
    }
}

static void run_on_ctxt(void (*fn)(void*), void *arg,
                        void (*gfn)(void*), void *garg) {
    Ctxt *ctxt = malloc(sizeof(*ctxt));
    if (!ctxt) {
        fprintf(stderr, "libgopy.ext internal error: ctxt NULL.\n");
        abort();
    }
    ctxt->fn   = fn;
    ctxt->arg  = arg;
    ctxt->gfn  = gfn;
    ctxt->garg = garg;
    sem_init(&ctxt->sem, 0, 0);
    send_ctxt(ctxt);
    while (sem_wait(&ctxt->sem) != 0) {
        if (errno == EINTR) continue;
        fprintf(stderr, "libgopy.ext internal error: sem_wait failed.\n");
        abort();
    }
    sem_destroy(&ctxt->sem);
    free(ctxt);
}

extern void main_main(void) __asm__ ("main.main");
extern void main_main(void) {
    __go_go(ctxtDispatcher, NULL);
    while (1) {
        simple_cgocall(g0_main, NULL);
        // g0_main only returns when we want to run code in a goroutine ...
        s->gfn(s->garg);
    }
}

static void mainstart(void *arg __attribute__((unused))) {
    runtime_main();
}

static void activate_go(void) {
    struct sigaction handlers[_NSIG];

    s->in_go = 1;

    // Swap over to the Go signal handlers
    store_signal_handlers(handlers);
    restore_signal_handlers(s0->handlers);

    __splitstack_getcontext(&s->stack[0]);
    getcontext(&s->c);

    if (s->in_go) {
        __splitstack_setcontext(&s->g0_stack[0]);
        setcontext(&s->g0c);
    }

    // Swap back to our own signal handlers
    restore_signal_handlers(handlers);
}

static void run_on_g0(void (*f)(void*), void *p) {
    if (!s) return run_on_ctxt(f, p, NULL, NULL);
    s->fn = f;
    s->arg = p;
    activate_go();
}

static void run_on_g(void (*f)(void*), void *p) {
    if (!s) return run_on_ctxt(NULL, NULL, f, p);
    s->gfn = f;
    s->garg = p;
    activate_go();
}

static void go_main(void) {
    // 2 because Go expects to find a null terminated list of env strings at the
    // end of argv ...
    char *argv[2] = {0};

    runtime_check();
    runtime_args(0, argv);
    runtime_osinit();
    runtime_schedinit();
    __go_go(mainstart, NULL);
    runtime_mstart(runtime_m());
}

static void state_init(void (*fn)(void)) {
    size_t ss;

    // Create a new state object
    s = calloc(1, sizeof(S));

    // we can't start the go runtime from the current context, so we need to
    // creeate a new one ...
    getcontext(&s->gmc);
    s->gmc.uc_stack.ss_sp = malloc(2 * 1024 * 1024);
    s->gmc.uc_stack.ss_size = 2 * 1024 * 1024;
    s->gmc.uc_link = NULL;
    makecontext(&s->gmc, fn, 0);
    __splitstack_makecontext(s->gmc.uc_stack.ss_size, &s->gm_stack[0], &ss);

    // we are about to go into go
    s->in_go = 1;

    // If this isn't s0, then we need to restore the s0 signal handlers, as
    // those are the ones Go expects to be running.
    if (s0) restore_signal_handlers(s0->handlers);

    // once we have started the go runtime in it's own context, we need to be
    // able to get back to this one to return from this function
    __splitstack_getcontext(&s->stack[0]);
    getcontext(&s->c);

    if (s->in_go) {
        // actually jump into go
        __splitstack_setcontext(&s->gm_stack[0]);
        setcontext(&s->gmc);
    }

    // save the signal handlers that Go is using, so that we can reactivate them
    // when switching into Go context ...
    store_signal_handlers(s->handlers);
}

static void cgocallback_g(void *_a) {
    struct {
        void (*fn)(void*);
        void *arg;
    } *a = _a;
    simple_cgocallback(a->fn, NULL, a->arg);
}

static void cgocallback_wrapper(void (*fn)(void*), void (*ef)(void*),
                                void *param) {
    struct {
        void (*fn)(void*);
        void *arg;
    } a;
    if (s && s->in_go) return simple_cgocallback(fn, NULL, param);
    a.fn = fn;
    a.arg = param;
    run_on_g0(cgocallback_g, &a);
}

static void base_init(void) {
    struct sigaction handlers[_NSIG];

    cgocallback = cgocallback_wrapper;

    // we need to save all the signal handlers, so we can stop Go from co-opting
    // them.
    store_signal_handlers(handlers);

    state_init(go_main);

    // store the base state so that all states can look at it
    s0 = s;

    // we should restore the signal handlers now.
    restore_signal_handlers(handlers);

    // the runtime is now up and running in it's own context, and we can access
    // it using the g0c context - time to return
    base_init_done = 1;
}

typedef void (*ifunc)(void);

struct ie {
    ifunc f;
    struct ie *n;
};

static struct ie *init_done = NULL;

static void do_init(ifunc f) {
    struct ie *i = init_done;
    while (i) {
        if ((i)->f == f) return;
        i = i->n;
    }
    f();
    i = malloc(sizeof(*i));
    i->f = f;
    i->n = init_done;
    init_done = i;
}

static void do_inits(ifunc funcs[]) {
    int i = 0;
    while (funcs[i]) {
        do_init(funcs[i++]);
    }
}

extern void _init_go(ifunc funcs[]) {
    int i = 0;

    if (!base_init_done) base_init();

    run_on_g((void (*)(void*))do_inits, funcs);
}

extern ifunc py_init_funcs[];

extern void main_init(void) {
    // We need to make sure that we don't get switched off to another thread,
    // otherwise Python will get very confused when we return from a function
    // call on a different thread to the one that called it!
    runtime_LockOSThread();

    do_inits(py_init_funcs);
}
