package models

type VisualizerTheme struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

type Visualizer struct {
	ID                    string            `json:"id"`
	Name                  string            `json:"name"`
	Description           string            `json:"description"`
	Category              string            `json:"category"`
	Tags                  []string          `json:"tags"`
	PreviewType           string            `json:"previewType"`
	ImageURLTemplate      string            `json:"imageUrlTemplate"`
	MarkdownTemplate      string            `json:"markdownTemplate"`
	WebsiteURL            string            `json:"websiteUrl"`
	RepositoryURL         string            `json:"repositoryUrl"`
	Themes                []VisualizerTheme `json:"themes"`
	DefaultTheme          string            `json:"defaultTheme"`
	RequiresUsername      bool              `json:"requiresUsername"`
	RequiresExternalSetup bool              `json:"requiresExternalSetup"`
	SetupExplanation      string            `json:"setupExplanation,omitempty"`
	ActionWorkflowYaml    string            `json:"actionWorkflowYaml,omitempty"`
	PrivacyNotice         string            `json:"privacyNotice,omitempty"`
	AddedAt               string            `json:"addedAt"`
	PopularityRank        int               `json:"popularityRank"`
	Enabled               bool              `json:"enabled"`
	EstimatedHeight       int               `json:"estimatedHeight"`
	VoteCount             int               `json:"voteCount"`
	UserVoted             bool              `json:"userVoted"`
	GitHubStars           int               `json:"githubStars"`
	LastRefreshedAt       string            `json:"lastRefreshedAt,omitempty"`
	NextRefreshAt         string            `json:"nextRefreshAt,omitempty"`
}

type VoteRequest struct {
	DeviceID string `json:"deviceId"`
}

type VoteResponse struct {
	Voted     bool `json:"voted"`
	VoteCount int  `json:"voteCount"`
}
