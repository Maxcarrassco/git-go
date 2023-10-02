package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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
		return errors.New(fmt.Sprintf("Error reading compress file: %s\n", err))
	}
	b := bytes.NewReader(file)
	data, err := zlib.NewReader(b)
	defer data.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("Error uncompressing file: %s\n", err))
	}
	if flag == "-p" {
		data, err := io.ReadAll(data)
		if err != nil {
			return errors.New(fmt.Sprintf("Error reading uncompress file: %s\n", err))
		}
		str := strings.Split(string(data), "\x00")
		fmt.Print(str[1])
	}
	return nil
}

func GitHashObject(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return errors.New(fmt.Sprintf("Error reading compress file: %s\n", err))
	}
	f_len := len(file)
	var buf bytes.Buffer
	content := string(file)
	data := []byte(fmt.Sprintf("blob %d\x00%s", f_len, content))
	compressData := zlib.NewWriter(&buf)
	_, err = compressData.Write(data)
	if err != nil {
		return errors.New(fmt.Sprintf("Error writing file: %s\n", err))
	}
	compressData.Close()
	sum := sha1.Sum(data)
	key := strings.ReplaceAll(string(fmt.Sprintf("% x\n", sum)), " ", "")
	key = strings.ReplaceAll(key, "\n", "")
	key = strings.ReplaceAll(key, "$", "")
	filePath := fmt.Sprintf(".git/objects/%s/%s", key[:2], key[2:])
	os.MkdirAll(".git/objects/"+key[:2], 0755)
	f, err := os.Create(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("Error writing file: %s\n", err))
	}
	_, err = f.Write(buf.Bytes())
	if err != nil {
		return errors.New(fmt.Sprintf("Error writing file: %s\n", err))
	}
	fmt.Println(key)
	return nil
}

func main() {
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
	case "hash-object":
		err := GitHashObject(os.Args[3])
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
