test:
    go test -v --failfast ./pkg/...

test-race:
    go test -v --race --failfast ./pkg/...

fmt:
    go tool gofumpt -w -l ./pkg/

update:
    go get -u ./...

dep:
    go mod tidy

clean:
    rm -rf ~/Library/Caches/cue/github.com/octohelm/cuemod-versioned-example*

internal_fork := "go run ./internal/cmd/internalfork"

fork-go-internal:
    {{ internal_fork }} \
    	-p cmd/go/internal/modload \
    	-p cmd/go/internal/modfetch \
    	-p internal/godebug \
    	./internal/gomod/internal
