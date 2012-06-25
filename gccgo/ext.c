// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*****************************************************************************
******************************************************************************
**
** Includes
**
******************************************************************************
*****************************************************************************/

#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <ucontext.h>
#include <sys/syscall.h>
#include <errno.h>
#include <semaphore.h>
#include <pthread.h>

/*****************************************************************************
******************************************************************************
**
** Types
**
******************************************************************************
*****************************************************************************/

typedef struct ctxt {
    struct ctxt *next;
    void (*fn)(void*);
    void *arg;
    void (*gfn)(void*);
    void *garg;
    sem_t sem;
} Ctxt;

typedef struct state {
    void *gm_stack[10];
    void *stack[10];
    void *g0_stack[10];
    void (*fn)(void*);
    void *arg;
    void (*gfn)(void*);
    void *garg;
    struct sigaction handlers[_NSIG];
    ucontext_t c, g0c;
    int in_go;
} S;

typedef void (*ifunc)(void);

struct ie {
    ifunc f;
    struct ie *n;
};

/*****************************************************************************
******************************************************************************
**
** Externs
**
******************************************************************************
*****************************************************************************/

extern void __splitstack_getcontext(void *context[10]);
extern void __splitstack_setcontext(void *context[10]);

// Functions in main goPy library
extern void simple_cgocall(void (*)(void*), void*);
extern void simple_cgocallback(void (*)(void*), void*);

// Function pointer in main goPy library
extern void (*cgocallback)(void (*)(void*), void*);

// C functions in the Go runtime
extern void runtime_check(void);
extern void runtime_osinit(void);
extern void runtime_schedinit(void);
extern void *__go_go(void (*fn)(void *), void *);
extern void runtime_mstart(void *);
extern void *runtime_m(void);
extern void runtime_main(void);

// "Go" functions in the Go runtime
extern void runtime_entersyscall(void) __asm__("syscall.Entersyscall");
extern void runtime_exitsyscall(void) __asm__("syscall.Exitsyscall");
extern void runtime_LockOSThread(void) __asm__("runtime.LockOSThread");
extern void runtime_UnlockOSThread(void) __asm__("runtime.UnlockOSThread");

// "Go" functions normally implemented by a Go program
extern void main_init(void) __asm__ ("__go_init_main");
extern void main_main(void) __asm__ ("main.main");

// The list of init functions that need to be called to initialise the gopy
// library, exported by pyext.c
extern ifunc py_init_funcs[];

/*****************************************************************************
******************************************************************************
**
** Static Variables
**
******************************************************************************
*****************************************************************************/

static __thread S *s;
static S *s0;

static ucontext_t gmc;

static Ctxt *ctxt_head = NULL;
static Ctxt *ctxt_tail = NULL;
static sem_t ctxt_sem;
static pthread_mutex_t ctxt_mutex = PTHREAD_MUTEX_INITIALIZER;

static volatile int base_init_done = 0;

static struct ie *init_done = NULL;

/*****************************************************************************
******************************************************************************
**
** Interface Functions
**
******************************************************************************
*****************************************************************************/

// This function stores the currently registered signal handlers into the given
// list.
static void store_signal_handlers(struct sigaction handlers[_NSIG]) {
    int i;
    for (i = 0; i < _NSIG; i++) {
        sigaction(i, NULL, &handlers[i]);
    }
}

// This function sets the signal handlers to be the ones saved in the given
// list.
static void restore_signal_handlers(struct sigaction handlers[_NSIG]) {
    int i;
    for (i = 0; i < _NSIG; i++) {
        sigaction(i, &handlers[i], NULL);
    }
}

// This function is called from inside the Go runtime to receive a context
// structure
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

// This function is called from outside the Go runtime to send a context
// structure
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

// This function sends a context to the dispatcher runing inside the Go runtime
// (using the send_ctxt/recv_ctxt functions), asking the handler to run fn(arg)
// in syscall context and/or gfn(garg) in goroutine context.
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

// This function is the second entry/exit point for the context jumping in the
// main thread.  It is called by run_on_g0 and run_on_g to jump into the Go
// runtime when they are called on the main thread.
static void activate_go(void (*fn)(void*), void *arg,
                        void (*gfn)(void*), void *garg) {
    struct sigaction handlers[_NSIG];

    // Setup state to jump into Go
    s->in_go = 1;
    s->fn    = fn;
    s->arg   = arg;
    s->gfn   = gfn;
    s->garg  = garg;

    // Swap over to the Go signal handlers
    store_signal_handlers(handlers);
    restore_signal_handlers(s0->handlers);

    // Setup the return point
    __splitstack_getcontext(&s->stack[0]);
    getcontext(&s->c);

    if (s->in_go) {
        // Jump into Go context
        __splitstack_setcontext(&s->g0_stack[0]);
        setcontext(&s->g0c);
    }

    // Swap back to our own signal handlers
    restore_signal_handlers(handlers);
}

