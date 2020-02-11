package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"os/user"
	"path"
	"strconv"
)

type SeenFile struct {
	Path  string
	Items []string
}

func NewSeenFile(params ...string) *SeenFile {
	var seenPath string

	if len(params) == 0 {
		user, err := user.Current()
		panicOnError(err)

		homeDir := user.HomeDir

		seenPath = path.Join(homeDir, ".config", "transmission", "seen")
	} else {
		seenPath = params[0]
	}

	file, err := os.OpenFile(seenPath, os.O_RDONLY|os.O_CREATE, 0600)
	panicOnError(err)

	defer file.Close()

	var items []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		items = append(items, scanner.Text())
	}

	fmt.Println("SEEN " + strconv.Itoa(len(items)) + " items")

	seen := SeenFile{
		Path:  seenPath,
		Items: items,
	}

	return &seen
}

func (seen *SeenFile) Add(link string) {
	hash := sha256sum(link)

	if seen.IsPresent(link) {
		return
	}

	seen.Items = append(seen.Items, hash)

	file, err := os.OpenFile(seen.Path, os.O_APPEND|os.O_WRONLY, 0600)
	panicOnError(err)
	defer file.Close()

	_, err = file.Write([]byte(hash + "\n"))
	panicOnError(err)

	file.Close()
}

func (seen *SeenFile) IsPresent(link string) bool {
	for _, item := range seen.Items {
		if item == sha256sum(link) {
			return true
		}
	}

	return false
}

func (seen *SeenFile) Count() int {
	return len(seen.Items)
}

func (seen *SeenFile) Clear() {
	seen.Items = []string{}

	err := os.Truncate(seen.Path, 0)
	panicOnError(err)
}

func sha256sum(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
