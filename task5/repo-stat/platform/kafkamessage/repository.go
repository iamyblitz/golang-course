package kafkamessage

type CollectTask struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type CollectResult struct {
	Owner       string `json:"owner"`
	Repo        string `json:"repo"`
	FullName    string `json:"full_name,omitempty"`
	Description string `json:"description,omitempty"`
	Stars       int64  `json:"stars,omitempty"`
	Forks       int64  `json:"forks,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	Error       string `json:"error,omitempty"`
}
