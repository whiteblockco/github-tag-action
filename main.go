package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

type VersionTag struct {
	ref   *plumbing.Reference
	Tag   string
	Major int
	Minor int
	Patch int
	Pre   string
	Build string
}

func (v *VersionTag) String() string {
	ret := ""

	ret += v.Tag // append tag name

	ret += fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch) // append body

	if len(v.Pre) > 0 {
		ret += "-" + v.Pre //append pre-release
	}

	if len(v.Build) > 0 {
		ret += "+" + v.Build // append build-metadata
	}

	return ret
}

func VersionFromString(str string) (*VersionTag, error) {
	if !semverRegex.MatchString(str) {
		return nil, fmt.Errorf("invalid tag format: <%s>", str)
	}

	semver := semverRegex.FindStringSubmatch(str)

	tag := semver[1]
	major, _ := strconv.Atoi(semver[2])
	minor, _ := strconv.Atoi(semver[3])
	patch, _ := strconv.Atoi(semver[4])
	pre := semver[5]
	build := semver[6]

	return &VersionTag{
		Tag:   tag,
		Major: major,
		Minor: minor,
		Patch: patch,
		Pre:   pre,
		Build: build,
	}, nil
}

var (
	releaseBranchRegex = regexp.MustCompile("release/(0|[1-9]\\d*)\\.(0|[1-9]\\d*)")
	semverRegex        = regexp.MustCompile("^([a-z]*)(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$")
	//preRegex           = regexp.MustCompile("([a-zA-Z]+\\.)?(0|[1-9]\\d*)")
)

func main() {
	r, _ := git.PlainOpen("./")

	h, err := r.Head()
	if err != nil {
		panic(err)
	}

	if !h.Name().IsBranch() {
		panic("release workflow must be branch")
	}

	branchName := h.Name().String()

	if !releaseBranchRegex.MatchString(branchName) {
		panic(fmt.Errorf("not matching branch name pattern: wanted %s, got %s ", releaseBranchRegex.String(), branchName))
	}

	major, _ := strconv.Atoi(releaseBranchRegex.FindStringSubmatch(branchName)[1])
	minor, _ := strconv.Atoi(releaseBranchRegex.FindStringSubmatch(branchName)[2])

	tags, err := r.Tags()
	if err != nil {
		panic(err)
	}

	latest := &VersionTag{
		ref:   h,
		Tag:   "v",
		Major: major,
		Minor: minor,
		Patch: 0,
		Pre:   "",
		Build: "",
	}

	err = tags.ForEach(func(ref *plumbing.Reference) error {
		current, err := parseTag(ref)
		if err != nil {
			return nil
		}

		// not a tag of this release
		if current.Major != major || current.Minor != minor {
			return nil
		}

		if isNewerVersion(latest, current) {
			latest = current
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	latest.Patch++

	message, err := summeryCommitMessage(r, latest)
	if err != nil {
		message = fmt.Sprintf("Failed summery commit messages <%s>", err)
	}

	opts := &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  "whiteblock",
			Email: "developer@whiteblock.co",
			When:  kst(),
		},
		Message: message,
		SignKey: nil,
	}
	// Summery commit messages to write description of tag

	c, err := getHeadCommit(r)
	if err != nil {
		panic(err)
	}

	err = opts.Validate(r, c.Hash)
	if err != nil {
		panic(err)
	}


	_, err = r.CreateTag(latest.String(), c.Hash, opts)
	if err != nil {
		panic(err)
	}

	Info("Latest commit: ", c)
	refSpec := fmt.Sprintf("+refs/tags/%s:refs/tags/%s", latest.String(), latest.String())
	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: "USER_NAME", // this can be anything except an empty string
			Password: os.Getenv("REPO_TOKEN"),
		},
		RefSpecs: []config.RefSpec{config.RefSpec(refSpec)},
	})
	if err != nil {
		panic(err)
	}

	Info("Success to bump version: %s", latest.String())
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	log.Printf("[INFO]"+format, args...)
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	log.Printf("[WARN]"+format, args...)
}

func parseTag(ref *plumbing.Reference) (*VersionTag, error) {
	version, err := VersionFromString(ref.Name().Short())
	if err != nil {
		return nil, err
	}

	version.ref = ref

	return version, nil
}

func getHeadCommit(r *git.Repository) (*object.Commit, error) {
	head, err := r.Head()
	if err != nil {
		return nil, err
	}

	cIter, err := r.Log(&git.LogOptions{From: head.Hash()})
	if err != nil {
		return nil, err
	}

	return cIter.Next()
}

func summeryCommitMessage(r *git.Repository, prevLatestTag *VersionTag) (string, error) {
	head, err := r.Head()
	if err != nil {
		Warning("Failed to get head reference: %s", err.Error())
		return "", err
	}

	cIter, err := r.Log(&git.LogOptions{From: head.Hash()})
	if err != nil {
		Warning("Failed to iterate git log")
		return "", err
	}

	var messages []string
	commit, err := cIter.Next()
	if err != nil {
		Warning("Failed to get head commit: %s", err.Error())
		return "", err
	}

	h, err := r.ResolveRevision(plumbing.Revision(prevLatestTag.ref.Hash().String()))
	if err != nil {
		Warning("Failed to get latest tag: %s", err.Error())
		return "", err
	}
	obj, err := r.Object(plumbing.AnyObject, *h)
	prevTagCommit := obj.(*object.Commit)
	for commit != nil && commit.Hash != prevTagCommit.Hash {
		messages = append(messages, commit.Message)
		commit, err = cIter.Next()
		if err != nil {
			break
		}
	}

	summery := ""
	for i := range messages {
		summery += "* " + messages[i]
	}
	if summery == "" {
		summery = "Nothing new, just for tagging."
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
				if old.Build > new.Build {
					return false
				}
			}
		}
	}
	return true
}

func kst() time.Time {
	loc, _ := time.LoadLocation("Asia/Seoul")
	return time.Now().In(loc)
}
