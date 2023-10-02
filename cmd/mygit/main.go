package main

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"os"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...

func GitInit() error {
	for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return errors.New(fmt.Sprintf("Error creating directory: %s\n", err))
		}
	}
	headFileContents := []byte("ref: refs/heads/master\n")
	if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
		return errors.New(fmt.Sprintf("Error writing file: %s\n", err))
	}
	fmt.Println("Initialized git directory")
	return nil
}

func GitCatFile(hash, flag string) error {
	filePath := fmt.Sprintf(".git/objects/%s/%s", hash[:2], hash[2:])
	file, err := os.ReadFile(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("Error reading file: %s\n", err))
	}
	b := bytes.NewReader(file)
	data, err := zlib.NewReader(b)
	if err != nil {
		return errors.New(fmt.Sprintf("Error uncompressing file: %s\n", err))
	}
	if flag == "-p" {
		io.Copy(os.Stdout, data)
		data.Close()
	}
	return nil
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		err := GitInit()
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}
	case "cat-file":
		err := GitCatFile(os.Args[3], os.Args[2])
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
