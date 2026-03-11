module github.com/octohelm/cuekit

go 1.26.0

tool (
	github.com/octohelm/cuekit/internal/cmd/fmt
	github.com/octohelm/cuekit/internal/cmd/internalfork
)

replace cuelang.org/go => github.com/morlay/cue v0.16.1-0.20260311041243-6c51b2341218

// +gengo:import:group=0_controlled
require (
	github.com/octohelm/gengo v0.0.0-20260224022252-ec6c2fc2f701
	github.com/octohelm/unifs v0.0.0-20260224071133-c301ec0d3226
	github.com/octohelm/x v0.0.0-20260224043023-b48075b44477
)

require (
	cuelang.org/go v0.16.0
	github.com/fatih/color v1.18.0
	github.com/go-json-experiment/json v0.0.0-20260214004413-d219187c3433
	github.com/gobwas/glob v0.2.3
	golang.org/x/mod v0.33.0
	golang.org/x/sync v0.20.0
	golang.org/x/telemetry v0.0.0-20260306145045-e526e8a188f5
	golang.org/x/tools v0.42.0
)

require (
	cuelabs.dev/go/oci/ociregistry v0.0.0-20251212221603-3adeb8663819 // indirect
	github.com/cockroachdb/apd/v3 v3.2.1 // indirect
	github.com/emicklei/proto v1.14.3 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/octohelm/courier v0.0.0-20260224022830-37ae4d696763 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/protocolbuffers/txtpbfmt v0.0.0-20260217160748-a481f6a22f94 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	mvdan.cc/gofumpt v0.9.2 // indirect
)
