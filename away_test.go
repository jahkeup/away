package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Create testing directory for use in Away tests
//
// paths - A list of the paths that should exist under the src path to
// test putting away.
//
// Returns the directories that were created for testing usage:
//
// src - To be used at the Away Source root
// tgt - To be used as the Away Target root
//
// And a function to clean up these directories at the conclusion of
// the test.
func mkTestDir(t *testing.T, paths []string) (src, tgt string, cleanup func()) {
	root := path.Join(os.TempDir(), "away"+strconv.Itoa(rand.Intn(500)))
	src = path.Join(root, "src")
	tgt = path.Join(root, "tgt")

	os.MkdirAll(root, 0775)
	os.Mkdir(src, 0775)
	os.Mkdir(tgt, 0775)

	for _, p := range paths {
		if strings.Contains(p, "/") {
			// Create the parent directory for the file
			os.MkdirAll(path.Join(src, filepath.Dir(p)), 0775)
		}

		// Write filename into the test file for debug
		ioutil.WriteFile(path.Join(src, p), []byte(p), 0775)
	}

	return src, tgt, func() {
		lister := func(p string, _ os.FileInfo, _ error) error {
			info, _ := os.Lstat(p)
			if info.Mode()&os.ModeSymlink == os.ModeSymlink {
				dest, _ := os.Readlink(p)
				t.Logf("%s -> %s\n", p, dest)
			} else {
				t.Logf("%s\n", p)
			}
			return nil
		}
		t.Logf("Test directory:")
		filepath.Walk(root, lister)
		// Helper to cleanup the directory
		os.RemoveAll(root)
	}
}

func TestFileDiscovery(t *testing.T) {
	files := []string{"A", "B/Bchild"}
	src, tgt, cleanup := mkTestDir(t, files)
	defer cleanup()
	away := NewAway(src, tgt)

	res, err := away.walk(src)
	if err != nil {
		t.Error("Error occurred during the walk, must be raining")
	}
	for _, f := range files {
		_, ok := res[f]
		if !ok {
			t.Errorf("File %s was not in discovered list", f)
		}
	}
}

func TestStow(t *testing.T) {
	files := []string{"A", "B"}
	src, tgt, cleanup := mkTestDir(t, files)
	defer cleanup()
	away := NewAway(src, tgt)
	away.Plan()
	err := away.Stow()
	if err != nil {
		t.Fatalf("Error during test stow: %s", err)
	}

	srcMap, _ := away.walk(src)
	tgtMap, _ := away.walk(tgt)
	t.Logf("%v", srcMap)
	t.Logf("%v", tgtMap)
	if len(srcMap) != len(tgtMap) {
		t.Error("They don't match up..")
	}
	for f, _ := range srcMap {
		_, ok := tgtMap[f]
		if !ok {
			t.Errorf("File %s missing from the target", f)
		}
	}
}

func TestFilePlan(t *testing.T) {
	files := []string{"A", "B"}
	src, tgt, cleanup := mkTestDir(t, files)
	defer cleanup()
	os.Symlink(filepath.Join(src, "A"), filepath.Join(tgt, "A"))
}

func TestIsEffectiveSymlink(t *testing.T) {
	files := []string{"A", "B"}
	src, tgt, cleanup := mkTestDir(t, files)
	defer cleanup()

	linkTarget := filepath.Join(src, "A")
	os.Symlink(linkTarget, filepath.Join(tgt, "A"))
	// Correct situation, link points to the same thing as check (think re-Stow)
	if !IsEffectiveLink(filepath.Join(tgt, "A"), linkTarget) {
		t.Errorf("Expected symlink target %s to match", linkTarget)
	}

	// Regular file
	ioutil.WriteFile(filepath.Join(tgt, "C"), []byte("regular file"), 0444)
	if IsEffectiveLink(filepath.Join(tgt, "C"), filepath.Join(src, "C")) {
		t.Error("No way that's a symlink")
	}

	// Directory
	os.Mkdir(filepath.Join(tgt, "D"), 0775)
	if IsEffectiveLink(filepath.Join(tgt, "D"), filepath.Join(src, "D")) {
		t.Error("No way that's a symlink")
	}

	// Non-existent file
	if IsEffectiveLink(filepath.Join(tgt, "E"), filepath.Join(src, "E")) {
		t.Error("No way that's a symlink")
	}

}