// This function arranges for the given function to be run inside the Go runtime
// in syscall context - i.e. it will be called from g0_main.
static void run_on_g0(void (*f)(void*), void *p) {
    if (!s) return run_on_ctxt(f, p, NULL, NULL);
    activate_go(f, p, NULL, NULL);
}

// This function arranges for the given function to be run inside the Go runtime
// in normal goroutine context - i.e. it will be called from either ctxtHandler,
// or main_main.
static void run_on_g(void (*f)(void*), void *p) {
    if (!s) return run_on_ctxt(NULL, NULL, f, p);
    activate_go(NULL, NULL, f, p);
}

/*****************************************************************************
******************************************************************************
**
** Run Inside Go Runtime
**
******************************************************************************
*****************************************************************************/

// This function is run in "syscall" context, and comprises half of the code
// running, context jumping loop.
static void g0_main(void *_) {
    // Setup return entry point
    __splitstack_getcontext(&s->g0_stack[0]);
    getcontext(&s->g0c);

    // in_go + gfn == return to call gfn from goroutine
    if (s->in_go && s->gfn) return;

    // in_go + fn == we have come back to run some go code on g0 ...
    if (s->in_go && s->fn) s->fn(s->arg);

    // We are now leaving Go context
    s->in_go = 0;

    // jump back to whatever non-go code got us here
    __splitstack_setcontext(&s->stack[0]);
    setcontext(&s->c);
}

// The main function for a proxy thread.  This function runs in normal goroutine
// context.  It forms part of a code running, context jumping loop with g0_main
// - but it also initialises that loop with the functions in the provided
// context.
static void ctxtHandler(void *arg) {
    Ctxt *ctxt = arg;

    // We are using TLS storage, so we need to keep this goroutine on the same
    // thread.
    runtime_LockOSThread();

    // From here to exitsyscall is conceptually a simple_cgocall to a setup
    // function for the code running loop - except that then we would be jumping
    // between two functions trying to use the same stack ...
    runtime_entersyscall();

    // Create a new state object
    s = calloc(1, sizeof(S));

    // Swap over to the Go signal handlers
    store_signal_handlers(s->handlers);
    restore_signal_handlers(s0->handlers);

    // we are about to go into go
    s->in_go = 1;

    // Copy fn/gfn from ctxt to s
    s->fn   = ctxt->fn;
    s->arg  = ctxt->arg;
    s->gfn  = ctxt->gfn;
    s->garg = ctxt->garg;

    // Store the current context, this is where we will jump back to
    __splitstack_getcontext(&s->stack[0]);
    getcontext(&s->c);

    // This exitsyscall marks the end of the conceptual simple_cgocall to a
    // setup function
    runtime_exitsyscall();

    // !in_go means we are done
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
        // We are now leaving Go context
        s->in_go = 0;
    }
}

// This function runs in a goroutine, receives contexts from outside the Go
// runtime, and fires off ctxtHandler in a new goroutine to handle each context.
// In go it would look something like:
//
// for {
//     go ctxtHandler(<- ctxtChannel)
// }
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

// This function calls the given init function, unless it has already been
// called - which it tracks using the init_done list.
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

// This function is called by _init_go to call the init functions in funcs
static void do_inits(ifunc funcs[]) {
    int i = 0;
    while (funcs[i]) {
        do_init(funcs[i++]);
    }
}

// This function is just to unbox the arguments to simple_cgocallback and call
// it, but it does so from Go runtime syscall context
static void cgocallback_g(void *_a) {
    struct {
        void (*fn)(void*);
        void *arg;
    } *a = _a;
    simple_cgocallback(a->fn, a->arg);
}

/*****************************************************************************
******************************************************************************
**
** Go "main" Package Functions
**
******************************************************************************
*****************************************************************************/

// This function is the Go "main" function from the main package that is called
// by the Runtime at then end of startup.  It forms a code running, context
// jumping loop with g0_main in the special Go runtime context created on the
// "main" thread.
extern void main_main(void) {
    __go_go(ctxtDispatcher, NULL);
    while (1) {
        simple_cgocall(g0_main, NULL);
        // g0_main only returns when we want to run code in a goroutine ...
        s->gfn(s->garg);
        // We are now leaving Go context
        s->in_go = 0;
    }
}

// This function is the equivalent of the init() functions from the main
// package, it is run as part of the Go runtime initialisation.
extern void main_init(void) {
    // We need to make sure that we don't get switched off to another thread,
    // otherwise Python will get very confused when we return from a function
    // call on a different thread to the one that called it!
    runtime_LockOSThread();

    do_inits(py_init_funcs);
}

