package utils

const (
	// Application name for registration with keepassxc http api.
	ApplicationName = "keepassxc-http-tools-go"
	// Application name for the binary & config.
	ApplicationNameShort = "kpht"
	// Prefix for environment variables that can replace config values (results in "KPHT_")
	ConfigEnvPrefix = ApplicationNameShort
	// Default file name to look for in user's config directory.
	ConfigFileNameDefault = ApplicationNameShort + ".yaml"
	// Config key path for association name.
	ConfigKeypathAssocName = "assoc.name"
	// Config key path for association key, stored in base64.
	ConfigKeypathAssocKey = "assoc.key"
	// Config key path for the formatter settings to use to identify keepassxc entries.
	ConfigKeypathEntryIdentifier = "entryIdentifier"
	// Config key path for the formatter settings to select the field to copy.
	ConfigKeypathClipDefaultCopy = "clip.defaultCopy"
	// Config key path for the formatter settings override to select the field to copy for specific entries.
	ConfigKeypathClipCopy = "clip.copy"
	// StringFields need this Prefix (incl. at least one space) to be returned by keepassxc http api.
	StringFieldKeyPrefix = "KPH: "
	// File name of the socket file of keepassxc http api.
	SocketFileName = "org.keepassxc.KeePassXC.BrowserServer"
	// Environment variable at Windows for user name.
	WindowsEnvVarUsername = "USERNAME"
	// Environment variable at Darwin for temporary files.
	DarwinEnvTmpDir = "TMPDIR"
	// Environment variable at Linux for runtime files dir.
	LinuxEnvXdgRuntimeDir = "XDG_RUNTIME_DIR"
	// Environment variable at Linux for user's config files dir.
	LinuxEnvXdgConfigHome = "XDG_CONFIG_HOME"
	// Default subdir for user's config files.
	LinuxDefaultXdgConfigHomeSubdir = ".config"
	// Environment variable at Linux for temporary files.
	LinuxEnvTmpDir = "TMPDIR"
	// Linux temporary files default dir.
	LinuxEnvTmpDirDefault = "/tmp"
	// Linux snaps subdir in user home.
	LinuxSnapCommonSubdir = "snap/keepassxc/common/"
	// URL string for entries to be found by this tool.
	ScriptIndicatorUrl = "script://keepassxc.go"
)
