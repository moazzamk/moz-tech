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
	"github.com/moazzamk/moz-tech/service"
	"github.com/moazzamk/moz-tech/action"
	"gopkg.in/olivere/elastic.v3"
	"github.com/moazzamk/moz-tech/structures"
)


var (
	esClient *elastic.Client
	config *structures.Dictionary
	qc      *que.Client
	pgxpool *pgx.ConnPool
)

// indexURLJob would do whatever indexing is necessary in the background
func scanSkills(j *que.Job) error {
	var ir moz_tech.ScanSkillsRequest
	if err := json.Unmarshal(j.Args, &ir); err != nil {
		return err
	}

	storage := service.NewStorage(esClient)
	skillParser := service.NewSkillParser(storage)
	crawlAction := action.NewCrawlTagsAction(&skillParser, config, storage)

	crawlAction.Run()

	log.Println("IndexRequest", "Processing Scan skills!")

	return nil
}

func scanJobs(j *que.Job) error {
	storage := service.NewStorage(esClient)
	skillParser := service.NewSkillParser(storage)
	crawlAction := action.NewCrawlTagsAction(&skillParser, config, storage)

	crawlAction.Run()


	log.Println("IndexRequest", "Processing IndexRequest! (not really)")

	return nil

}

/**
 * Entry point for background workers
 */
func main() {
	var err error

	esUrl, pgUrl := moz_tech.GetConfigVars()
	esClient, err = elastic.NewClient(
		elastic.SetURL(esUrl),
		elastic.SetMaxRetries(10),
		elastic.SetSniff(false),
		elastic.SetScheme(`https`))

	if err != nil {
		panic(err)
	}

	pgxpool, qc, err :=	 moz_tech.SetupDb(pgUrl)
	if err != nil {
		log.Println("PSQL connection failed on workers")
	}
	defer pgxpool.Close()

	wm := que.WorkMap{
		moz_tech.ScanSkillsJob: scanSkills,
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
