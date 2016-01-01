package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

type Executor interface {
	Execute(*Plan) error
	Describe(*Plan) []byte
}

type Stowaway struct{}

func (s Stowaway) Execute(p *Plan) {
	os.MkdirAll(p.Dest, 0775)
}

func (s Stowaway) Describe(p *Plan) []byte {
	buf := &bytes.Buffer{}
	if p.Options.LinkFilesOnly {
		log.Info("This plan will create all directories between the destination node")
	}

	for _, n := range p.Nodes {
		dir := filepath.Dir(n.Rel(p.Src))
		if dir != "." && p.Options.LinkFilesOnly {
			fmt.Fprintf(buf, "MKDIR: %s", dir)
			continue
		}
		fmt.Fprintf(buf, "LINK:  %s => %s", n.PlannedPath(p), n.Rel(p.Dest))
	}

	return buf.Bytes()
}
