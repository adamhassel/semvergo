package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/adamhassel/semvergo/pkg/flags"
	"github.com/adamhassel/semvergo/pkg/semver"
)

var incMajor, incMinor, incPatch, usetags, usebranch flags.Bool
var version, prefix, suffix, prefixSeparator, suffixSeparator, gitdir flags.String

func init() {
	flag.Var(&version, "v", "version string to use")
	flag.Var(&incMajor, "major", "increment major version")
	flag.Var(&incMinor, "minor", "increment minor version")
	flag.Var(&incPatch, "patch", "increment patch version. This is the default if no other increments are set.")
	flag.Var(&prefix, "prefix", "prefix to add to semver string")
	flag.Var(&suffix, "suffix", "suffix to add to semver string")
	flag.Var(&prefixSeparator, "prefix-sep", "prefix separator used to separate prefix from  semver string. Used both for parsing and constructing. Default is empty")
	flag.Var(&suffixSeparator, "suffix-sep", "suffix separator used to separate semver string from suffix. Used both for parsing and constructing. Default is '-'")

	flag.Var(&usetags, "tags", "use latest tag on git repository as version string")
	flag.Var(&usebranch, "branch", "use branch name as suffix. When used with -tags, the version number used as input is the latest tag suffixed with the branch name")
	flag.Var(&gitdir, "gitdir", "git directory. Default is current directory.")
}

func main() {
	flag.Parse()

	var sv semver.SemVer

	if !suffixSeparator.IsSet() {
		suffixSeparator.Set("-")
	}

	sv.Presep(prefixSeparator.String())
	sv.Sufsep(suffixSeparator.String())

	switch {
	case usetags.Bool():
		dir := gitdir.String()

		if dir == "" {
			var err error
			dir, err = os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
		}
		repo, err := git.PlainOpen(dir)
		if err != nil {
			log.Fatal(err)
		}
		sv, err = getLatestGitVersion(repo, usebranch.Bool(), suffixSeparator.String())
	case version.IsSet() && version.String() != "":
		var err error
		sv, err = semver.ParseSeparated(version.String(), prefixSeparator.String(), suffixSeparator.String())
		if err != nil {
			log.Fatal(err)
		}
	}

	if incMajor.IsSet() && incMajor.Bool() {
		sv.IncrementMajor()
	}
	if incMinor.IsSet() && incMinor.Bool() {
		sv.IncrementMinor()
	}
	if incPatch.IsSet() && incPatch.Bool() {
		sv.IncrementPatch()
	}

	if !incMajor.IsSet() && !incMinor.IsSet() {
		sv.IncrementPatch()
	}

	if suffix.IsSet() {
		sv.Suffix(suffix.String())
	}
	if prefix.IsSet() {
		sv.Prefix(prefix.String())
	}

	fmt.Printf(sv.String())
}

func currentBranch(repo *git.Repository) string {
	head, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}
	return head.Name().Short()
}

func branchTags(repo *git.Repository) []string {
	tags, err := repo.Tags()
	if err != nil {
		log.Fatal(err)
	}
	var rv []string
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		rv = append(rv, ref.Name().Short())
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return rv
}

// getLatestGitVersion returns the latest version tag from the repository's tags. If `branch` is true, will only look at version tags suffixed with the branch name. sufsep is the suffix separator.
func getLatestGitVersion(repo *git.Repository, branch bool, sufsep string) (semver.SemVer, error) {
	tags := branchTags(repo)

	thisbranch := currentBranch(repo)
	var vs []semver.SemVer
	for _, tag := range tags {
		v, err := semver.ParseSeparated(tag, "", sufsep)
		if err != nil {
			continue
		}
		_, suffix := v.PreSuffix()
		if branch && suffix != thisbranch {
			continue
		}
		vs = append(vs, v)
	}
	rv := semver.Max(vs)
	if branch {
		rv.Sufsep(sufsep)
		rv.Suffix(thisbranch)
	}
	return rv, nil
}
