package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

// ConfigFlag is an int that represents a boolean config flag.
type ConfigFlag int

// The values for a ConfigFlag, Enabled maps to 1 in C Python, and Disabled maps
// to 0.
const (
	defaultFlag ConfigFlag = iota
	Enabled
	Disabled
)

// Enable sets the ConfigFlag to Enabled.
func (f *ConfigFlag) Enable() {
	*f = Enabled
}

// Disable sets the ConfigFlag to Disabled.
func (f *ConfigFlag) Disable() {
	*f = Disabled
}

// Set sets the ConfigFlag based on the value of v.
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

type IntMaxStrDigits int

const UnlimitedMaxStrDigits IntMaxStrDigits = -1

func (i *IntMaxStrDigits) SetUnlimited() {
	*i = UnlimitedMaxStrDigits
}

func (i *IntMaxStrDigits) Set(max int) error {
	if max < 640 {
		return fmt.Errorf("max digits must be >= 640")
	}
	*i = IntMaxStrDigits(max)
	return nil
}

func (i IntMaxStrDigits) apply(v *C.int) {
	if i == 0 {
		// no explicit value set, leave C version at default value
		return
	}
	if i == UnlimitedMaxStrDigits {
		// for the C version, unlimited means set to 0
		*v = 0
		return
	}
	// set other values directly
	*v = C.int(i)
}

// PreConfig is a structure used to preinitialize Python.
//
// It corresponds to PyPreConfig in C Python.
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

// PythonPreConfig returns a PreConfig struct that is ready to be used for
// pre-initialisation. Optionally with the provided command line arguments.
func PythonPreConfig(args ...string) *PreConfig {
	return &PreConfig{
		Args:     args,
		Isolated: Disabled,
	}
}

// IsolatedPreConfig returns a PreConfig struct that is ready to used for
// pre-initialisation in isolated mode.
func IsolatedPreConfig() *PreConfig {
	return &PreConfig{
		Isolated: Enabled,
	}
}

// PreInitialize performs Python pre-initialization based on the PreConfig
// settings.
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
	if s == nil {
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

// Config contains most of the parameters used to configure Python.
//
// It corresponds to PyConfig in C Python.
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
	IntMaxStrDigits       IntMaxStrDigits
	CpuCount              int
	Isolated              ConfigFlag
	// LegacyWindowStdio ConfigFlag - TODO: Windows only
	MallocStats       ConfigFlag
	PlatLibDir        string
	PythonPathEnv     string
	ModuleSearchPath  []string
	OptimizationLevel int
	// orig_argv - this doesn't seem useful?
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
	Verbose             int
	WarnOptions         []string
	WriteBytecode       ConfigFlag
	XOptions            []string
}

// PythonConfig returns a Config struct that is ready to be used for
// initialisation. Optionally with the provided command line arguments.
func PythonConfig(args ...string) *Config {
	return &Config{
		Args:     args,
		Isolated: Disabled,
	}
}

// IsolatedConfig returns a Config struct that is ready to used for
// initialisation in isolated mode.
func IsolatedConfig() *Config {
	return &Config{
		Isolated: Enabled,
	}
}

