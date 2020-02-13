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

	self := SeenFile{
		Path:  seenPath,
		Items: items,
	}

	return &self
}

func (self *SeenFile) Add(link string) {
	hash := sha256sum(link)

	if self.IsPresent(link) {
		return
	}

	self.Items = append(self.Items, hash)

	file, err := os.OpenFile(self.Path, os.O_APPEND|os.O_WRONLY, 0600)
	panicOnError(err)
	defer file.Close()

	_, err = file.Write([]byte(hash + "\n"))
	panicOnError(err)

	file.Close()
}

func (self *SeenFile) IsPresent(link string) bool {
	for _, item := range self.Items {
		if item == sha256sum(link) {
			return true
		}
	}

	return false
}

func (self *SeenFile) Count() int {
	return len(self.Items)
}

func (self *SeenFile) Clear() {
	self.Items = []string{}

	err := os.Truncate(self.Path, 0)
	panicOnError(err)
}

func sha256sum(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
