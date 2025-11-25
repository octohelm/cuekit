test:
    go test -count=1 -v --failfast ./internal/gomod
    go test -count=1 -v --failfast ./pkg/...

test-race:
    go test -count=1 -v --race --failfast ./pkg/...

fmt:
    go fix ./pkg/...
    go fix ./internal/gomod
    go tool fmt .

update:
    go get -u ./...

dep:
    go mod tidy

clean:

internal_fork := "go tool internalfork"

fork-go-internal:
    {{ internal_fork }} ./internal/gomod/internal/go
