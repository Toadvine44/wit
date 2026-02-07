package objects

import (
	"crypto/sha1"
	"fmt"
	"sort"
)

// FileMode represents the mode/permissions of a tree entry.
type FileMode uint32

const (
	ModeFile    FileMode = 0100644 // regular file
	ModeExec    FileMode = 0100755 // executable file
	ModeSymlink FileMode = 0120000 // symbolic link
	ModeDir     FileMode = 0040000 // directory (tree)
)

// TreeEntry represents a single entry in a tree object.
// Each entry points to either a blob (file) or another tree (subdirectory).
type TreeEntry struct {
	Mode FileMode
	Name string
	Hash [20]byte
}

// Tree represents a directory listing in the object store.
// A tree contains references to blobs (files) and other trees (subdirectories).
type Tree struct {
	Entries []TreeEntry
}

// NewTree creates a new empty Tree.
func NewTree() *Tree {
	return &Tree{
		Entries: make([]TreeEntry, 0),
	}
}

// AddEntry adds a new entry to the tree.
func (t *Tree) AddEntry(mode FileMode, name string, hash [20]byte) {
	t.Entries = append(t.Entries, TreeEntry{
		Mode: mode,
		Name: name,
		Hash: hash,
	})
}

// Sort sorts tree entries by name (required for consistent hashing).
func (t *Tree) Sort() {
	sort.Slice(t.Entries, func(i, j int) bool {
		return t.Entries[i].Name < t.Entries[j].Name
	})
}

// Hash computes the SHA-1 hash of the tree in git format.
func (t *Tree) Hash() [20]byte {
	t.Sort()
	var content []byte
	for _, entry := range t.Entries {
		line := fmt.Sprintf("%o %s\x00", entry.Mode, entry.Name)
		content = append(content, []byte(line)...)
		content = append(content, entry.Hash[:]...)
	}
	header := fmt.Sprintf("tree %d\x00", len(content))
	full := append([]byte(header), content...)
	return sha1.Sum(full)
}

// HashString returns the hex-encoded hash of the tree.
func (t *Tree) HashString() string {
	hash := t.Hash()
	return fmt.Sprintf("%x", hash)
}

// Type returns the object type identifier.
func (t *Tree) Type() ObjectType {
	return TypeTree
}
