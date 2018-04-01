// Copyright 2018 Maxim R. <epiqm.github.io>

// Pack or unpack files and/or directories.
// Create or extract binary packages created with PAQ.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// CLI information constants
const (
	SIG      string = "EPAQ"
	VER      string = "0"
	VERI     string = "1"
	DESC     string = "pack or unpack files and/or directories"
	USAGE    string = "Usage:\n\tpaq command [path..] [options..]"
	COMMANDS string = "Commands:\n\t" +
		"pack    create a package\n\t" +
		"unpack  extract a package\n\t" +
		"scan    analyze package and output meta\n\t" +
		"help    show instructions\n\t" +
		"version display release numbers"
	OPTIONS string = "Options:\n\t-o     package filename or extract directory"
	HELP    string = "Use -h or help to display available commands and options."
	NL      string = "\n"
	SP      string = " "
	NLL     string = "\n\n"
	NR      string = "\n\t"
	TR      string = "\r"
	SEP     string = string(filepath.Separator)
	SUCCESS string = "done."
)

// Type Package contains specified filename
// with data used for meta table generation.
type Package struct {
	Name  string
	Files []Meta
}

// Type Meta contains filename, start offset,
// and total size (in bytes).
type Meta struct {
	Name   string
	Offset int
	Size   int
}

var (
	// Regular expression to parse file meta table
	// that is used for extraction.
	MetaRegex = regexp.MustCompile(`\:\/\/(.?\S+)+\>\>(\d*\.?\d+)+`)

	// RawPackage object contains filenames and meta tables.
	RawPackage Package
)

