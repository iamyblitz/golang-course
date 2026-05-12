package dto

type RepositoryInfoResponse struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int64  `json:"stars"`
	Forks       int64  `json:"forks"`
	CreatedAt   string `json:"created_at"`
}

type RepositoriesInfoResponse struct {
	Repositories []RepositoryInfoResponse `json:"repositories"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
