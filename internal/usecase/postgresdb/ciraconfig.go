package postgresdb

import "github.com/open-amt-cloud-toolkit/console/pkg/postgres"

// DomainRepo -.
type CIRARepo struct {
	*postgres.DB
}

// New -.
func NewCIRARepo(pg *postgres.DB) *CIRARepo {
	return &CIRARepo{pg}
}
