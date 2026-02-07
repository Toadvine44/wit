package objects

// ObjectType represents the type of a git object.
type ObjectType string

const (
	TypeBlob   ObjectType = "blob"
	TypeTree   ObjectType = "tree"
	TypeCommit ObjectType = "commit"
)

// Object is the interface implemented by all git objects.
type Object interface {
	Hash() [20]byte
	HashString() string
	Type() ObjectType
}
