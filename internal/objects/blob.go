package objects

import (
	"crypto/sha1"
	"fmt"
)

// Blob represents a file's contents in the object store.
// In git's model, a blob is simply the compressed contents of a file,
// identified by the SHA-1 hash of "blob <size>\0<content>".
type Blob struct {
	Data []byte
}

// NewBlob creates a new Blob from raw file data.
func NewBlob(data []byte) *Blob {
	return &Blob{Data: data}
}

// Hash computes the SHA-1 hash of the blob in git format.
func (b *Blob) Hash() [20]byte {
	header := fmt.Sprintf("blob %d\x00", len(b.Data))
	content := append([]byte(header), b.Data...)
	return sha1.Sum(content)
}

// HashString returns the hex-encoded hash of the blob.
func (b *Blob) HashString() string {
	hash := b.Hash()
	return fmt.Sprintf("%x", hash)
}

// Type returns the object type identifier.
func (b *Blob) Type() ObjectType {
	return TypeBlob
}

// Size returns the size of the blob data in bytes.
func (b *Blob) Size() int {
	return len(b.Data)
}
