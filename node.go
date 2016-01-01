package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/juju/errors"
)

var (
	ErrNodeExists     = errors.New("Conflicting path exists at destination")
	ErrNodeBadSymlink = errors.New("Cannot read destination node symlink")
	ErrNodeIncSymlink = errors.New("Destination symlink points elsewhere")
)

type Node struct {
	Path string
	Info os.FileInfo
}

func NewNode(path string) (node *Node, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Annotatef(err, "Could not resolve absolute path for %s", path)
	}
	// Lstat enables symlinks to be nodes themselves so links can point
	// to links.
	lstat, err := os.Lstat(absPath)
	if err != nil {
		return nil, errors.Annotatef(err, "Cannot Lstat path %s", path)
	}
	return &Node{
		Path: absPath,
		Info: lstat,
	}, nil
}

// Determine Node's path in terms of the plan
func (n *Node) PlannedPath(plan *Plan) string {
	return filepath.Join(plan.Dest, n.Rel(plan.Src))
}

// Resolve a relative path for this node from path
func (n *Node) Rel(path string) string {
	rel, _ := filepath.Rel(path, n.Path)
	return rel
}

// Is planned path a conflict?
func (n *Node) CheckNode(plan *Plan) error {
	path := n.PlannedPath(plan)
	// Check if destination exists
	lstat, err := os.Lstat(path)
	if err != nil {
		return nil // path doesn't exist
	}

	if lstat.Mode().IsRegular() || lstat.Mode().IsDir() {
		// This is already a file so we can't link here
		return ErrNodeExists
	}

	if lstat.Mode()&os.ModeSymlink == os.ModeSymlink {
		linkTarget, err := filepath.EvalSymlinks(path)
		if err != nil {
			return ErrNodeBadSymlink
		}
		// Resolve the absolute target for comparison
		linkAbsTarget, err := filepath.Abs(linkTarget)
		if err != nil {
			return ErrNodeBadSymlink
		}
		if linkAbsTarget != n.Path {
			return errors.Wrap(ErrNodeIncSymlink,
				errors.Errorf("Symlink exists with different path: actual - %s, desired -%s'",
					linkAbsTarget, n.Path))
		}
	}
	return nil
}

func (n *Node) String() string {
	return fmt.Sprintf("#<Node Path: %s>", n.Path)
}
