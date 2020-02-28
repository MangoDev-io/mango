package data

import (
	"time"

	"github.com/go-pg/pg"
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

// Ping pings the database
func (s *DatabaseService) Ping() error {
	i := 0

	_, err := s.QueryOne(pg.Scan(&i), "SELECT 1")
	return err
}
