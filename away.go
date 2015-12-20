package main

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

type AwayOptions struct{}

type Away struct {
	Options AwayOptions

	Target string
	Source string

	// Map of Path to Path that will be from the symlink
	plan map[string]string
}

func NewAway(source, target string) *Away {
	return &Away{
		Source: source,
		Target: target,

		// Unused options
		Options: AwayOptions{},
	}
}

func (a *Away) Plan() {
	pathMap, _ := a.walk(a.Source)
	a.plan = pathMap
}

func (a *Away) Stow() (err error) {
	if len(a.plan) == 0 {
		log.Errorln("Plan has not been made, run #Plan first")
		return
	}

	links := make(map[string]string)
	for srcfile, _ := range a.plan {
		target := filepath.Join(a.Source, srcfile)
		link := filepath.Join(a.Target, srcfile)

		_, err = os.Lstat(link)
		if err == nil {
			if IsEffectiveLink(link, target) {
				continue // our work here is done
			}
			log.Errorf("File '%s' occupies desired link space!", link)
			return err
		}
		links[link] = target
	}

	// Make all them links!
	for link, target := range links {
		os.Symlink(target, link)
	}
	return nil
}

// Discover the source files
func (a Away) walk(dir string) (pathMap map[string]string, err error) {
	pathMap = make(map[string]string)
	walker := func(p string, _info os.FileInfo, _err error) error {
		stat, err := os.Lstat(p)
		if err != nil {
			log.Fatalf("Error inspecting file %s: %s", p, err)
		}
		relpath, _ := filepath.Rel(dir, p)

		// Skip the top level directory itself
		if relpath == "." {
			return nil
		}

		if stat.Mode()&os.ModeSymlink == os.ModeSymlink {
			dest, err := os.Readlink(p)
			if err != nil {
				return err
			}
			// This file is a symlink and points elsewhere
			pathMap[relpath] = dest
		} else {
			// Retain the relative path to the source
			pathMap[relpath] = relpath
		}
		return nil
	}
	filepath.Walk(dir, walker)
	return pathMap, nil
}

// Compare target to current Symlink at name, true if its the same
func IsEffectiveLink(name, target string) bool {
	// Probably should exist..
	stat, err := os.Lstat(name)
	if err != nil {
		return false
	}

	// Current file at name is not a Symlink so there's no way that its
	// the same effective target
	if stat.Mode().IsRegular() {
		return false
	}

	dest, err := os.Readlink(name)
	if err != nil {
		return false
	}

	// It is a Symlink but not the same target
	if dest != target {
		return false
	}
	return true
}
