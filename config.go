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

func (f ConfigFlag) apply(v *C.int) {
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

func (f AllocatorMode) apply(v *C.int) {
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

	c.Allocator.apply(&cfg.allocator)
	c.ConfigureLocale.apply(&cfg.configure_locale)
	c.CoerceCLocale.apply(&cfg.coerce_c_locale)
	c.CoerceCLocaleWarn.apply(&cfg.coerce_c_locale_warn)
	c.DevMode.apply(&cfg.dev_mode)
	c.Isolated.apply(&cfg.isolated)
	// c.LegacyWindowsFSEncoding.set(&cfg.legacy_windows_fs_encoding)
	c.ParseArgv.apply(&cfg.parse_argv)
	c.UseEnvironment.apply(&cfg.use_environment)
	c.UTF8Mode.apply(&cfg.utf8_mode)

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

func setString(s string, cfg *C.PyConfig, target **C.wchar_t) error {
	if len(s) == 0 {
		return nil
	}
	src := C.CString(s)
	defer C.free(unsafe.Pointer(src))
	return status2Err(C.PyConfig_SetBytesString(cfg, target, src))
}

func setStringList(s []string, cfg *C.PyConfig, target *C.PyWideStringList) error {
	if len(s) == 0 {
		return nil
	}

	items := (**C.wchar_t)(C.PyMem_RawMalloc(C.size_t(len(s) * int(unsafe.Sizeof((*C.wchar_t)(nil))))))
	src := unsafe.Slice(items, len(s))

	for i, arg := range s {
		cstr := C.CString(arg)
		defer C.free(unsafe.Pointer(cstr))
		src[i] = nil // PyConfig_SetBytesString will free the "old" entry
		err := status2Err(C.PyConfig_SetBytesString(cfg, &src[i], cstr))
		if err != nil {
			// TODO: free memory for stringlist and strings
			return err
		}
	}

	target.length = C.Py_ssize_t(len(src))
	target.items = items

	return nil
}

type CheckHashPYCsMode string

const (
	AlwaysCheckHash CheckHashPYCsMode = "always"
	NeverCheckHash  CheckHashPYCsMode = "never"
)

type EncodingErrorHandler string

const (
	StrictErrorHandler          EncodingErrorHandler = "strict"
	SurrogateEscapeErrorHandler EncodingErrorHandler = "surrogateescape"
	SurrogatePassErrorHandler   EncodingErrorHandler = "surrogatepass"
)

type Config struct {
	Args                  []string
	SafePath              ConfigFlag
	BaseExecPrefix        string
	BaseExecutable        string
	BasePrefix            string
	BufferedStdio         ConfigFlag
	BytesWarning          ConfigFlag
	WarnDefaultEncoding   ConfigFlag
	CodeDebugRanges       ConfigFlag
	CheckHashPVCsMode     CheckHashPYCsMode
	ConfigureCStdio       ConfigFlag
	DevMode               ConfigFlag
	DumpRefs              ConfigFlag
	ExecPrefix            string
	Executable            string
	FaultHandler          ConfigFlag
	FilesystemEncoding    string
	FilesystemErrors      EncodingErrorHandler
	HashSeed              uint64
	UseHashSeed           ConfigFlag
	Home                  string
	ImportTime            ConfigFlag
	Inspect               ConfigFlag
	InstallSignalHandlers ConfigFlag
	Interactive           ConfigFlag
	// int_max_str_digits
	// cpu_count
	Isolated ConfigFlag
	// LegacyWindowStdio ConfigFlag - TODO: Windows only
	MallocStats         ConfigFlag
	PlatLibDir          string
	PythonPathEnv       string
	ModuleSearchPath    []string
	ModuleSearchPathSet ConfigFlag // TODO: just use nil vs {} for set?
	OptimizationLevel   int
	// orig_argv
	ParseArgv          ConfigFlag
	ParserDebug        ConfigFlag
	PathConfigWarnings ConfigFlag
	Prefix             string
	ProgramName        string
	PyCachePrefix      string
	Quiet              ConfigFlag
	RunCommand         string
	RunFilename        string
	RunModule          string
	// RunPreSite          string // TODO: debug build only
	ShowRefCount        ConfigFlag
	SiteImport          ConfigFlag
	SkipSourceFirstTime ConfigFlag
	StdioEncoding       string
	StdioErrors         EncodingErrorHandler
	TraceMalloc         ConfigFlag
	PerfProfiling       ConfigFlag
	UseEnvironment      ConfigFlag
	UserSiteDirectory   ConfigFlag
	Verbose             ConfigFlag // TODO: this isn't actually a bool flag
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

	c.SafePath.apply(&cfg.safe_path)
	setString(c.BaseExecPrefix, &cfg, &cfg.base_exec_prefix)
	setString(c.BaseExecutable, &cfg, &cfg.base_executable)
	setString(c.BasePrefix, &cfg, &cfg.base_prefix)
	c.BufferedStdio.apply(&cfg.buffered_stdio)
	c.BytesWarning.apply(&cfg.bytes_warning)
	c.WarnDefaultEncoding.apply(&cfg.warn_default_encoding)
	c.CodeDebugRanges.apply(&cfg.code_debug_ranges)
	setString(string(c.CheckHashPVCsMode), &cfg, &cfg.check_hash_pycs_mode)
	c.ConfigureCStdio.apply(&cfg.configure_c_stdio)
	c.DevMode.apply(&cfg.dev_mode)
	c.DumpRefs.apply(&cfg.dump_refs)
	setString(c.ExecPrefix, &cfg, &cfg.exec_prefix)
	setString(c.Executable, &cfg, &cfg.executable)
	c.FaultHandler.apply(&cfg.faulthandler)
	setString(c.FilesystemEncoding, &cfg, &cfg.filesystem_encoding)
	setString(string(c.FilesystemErrors), &cfg, &cfg.filesystem_errors)
	if c.HashSeed > 0 {
		cfg.hash_seed = C.ulong(c.HashSeed)
	}
	c.UseHashSeed.apply(&cfg.use_hash_seed)
	setString(c.Home, &cfg, &cfg.home)
	c.ImportTime.apply(&cfg.import_time)
	c.Inspect.apply(&cfg.inspect)
	c.InstallSignalHandlers.apply(&cfg.install_signal_handlers)
	c.Interactive.apply(&cfg.interactive)
	c.Isolated.apply(&cfg.isolated)
	// c.LegacyWindowStdio.set(&cfg.legacy_windows_stdio) - TODO: windows
	c.MallocStats.apply(&cfg.malloc_stats)
	setString(c.PlatLibDir, &cfg, &cfg.platlibdir)
	setString(c.PythonPathEnv, &cfg, &cfg.pythonpath_env)
	setStringList(c.ModuleSearchPath, &cfg, &cfg.module_search_paths)
	c.ModuleSearchPathSet.apply(&cfg.module_search_paths_set)
	if c.OptimizationLevel > 0 {
		cfg.optimization_level = C.int(c.OptimizationLevel)
	}
	c.ParseArgv.apply(&cfg.parse_argv)
	c.ParserDebug.apply(&cfg.parser_debug)
	c.PathConfigWarnings.apply(&cfg.pathconfig_warnings)
	setString(c.Prefix, &cfg, &cfg.prefix)
	setString(c.ProgramName, &cfg, &cfg.program_name)
	setString(c.PyCachePrefix, &cfg, &cfg.pycache_prefix)
	c.Quiet.apply(&cfg.quiet)
	setString(c.RunCommand, &cfg, &cfg.run_command)
	setString(c.RunFilename, &cfg, &cfg.run_filename)
	setString(c.RunModule, &cfg, &cfg.run_module)
	// setString(c.RunPreSite, &cfg, &cfg.run_presite)
	c.ShowRefCount.apply(&cfg.show_ref_count)
	c.SiteImport.apply(&cfg.site_import)
	c.SkipSourceFirstTime.apply(&cfg.skip_source_first_line)
	setString(c.StdioEncoding, &cfg, &cfg.stdio_encoding)
	setString(string(c.StdioErrors), &cfg, &cfg.stdio_errors)
	c.TraceMalloc.apply(&cfg.tracemalloc)
	c.PerfProfiling.apply(&cfg.perf_profiling)
	c.UseEnvironment.apply(&cfg.use_environment)
	c.UserSiteDirectory.apply(&cfg.user_site_directory)
	c.Verbose.apply(&cfg.verbose)
	c.WriteBytecode.apply(&cfg.write_bytecode)

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
