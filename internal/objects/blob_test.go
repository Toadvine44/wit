package objects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBlob_Success(t *testing.T) {
	_, err := NewBlob("./test/test_file.md")

	// expectedBlob := &Blob{
	// 	filename: "./test/test_file.md",
	// 	Sha1:     [20]byte{},
	// }

	assert.Nil(t, err)
	// assert.Equal(t, expectedBlob, b)
}
