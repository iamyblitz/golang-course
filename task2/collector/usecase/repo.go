package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	pb "github.com/iamyblitz/golang-course/task2/proto"
)

var ErrRepoNotFound = errors.New("repository not found")

type RepoUsecase struct {
	client *http.Client
}

func NewRepoUsecase() *RepoUsecase {
	return &RepoUsecase{
		client: &http.Client{},
	}
}

type githubRepoResponse struct {
	FullName        string `json:"full_name"`
	Description     string `json:"description"`
	StargazersCount int32  `json:"stargazers_count"`
	ForksCount      int32  `json:"forks_count"`
	CreatedAt       string `json:"created_at"`
}

func (u *RepoUsecase) GetRepo(ctx context.Context, owner, repo string) (*pb.RepoResponse, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build github request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "task2-collector")

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request github api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrRepoNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	var ghResp githubRepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&ghResp); err != nil {
		return nil, fmt.Errorf("decode github response: %w", err)
	}

	name := ghResp.FullName
	if name == "" {
		name = owner + "/" + repo
	}

	return &pb.RepoResponse{
		Name:        name,
		Description: ghResp.Description,
		Stars:       ghResp.StargazersCount,
		Forks:       ghResp.ForksCount,
		CreatedAt:   ghResp.CreatedAt,
	}, nil
}
