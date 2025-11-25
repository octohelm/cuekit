package cfg

import (
	"os"
	"path/filepath"
	"sync"

	syncx "github.com/octohelm/x/sync"
)

var OrigEnv []string

var (
	GOROOT string

	GOROOTbin string
	GOROOTsrc string
	GOBIN     string

	GOMODCACHE string

	GOCACHEPROG = envOr("GOCACHEPROG", "")
	GOFIPS140   = envOr("GOFIPS140", "")
	GOPROXY     = envOr("GOPROXY", "https://proxy.golang.org,direct")
	GOSUMDB     = envOr("GOSUMDB", "sum.golang.org")
	GOPRIVATE   = envOr("GOPRIVATE", "")
	GONOPROXY   = envOr("GONOPROXY", GOPRIVATE)
	GONOSUMDB   = envOr("GONOSUMDB", GOPRIVATE)
	GOINSECURE  = envOr("GOINSECURE", "")
	GOAUTH      = envOr("GOAUTH", "netrc")
)

func init() {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	GOMODCACHE = filepath.Join(userCacheDir, "cue/modcache")
}

func envOr(key, def string) string {
	val := Getenv(key)
	if val == "" {
		val = def
	}
	return val
}

var c = syncx.Map[string, func() string]{}

func Getenv(key string) string {
	get, _ := c.LoadOrStore(key, sync.OnceValue(func() string {
		return os.Getenv(key)
	}))
	return get()
}
