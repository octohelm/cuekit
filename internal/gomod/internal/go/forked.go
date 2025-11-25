//go:fork cmd/go/internal/modload
//go:fork cmd/go/internal/modfetch
//go:fork internal/goroot
//go:fork:replace internal/syscall/windows syscall/windows
//go:fork:replace internal/singleflight    golang.org/x/sync/singleflight
//go:fork:replace cmd/go/internal/cfg      github.com/octohelm/cuekit/internal/gomod/internal/cfg
//go:fork:replace internal/godebug         github.com/octohelm/cuekit/internal/gomod/internal/cfg/godebug
//go:fork:replace cmd/go/internal/fips140  github.com/octohelm/cuekit/internal/gomod/internal/cfg/fips140
package forked
