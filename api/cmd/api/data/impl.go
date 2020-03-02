package data

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/haardikk21/algorand-asset-manager/api/cmd/api/models"
)

// WaitUntilReady returns when the database service has connected
func (s *DatabaseService) WaitUntilReady() *DatabaseService {
	pingCounter := 0
	err := s.Ping()
	for err != nil {
		pingCounter++
		if pingCounter > s.config.MaxConnectionRetries {
			panic("Could not connect to DB")
		}
		time.Sleep(time.Second * time.Duration(pingCounter))
		err = s.Ping()
	}

	return s
}

// Instantiate creates DB tables if they don't exist
func (s *DatabaseService) Instantiate() error {
	for _, model := range []interface{}{
		(*models.OwnedAssets)(nil),
	} {
		err := s.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// Ping pings the database
func (s *DatabaseService) Ping() error {
	i := 0

	_, err := s.QueryOne(pg.Scan(&i), "SELECT 1")
	return err
}
