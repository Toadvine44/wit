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
}

func NewBlob(filename string) (*Blob, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	b := &Blob{}
	b.filename = filename
	b.Sha1 = b.hash(data)
	b.writeToDisk(data)
	return b, nil
}

func compressFile(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		return nil, err
	}

	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}
	w.Close()
	return buf.Bytes(), nil
}

func (b *Blob) writeToDisk(data []byte) error {
	err := os.Mkdir(fmt.Sprintf("./.wit/objects/%x", string(b.Sha1[:1])), 0755)
	outFile, err := os.Create(fmt.Sprintf("./.wit/objects/%x/%x", string(b.Sha1[:1]), string(b.Sha1[1:])))
	if err != nil {
		return err
	}
	defer outFile.Close()
	compressed, err := compressFile(data)
	if err != nil {
		return err
	}
	_, err = io.Copy(outFile, bytes.NewReader(compressed))
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
