package fsadapter

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/afero"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/helper/chroot"
)

const (
	defaultDirectoryMode = 0755
	defaultCreateMode    = 0666
)

// AdapterFs holds an afero Fs interface for adaptation to billy.Filesystem
type AdapterFs struct {
	fs afero.Fs
}

func New(fs afero.Fs) billy.Filesystem {
	return chroot.New(&AdapterFs{fs}, "/")
}

func (fs *AdapterFs) Create(filename string) (billy.File, error) {
	return fs.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, defaultCreateMode)
}

func (fs *AdapterFs) OpenFile(filename string, flag int, perm os.FileMode) (billy.File, error) {
	if flag&os.O_CREATE != 0 {
		if err := fs.createDir(filename); err != nil {
			return nil, err
		}
	}

	f, err := fs.fs.OpenFile(filename, flag, perm)
	if err != nil {
		return nil, err
	}

	mutexFile := &file{
		File: f,
	}

	return mutexFile, err
}

func (fs *AdapterFs) createDir(fullpath string) error {
	dir := filepath.Dir(fullpath)
	if dir != "." {
		if err := fs.fs.MkdirAll(dir, defaultDirectoryMode); err != nil {
			return err
		}
	}

	return nil
}

func (fs *AdapterFs) ReadDir(path string) ([]os.FileInfo, error) {
	l, err := afero.ReadDir(fs.fs, path)
	if err != nil {
		return nil, err
	}

	var s = make([]os.FileInfo, len(l))
	for i, f := range l {
		s[i] = f
	}

	return s, nil
}

func (fs *AdapterFs) Rename(from, to string) error {
	if err := fs.createDir(to); err != nil {
		return err
	}

	return os.Rename(from, to)
}

func (fs *AdapterFs) MkdirAll(path string, perm os.FileMode) error {
	return fs.fs.MkdirAll(path, defaultDirectoryMode)
}

func (fs *AdapterFs) Open(filename string) (billy.File, error) {
	return fs.OpenFile(filename, os.O_RDONLY, 0)
}

func (fs *AdapterFs) Stat(filename string) (os.FileInfo, error) {
	return fs.fs.Stat(filename)
}

func (fs *AdapterFs) Remove(filename string) error {
	return fs.fs.Remove(filename)
}

func (fs *AdapterFs) TempFile(dir, prefix string) (billy.File, error) {
	if err := fs.createDir(dir + string(os.PathSeparator)); err != nil {
		return nil, err
	}

	f, err := afero.TempFile(fs.fs, dir, prefix)
	if err != nil {
		return nil, err
	}
	return &file{File: f}, nil
}

func (fs *AdapterFs) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (fs *AdapterFs) RemoveAll(path string) error {
	return fs.fs.RemoveAll(filepath.Clean(path))
}

func (fs *AdapterFs) Lstat(filename string) (os.FileInfo, error) {

	info, success := fs.fs.(afero.Lstater)
	if success {
		s, _, err := info.LstatIfPossible(filename)
		if err != nil {
			return nil, err
		}

		return s, nil

	}

	return fs.fs.Stat(filename)
}

func (fs *AdapterFs) Symlink(target, link string) error {
	if err := fs.createDir(link); err != nil {
		return err
	}

	// TODO afero does not support symlinks
	return nil
}

func (fs *AdapterFs) Readlink(link string) (string, error) {

	// TODO afero does not support symlinks
	return "", nil
}

// Capabilities implements the Capable interface.
func (fs *AdapterFs) Capabilities() billy.Capability {
	return billy.DefaultCapabilities
}

// file is a wrapper for an os.File which adds support for file locking.
type file struct {
	afero.File
	m sync.Mutex
}

func (f *file) Lock() error {
	return nil
}

func (f *file) Unlock() error {
	return nil
}
