package objects

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"time"
)

// Signature represents an author or committer with timestamp.
type Signature struct {
	Name  string
	Email string
	Time  time.Time
}

// String formats the signature in git format.
func (s Signature) String() string {
	return fmt.Sprintf("%s <%s> %d %s",
		s.Name,
		s.Email,
		s.Time.Unix(),
		s.Time.Format("-0700"),
	)
}

// Commit represents a commit object in the object store.
// A commit points to a tree (the snapshot), has parent commits,
// and contains metadata about who made the commit and when.
type Commit struct {
	Tree      [20]byte   // hash of the root tree
	Parents   [][20]byte // hashes of parent commits (empty for initial commit)
	Author    Signature
	Committer Signature
	Message   string
}

// NewCommit creates a new Commit with the given tree and message.
func NewCommit(tree [20]byte, message string, author Signature) *Commit {
	return &Commit{
		Tree:      tree,
		Parents:   make([][20]byte, 0),
		Author:    author,
		Committer: author,
		Message:   message,
	}
}

// AddParent adds a parent commit hash.
func (c *Commit) AddParent(parent [20]byte) {
	c.Parents = append(c.Parents, parent)
}

// Content builds the commit content in git format (without the header).
func (c *Commit) Content() []byte {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("tree %x\n", c.Tree))

	for _, parent := range c.Parents {
		builder.WriteString(fmt.Sprintf("parent %x\n", parent))
	}

	builder.WriteString(fmt.Sprintf("author %s\n", c.Author))
	builder.WriteString(fmt.Sprintf("committer %s\n", c.Committer))
	builder.WriteString("\n")
	builder.WriteString(c.Message)

	return []byte(builder.String())
}

// Hash computes the SHA-1 hash of the commit in git format.
func (c *Commit) Hash() [20]byte {
	content := c.Content()
	header := fmt.Sprintf("commit %d\x00", len(content))
	full := append([]byte(header), content...)
	return sha1.Sum(full)
}

// HashString returns the hex-encoded hash of the commit.
func (c *Commit) HashString() string {
	hash := c.Hash()
	return fmt.Sprintf("%x", hash)
}

// Type returns the object type identifier.
func (c *Commit) Type() ObjectType {
	return TypeCommit
}
