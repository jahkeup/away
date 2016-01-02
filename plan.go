package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/juju/errors"
)

type PlanOptions struct {
	DryRun        bool
	LinkFilesOnly bool
	DirMode       os.FileMode
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
	if !filepath.IsAbs(src) {
		src, _ = filepath.Abs(src)
	}
	if !filepath.IsAbs(dest) {
		dest, _ = filepath.Abs(dest)
	}

	// If DirMode is not set, this defaults to 770
	if options.DirMode == 0 {
		options.DirMode = 0770
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

	// Store directories as nodes and skip its children unless the user
	// wants to only link in files.
	if info.IsDir() {
		if node.CheckNode(p) == ErrNodeExists {
			return nil
		}
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

func (p *Plan) String() string {
	return fmt.Sprintf("#<Plan Src: %s, Dest: %s, Nodes: %s>", p.Src, p.Dest, p.Nodes)
}
