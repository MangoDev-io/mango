package data

import (
	"github.com/go-pg/pg"
	"github.com/haardikk21/algorand-asset-manager/api/cmd/api/config"
)

type DatabaseService struct {
	*pg.DB
	config config.DatabaseConfig
}

func NewDatabaseService(dbConfig config.DatabaseConfig) *DatabaseService {
	pgDB := pg.Connect(&pg.Options{
		Addr:     dbConfig.PostgreSQLHost,
		User:     dbConfig.PostgreSQLUsername,
		Password: dbConfig.PostgreSQLPassword,
		Database: dbConfig.PostgreSQLDatabase,
	})

	service := &DatabaseService{pgDB, dbConfig}
	return service
}