/*****************************************************************************
******************************************************************************
**
** Callback Entry Point
**
******************************************************************************
*****************************************************************************/

// This function either calls simple_cgocallback directly if already inside the
// Go runtime, other wise it arranges for it to be called from Go syscall
// context via run_on_g0 and cgocallback_g.
static void cgocallback_wrapper(void (*fn)(void*), void *param) {
    struct {
        void (*fn)(void*);
        void *arg;
    } a;
    if (s && s->in_go) return simple_cgocallback(fn, param);
    a.fn = fn;
    a.arg = param;
    run_on_g0(cgocallback_g, &a);
}

/*****************************************************************************
******************************************************************************
**
** Go Runtime Startup Functions
**
******************************************************************************
*****************************************************************************/

// This function is called as part of the Go runtime startup, it is called by
// go_main as the function to run on the first goroutine - all it has to do is
// call runtime_main().
static void mainstart(void *arg __attribute__((unused))) {
    runtime_main();
}

// This function starts the Go runtime - in a normal Go program this would be
// the C main function.  Here it is the first function run in a specially
// created context.
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

// This function prepares the new context that the Go runtime will run in, keeps
// a copy of the signal handlers setup by Go, and provides one of the entry/exit
// points for the context jumping in the main thead's loop.
static void base_init(void) {
    size_t ss;
    struct sigaction handlers[_NSIG];

    // We need to replace the cgocallback used in gopy with our wrapper - so
    // that we can get into the Go runtime.
    cgocallback = cgocallback_wrapper;

    // we need to save all the signal handlers, so we can stop Go from co-opting
    // them.
    store_signal_handlers(handlers);

    // Create a new state object
    s = calloc(1, sizeof(S));

    // we can't start the go runtime from the current context, so we need to
    // create a new one ...
    getcontext(&gmc);
    gmc.uc_stack.ss_sp = malloc(2 * 1024 * 1024);
    gmc.uc_stack.ss_size = 2 * 1024 * 1024;
    gmc.uc_link = NULL;
    makecontext(&gmc, go_main, 0);
    __splitstack_makecontext(gmc.uc_stack.ss_size, &s->gm_stack[0], &ss);

    // we are about to go into go
    s->in_go = 1;

    // once we have started the go runtime in it's own context, we need to be
    // able to get back to this one to return from this function
    __splitstack_getcontext(&s->stack[0]);
    getcontext(&s->c);

    if (s->in_go) {
        // actually jump into go
        __splitstack_setcontext(&s->gm_stack[0]);
        setcontext(&gmc);
    }

    // save the signal handlers that Go is using, so that we can reactivate them
    // when switching into Go context ...
    store_signal_handlers(s->handlers);

    // store the base state so that all threads can get at the signal handlers
    s0 = s;

    // we should restore the signal handlers now.
    restore_signal_handlers(handlers);

    // the runtime is now up and running in it's own context, and we can access
    // it using the g0c context - time to return
    base_init_done = 1;
}

/*****************************************************************************
******************************************************************************
**
** External Interface
**
******************************************************************************
*****************************************************************************/

// This is the function called by generated initXXX CPython extension module
// init function.  funcs is the list of init functions that need to be called
// for the module in question.
extern void _init_go(ifunc funcs[]) {
    int i = 0;

    if (!base_init_done) base_init();

    run_on_g((void (*)(void*))do_inits, funcs);
}

// This function initialises the Go runtime without creating a new context, and
// then calls the given list of init funcs.  This functions is used to make a Go
// executable that can load extension modules.
extern int _init_go_main(int argc, char *argv[], ifunc funcs[]) {
    // We need to replace the cgocallback used in gopy with our wrapper - so
    // that we can call back into the Go runtime from Python threads.
    cgocallback = cgocallback_wrapper;

    // Create a new state object
    s = calloc(1, sizeof(S));

    // we are about to go into go
    s->in_go = 1;

    // Once in Go, call do_inits in a goroutine, with funcs as the argument
    s->gfn  = (void (*)(void *))do_inits;
    s->garg = funcs;

    // we set the base_init_done flag to stop _init_go from trying to start the
    // Go runtime for a second time.
    base_init_done = 1;

    // once we have started the go runtime, we can only get back by jumping back
    // to this context (since that is the way the main code loop works).
    __splitstack_getcontext(&s->stack[0]);
    getcontext(&s->c);

    // !in_go means we are back from Go - exit.
    if (!s->in_go) exit(0);

    // Initialise the Go runtime, which will take over ownership of this thread.
    runtime_check();
    runtime_args(argc, argv);
    runtime_osinit();
    runtime_schedinit();
    __go_go(mainstart, NULL);
    runtime_mstart(runtime_m());

    // We should never get here - force a SIGSEGV.
    while (1) *(int *)0 = 0;
}
