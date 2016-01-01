package main

import (
	"os"
	"path/filepath"
	"testing"
)

// Helper - runs mkTestDir and creates a plan using the resulting src
// and dest directories
func mkTestPlan(t *testing.T, paths []string) (plan *Plan, cleanup func()) {
	src, dest, cleanup := mkTestDir(t, paths)
	plan, err := NewPlan(src, dest, PlanOptions{})
	if err != nil {
		t.Fatalf("Could not create plan for testing: %s", err)
	}
	return plan, cleanup
}

func TestWalker(t *testing.T) {
	paths := []string{"A/B/c", "C/D/e"}
	plan, cleanup := mkTestPlan(t, paths)
	defer cleanup()

	err := plan.FindNodes()
	if err != nil {
		t.Fatal(err)
	}

	if len(plan.Nodes) != len(paths) {
		t.Errorf("There should only be the minimum nodes necessary to link: %d != %d", len(plan.Nodes), len(paths))
		for _, n := range plan.Nodes {
			t.Log(n)
		}
	}
}

func TestWalkerWithExistingDirs(t *testing.T) {
	path := "A/B/c"
	plan, cleanup := mkTestPlan(t, []string{path})
	defer cleanup()

	// Make all directories
	os.MkdirAll(filepath.Join(plan.Dest, "A", "B"), 0777)
	// Should find only 'c' needs to be linked
	err := plan.FindNodes()
	if err != nil {
		t.Fatal(err)
	}

	if len(plan.Nodes) != 1 {
		t.Error("Should have a single node")
		t.FailNow()
	}

	rel, _ := filepath.Rel(plan.Src, plan.Nodes[0].Path)
	if rel != "A/B/c" {
		t.Errorf("Node was not the lowest required: %s", plan.Nodes[0])
	}
}

func TestWalkerWithLinkFilesOnly(t *testing.T) {
	path := "A/B/c"
	plan, cleanup := mkTestPlan(t, []string{path})
	defer cleanup()

	// This is the same as directories existing in the destination
	plan.Options.LinkFilesOnly = true
	// Should find only 'c' needs to be linked
	err := plan.FindNodes()
	if err != nil {
		t.Fatal(err)
	}

	if len(plan.Nodes) != 1 {
		t.Error("Should have a single node")
		t.FailNow()
	}

	rel, _ := filepath.Rel(plan.Src, plan.Nodes[0].Path)
	if rel != "A/B/c" {
		t.Errorf("Node was not the lowest required: %s", plan.Nodes[0])
	}
}

func TestDirExists(t *testing.T) {
	src, _, cleanup := mkTestDir(t, []string{"A/file"})
	defer cleanup()

	if DirExist(filepath.Join(src, "A/file")) {
		t.Error("Its a file not a directory.")
	}

	if !DirExist(filepath.Join(src, "A")) {
		t.Errorf("This directory exists: %s", filepath.Join(src, "A"))
	}

	if DirExist("bogus") {
		t.Error("Bogus path should not be a dir.")
	}
}

func TestBadPlan(t *testing.T) {
	_, err := NewPlan("bogus", "/tmp", PlanOptions{})
	if err == nil {
		t.Error("Plan should have errored on a non-existent source")
	}

	_, err = NewPlan("/tmp", "bogus", PlanOptions{})
	if err == nil {
		t.Error("Plan should have errored on non-existent dest")
	}
}
