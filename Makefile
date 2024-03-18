test:
	go test -v --failfast ./pkg/...

fmt:
	goimports -w -l ./pkg

clean:
	rm -rf ~/Library/Caches/cue/github.com/octohelm/cuemod-versioned-example*


INTERNAL_FORK = go run ./internal/cmd/internalfork

fork.go.internal:
	$(INTERNAL_FORK) \
		-p cmd/go/internal/modload \
		-p cmd/go/internal/modfetch \
		-p internal/godebug \
		./internal/gomod/internal

fork.cue.internal:
	$(INTERNAL_FORK) \
		-p cuelang.org/go/internal/mod/modload \
		./internal/cue/internal

