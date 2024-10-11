package mocks

import "github.com/open-amt-cloud-toolkit/console/pkg/db"

type MockSQLDB struct { //nolint:revive // Ignore stutter since the S is part of "S"QL not mock"s"
	*db.SQL
}

func NewMockSQLDB() *db.SQL {
	return &db.SQL{}
}
