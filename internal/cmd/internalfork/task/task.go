package task

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

type Task struct {
	pkg *packages.Package

	replaces map[string]string
	prefixes map[string]bool
	packages map[string]*packages.Package
	l        *slog.Logger
}

func (t *Task) Init(dir string) error {
	t.l = slog.Default()

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedTypes,
		Dir:  dir,
	}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return err
	}

	if len(pkgs) == 0 {
		return fmt.Errorf("no packages found in %s", dir)
	}

	pkg := pkgs[0]

	t.pkg = pkg

	for _, file := range pkg.Syntax {
		for _, commentGroup := range file.Comments {
			for _, c := range commentGroup.List {
				if strings.HasPrefix(c.Text, "//go:fork") {
					parts := strings.Fields(c.Text)
					switch parts[0] {
					case "//go:fork":
						entryParts := parts[1:]
						if len(entryParts) != 1 {
							return fmt.Errorf("invalid fork entries: %s", c.Text)
						}
						if t.prefixes == nil {
							t.prefixes = make(map[string]bool)
						}
						t.prefixes[entryParts[0]] = true
					case "//go:fork:replace":
						replaceParts := parts[1:]
						if len(replaceParts) != 2 {
							return fmt.Errorf("invalid fork entries: %s", c.Text)
						}
						if t.replaces == nil {
							t.replaces = make(map[string]string)
						}
						t.replaces[replaceParts[0]] = replaceParts[1]
					}
				}
			}
		}
	}

	return nil
}

func (t *Task) Commit() error {
	for _, pkg := range t.packages {
		if err := t.commit(pkg); err != nil {
			return err
		}
	}

	return nil
}

func (t *Task) commit(pkg *packages.Package) error {
	forkedPkgPath := t.replaceInternal(pkg.PkgPath)

	for _, file := range pkg.Syntax {
		filename := pkg.Fset.File(file.Pos()).Name()

		for _, i := range file.Imports {
			importPath, _ := strconv.Unquote(i.Path.Value)

			if replaced, ok := t.replaces[importPath]; ok {
				_ = astutil.RewriteImport(
					pkg.Fset, file,
					importPath, replaced,
				)
			} else if _, ok := t.packages[importPath]; ok {
				_ = astutil.RewriteImport(
					pkg.Fset, file,
					importPath, filepath.Join(t.pkg.PkgPath, t.replaceInternal(importPath)),
				)
			}
		}

		dest := filepath.Join(t.pkg.Dir, forkedPkgPath, filepath.Base(filename))
		buf := bytes.NewBuffer(nil)
		if err := format.Node(buf, pkg.Fset, file); err != nil {
			return err
		}
		if err := writeFile(dest, buf.Bytes()); err != nil {
			return err
		}
	}

	t.l.Info("commited", slog.String("forked", forkedPkgPath))

	return nil
}

func (t *Task) Scan() error {
	for importPath := range t.prefixes {
		if err := t.scan(importPath); err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) scan(importPath string) error {
	// skip not internal pkg
	if !t.isInternalPkg(importPath) {
		return nil
	}

	if _, ok := t.packages[importPath]; ok {
		return nil
	}

	if t.packages == nil {
		t.packages = make(map[string]*packages.Package)
	}

	// to avoid loop
	t.packages[importPath] = &packages.Package{}

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles,
	}

	pkgs, err := packages.Load(cfg, importPath)
	if err != nil {
		return err
	}

	if len(pkgs) == 0 {
		return fmt.Errorf("no packages found in %s", importPath)
	}

	pkg := pkgs[0]

	files, err := os.ReadDir(pkg.Dir)
	if err != nil {
		return fmt.Errorf("failed to read dir %s: %w", pkg.Dir, err)
	}

	pkg.Fset = token.NewFileSet()

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".go") || strings.HasSuffix(f.Name(), "_test.go") {
			continue
		}

		filePath := filepath.Join(pkg.Dir, f.Name())

		fileAST, err := parser.ParseFile(pkg.Fset, filePath, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("failed to parse file %s: %w", filePath, err)
		}

		pkg.Syntax = append(pkg.Syntax, fileAST)
	}

	for _, file := range pkg.Syntax {
		for _, spec := range file.Imports {
			ip, _ := strconv.Unquote(spec.Path.Value)
			if err := t.scan(ip); err != nil {
				return err
			}
		}
	}

	t.packages[importPath] = pkg

	t.l.Info("scanned", slog.String("pkg", pkg.PkgPath))

	return nil
}

func (t *Task) isInternalPkg(importPath string) bool {
	if t.replaces != nil {
		if _, ok := t.replaces[importPath]; ok {
			return false
		}
	}
	if strings.Contains(importPath, "internal/") {
		return true
	}
	return strings.HasSuffix(importPath, "internal") || strings.HasPrefix(importPath, "internal")
}

func (t *Task) replaceInternal(p string) string {
	if strings.HasSuffix(p, "internal") {
		return filepath.Join(filepath.Dir(p), "./internals")
	}

	return strings.Replace(p, "internal/", "internals/", -1)
}
