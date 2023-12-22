package acid

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
)

type DigestFS struct {
	FS         fs.FS
	Map        map[string]string
	ReverseMap map[string]string
}

var assetsWithDigests *DigestFS

func FS() http.FileSystem {
	if assetsWithDigests == nil {
		panic("AssetsWithDigests not initialized")
	}
	return http.FS(assetsWithDigests)
}

func CalculateDigests(fsys fs.FS, path string) error {
	var err error
	assetsWithDigests, err = newDigestFS(fsys, path)
	return err
}

func newDigestFS(fsys fs.FS, prefix string) (*DigestFS, error) {
	cfs := &DigestFS{
		FS:         fsys,
		Map:        make(map[string]string),
		ReverseMap: make(map[string]string),
	}

	err := fs.WalkDir(fsys, prefix, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			newPath, err := hashedFilePath(fsys, path)
			if err != nil {
				return err
			}
			cfs.Map[newPath] = path
			cfs.ReverseMap[path] = newPath
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return cfs, nil
}

func (c *DigestFS) FindFile(name string) string {
	if c.Map == nil {
		panic("digestfs not initialized")
	}

	if newName, ok := c.Map[name]; ok {
		return newName
	}
	return name
}

func (c *DigestFS) Open(name string) (fs.File, error) {
	return c.FS.Open(c.FindFile(name))
}

func fileDigest(fsys fs.FS, filePath string) (string, error) {
	data, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

func hashedFilePath(fsys fs.FS, path string) (string, error) {
	hash, err := fileDigest(fsys, path)
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(path)
	ext := filepath.Ext(path)
	base := filepath.Base(path)
	nameWithoutExt := base[0 : len(base)-len(ext)]

	newName := fmt.Sprintf("%s-%s%s", nameWithoutExt, hash, ext)
	return filepath.Join(dir, newName), nil
}
