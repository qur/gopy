package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

type ConfigFlag int

const (
	defaultFlag ConfigFlag = iota
	Enabled
	Disabled
)

func (f *ConfigFlag) Set(v bool) {
	if v {
		*f = Enabled
	} else {
		*f = Disabled
	}
}

func (f ConfigFlag) set(v *C.int) {
	switch f {
	case defaultFlag:
	case Enabled:
		*v = 1
	case Disabled:
		*v = 0
	default:
		panic(fmt.Sprintf("invalid ConfigFlag value: %d", f))
	}
}

type AllocatorMode int

const (
	AllocatorNotSet        AllocatorMode = C.PYMEM_ALLOCATOR_NOT_SET
	AllocatorDefault       AllocatorMode = C.PYMEM_ALLOCATOR_DEFAULT
	AllocatorDebug         AllocatorMode = C.PYMEM_ALLOCATOR_DEBUG
	AllocatorMalloc        AllocatorMode = C.PYMEM_ALLOCATOR_MALLOC
	AllocatorMallocDebug   AllocatorMode = C.PYMEM_ALLOCATOR_MALLOC_DEBUG
	AllocatorPyMalloc      AllocatorMode = C.PYMEM_ALLOCATOR_PYMALLOC
	AllocatorPyMallocDebug AllocatorMode = C.PYMEM_ALLOCATOR_PYMALLOC_DEBUG
	AllocatorMiMalloc      AllocatorMode = C.PYMEM_ALLOCATOR_MIMALLOC
	AllocatorMiMallocDebug AllocatorMode = C.PYMEM_ALLOCATOR_MIMALLOC_DEBUG
)

func (f AllocatorMode) set(v *C.int) {
	*v = C.int(f)
}

type PreConfig struct {
	Args              []string
	Allocator         AllocatorMode
	ConfigureLocale   ConfigFlag
	CoerceCLocale     ConfigFlag
	CoerceCLocaleWarn ConfigFlag
	DevMode           ConfigFlag
	Isolated          ConfigFlag
	// LegacyWindowsFSEncoding ConfigFlag - TODO: only available on Windows ...
	ParseArgv      ConfigFlag
	UseEnvironment ConfigFlag
	UTF8Mode       ConfigFlag
}

func PythonPreConfig(args ...string) *PreConfig {
	return &PreConfig{
		Args:     args,
		Isolated: Disabled,
	}
}

func IsolatedPreConfig() *PreConfig {
	return &PreConfig{
		Isolated: Enabled,
	}
}

func (c *PreConfig) PreInitialize() error {
	cfg := C.PyPreConfig{}
	if c.Isolated == Enabled {
		C.PyPreConfig_InitIsolatedConfig(&cfg)
	} else {
		C.PyPreConfig_InitPythonConfig(&cfg)
	}

	c.Allocator.set(&cfg.allocator)
	c.ConfigureLocale.set(&cfg.configure_locale)
	c.CoerceCLocale.set(&cfg.coerce_c_locale)
	c.CoerceCLocaleWarn.set(&cfg.coerce_c_locale_warn)
	c.DevMode.set(&cfg.dev_mode)
	c.Isolated.set(&cfg.isolated)
	// c.LegacyWindowsFSEncoding.set(&cfg.legacy_windows_fs_encoding)
	c.ParseArgv.set(&cfg.parse_argv)
	c.UseEnvironment.set(&cfg.use_environment)
	c.UTF8Mode.set(&cfg.utf8_mode)

	if c.Args == nil {
		return status2Err(C.Py_PreInitialize(&cfg))
	}

	argv := make([]*C.char, len(c.Args))

	for i, arg := range c.Args {
		argv[i] = C.CString(arg)
		defer C.free(unsafe.Pointer(argv[i]))
	}

	return status2Err(C.Py_PreInitializeFromBytesArgs(&cfg, C.Py_ssize_t(len(argv)), &argv[0]))
}

func setArgs(args []string, cfg *C.PyConfig) error {
	argv := make([]*C.char, len(args))

	for i, arg := range args {
		argv[i] = C.CString(arg)
		defer C.free(unsafe.Pointer(argv[i]))
	}

	return status2Err(C.PyConfig_SetBytesArgv(cfg, C.Py_ssize_t(len(argv)), &argv[0]))
}

