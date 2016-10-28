package main

import (
	"github.com/moazzamk/moz-tech/app"
	"log"
	"github.com/bgentry/que-go"
	"os"
	"os/signal"
	"syscall"
	"encoding/json"
	"github.com/jackc/pgx"
)


var (
	qc      *que.Client
	pgxpool *pgx.ConnPool
)

// indexURLJob would do whatever indexing is necessary in the background
func indexURLJob(j *que.Job) error {
	var ir qe.IndexRequest
	if err := json.Unmarshal(j.Args, &ir); err != nil {
		return errors.Wrap(err, "Unable to unmarshal job arguments into IndexRequest: "+string(j.Args))
	}

	log.Println("IndexRequest", "Processing IndexRequest! (not really)")

	return nil
}

/**
 * Entry point for background workers
 */
func main() {
	config := moz_tech.NewAppConfig(`config/config.txt`)
	pgxpool, qc, err := moz_tech.SetupDb(config["pgsql_url"])
	if err != nil {
		log.Println("PSQL connection failed on workers")
	}
	defer pgxpool.Close()

	wm := que.WorkMap{
		moz_tech.IndexRequestJob: indexURLJob,
	}

	workers := que.NewWorkerPool(qc, wm, 2)

	// Catch signal so we can shutdown gracefully
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go workers.Start()

	sig := <-sigCh
	log.Println("signal", sig, "Signal received. Shutting down.")

	workers.Shutdown()
}
