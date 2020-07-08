package main

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type VersionTag struct {
	Tag         object.Tag
	Major       int
	Minor       int
	Patch       int
	BuildNumber int
}

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

func parseTag(tag object.Tag) (VersionTag, error) {
	arr := strings.Split(tag.Name, ".")
	if len(arr) != 3 {
		return VersionTag{}, errors.New(fmt.Sprintf("Invalid tag format: <%s>", tag.Name))
	}

	major, _ := strconv.Atoi(arr[0])
	minor, _ := strconv.Atoi(arr[1])
	patch, _ := strconv.Atoi(arr[2])
	buildNumber := 0
	if strings.Index(arr[2], "-") != -1 {
		str := strings.Split(arr[2], "-")
		patch, _ = strconv.Atoi(str[0])
		buildNumber, _ = strconv.Atoi(str[1])
	}
	return VersionTag{tag, major, minor, patch, buildNumber}, nil
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

func summeryCommitMessage(r *git.Repository, prevLatestTag *VersionTag) (string, error) {
	head, err := r.Head()
	if err != nil {
		Warning("Failed to get head reference: %s", err.Error())
		return "", err
	}

	cIter, err := r.Log(&git.LogOptions{From: head.Hash()})
	var messages []string
	commit, err := cIter.Next()
	if err != nil {
		Warning("Failed to get head commit: %s", err.Error())
		return "", err
	}
	prevTagCommit, err := prevLatestTag.Tag.Commit()
	if err != nil {
		Warning("Failed to get latest tag: %s", err.Error())
		return "", err
	}
	for commit != nil && commit.Hash != prevTagCommit.Hash {
		messages = append(messages, commit.Message)
		commit, err = cIter.Next()
		if err != nil {
			break
		}
	}

	summery := ""
	for i := range messages {
		summery += "*" + messages[i]
	}
	if summery == "" {
		summery = "Nothing new, Just for tagging."
	}
	return summery, nil
}

func isNewerVersion(old, new *VersionTag) bool {
	if old.Major > new.Major {
		return false
	} else if old.Major == new.Major {
		if old.Minor > new.Minor {
			return false
		} else if old.Minor == new.Minor {
			if old.Patch > new.Patch {
				return false
			} else if old.Patch == new.Patch {
				if old.BuildNumber > new.BuildNumber {
					return false
				}
			}
		}
	}
	return true
}

func getLatestTag(tagIter *object.TagIter) (VersionTag, error) {
	latestTag := VersionTag{
		Major:       0,
		Minor:       0,
		Patch:       0,
		BuildNumber: 0,
	}
	if err := tagIter.ForEach(func(t *object.Tag) error {
		tmpTag, err := parseTag(*t)
		if err != nil {
			return err
		}
		if isNewerVersion(&latestTag, &tmpTag) {
			latestTag = tmpTag
		}
		return nil
	}); err != nil {
		return latestTag, err
	}
	return latestTag, nil
}

func koreanTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Seoul")
	return time.Now().In(loc)
}

func main() {
	r, _ := git.PlainOpen("./")

	tags, err := r.TagObjects()
	ExitIfError(err)

	latestTag, err := getLatestTag(tags)
	ExitIfError(err)

	Info("::debug:: latestTag: %s", latestTag.Tag.Name)
	if latestTag.BuildNumber == 0 {
		latestTag.Patch++
	}
	latestTag.BuildNumber++
	// Increase build number

	message, err := summeryCommitMessage(r, &latestTag)
	if err != nil {
		message = fmt.Sprintf("Failed summery commit messages <%s>", err)
	}
	opts := &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  "whiteblock",
			Email: "developer@whiteblock.co",
			When:  koreanTime(),
		},
		Message: message,
		SignKey: nil,
	}
	// Summery commit messages to write description of tag

	c, err := getHeadCommit(r)
	ExitIfError(err)
	err = opts.Validate(r, c.Hash)
	newTag := fmt.Sprintf("%d.%d.%d-%d", latestTag.Major, latestTag.Minor, latestTag.Patch, latestTag.BuildNumber)
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
	Info("Success to bump version: %s", newTag)
}
