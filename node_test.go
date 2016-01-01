package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNodePathPlanning(t *testing.T) {
	path := "A/file"
	src, dest, cleanup := mkTestDir(t, []string{path})
	defer cleanup()
	plan, _ := NewPlan(src, dest, PlanOptions{})

	nodePath := filepath.Join(src, path)
	expectedPath := filepath.Join(dest, path)

	node, err := NewNode(nodePath)
	if err != nil {
		t.Logf("Error occurred", err)
		t.Error(err)
	}
	plannedPath := node.PlannedPath(plan)
	if expectedPath != plannedPath {
		t.Logf("Expected path %s to be %s", expectedPath, plannedPath)
		t.Fail()
	}
}

func TestNodeRel(t *testing.T) {
	path := "A/file"
	src, _, cleanup := mkTestDir(t, []string{path})
	defer cleanup()

	node, _ := NewNode(filepath.Join(src, path))
	if node.Rel(src) != path {
		t.Error("Node should have returned its relative path", node.Rel(src))
	}
}

func TestNodeCheck(t *testing.T) {
	path := "file"

	// Check with clean destination
	plan, cleanup := mkTestPlan(t, []string{path})
	node, _ := NewNode(filepath.Join(plan.Src, path))
	err := node.CheckNode(plan)
	if err != nil {
		t.Error("Node should be okay without files in the target")
	}
	cleanup()

	// Check regular files
	plan, cleanup = mkTestPlan(t, []string{path})
	node, _ = NewNode(filepath.Join(plan.Src, path))

	// Make regular file at the PlannedPath
	os.MkdirAll(filepath.Dir(node.PlannedPath(plan)), 0777)
	ioutil.WriteFile(node.PlannedPath(plan), []byte("conflict!"), 0666)
	err = node.CheckNode(plan)
	if err == nil {
		t.Error("Node should have an error when a file exists at the Planned path")
	}
	cleanup()

	// Check bad symlinks
	plan, cleanup = mkTestPlan(t, []string{path})
	node, _ = NewNode(filepath.Join(plan.Src, path))
	// Make regular file at the PlannedPath
	os.MkdirAll(filepath.Dir(node.PlannedPath(plan)), 0777)
	os.Symlink("/tmp/bad_link", node.PlannedPath(plan))
	err = node.CheckNode(plan)
	if err == nil {
		t.Error("Node should have an error when a file exists at the Planned path")
	}
	cleanup()

	// Check existing but unplanned symlinks
	plan, cleanup = mkTestPlan(t, []string{path})
	node, _ = NewNode(filepath.Join(plan.Src, path))
	// Make regular file at the PlannedPath
	os.MkdirAll(filepath.Dir(node.PlannedPath(plan)), 0777)
	os.Symlink("/usr/bin/false", node.PlannedPath(plan))
	err = node.CheckNode(plan)
	if err == nil {
		t.Error("Node should have an error when a file exists at the Planned path")
	}
	cleanup()

	// Check existing but unplanned symlinks
	plan, cleanup = mkTestPlan(t, []string{path})
	node, _ = NewNode(filepath.Join(plan.Src, path))
	// Make regular file at the PlannedPath
	os.MkdirAll(filepath.Dir(node.PlannedPath(plan)), 0777)
	os.Symlink(node.Path, node.PlannedPath(plan))
	err = node.CheckNode(plan)
	if err != nil {
		t.Errorf("Node should be able to determine its been linked properly before: %s", err)
	}
	cleanup()
}

func TestNewNode(t *testing.T) {
	_, err := NewNode("/tmp")
	if err != nil {
		t.Error("Should have been able to create a new Node for /tmp")
	}
	_, err = NewNode("bogus")
	if err == nil {
		t.Error("Should have errored trying to create Node for bogus path")
	}
}
