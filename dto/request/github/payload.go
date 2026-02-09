package github

type PushPayload struct {
	Ref        string `json:"ref"`
	After      string `json:"after"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
	HeadCommit struct {
		Message   string `json:"message"`
		ID        string `json:"id"`
		Timestamp string `json:"timestamp"`
	}
	Commits []struct {
		ID        string   `json:"id"`
		Timestamp string   `json:"timestamp"`
		Added     []string `json:"added"`
		Modified  []string `json:"modified"`
		Removed   []string `json:"removed"`
	} `json:"commits"`
}

type PushRequestPayload struct {
	Action     string `json:"action"`
	Number     int    `json:"number"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	PullRequest struct {
		Merged  bool   `json:"merged"`
		MergeAt string `json:"merge_at"`
	} `json:"pull_request"`
}
