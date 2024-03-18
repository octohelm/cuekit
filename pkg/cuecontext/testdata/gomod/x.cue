package gomod

import (
	v0 "github.com/octohelm/cuemod-versioned-example/cuepkg"
	v2 "github.com/octohelm/cuemod-versioned-example/v2/cuepkg"
)

name: "hello"

deps: {
	"v0": "\(v0.#Version)"
	"v2": "\(v2.#Version)"
}