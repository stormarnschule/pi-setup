package main

import (
	"bytes"
	"crypto/md5"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	logger  = log.New(os.Stdout, "", 0)
	binDir  string
	options *flags
)

func escapePath(s string) string {
	s = strings.Replace(s, ".", "_", -1)
	s = strings.Replace(s, string(os.PathSeparator), "_", -1)
	s = strings.Replace(s, ":", "_", -1)
	return s
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		logger.Fatalln(err)
	}
	return info.Mode().IsRegular()
}

func removeExt(s string) string {
	oldExt := filepath.Ext(s)
	return s[:len(s)-len(oldExt)]
}

func main() {
	options = parseFlags(os.Args[1:])

	binDir = filepath.Join(filepath.Dir(options.source), "bin")
	binFilePath := filepath.Join(binDir, removeExt(filepath.Base(options.source)))
	hashFilePath := binFilePath + ".hash"

	logger.Println(options.source, "->", binFilePath)

	// compile
	err := compile(options.source, binFilePath, hashFilePath)
	if err != nil {
		logger.Fatalln(err)
	}

	// run
	cmd := exec.Command(binFilePath)
	logger.Println("--- Running ---------------------------------------")
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		logger.Fatalln(err)
	}
}

func compile(srcFile, binFile, hashFilePath string) error {
	// create libFlags
	var libFlags string
	for i, lib := range options.libraries {
		if i != 0 {
			libFlags += " "
		}
		libFlags += "-l" + lib
	}

	// check file changed
	srcFileHash := hashFile(srcFile)
	oldHash, err := ioutil.ReadFile(hashFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			// error reading and hashFile exists
			return err
		}
		// file does not exist, go on
	} else {
		// read oldHash successfully
		if bytes.Equal(srcFileHash, oldHash) {
			// no compilation required
			return nil
		}
		// old hash is not new hash, go on
	}

	// compilation required
	err = os.MkdirAll(binDir, os.ModePerm)
	if err != nil {
		return err
	}

	cmdLine := strings.Join([]string{"g++", srcFile, "-o " + binFile, libFlags}, " ")
	logger.Println("Compiling:", cmdLine)

	cmd := exec.Command("sh", "-c", cmdLine)
	cmd.Stdout = os.Stdout // pipe g++ output

	// run command
	err = cmd.Run()
	if err != nil {
		return err
	}

	// save new hash
	return ioutil.WriteFile(hashFilePath, srcFileHash, os.ModePerm)
}

func hashFile(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		logger.Fatalln(err)
	}
	defer f.Close()

	hasher := md5.New()
	_, err = io.Copy(hasher, f)

	if err != nil {
		logger.Fatalln(err)
	}
	return hasher.Sum(nil)
}