// Packs files/directories into single binary.
// Creates meta table and appends binary content.
func Pack(ofname string, fplist []string) (result bool) {
	var pind int = 1           // assuming its first package in current directory
	var efname string = ofname // save original package name for iterations

	for _, v := range fplist {
		// open file to get stats
		finfo, err := os.Stat(v)
		if err != nil {
			fmt.Println(err)
			return
		}

		// check if argument is a directory
		if finfo.IsDir() {
			dirFiles, err := ioutil.ReadDir(v)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, f := range dirFiles {
				// strip suffix slash
				if strings.HasSuffix(v, SEP) {
					v = v[:len(v)-1]
				}

				// append each file in directory
				fplist = append(fplist, fmt.Sprintf("%s%s%s", v, SEP, f.Name()))
			}
		}
	}

	// check packages with same filename
	for {
		if _, err := os.Stat(ofname); os.IsNotExist(err) {
			break
		}
		pind++
		ofname = fmt.Sprintf("%s%d", efname, pind)
	}

	// create package
	fo, err := os.Create(ofname)
	if err != nil {
		fmt.Println(err)
		return
	}

	// free package after writing
	defer func() {
		if err := fo.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	// write signature and packaging method version
	// (SIG + VER + VERI + SP)
	var aheader string = SIG + VER + VERI
	if _, err := fo.Write([]byte(aheader)); err != nil {
		fmt.Println(err)
		return
	}

	// generate meta table
	var ameta string
	for _, v := range fplist {
		// open file to get stats
		finfo, err := os.Stat(v)
		if err != nil {
			fmt.Println(err)
			return
		}

		// skip if its a directory
		if finfo.IsDir() {
			continue
		}

		// create meta for every single unit
		ameta = ameta + fmt.Sprintf("://%s>>%d?", finfo.Name(), finfo.Size())
	}

	// first pass:
	// write meta string to package
	if _, err := fo.Write([]byte(ameta + ";;")); err != nil {
		fmt.Println(err)
		return
	}

	// second pass:
	// copy files/directories content to package
	for _, v := range fplist {
		// open file to get stats
		finfo, err := os.Stat(v)
		if err != nil {
			fmt.Println(err)
			return
		}

		// skip if its a directory
		if finfo.IsDir() {
			continue
		}

		// open file for reading
		fi, err := os.Open(v)
		if err != nil {
			fmt.Println(err)
			return
		}

		// close file after complete
		defer func() {
			if err := fi.Close(); err != nil {
				fmt.Println(err)
				return
			}
		}()

		buf := make([]byte, 1024) // create a buffer
		for {
			n, err := fi.Read(buf) // read file to buffer

			// read until the end
			if err != nil && err != io.EOF {
				fmt.Println(err)
				return
			}

			// break when theres nothing to read
			if n == 0 {
				break
			}

			// write content to package or output an error
			if _, err := fo.Write(buf[:n]); err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	return true
}

func Scan(fplist []string, RawPackage *Package) (result bool) {
	for _, v := range fplist {
		// get file stats
		finfo, err := os.Stat(v)
		if err != nil {
			fmt.Println(err)
			return
		}

		// output filename
		fmt.Println(fmt.Sprintf("%s:", finfo.Name()))

		// open file for reading
		fi, err := os.Open(v)
		if err != nil {
			fmt.Println(err)
			return
		}

		// close file after complete
		defer func() {
			if err := fi.Close(); err != nil {
				fmt.Println(err)
				return
			}
		}()

		fi.Seek(6, os.SEEK_SET) // offset after tool name and version

		buf := make([]byte, 1)

		var meta string
		var headerlen int = 6

		for {
			n, err := fi.Read(buf) // read a byte from buffer

			// return on error or end of the file
			if err != nil && err != io.EOF {
				fmt.Println(err)
				return
			}
			if n == 0 {
				// break when no bytes left
				break
			}

			// calculate header length
			headerlen += n

			meta += string(buf[:n])             // append meta content
			match := strings.Index(meta, "?;;") // meta end separator

			// calculate offset of data
			if match > -1 {
				meta = meta[:match+3]
				break
			}
		}

		// split meta content using delimiter
		fsplit := strings.Split(meta, "?")
		if len(fsplit) > 0 {
			poffset := headerlen

			var fid int     // index of file
			var foffset int // offset in bytes
			var fsize int   // total size in bytes

			// iteratre through single meta table element
			for _, mv := range fsplit {
				// find element that matches expression
				metaInfo := MetaRegex.FindAllStringSubmatch(mv, -1)

				if len(metaInfo) == 1 {
					// calculate new offset to get next element
					if fid < 1 {
						foffset = poffset
					} else if fid >= 1 {
						foffset += fsize
					}

					// convert size to bytes integer
					fsize, err = strconv.Atoi(string(metaInfo[0][2]))
					if err != nil {
						fmt.Println(err)
						return
					}

					fid += 1

					// fill single meta element
					fm := Meta{}
					fm.Name = string(metaInfo[0][1])
					fm.Offset = foffset
					fm.Size = fsize

					// append to a raw package
					RawPackage.Files = append(RawPackage.Files, fm)

					// output filename, size, and offset
					fmt.Println(fmt.Sprintf("  %s (%d bytes, offset %d)",
						fm.Name, fm.Size, fm.Offset))
				}
			}
		}
	}

	return true
}

func Unpack(fplist []string) (result bool) {
	Scan(fplist, &RawPackage) // scan the list and create meta

	// go through the specified files
	for _, fa := range RawPackage.Files {

		// create a file for element contents
		fo, err := os.Create(fa.Name)
		if err != nil {
			fmt.Println(err)
			return
		}

		// close file when done
		defer func() {
			if err := fo.Close(); err != nil {
				fmt.Println(err)
				return
			}
		}()

		// open file for output
		fi, err := os.Open(fplist[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		// close file after output
		defer func() {
			if err := fi.Close(); err != nil {
				fmt.Println(err)
				return
			}
		}()

		// set offset location to start reading
		fi.Seek(int64(fa.Offset), os.SEEK_SET)

		// validate final size
		var bufsize int = 1024
		if bufsize > fa.Size {
			bufsize = fa.Size
		}

		// create a buffer
		buf := make([]byte, bufsize)
		for {
			// read n bytes from buffer
			n, err := fi.Read(buf)

			// check for error or end of package
			if err != nil && err != io.EOF {
				fmt.Println(err)
				return
			}

			// write n bytes from buffer
			if _, err := fo.Write(buf[:n]); err != nil {
				fmt.Println(err)
				return
			}

			// if processed or overflow
			if n == 0 || n >= fa.Size {
				break
			}
		}
	}

	return true
}

func main() {
	args := os.Args[1:] // get the arguments

	var cmd string = "pack"          // a default command
	var ofname string = "package.pq" // default package name
	var fplist []string              // list of filenames

	if len(args) > 0 {
		for k, v := range args {
			if v == "-h" || v == "--help" || v == "help" {
				fmt.Println(DESC + NLL +
					USAGE + NLL +
					COMMANDS + NLL +
					OPTIONS)
				return
			}
			if v == "-v" || v == "--version" || v == "version" {
				fmt.Println(VER + "." + VERI)
				return
			}
			if v == "pack" && k == 0 {
				continue
			}
			if v == "unpack" && k == 0 {
				cmd = "unpack"
				continue
			}
			if v == "scan" && k == 0 {
				cmd = "scan"
				continue
			}
			if string(v[0]) == "-" {
				continue
			}
			if len(args) > 2 && k > 1 {
				optn := args[k-1]
				if optn == "-o" || optn == "--output" {
					ofname = v
					continue
				}
			}

			fplist = append(fplist, v)
		}

		if cmd == "pack" {
			if Pack(ofname, fplist) == true {
				fmt.Println(SUCCESS)
			} else {
				// remove temp package
				os.Remove(ofname)
			}
		} else if cmd == "scan" {
			if Scan(fplist, &RawPackage) == true {
				fmt.Println(SUCCESS)
			}
		} else if cmd == "unpack" {
			if Unpack(fplist) == true {
				fmt.Println(SUCCESS)
			}
		}

		return
	}

	// output description and usage
	fmt.Println(DESC + NLL +
		USAGE + NLL +
		TR + HELP)
}
