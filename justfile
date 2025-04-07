test:
    go test -v --failfast ./pkg/...

fmt:
    go tool gofumpt -w -l ./pkg/

clean:
    rm -rf ~/Library/Caches/cue/github.com/octohelm/cuemod-versioned-example*

internal_fork := "go run ./internal/cmd/internalfork"

fork-go-internal:
    {{ internal_fork }} \
    	-p cmd/go/internal/modload \
    	-p cmd/go/internal/modfetch \
    	-p internal/godebug \
    	./internal/gomod/internal

fork-cue-internal:
    {{ internal_fork }} \
    	-p cuelang.org/go/internal/mod/modload \
    	./internal/cue/internal
