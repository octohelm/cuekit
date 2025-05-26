package x

import (
	"strconv"
	"strings"
	"regexp"
)

local: #Def & {
	cwd: string
}

def: #Def & {
	a: string
}

actions: {
	x: X = #Build & {
		cwd:     local.cwd
		version: def.a

		ldflags: [
			"-s", "-w",
			"-X",
			"\(X.module)/internal/version.version=\(X.version)"]
	}
}

#Build: X = {
	#Base

	_info: #Info & {
		target: wd: X.cwd
	}

	module: string | *_info.output.module

	_build: #Exec & {
		cmd: [
			"go", "build",
			"-ldflags", strconv.Quote(strings.Join(X.ldflags | *["-s", "-w"], " ")),
		]
	}

	output: _build.output
}

#Info: {
	target: #File & {
		filename: "x"
	}

	_read: #Def & {
		"cwd":      target
		"contents": string
	}

	output: #Def & {
		module:    regexp.FindSubmatch(#"module (.+)\n"#, _read.contents)[1]
		goversion: regexp.FindSubmatch(#"\ngo (.+)\n?"#, _read.contents)[1]
	}
}

#Base: {
	cwd:     string
	version: string
	ldflags?: [...string]
}

#Exec: "$$type": name: "Exec"
#Exec: "$$task": true
#Exec: {
	cmd: [...string]
	output: string @generated()
}

#Def: "$$type": name: "Def"
#Def: "$$task": true
#Def: {
	...
}

#File: $$type: name: "File"
#File: {
	wd!:       string
	filename!: string
}
