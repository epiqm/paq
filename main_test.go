// Copyright 2018 Maxim R. <epiqm.github.io>

// Packs files into a package then unpacks.
// Compares extracted files with original ones.
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPack(t *testing.T) {
	// get working directory
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "paq") {
		wd = filepath.Dir(wd)
	}

	var packageName string = "test.pq" // create package name
	var filesPath []string             // create files hashmap

	// remove previous package
	os.Remove(packageName)

	// append files to package
	filesPath = append(filesPath, fmt.Sprintf("%s/tests/note.txt", wd))
	filesPath = append(filesPath, fmt.Sprintf("%s/tests/rabbit.png", wd))

	// test Pack function
	result := Pack(packageName, filesPath)
	if result != true {
		t.Error("Failed to create package.")
	}

	// total files size
	var totalFilesSize int64
	for _, v := range filesPath {
		finfo, _ := os.Stat(v)
		totalFilesSize += finfo.Size()
	}

	// package size
	pkginfo, _ := os.Stat(packageName)
	pkgSize := pkginfo.Size()

	// compare package size
	if pkgSize < totalFilesSize {
		t.Error("Package appears to be corrupt.")
	}
}

func TestScan(t *testing.T) {
	var packageNames []string
	packageNames = append(packageNames, "test.pq")

	// test scan function
	result := Scan(packageNames, &RawPackage)
	if result != true {
		t.Error("Failed to get package meta info.")
	}
}

func TestUnpack(t *testing.T) {
	// get working directory
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "paq") {
		wd = filepath.Dir(wd)
	}

	var packageNames []string

	// append package that was created
	packageNames = append(packageNames, "test.pq")

	// test unpack function
	result := Unpack(packageNames, "./")
	if result != true {
		t.Error("Failed to extract package.")
	}

	var filesPath []string // original files hashmap
	filesPath = append(filesPath, fmt.Sprintf("%s/tests/note.txt", wd))
	filesPath = append(filesPath, fmt.Sprintf("%s/tests/rabbit.png", wd))

	var extractPath []string // extracted files hashmap
	extractPath = append(filesPath, fmt.Sprintf("%s/note.txt", wd))
	extractPath = append(filesPath, fmt.Sprintf("%s/rabbit.png", wd))

	// compare extracted files to original ones
	for k, v := range filesPath {
		f1, err := ioutil.ReadFile(v)
		if err != nil {
			t.Error("Failed to read original file.")
		}

		f2, err := ioutil.ReadFile(extractPath[k])
		if err != nil {
			t.Error("Failed to read extracted file.")
		}

		if bytes.Equal(f1, f2) == false {
			t.Error("Extracted content does not match original file.")
		}
	}

	// remove temporary files
	os.Remove(packageNames[0])
	os.Remove("note.txt")
	os.Remove("rabbit.png")
}
