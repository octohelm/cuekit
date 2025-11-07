package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"

	cueformat "cuelang.org/go/cue/format"

	"github.com/octohelm/unifs/pkg/filesystem"
)

func Source(ctx context.Context, registry Registry) (filesystem.FileSystem, error) {
	sources := map[string]*source{}

	for t := range registry.Tasks() {
		nt := t.(*named)

		s, ok := sources[nt.decl.PkgPath]
		if !ok {
			s = &source{
				pkgName: filepath.Base(nt.decl.PkgPath),
				imports: map[string]string{},
			}
			sources[nt.decl.PkgPath] = s
		}

		if err := s.writeDecl(nt); err != nil {
			return nil, err
		}
	}

	fs := filesystem.NewMemFS()

	for pathPath, s := range sources {
		code, err := s.Source()
		if err != nil {
			return nil, err
		}

		if err := writeFile(ctx, fs, path.Join(pathPath, s.pkgName+".cue"), code); err != nil {
			return nil, err
		}
	}

	return fs, nil
}

func writeFile(ctx context.Context, fs filesystem.FileSystem, filename string, data []byte) error {
	if err := filesystem.MkdirAll(ctx, fs, filepath.Dir(filename)); err != nil {
		return err
	}
	file, err := fs.OpenFile(ctx, filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return err
	}
	return nil
}

type source struct {
	pkgName string
	imports map[string]string
	bytes.Buffer
}

func (s *source) writeDecl(nt *named) error {
	for k, v := range nt.decl.Imports {
		s.imports[k] = v
	}

	s.WriteString("\n")

	return nt.WriteDeclTo(s)
}

func (s *source) Source() ([]byte, error) {
	b := bytes.NewBufferString("package " + s.pkgName)

	if len(s.imports) > 0 {
		_, _ = fmt.Fprintf(b, `

import (
`)

		for _, p := range slices.Sorted(maps.Keys(s.imports)) {
			_, _ = fmt.Fprintf(b, `%s %q
`, s.imports[p], p)
		}

		_, _ = fmt.Fprintf(b, `)
`)
	}

	_, _ = io.Copy(b, s)

	data, err := cueformat.Source(b.Bytes(), cueformat.Simplify())
	if err != nil {
		return nil, fmt.Errorf(`format invalid: %w
%s`, err, b.Bytes())
	}
	return data, nil
}
