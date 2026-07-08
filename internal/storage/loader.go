package storage

import (
	"fmt"
	"io/fs"
	"path"

	"github.com/mateom/vaultsh/internal/filesystem"
)

func Load(source fs.FS) (*filesystem.Directory, error) {
	root := filesystem.NewDirectory("")
	directories := map[string]*filesystem.Directory{".": root}

	err := fs.WalkDir(source, ".", func(name string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if name == "." {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("stat file %s: %w", name, err)
		}

		parent, found := directories[path.Dir(name)]
		if !found {
			return fmt.Errorf("parent directory not loaded: %s", path.Dir(name))
		}

		if entry.IsDir() {
			directory := filesystem.NewDirectoryWithModTime(entry.Name(), info.ModTime())
			if err := parent.Add(directory); err != nil {
				return fmt.Errorf("add directory %s: %w", name, err)
			}
			directories[name] = directory
			return nil
		}

		content, err := fs.ReadFile(source, name)
		if err != nil {
			return fmt.Errorf("read file %s: %w", name, err)
		}
		file := filesystem.NewFileWithModTime(entry.Name(), string(content), info.ModTime())
		if err := parent.Add(file); err != nil {
			return fmt.Errorf("add file %s: %w", name, err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("load content: %w", err)
	}

	return root, nil
}