type Config struct {
	Args     []string
	SafePath ConfigFlag
	// base_exec_prefix
	// base_executable
	// base_prefix
	BufferedStdio       ConfigFlag
	BytesWarning        ConfigFlag
	WarnDefaultEncoding ConfigFlag
	CodeDebugRanges     ConfigFlag
	// check_hash_pycs_mode
	ConfigureCStdio ConfigFlag
	DevMode         ConfigFlag
	DumpRefs        ConfigFlag
	// exec_prefix
	// executable
	FaultHandler ConfigFlag
	// filesystem_encoding
	// filesystem_errors
	HashSeed    uint64
	UseHashSeed ConfigFlag
	// home
	ImportTime            ConfigFlag
	Inspect               ConfigFlag
	InstallSignalHandlers ConfigFlag
	Interactive           ConfigFlag
	// int_max_str_digits
	// cpu_count
	Isolated ConfigFlag
	// LegacyWindowStdio ConfigFlag - TODO: Windows only
	MallocStats ConfigFlag
	// platlibdir
	// pythonpath_env
	// module_search_paths
	ModuleSearchPathSet ConfigFlag
	OptimizationLevel   int
	// orig_argv
	ParseArgv          ConfigFlag
	ParserDebug        ConfigFlag
	PathConfigWarnings ConfigFlag
	// prefix
	// program_name
	// pycache_prefix
	Quiet ConfigFlag
	// run_command
	// run_filename
	// run_module
	// run_presite
	ShowRefCount        ConfigFlag
	SiteImport          ConfigFlag
	SkipSourceFirstTime ConfigFlag
	// stdio_encoding
	TraceMalloc       ConfigFlag
	PerfProfiling     ConfigFlag
	UseEnvironment    ConfigFlag
	UserSiteDirectory ConfigFlag
	Verbose           ConfigFlag // TODO: this isn't actually a bool flag
	// warnoptions
	WriteBytecode ConfigFlag
	// xoptions
}

func PythonConfig(args ...string) *Config {
	return &Config{
		Args:     args,
		Isolated: Disabled,
	}
}

func IsolatedConfig() *Config {
	return &Config{
		Isolated: Enabled,
	}
}

func (c *Config) Initialize() error {
	cfg := C.PyConfig{}
	if c.Isolated == Enabled {
		C.PyConfig_InitIsolatedConfig(&cfg)
	} else {
		C.PyConfig_InitPythonConfig(&cfg)
	}
	defer C.PyConfig_Clear(&cfg)

	if c.Args != nil {
		if err := setArgs(c.Args, &cfg); err != nil {
			return err
		}
	}

	c.SafePath.set(&cfg.safe_path)
	c.BufferedStdio.set(&cfg.buffered_stdio)
	c.BytesWarning.set(&cfg.bytes_warning)
	c.WarnDefaultEncoding.set(&cfg.warn_default_encoding)
	c.CodeDebugRanges.set(&cfg.code_debug_ranges)
	c.ConfigureCStdio.set(&cfg.configure_c_stdio)
	c.DevMode.set(&cfg.dev_mode)
	c.DumpRefs.set(&cfg.dump_refs)
	c.FaultHandler.set(&cfg.faulthandler)
	if c.HashSeed > 0 {
		cfg.hash_seed = C.ulong(c.HashSeed)
	}
	c.UseHashSeed.set(&cfg.use_hash_seed)
	c.ImportTime.set(&cfg.import_time)
	c.Inspect.set(&cfg.inspect)
	c.InstallSignalHandlers.set(&cfg.install_signal_handlers)
	c.Interactive.set(&cfg.interactive)
	c.Isolated.set(&cfg.isolated)
	// c.LegacyWindowStdio.set(&cfg.legacy_windows_stdio) - TODO: windows
	c.MallocStats.set(&cfg.malloc_stats)
	c.ModuleSearchPathSet.set(&cfg.module_search_paths_set)
	if c.OptimizationLevel > 0 {
		cfg.optimization_level = C.int(c.OptimizationLevel)
	}
	c.ParseArgv.set(&cfg.parse_argv)
	c.ParserDebug.set(&cfg.parser_debug)
	c.PathConfigWarnings.set(&cfg.pathconfig_warnings)
	c.Quiet.set(&cfg.quiet)
	c.ShowRefCount.set(&cfg.show_ref_count)
	c.SiteImport.set(&cfg.site_import)
	c.SkipSourceFirstTime.set(&cfg.skip_source_first_line)
	c.TraceMalloc.set(&cfg.tracemalloc)
	c.PerfProfiling.set(&cfg.perf_profiling)
	c.UseEnvironment.set(&cfg.use_environment)
	c.UserSiteDirectory.set(&cfg.user_site_directory)
	c.Verbose.set(&cfg.verbose)
	c.WriteBytecode.set(&cfg.write_bytecode)

	return status2Err(C.Py_InitializeFromConfig(&cfg))
}

func (c *Config) InitAndLock() (*Lock, error) {
	// Lock the current goroutine to the current OS thread, until we have
	// released the GIL (as CPython uses per-thread state)
	runtime.LockOSThread()

	// Initialize the default Python interpreter
	if err := c.Initialize(); err != nil {
		runtime.UnlockOSThread()
		return nil, err
	}

	// Immediately release the GIL (and thus "deactivate" any per-thread state
	// associated with the current thread
	C.PyEval_SaveThread()

	// We can now unlock the current goroutine from the current OS thread, as
	// there is no active per-thread state
	runtime.UnlockOSThread()

	// Now that Python is setup, we can return a locked Lock, ready for the
	// calling code to use
	return NewLock(), nil
}
