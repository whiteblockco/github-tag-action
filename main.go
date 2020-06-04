package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var versionFormatWithBuildNumber = "%d.%d.%d-%d"

// ExitIfError should be used to naively panics if an error is not nil.
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

func parseTag(tag string) (int, int, int, int) {
	arr := strings.Split(tag, ".")
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

func isNewerVersion(old, new string) bool {
	oMajor, oMinor, oPatch, oBuildNumber := parseTag(old)
	nMajor, nMinor, nPatch, nBuildNumber := parseTag(new)
	if oMajor > nMajor {
		return false
	} else if oMajor == nMajor {
		if oMinor > nMinor {
			return false
		} else if oMinor == nMinor {
			if oPatch > nPatch {
				return false
			} else if oPatch == nPatch {
				if oBuildNumber > nBuildNumber {
					return false
				}
			}
		}
	}
	return true
}

func getLatestTag(tagIter *object.TagIter) string {
	latestTag := "0.0.0-0"
	tagIter.ForEach(func(t *object.Tag) error {
		if isNewerVersion(latestTag, t.Name) {
			major, minor, patch, buildNumber := parseTag(t.Name)
			latestTag = fmt.Sprintf(versionFormatWithBuildNumber, major, minor, patch, buildNumber)
		}
		return nil
	})
	return latestTag
}

func koreanTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Seoul")
	return time.Now().In(loc)
}

func main() {
	r, _ := git.PlainOpen("./")

	tags, err := r.TagObjects()
	ExitIfError(err)

	latestTag := getLatestTag(tags)
	major, minor, patch, buildNumber := parseTag(latestTag)
	buildNumber++
	newTag := fmt.Sprintf("%d.%d.%d-%d", major, minor, patch, buildNumber)
	fmt.Println(fmt.Sprintf("New tag: <%s>", newTag))
	c, err := getHeadCommit(r)
	ExitIfError(err)

	opts := &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  "whiteblock",
			Email: "developer@whiteblock.co",
			When: koreanTime(),
		},
		Message: "message",
		SignKey: nil,
	}
	err = opts.Validate(r, c.Hash)
	_, err = r.CreateTag(newTag, c.Hash, opts)
	ExitIfError(err)
	fmt.Println("Latest commit: ", c)
	refSpec := fmt.Sprintf("+refs/tags/%s:refs/tags/%s", newTag, newTag)
	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: "USER_NAME", // this can be anything except an empty string
			Password: os.Getenv("REPO_TOKEN"),
		},
		RefSpecs: []config.RefSpec{config.RefSpec(refSpec)},
	})
	ExitIfError(err)
}
