package module

import (
	"cuelang.org/go/mod/module"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func OSDirFS(dir string) fs.FS {
	return module.OSDirFS(dir)
}

func BasePath(mpath string) string {
	base, _, ok := module.SplitPathVersion(mpath)
	if ok {
		return base
	}
	return mpath
}

func SplitPathVersion(path string) (prefix, version string, ok bool) {
	return module.SplitPathVersion(path)
}

func SourceLocOfOSDir(dir string) SourceLoc {
	return SourceLoc{
		FS:  OSDirFS(dir),
		Dir: ".",
	}
}

func WalkCueFile(root string, fromPath string) ([]string, error) {
	files := make([]string, 0)

	walkSubDir := strings.HasSuffix(fromPath, "/...")

	if walkSubDir {
		fromPath = fromPath[0 : len(fromPath)-4]
	}

	start := filepath.Join(root, fromPath)

	err := filepath.Walk(start, func(path string, info os.FileInfo, err error) error {
		if path == start {
			return nil
		}

		// skip cue.mod
		if isSubDirFor(path, "cue.mod") {
			return filepath.SkipDir
		}

		if info.IsDir() {
			// skip sub dir which is cuemod root
			if _, err := os.Stat(filepath.Join(path, fileModule)); err == nil {
				return filepath.SkipDir
			}
			if walkSubDir {
				return nil
			}
			return filepath.SkipDir
		}

		if filepath.Ext(path) == ".cue" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func isSubDirFor(targetPath string, root string) bool {
	targetPath = targetPath + "/"
	root = root + "/"
	return strings.HasPrefix(targetPath, root)
}
