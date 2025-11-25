package main

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/octohelm/cuekit/internal/cmd/internalfork/task"
)

func main() {
	flag.Parse()

	t := &task.Task{}
	if err := t.Init(flag.Args()[0]); err != nil {
		slog.Default().Error(fmt.Sprintf("failed to init: %s", err))
		return
	}

	if err := t.Scan(); err != nil {
		slog.Default().Error(fmt.Sprintf("failed to scan: %s", err))
		return
	}

	if err := t.Commit(); err != nil {
		slog.Default().Error(fmt.Sprintf("failed to commit: %s", err))
		return
	}
}
