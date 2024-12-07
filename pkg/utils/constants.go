package utils

const (
	// Application name for registration with keepassxc http api.
	ApplicationName = "keepassxc-http-tools-go"
	// StringFields need this Prefix (incl. at least one space) to be returned by keepassxc http api.
	StringFieldKeyPrefix = "KPH: "
	// File name of the socket file of keepassxc http api.
	SocketFileName = "org.keepassxc.KeePassXC.BrowserServer"
	// Environment variable at Windows for user name.
	WindowsEnvVarUsername = "USERNAME"
	// Environment variable at Darwin for temporary files.
	DarwinEnvTmpDir = "TMPDIR"
	// Environment variable at Linux for runtime files.
	LinuxEnvXdgRuntimeDir = "XDG_RUNTIME_DIR"
	// Environment variable at Linux for temporary files.
	LinuxEnvTmpDir = "TMPDIR"
	// Linux temporary files default dir.
	LinuxEnvTmpDirDefault = "/tmp"
	// Linux snaps subdir in user home.
	LinuxSnapCommonSubdir = "snap/keepassxc/common/"
	// URL string for entries to be found by this tool.
	ScriptIndicatorUrl = "script://keepassxc.go"
)
