package ghcli

type RepositoryInput struct {
	Name        string
	Description string
	Private     bool
}

type Repository struct {
	ID          int64
	Name        string
	Description string
	Private     bool
}
