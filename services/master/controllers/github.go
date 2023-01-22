package controllers

import (
	"strings"

	"os"
	"path/filepath"

	"context"
	_ "io/ioutil"
	"time"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v41/github"

	"github.com/sdslabs/gasper/lib/utils"
	"golang.org/x/oauth2"
)

// Endpoint to create repository in GitHub
func CreateRepository(c *gin.Context) {
	//TODO: not able to receive params from GCTL query, request body is empty
	//fmt.Println("--------------\nRequest Body", c.Request.Body, "\n----------------------")
	//Endpoint works pefectly when called directly, not from GCTL
	filters := utils.QueryToFilter(c.Request.URL.Query())
	repoName, pathToApplication := filters["name"], filters["path"]
	repository, err := CreateRepositoryGithub(repoName.(string))

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}
	_, err = GitPush(pathToApplication.(string), repository.GetCloneURL())

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": repository.GetCloneURL(),
	})
}

//Needs an .env file with USERNAME and PAT
//TODO: Shift env variables to config.toml
func GoDotEnvVariable(key string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	environmentPath := filepath.Join(strings.TrimRight(dir, "/tmp"), ".env")
	err = godotenv.Load(environmentPath)

	if err != nil {
		panic(err)
	}

	return os.Getenv(key)
}

func CreateRepositoryGithub(repoName string) (*github.Repository, error) {
	tc := oauth2.NewClient(
		context.Background(),
		oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: GoDotEnvVariable("PAT"), //PAT
			},
		),
	)
	client := github.NewClient(tc)
	repo := &github.Repository{
		Name:    github.String(repoName),
		Private: github.Bool(true),
	}
	repo, _, err := client.Repositories.Create(context.Background(), "", repo)
	return repo, err
}

func GitInit(directoryPath string) (*git.Repository, error) {
	var (
		err error
	)
	_, err = os.Stat(directoryPath)
	if err != nil {
		return nil, err
	}
	repository, err := git.PlainInit(directoryPath, false)
	return repository, err
}

func GitPush(pathToApplication string, repoURL string) (*git.Repository, error) {
	var firstInit bool
	repo, err := git.PlainOpen(pathToApplication)
	if err != nil {
		firstInit = true
		repo, err = GitInit(pathToApplication)
		if err != nil {
			return nil, err
		}
	} else {
		firstInit = false
	}
	_, err = repo.Remote("origin")
	if err != nil {
		_, err = repo.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{repoURL},
		})
		if err != nil {
			return nil, err
		}
	}
	w, _ := repo.Worktree()
	if firstInit {
		err = w.AddGlob(".")
		if err != nil {
			return nil, err
		}
	} else {
		_, err = w.Add(".")
		if err != nil {
			return nil, err
		}
	}

	_, _ = w.Commit("latest commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  GoDotEnvVariable("USERNAME"),
			Email: GoDotEnvVariable("EMAIL"),
			When:  time.Now(),
		},
	})

	auth := &http.BasicAuth{
		Username: GoDotEnvVariable("USERNAME"),
		Password: GoDotEnvVariable("PAT"),
	}

	if err != nil {
		return nil, err
	}
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	if err != nil {
		return nil, err
	}

	return repo, err
}
