package cfg

import (
	"runtime"
)

func ToolExeSuffix() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}
