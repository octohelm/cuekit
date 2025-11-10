module github.com/octohelm/cuekit

go 1.25.4

tool github.com/octohelm/cuekit/internal/cmd/fmt

// +gengo:import:group=0_controlled
require (
	github.com/octohelm/gengo v0.0.0-20251031090256-3806c0e4d0f8
	github.com/octohelm/unifs v0.0.0-20251110023907-a933885d117e
	github.com/octohelm/x v0.0.0-20251028032356-02d7b8d1c824
)

replace cuelang.org/go => github.com/morlay/cue v0.15.1-0.20251107024512-675234fd06ac

require (
	cuelang.org/go v0.15.0
	github.com/fatih/color v1.18.0
	github.com/go-json-experiment/json v0.0.0-20251027170946-4849db3c2f7e
	github.com/gobwas/glob v0.2.3
	golang.org/x/mod v0.29.0
	golang.org/x/sync v0.18.0
	golang.org/x/telemetry v0.0.0-20251106204201-5cc343d36a5d
	golang.org/x/tools v0.38.0
)

require (
	cuelabs.dev/go/oci/ociregistry v0.0.0-20250722084951-074d06050084 // indirect
	github.com/cockroachdb/apd/v3 v3.2.1 // indirect
	github.com/emicklei/proto v1.14.2 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/octohelm/courier v0.0.0-20251010073531-57524a0631a3 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/protocolbuffers/txtpbfmt v0.0.0-20251016062345-16587c79cd91 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/oauth2 v0.33.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	mvdan.cc/gofumpt v0.9.2 // indirect
)
