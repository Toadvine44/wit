package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupWitDir(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	origDir, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmpDir))
	require.NoError(t, os.MkdirAll(".wit/objects", 0755))

	t.Cleanup(func() {
		os.Chdir(origDir)
	})

	return tmpDir
}

// createTempFile writes content to a file inside dir and returns its path.
func createTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	return path
}

// --- compressFile tests ---

func TestCompressFile_RoundTrip(t *testing.T) {
	data := []byte("hello world, this is some test data for compression")
	compressed, err := compressFile(data)
	require.NoError(t, err)

	r, err := zlib.NewReader(bytes.NewReader(compressed))
	require.NoError(t, err)
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.Equal(t, data, decompressed)
}

func TestCompressFile_EmptyInput(t *testing.T) {
	compressed, err := compressFile([]byte{})
	require.NoError(t, err)

	r, err := zlib.NewReader(bytes.NewReader(compressed))
	require.NoError(t, err)
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.Empty(t, decompressed)
}

func TestCompressFile_OutputSmallerThanInput(t *testing.T) {
	// Repetitive data should compress well
	data := bytes.Repeat([]byte("This is a test file.\n"), 100)
	compressed, err := compressFile(data)
	require.NoError(t, err)
	assert.Less(t, len(compressed), len(data))
}

// --- hash tests ---

func TestHash_KnownValue(t *testing.T) {
	b := &Blob{}
	data := []byte("hello")

	header := fmt.Sprintf("blob %d\x00", len(data))
	expected := sha1.Sum(append([]byte(header), data...))

	got := b.hash(data)
	assert.Equal(t, expected, got)
}

func TestHash_EmptyData(t *testing.T) {
	b := &Blob{}
	data := []byte{}

	expected := sha1.Sum([]byte("blob 0\x00"))

	got := b.hash(data)
	assert.Equal(t, expected, got)
}

// --- writeToDisk tests ---

func TestWriteToDisk_CreatesFile(t *testing.T) {
	dir := setupWitDir(t)
	_ = dir

	b := &Blob{Sha1: sha1.Sum([]byte("test"))}
	data := []byte("file contents")

	err := b.writeToDisk(data)
	require.NoError(t, err)

	objPath := fmt.Sprintf(".wit/objects/%x/%x", string(b.Sha1[:1]), string(b.Sha1[1:]))
	_, err = os.Stat(objPath)
	assert.NoError(t, err, "object file should exist on disk")
}

func TestWriteToDisk_ContentIsZlibCompressed(t *testing.T) {
	setupWitDir(t)

	b := &Blob{Sha1: sha1.Sum([]byte("test"))}
	data := []byte("compressed content check")

	require.NoError(t, b.writeToDisk(data))

	objPath := fmt.Sprintf(".wit/objects/%x/%x", string(b.Sha1[:1]), string(b.Sha1[1:]))
	raw, err := os.ReadFile(objPath)
	require.NoError(t, err)

	r, err := zlib.NewReader(bytes.NewReader(raw))
	require.NoError(t, err)
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.Equal(t, data, decompressed)
}

// --- NewBlob tests ---

func TestNewBlob_Success(t *testing.T) {
	dir := setupWitDir(t)
	path := createTempFile(t, dir, "hello.txt", "hello world")

	b, err := NewBlob(path)
	require.NoError(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, path, b.filename)
}

func TestNewBlob_SetsCorrectHash(t *testing.T) {
	dir := setupWitDir(t)
	content := "hash me please"
	path := createTempFile(t, dir, "hashtest.txt", content)

	b, err := NewBlob(path)
	require.NoError(t, err)

	expected := sha1.Sum([]byte(fmt.Sprintf("blob %d\x00%s", len(content), content)))
	assert.Equal(t, expected, b.Sha1)
}

func TestNewBlob_WritesObjectFile(t *testing.T) {
	dir := setupWitDir(t)
	path := createTempFile(t, dir, "obj.txt", "object data")

	b, err := NewBlob(path)
	require.NoError(t, err)

	objPath := fmt.Sprintf(".wit/objects/%x/%x", string(b.Sha1[:1]), string(b.Sha1[1:]))
	_, err = os.Stat(objPath)
	assert.NoError(t, err, "object file should be written to disk")
}

func TestNewBlob_FileNotFound(t *testing.T) {
	b, err := NewBlob("/nonexistent/path/file.txt")
	assert.Error(t, err)
	assert.Nil(t, b)
}

func TestNewBlob_EmptyFile(t *testing.T) {
	dir := setupWitDir(t)
	path := createTempFile(t, dir, "empty.txt", "")

	b, err := NewBlob(path)
	require.NoError(t, err)
	assert.NotNil(t, b)

	expected := sha1.Sum([]byte("blob 0\x00"))
	assert.Equal(t, expected, b.Sha1)
}