func (c *Config) initialize() error {
	cfg := C.PyConfig{}
	if c.Isolated == Enabled {
		C.PyConfig_InitIsolatedConfig(&cfg)
	} else {
		C.PyConfig_InitPythonConfig(&cfg)
	}
	defer C.PyConfig_Clear(&cfg)

	// These values need to be set before we call any function that could
	// trigger pre-initialisation.
	c.DevMode.apply(&cfg.dev_mode)
	c.Isolated.apply(&cfg.isolated)
	c.ParseArgv.apply(&cfg.parse_argv)
	c.UseEnvironment.apply(&cfg.use_environment)

	// Setting the arguments needs to be the next thing that we do. As if we
	// call PyConfig_SetBytesArgv then it should be called before other methods.
	if c.Args != nil {
		if err := setArgs(c.Args, &cfg); err != nil {
			return err
		}
	}

	// Now we set the rest of the values, in declaration order.
	c.SafePath.apply(&cfg.safe_path)
	if err := setString(c.BaseExecPrefix, &cfg, &cfg.base_exec_prefix); err != nil {
		return err
	}
	if err := setString(c.BaseExecutable, &cfg, &cfg.base_executable); err != nil {
		return err
	}
	if err := setString(c.BasePrefix, &cfg, &cfg.base_prefix); err != nil {
		return err
	}
	c.BufferedStdio.apply(&cfg.buffered_stdio)
	c.BytesWarning.apply(&cfg.bytes_warning)
	c.WarnDefaultEncoding.apply(&cfg.warn_default_encoding)
	c.CodeDebugRanges.apply(&cfg.code_debug_ranges)
	if err := setString(string(c.CheckHashPVCsMode), &cfg, &cfg.check_hash_pycs_mode); err != nil {
		return err
	}
	c.ConfigureCStdio.apply(&cfg.configure_c_stdio)
	c.DumpRefs.apply(&cfg.dump_refs)
	if err := setString(c.ExecPrefix, &cfg, &cfg.exec_prefix); err != nil {
		return err
	}
	if err := setString(c.Executable, &cfg, &cfg.executable); err != nil {
		return err
	}
	c.FaultHandler.apply(&cfg.faulthandler)
	if err := setString(c.FilesystemEncoding, &cfg, &cfg.filesystem_encoding); err != nil {
		return err
	}
	if err := setString(string(c.FilesystemErrors), &cfg, &cfg.filesystem_errors); err != nil {
		return err
	}
	if c.HashSeed > 0 {
		cfg.hash_seed = C.ulong(c.HashSeed)
	}
	c.UseHashSeed.apply(&cfg.use_hash_seed)
	if err := setString(c.Home, &cfg, &cfg.home); err != nil {
		return err
	}
	c.ImportTime.apply(&cfg.import_time)
	c.Inspect.apply(&cfg.inspect)
	c.InstallSignalHandlers.apply(&cfg.install_signal_handlers)
	c.Interactive.apply(&cfg.interactive)
	c.IntMaxStrDigits.apply(&cfg.int_max_str_digits)
	if c.CpuCount > 0 {
		cfg.cpu_count = C.int(c.CpuCount)
	}
	// c.LegacyWindowStdio.set(&cfg.legacy_windows_stdio) - TODO: windows
	c.MallocStats.apply(&cfg.malloc_stats)
	if err := setString(c.PlatLibDir, &cfg, &cfg.platlibdir); err != nil {
		return err
	}
	if err := setString(c.PythonPathEnv, &cfg, &cfg.pythonpath_env); err != nil {
		return err
	}
	if err := setStringList(c.ModuleSearchPath, &cfg, &cfg.module_search_paths); err != nil {
		return err
	}
	if c.ModuleSearchPath != nil {
		cfg.module_search_paths_set = 1
	}
	if c.OptimizationLevel > 0 {
		cfg.optimization_level = C.int(c.OptimizationLevel)
	}
	c.ParserDebug.apply(&cfg.parser_debug)
	c.PathConfigWarnings.apply(&cfg.pathconfig_warnings)
	if err := setString(c.Prefix, &cfg, &cfg.prefix); err != nil {
		return err
	}
	if err := setString(c.ProgramName, &cfg, &cfg.program_name); err != nil {
		return err
	}
	if err := setString(c.PyCachePrefix, &cfg, &cfg.pycache_prefix); err != nil {
		return err
	}
	c.Quiet.apply(&cfg.quiet)
	if err := setString(c.RunCommand, &cfg, &cfg.run_command); err != nil {
		return err
	}
	if err := setString(c.RunFilename, &cfg, &cfg.run_filename); err != nil {
		return err
	}
	if err := setString(c.RunModule, &cfg, &cfg.run_module); err != nil {
		return err
	}
	// if err := setString(c.RunPreSite, &cfg, &cfg.run_presite)
	c.ShowRefCount.apply(&cfg.show_ref_count)
	c.SiteImport.apply(&cfg.site_import)
	c.SkipSourceFirstTime.apply(&cfg.skip_source_first_line)
	if err := setString(c.StdioEncoding, &cfg, &cfg.stdio_encoding); err != nil {
		return err
	}
	if err := setString(string(c.StdioErrors), &cfg, &cfg.stdio_errors); err != nil {
		return err
	}
	c.TraceMalloc.apply(&cfg.tracemalloc)
	c.PerfProfiling.apply(&cfg.perf_profiling)
	c.UserSiteDirectory.apply(&cfg.user_site_directory)
	if c.Verbose > 0 {
		cfg.verbose = C.int(c.Verbose)
	}
	if err := setStringList(c.WarnOptions, &cfg, &cfg.warnoptions); err != nil {
		return err
	}
	c.WriteBytecode.apply(&cfg.write_bytecode)
	if err := setStringList(c.XOptions, &cfg, &cfg.xoptions); err != nil {
		return err
	}

	return status2Err(C.Py_InitializeFromConfig(&cfg))
}

// Initialize the Python interpreter using the settings in the Config struct.
//
// You probably want InitAndLock, as it doesn't require the caller to worry
// about goroutines or threads.
func (c *Config) Initialize() error {
	if err := c.initialize(); err != nil {
		return err
	}

	if err := setupImporter(); err != nil {
		return fmt.Errorf("failed to setup importer: %s", err)
	}

	return nil
}

// Initialize the Python interpreter using the settings in the Config struct.
// Returns a Lock

// InitAndLock is a convenience function.  It initializes Python, enables thread
// support, and returns a locked Lock instance.
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
