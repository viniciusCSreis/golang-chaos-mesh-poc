package git

import (
	"io"
)

type RepoVersion string

func (r RepoVersion) String() string {
	return string(r)
}

type Git struct {
	Repos       Repositories
	NewRepoInfo func(url string, token string) RepoInfo
}

type Tag struct {
	Name string `json:"tag_name"`
}

type Tags []Tag

func (t Tags) Names() []string {
	var tags []string
	for i := range t {
		tags = append(tags, t[i].Name)
	}

	return tags
}

type RepoInfo interface {
	ZipUrl(version string) string
	TagsUrl() string
	LatestTagUrl() string
	TokenHeader() string
	Token() string
}

type Repositories interface {
	Zipball(info RepoInfo, version string) (io.ReadCloser, error)
	Tags(info RepoInfo) (Tags, error)
	LatestTag(info RepoInfo) (Tag, error)
}

type RepoProvider string

func (r RepoProvider) String() string {
	return string(r)
}

type RepoProviders map[RepoProvider]Git

func (re RepoProviders) Add(provider RepoProvider, git Git) {
	re[provider] = git
}

func (re RepoProviders) Resolve(provider RepoProvider) Git {
	return re[provider]
}

const (
	GitHub = "github"
)
