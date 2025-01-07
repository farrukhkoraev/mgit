package main

import (
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path"
)

const (
	BLOB   = "blob"
	TREE   = "tree"
	COMMIT = "commit"

	GIT_DIR = ".git"
)

func getBlobObject(data []byte) []byte {
	header := []byte(fmt.Sprintf("%s %d\x00", BLOB, len(data)))
	return append(header, data...)
}
func writeGitObject(objType string, data []byte) string {

	obj := getBlobObject(data)
	objSHA := sha1.Sum(obj)

	basePath := path.Join(GIT_DIR, "objects")

	err := os.Mkdir(
		path.Join(basePath, fmt.Sprintf("%x", objSHA[0])),
		0755,
	)
	wr := zlib.NewWriter(os.Stdout)
	defer wr.Close()

	if err != nil && !os.IsExist(err) {
		fmt.Printf("Failed to create directory %s\n", fmt.Sprintf("%x", objSHA[0]))
	}

	filepath := path.Join(basePath, fmt.Sprintf("%x", objSHA[0]), fmt.Sprintf("%x", objSHA[1:]))
	_, err = os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			defer file.Close()

			writer := zlib.NewWriter(file)
			defer writer.Close()

			_, err = writer.Write(obj)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		} else {
			fmt.Println(err.Error())
			os.Exit(1)

		}
	}

	return fmt.Sprintf("%x", objSHA)
}

func readGitObject(objSHA string) []byte {
	var data []byte
	filepath := path.Join(
		GIT_DIR,
		"objects",
		objSHA[:2],
		objSHA[2:],
	)

	_, err := os.Stat(filepath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	reader, err := zlib.NewReader(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer reader.Close()

	data, err = io.ReadAll(reader)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(data)

	return data
}
