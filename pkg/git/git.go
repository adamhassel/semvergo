package git

import (
	"log"

	ggit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/adamhassel/semvergo/pkg/semver"
)

func currentBranch(repo *ggit.Repository) string {
	head, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}
	return head.Name().Short()
}

func branchTags(repo *ggit.Repository) []string {
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

// LatestsGitVersionTag returns the latest version tag from the repository's tags. If `branch` is true, will only look at version tags suffixed with the branch name. sufsep is the suffix separator.
func LatestsGitVersionTag(repo *ggit.Repository, branch bool, sufsep string) (semver.SemVer, error) {
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
	rv := semver.MaxSlice(vs)
	if branch {
		rv.Sufsep(sufsep)
		rv.Suffix(thisbranch)
	}
	return rv, nil
}
