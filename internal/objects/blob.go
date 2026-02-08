package objects

import (
	"bytes"
	"fmt"
	// "hex"
	"io"
	"os"

	"compress/zlib"
	"crypto/sha1"
)

// Blob represents a file's contents in the object store.
// In git's model, a blob is simply the compressed contents of a file,
// identified by the SHA-1 hash of "blob <size>\0<content>".
type Blob struct {
	filename string
	Sha1     [20]byte
	// ZipData []byte
}

func NewBlob(filename string) (*Blob, error) {
	data, err := os.ReadFile(filename)
	fmt.Printf("Bytes: %d read from %s\n", len(data), filename)
	if err != nil {
		return nil, err
	}

	b := &Blob{}
	b.filename = filename
	b.Sha1 = b.hash(data)
	fmt.Printf("Sha hash: %v\n", fmt.Sprintf("%x", b.Sha1[:]))

	b.writeToDisk(data)
	return b, nil
}

func (b *Blob) writeToDisk(data []byte) error {
	r := bytes.NewReader(data)
	err := os.Mkdir(fmt.Sprintf("./.wit/objects/%s", string(b.Sha1[:2])), 0755)
	outFile, err := os.Create(fmt.Sprintf("./.wit/objects/%s/%s", string(b.Sha1[:2]), string(b.Sha1[2:])))
	if err != nil {
		return err
	}
	defer outFile.Close()
	w, err := zlib.NewWriterLevel(outFile, zlib.BestCompression)
	if err != nil {
		return err
	}
	written, err := io.Copy(w, r)
	fmt.Printf("Bytes written: %d\n", written)
	return err
}

func (b *Blob) hash(data []byte) [20]byte {
	header := fmt.Sprintf("blob %d\x00", len(data))
	content := append([]byte(header), data...)
	return sha1.Sum(content)
}

// func (b *Blob) Type() ObjectType {
// 	return TypeBlob
// }
//
// func (b *Blob) Size() int {
// 	return len(b.Data)
// }
