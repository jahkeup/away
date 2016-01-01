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
