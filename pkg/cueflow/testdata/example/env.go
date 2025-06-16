package example

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"cuelang.org/go/cue"
	"github.com/octohelm/cuekit/pkg/cueconvert"
	"github.com/octohelm/cuekit/pkg/cueflow"
	"github.com/octohelm/cuekit/pkg/cueflow/task"
)

func init() {
	Registry.Register(&EnvInterface{})
}

// EnvInterface of client
type EnvInterface struct {
	// to avoid added ok
	task.Task `json:"-"`

	RequiredEnv map[string]string `json:"-"`
	OptionalEnv map[string]string `json:"-"`
}

var _ cueflow.CueValueUnmarshaler = &EnvInterface{}

func (ei *EnvInterface) UnmarshalCueValue(v cue.Value) error {
	i, err := v.Fields(cue.All())
	if err != nil {
		return err
	}

	ei.RequiredEnv = make(map[string]string)
	ei.OptionalEnv = make(map[string]string)

	for i.Next() {
		envKey := i.Selector().Unquoted()

		// skip task name
		if strings.HasPrefix(envKey, "$$") {
			continue
		}

		var envVar string

		if i.Selector().Type()&cue.RequiredConstraint != 0 {
			ei.RequiredEnv[envKey] = envVar
		} else {
			ei.OptionalEnv[envKey] = envVar
		}
	}

	return nil
}

var _ cueconvert.OutputValuer = EnvInterface{}

func (ei EnvInterface) OutputValues() map[string]any {
	values := map[string]any{}

	for k, v := range ei.RequiredEnv {
		values[k] = v
	}

	for k, v := range ei.OptionalEnv {
		values[k] = v
	}

	return values
}

func (ei *EnvInterface) Do(ctx context.Context) error {
	clientEnvs := getClientEnvs()

	for key := range ei.RequiredEnv {
		if envVar, ok := clientEnvs[key]; ok {
			ei.RequiredEnv[key] = envVar
		} else {
			return fmt.Errorf("env var %s is required, but not defined", key)
		}
	}

	for key := range ei.OptionalEnv {
		if envVar, ok := clientEnvs[key]; ok {
			ei.OptionalEnv[key] = envVar
		}
	}

	return nil
}

var getClientEnvs = sync.OnceValue(func() map[string]string {
	clientEnvs := map[string]string{}

	for _, i := range os.Environ() {
		parts := strings.SplitN(i, "=", 2)
		if len(parts) == 2 {
			clientEnvs[parts[0]] = parts[1]
		} else {
			clientEnvs[parts[0]] = ""
		}
	}

	return clientEnvs
})
