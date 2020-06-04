package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// CheckIfError should be used to naively panics if an error is not nil.
func ExitIfError(err error) {
	if err == nil {
		return
	}
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func parseTag(t *object.Tag) (int, int, int, int) {
	arr := strings.Split(t.Name, ".")
	if len(arr) != 3 {
		// TODO: Make to error below
		return 0, 0, 0, 0
	}

	major, _ := strconv.Atoi(arr[0])
	minor, _ := strconv.Atoi(arr[1])
	var patch int
	buildNumber := 0
	if strings.Index(arr[2], "-") != -1 {
		str := strings.Split(arr[2], "-")
		patch, _ = strconv.Atoi(str[0])
		buildNumber, _ = strconv.Atoi(str[1])
	}
	return major, minor, patch, buildNumber
}

func getHeadCommit(r *git.Repository) (*object.Commit, error) {
	head, err := r.Head()
	ExitIfError(err)
	cIter, err := r.Log(&git.LogOptions{From: head.Hash()})
	ExitIfError(err)
	commit, err := cIter.Next()
	ExitIfError(err)
	return commit, err
}
func main() {
	r, _ := git.PlainOpen("./")

	tags, err := r.TagObjects()
	ExitIfError(err)

	var major, minor, patch, buildNumber int
	latestTag, err := tags.Next()
	ExitIfError(err)

	major, minor, patch, buildNumber = parseTag(latestTag)
	buildNumber++
	newTag := fmt.Sprintf("%d.%d.%d-%d", major, minor, patch, buildNumber)
	fmt.Println(newTag)
	c, err := getHeadCommit(r)
	ExitIfError(err)

	fmt.Println("Lastest commit: ", c)
	ref, err := r.CreateTag(newTag, c.Hash, nil)
	ExitIfError(err)

	fmt.Println(ref)
}
