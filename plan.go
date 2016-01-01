package main

import (
	"os"
	"path/filepath"

	"github.com/juju/errors"
)

type PlanOptions struct {
	DryRun        bool
	LinkFilesOnly bool
}

type Plan struct {
	Src   string
	Dest  string
	Nodes []*Node

	Options PlanOptions
}

func NewPlan(src, dest string, options PlanOptions) (*Plan, error) {
	if !DirExist(src) || !DirExist(dest) {
		return &Plan{}, errors.New("Source or destination directory does not exist")
	}
	return &Plan{
		Src:     src,
		Dest:    dest,
		Options: options,
	}, nil
}

func (p *Plan) FindNodes() error {
	return filepath.Walk(p.Src, p.walker)
}

// filepath.WalkerFunc for discovering the nodes based on current plan
// configuration
func (p *Plan) walker(path string, info os.FileInfo, ferr error) (err error) {
	node, err := NewNode(path)
	if err != nil {
		return err
	}

	// Delve deeper if this directory exists at the destination
	if info.IsDir() && (ErrNodeExists == node.CheckNode(p)) {
		return
	}

	// Store directories as nodes and skip its children unless the user
	// wants to only link in files.
	if info.IsDir() {
		if p.Options.LinkFilesOnly {
			return err
		}
		err = filepath.SkipDir
	}

	p.Nodes = append(p.Nodes, node)
	return err
}

// Check that directory exists
func DirExist(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}
