package moz_tech

import (
	"github.com/jackc/pgx"
	"github.com/bgentry/que-go"
	"fmt"
	"os"
)

const (
	ScanSkillsJob = "IndexRequests"
	ScanJobsJob = "ScanJobs"
)

type ScanSkillsRequest struct {
	URL string `json:"url"`
}

type ScanJobsRequest struct {
	Keywords string
	State string
	City string
}

func SetupDb(dbUrl string) (*pgx.ConnPool, *que.Client, error) {

	pgxpool, err := GetPgxPool(dbUrl)
	if err != nil {
		return nil, nil, err
	}

	qc := que.NewClient(pgxpool)

	return pgxpool, qc, nil
}

// GetPgxPool based on the provided database URL
func GetPgxPool(dbURL string) (*pgx.ConnPool, error) {
	pgxcfg, err := pgx.ParseURI(dbURL)
	if err != nil {
		return nil, err
	}

	pgxpool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:   pgxcfg,
		AfterConnect: que.PrepareStatements,
	})

	if err != nil {
		return nil, err
	}

	return pgxpool, nil
}

func GetConfigVars() (string, string) {
	var esUrl string
	var pgUrl string
	//esUrlKey := `SEARCHBOX_SSL_URL`
	//pgUrlKey := `HEROKU_POSTGRESQL_AQUA_URL`

	esUrlKey := `ELASTICSEARCH_URL`
	pgUrlKey := `DATABASE_URL`

	fmt.Println(`Starting worker`)

	if os.Getenv(esUrlKey) == `` {
		config := NewAppConfig(`config/config.txt`)
		esUrl, _ = config.Get(`es_url`)
		pgUrl, _ = config.Get(`psql_url`)

	} else {
		pgUrl = os.Getenv(pgUrlKey)
		esUrl = os.Getenv(esUrlKey)
	}

	return esUrl, pgUrl
}
