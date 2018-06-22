package contentaddressable

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"os"
	"path/filepath"
)

var (
	AlreadyClosed = errors.New("Already closed.")
	HasData       = errors.New("Destination file already has data.")
	DefaultSuffix = "-temp"
)

// File handles the atomic writing of a content addressable file.  It writes to
// a temp file, and then renames to the final location after Accept().
type File struct {
	Oid          string
	filename     string
	tempFilename string
	tempFile     *os.File
	hasher       hash.Hash
}

// NewFile initializes a content addressable file for writing.  It is identical
// to NewWithSuffix, except it uses DefaultSuffix as the suffix.
func NewFile(filename string) (*File, error) {
	return NewWithSuffix(filename, DefaultSuffix)
}

// NewWithSuffix initializes a content addressable file for writing.
// Data is written to a temporary file, and atomically renamed to the destination
// filename when Accept() is called. The *File OID is taken from the base name
// of the given filename.
func NewWithSuffix(filename, suffix string) (*File, error) {
	oid := filepath.Base(filename)
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	tempFilename := filename + suffix
	tempFile, err := os.OpenFile(tempFilename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return nil, err
	}

	caw := &File{
		Oid:          oid,
		filename:     filename,
		tempFilename: tempFilename,
		tempFile:     tempFile,
		hasher:       sha256.New(),
	}

	return caw, nil
}

// Write sends data to the temporary file.
func (w *File) Write(p []byte) (int, error) {
	if w.Closed() {
		return 0, AlreadyClosed
	}

	w.hasher.Write(p)
	return w.tempFile.Write(p)
}

// Accept verifies the written content SHA-256 signature matches the given OID.
// If it matches, the temp file is renamed to the destination filename.
// Returns a bool indicating whether the destination file was created (if not,
// someone else adding the same contents in parallel got there first), and
// an error that might have occurred during the rename.
func (w *File) Accept() (bool, error) {
	if w.Closed() {
		return false, AlreadyClosed
	}

	sig := hex.EncodeToString(w.hasher.Sum(nil))
	if sig != w.Oid {
		return false, fmt.Errorf("Content mismatch.  Expected OID %s, got %s", w.Oid, sig)
	}

	// Only bother renaming the temp file if the destination file doesn't already exist.
	// Since the SHA-256 must match, we can be confident that the contents are identical.
	if _, err := os.Stat(w.filename); err != nil {
		w.tempFile.Close()
		w.tempFile = nil

		// rename the temp file to the real file
		return true, os.Rename(w.tempFilename, w.filename)
	}

	return false, w.Close()
}

// Close cleans up the internal file objects.
func (w *File) Close() error {
	if w.tempFile != nil {
		if err := cleanupFile(w.tempFile); err != nil {
			return err
		}
		w.tempFile = nil
	}

	return nil
}

// Closed reports whether this file object has been closed.
func (w *File) Closed() bool {
	if w.tempFile == nil {
		return true
	}
	return false
}

func cleanupFile(f *os.File) error {
	err := f.Close()
	if err := os.RemoveAll(f.Name()); err != nil {
		return err
	}

	return err
}
