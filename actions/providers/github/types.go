package github

// Resource represents GitHub API Resource
type Resource struct {
	URL string `json:"html_url" description:"The URL to visit to view the resource"`
	APIEndpoint string `json:"url" description:"The API endpoint URL that represents this resource"`
}

// a CommitAuthor is an identity in Git
type CommitAuthor struct {
	Name string   `json:"name"`
	Email string  `json:"email"`
}

type Commit struct {
	Resource
	SHA string `json:"sha" description:"The SHA of the commit."`
	Message string `json:"message" description:"The commit message"`
	Author CommitAuthor `json:"author" description:"The git author of this commit"`
	IsDistinct bool `json:"distinct" description:"Whether this commit is distinct from any that were pushed before."`
}

// Repository is a GitHub repository
type Repository struct {
	Resource
	ID string `json:"id" description:"Github's ID for the repository'"`
	Name string `json:"name "description:"The name of the repository"`
	FullName string `json:"full_name" description:"The full name of the repository, including the user/organisation name. Example: 'markbates/pop'"`
	Description string `json:"description" description:"The repository description."`
	IsPrivate bool `json:"private" description:"Whether the repository is private."`
	IsFork bool `json:"fork" description:"Whether the repository is a fork."`
}

// Use"r is a GitHub user.
type User struct {
	Resource
	Username string `json:"login" description:"The user's public username'"`
	ID string `json:"id" description:"GitHub's ID for the user"`
	AvatarURL string `json:"avatar_url" description:"The URL to the user's avatar image.'"`
}