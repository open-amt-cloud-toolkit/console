package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/github"
	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
)

var (
	ErrGithub        = consoleerrors.CreateConsoleError("LatestReleaseHandler")
	ErrFailedToFetch = errors.New("repositoryError")
)

func RepositoryError(status string) error {
	return fmt.Errorf("failed to fetch latest release: %w: %s", ErrFailedToFetch, status)
}

// FetchLatestRelease fetches the latest release information from GitHub API
func FetchLatestRelease(c *gin.Context, repo string) (*github.Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	client := &http.Client{}

	req2, _ := http.NewRequestWithContext(c, http.MethodGet, url, http.NoBody)

	resp, err := client.Do(req2)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrGithub.Wrap("FetchLatestRelease", "http.Get", RepositoryError(resp.Status))
	}

	var release github.Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// LatestReleaseHandler is the Gin handler function to check for the latest release
func LatestReleaseHandler(c *gin.Context) {
	repo := Config.App.Repo

	release, err := FetchLatestRelease(c, repo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"current": Config.App.Version,
			"error":   err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"current": Config.App.Version,
		"latest": map[string]interface{}{
			"tag_name":     release.TagName,
			"name":         release.Name,
			"body":         release.Body,
			"prerelease":   release.Prerelease,
			"created_at":   release.CreatedAt,
			"published_at": release.PublishedAt,
			"html_url":     release.HTMLURL,
			"author":       release.Author,
			"assets":       release.Assets,
		},
	})
}
