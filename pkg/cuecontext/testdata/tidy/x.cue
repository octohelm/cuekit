package tidy

import (
	v0 "github.com/octohelm/cuemod-versioned-example/cuepkg"
	v2 "github.com/octohelm/cuemod-versioned-example/v2/cuepkg"
	"mem.octothelm.tech/x"
	"github.com/octohelm/kubepkg/cuepkg/kubepkg"
)

"kubepkg": kubepkg.#KubePkg & {}

deps: {
	"v0": "\(v0.#Version)"
	"v2": "\(v2.#Version)"
	"x": "\(x.#Version)"
}