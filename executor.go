package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

// Executor handler interface that are responsible for taking plans
// into action.
type Executor interface {
	// Run the plan with effect
	Execute(*Plan) error
	// Describe what the executor would do without effect
	Describe(*Plan) string
}

// Putaway files into the filesystem by linking
type Putaway struct{}

func (s Putaway) Execute(p *Plan) error {
	for _, n := range p.Nodes {
		dir := filepath.Dir(n.Rel(p.Src))
		if dir != "." && p.Options.LinkFilesOnly {
			os.MkdirAll(filepath.Join(p.Dest, dir), p.Options.DirMode)
		}
		os.Symlink(n.Path, n.PlannedPath(p))
	}
	return nil
}

func (s Putaway) Describe(p *Plan) string {
	buf := &bytes.Buffer{}
	if p.Options.LinkFilesOnly {
		log.Info("This plan will create all directories between the destination node")
	}
	for _, n := range p.Nodes {
		dir := filepath.Dir(n.Rel(p.Src))
		if dir != "." && p.Options.LinkFilesOnly {
			fmt.Fprintf(buf, "MKDIR: %s\n", filepath.Join(p.Dest, dir))
		}
		fmt.Fprintf(buf, "LINK:  %s => %s\n", n.PlannedPath(p), n.Path)
	}

	return buf.String()
}

// Undo the Putaway removing links that the plan can account for
type Unputaway struct{}

func (e Unputaway) Execute(p *Plan) error {
	for _, n := range p.Nodes {
		// Node should pass checks if its to be removed by Unputaway
		if err := n.CheckNode(p); err == nil {
			os.Remove(n.PlannedPath(p))
		} else {
			log.Warnf("SKIP %s - link is not managed by this away: %s\n", n.PlannedPath(p), err)
		}
	}
	return nil
}

func (e Unputaway) Describe(p *Plan) string {
	buf := &bytes.Buffer{}
	for _, n := range p.Nodes {
		// Node should pass checks if its to be removed by Unputaway
		if err := n.CheckNode(p); err == nil {
			fmt.Fprintf(buf, "RM:   %s\n", n.PlannedPath(p))
		} else {
			fmt.Fprintf(buf, "SKIP: %s - link is not managed by this away\n", n.PlannedPath(p))
		}
	}
	return buf.String()
}
