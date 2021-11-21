package fsutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Listing struct {
	Name    string
	Path    string
	Entries []Entry
}

type Entry struct {
	Name           string
	Path           string
	RelPath        string
	LastModifiedAt time.Time
	Size           Size
	IsDir          bool
}

type Size int64

func List(dir, root string) (Listing, error) {
	listing := Listing{
		Name:    filepath.Base(dir),
		Path:    dir,
		Entries: make([]Entry, 0),
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return listing, fmt.Errorf("list dir %v: %w", dir, err)
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return listing, fmt.Errorf("list dir %v: %w", dir, err)
		}

		path := filepath.Join(dir, entry.Name())
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return listing, fmt.Errorf("list dir %v: %w", dir, err)
		}

		listing.Entries = append(listing.Entries, Entry{
			Name:           entry.Name(),
			Path:           path,
			RelPath:        relPath,
			LastModifiedAt: info.ModTime(),
			Size:           Size(info.Size()),
			IsDir:          info.IsDir(),
		})
	}

	return listing, nil
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, fmt.Errorf("stat path %v: %w", path, err)
	}
}

func IsDir(path string) (bool, error) {
	stat, err := os.Stat(path)

	if err == nil {
		return stat.IsDir(), nil
	} else {
		return false, fmt.Errorf("stat path %v: %w", path, err)
	}
}

const (
	kilobyte Size = 1024
	megabyte      = kilobyte * 1024
	gigabyte      = megabyte * 1024
)

func (size Size) String() string {
	if size < kilobyte {
		return fmt.Sprintf("%d", int64(size))
	} else if size < megabyte {
		return fmt.Sprintf("%.1fK", float64(size)/float64(kilobyte))
	} else if size < gigabyte {
		return fmt.Sprintf("%.1fM", float64(size)/float64(megabyte))
	} else {
		return fmt.Sprintf("%.1fG", float64(size)/float64(gigabyte))
	}
}
