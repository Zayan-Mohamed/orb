package filesystem

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Zayan-Mohamed/orb/pkg/protocol"
)

var (
	ErrPathTraversal    = errors.New("path traversal attempt detected")
	ErrSymlinkEscape    = errors.New("symlink points outside shared directory")
	ErrInvalidPath      = errors.New("invalid path")
	ErrPermissionDenied = errors.New("permission denied")
)

// SecureFilesystem provides sandboxed filesystem operations
type SecureFilesystem struct {
	rootPath string
	readOnly bool
}

// NewSecureFilesystem creates a new secure filesystem handler
func NewSecureFilesystem(rootPath string, readOnly bool) (*SecureFilesystem, error) {
	// Resolve to absolute path
	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve root path: %w", err)
	}

	// Verify root exists and is a directory
	info, err := os.Stat(absRoot)
	if err != nil {
		return nil, fmt.Errorf("root path does not exist: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("root path is not a directory")
	}

	return &SecureFilesystem{
		rootPath: absRoot,
		readOnly: readOnly,
	}, nil
}

// sanitizePath ensures the path is within the root directory
// This prevents path traversal attacks
func (fs *SecureFilesystem) sanitizePath(path string) (string, error) {
	// Clean the path (removes .., ., etc.)
	cleaned := filepath.Clean(path)

	// Remove leading slash to make it relative
	cleaned = strings.TrimPrefix(cleaned, string(filepath.Separator))

	// Join with root
	fullPath := filepath.Join(fs.rootPath, cleaned)

	// Resolve any symlinks
	resolved, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		// Path doesn't exist yet (for create operations)
		// Check parent directory instead
		parent := filepath.Dir(fullPath)
		resolved, err = filepath.EvalSymlinks(parent)
		if err != nil {
			return "", fmt.Errorf("invalid path: %w", err)
		}
		// Reconstruct the full path with the original filename
		resolved = filepath.Join(resolved, filepath.Base(fullPath))
	}

	// Ensure resolved path is still within root
	if !strings.HasPrefix(resolved, fs.rootPath) {
		return "", ErrPathTraversal
	}

	return resolved, nil
}

// List returns directory contents
func (fs *SecureFilesystem) List(path string) (*protocol.ListResponse, error) {
	safePath, err := fs.sanitizePath(path)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(safePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	files := make([]protocol.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue // Skip entries we can't stat
		}

		// Check if symlink points outside root
		if info.Mode()&os.ModeSymlink != 0 {
			linkPath := filepath.Join(safePath, entry.Name())
			target, err := filepath.EvalSymlinks(linkPath)
			if err != nil || !strings.HasPrefix(target, fs.rootPath) {
				// Skip symlinks that point outside or are broken
				continue
			}
		}

		files = append(files, protocol.FileInfo{
			Name:    entry.Name(),
			Size:    info.Size(),
			Mode:    uint32(info.Mode()),
			ModTime: info.ModTime().Unix(),
			IsDir:   info.IsDir(),
		})
	}

	return &protocol.ListResponse{Files: files}, nil
}

// Stat returns file information
func (fs *SecureFilesystem) Stat(path string) (*protocol.StatResponse, error) {
	safePath, err := fs.sanitizePath(path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(safePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	return &protocol.StatResponse{
		Info: protocol.FileInfo{
			Name:    info.Name(),
			Size:    info.Size(),
			Mode:    uint32(info.Mode()),
			ModTime: info.ModTime().Unix(),
			IsDir:   info.IsDir(),
		},
	}, nil
}

// Read reads file contents
func (fs *SecureFilesystem) Read(path string, offset, length int64) (*protocol.ReadResponse, error) {
	safePath, err := fs.sanitizePath(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(safePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Warning: failed to close file: %v", err)
		}
	}()

	// Get file size
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Validate offset
	if offset < 0 || offset > info.Size() {
		return nil, errors.New("invalid offset")
	}

	// Seek to offset
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek: %w", err)
	}

	// Calculate read length
	if length <= 0 || offset+length > info.Size() {
		length = info.Size() - offset
	}

	// Limit read size to prevent memory exhaustion
	const maxReadSize = 10 * 1024 * 1024 // 10 MB
	if length > maxReadSize {
		length = maxReadSize
	}

	// Read data
	data := make([]byte, length)
	n, err := io.ReadFull(file, data)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return &protocol.ReadResponse{Data: data[:n]}, nil
}

// Write writes data to a file
func (fs *SecureFilesystem) Write(path string, offset int64, data []byte) (*protocol.WriteResponse, error) {
	if fs.readOnly {
		return nil, ErrPermissionDenied
	}

	safePath, err := fs.sanitizePath(path)
	if err != nil {
		return nil, err
	}

	// Open or create file
	file, err := os.OpenFile(safePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Warning: failed to close file: %v", err)
		}
	}()

	// Seek to offset
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek: %w", err)
	}

	// Write data
	n, err := file.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &protocol.WriteResponse{BytesWritten: int64(n)}, nil
}

// Delete removes a file or directory
func (fs *SecureFilesystem) Delete(path string) error {
	if fs.readOnly {
		return ErrPermissionDenied
	}

	safePath, err := fs.sanitizePath(path)
	if err != nil {
		return err
	}

	// Prevent deleting the root directory
	if safePath == fs.rootPath {
		return errors.New("cannot delete root directory")
	}

	if err := os.RemoveAll(safePath); err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	return nil
}

// Rename renames a file or directory
func (fs *SecureFilesystem) Rename(oldPath, newPath string) error {
	if fs.readOnly {
		return ErrPermissionDenied
	}

	safeOldPath, err := fs.sanitizePath(oldPath)
	if err != nil {
		return err
	}

	safeNewPath, err := fs.sanitizePath(newPath)
	if err != nil {
		return err
	}

	// Prevent renaming the root directory
	if safeOldPath == fs.rootPath || safeNewPath == fs.rootPath {
		return errors.New("cannot rename root directory")
	}

	if err := os.Rename(safeOldPath, safeNewPath); err != nil {
		return fmt.Errorf("failed to rename: %w", err)
	}

	return nil
}

// Mkdir creates a directory
func (fs *SecureFilesystem) Mkdir(path string, perm uint32) error {
	if fs.readOnly {
		return ErrPermissionDenied
	}

	safePath, err := fs.sanitizePath(path)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(safePath, os.FileMode(perm)); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// IsReadOnly returns whether the filesystem is read-only
func (fs *SecureFilesystem) IsReadOnly() bool {
	return fs.readOnly
}

// RootPath returns the root path
func (fs *SecureFilesystem) RootPath() string {
	return fs.rootPath
}
