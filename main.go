package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"os"
)

type InitCmd struct{}
type HashObjectCmd struct {
	W     bool   `short:"w" help:"Write blob"`
	StdIn bool   `name:"stdin" help:"Get input from stdin"`
	File  string `arg:"" help:"File to read from" optional:""`
}
type CatFileCmd struct {
	Type   string `arg:"" required:""`
	Object string `arg:"" required:""`
}

type CLIArg struct {
	Init       *InitCmd       `cmd:"" help:"Create an enpty Git repo"`
	HashObject *HashObjectCmd `cmd:"" help:"Compute object ID and optionally creates a blob from a file"`
	CatFile    *CatFileCmd    `cmd:"" help:"Provide content or type and size information for repository objects"`
}

func parse() *CLIArg {
	arg := &CLIArg{}
	kong.Parse(arg)
	return arg
}

func main() {
	args := os.Args[1:]
	switch args[0] {
	case "hash-object":
		fmt.Println(writeGitObject(BLOB, []byte("Hello. It's a content of blob object.")))
	case "cat-file":
		fmt.Println(string(readGitObject(args[1])))
	}

}
