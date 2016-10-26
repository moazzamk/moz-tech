package main

import (

	"github.com/moazzamk/moz-tech/app"
	"github.com/moazzamk/moz-tech/service"
	"gopkg.in/olivere/elastic.v3"
	"github.com/moazzamk/moz-tech/action"
)

/*

- find jobs that pay the most
- find companies that pay the most
- find technologies that relate to each other

*/

func main() {
	config := moz_tech.NewAppConfig(`config/config.txt`)
	esUrl, _ := config.Get(`es_url`)

	client, err := elastic.NewClient(
		elastic.SetURL(esUrl),
		elastic.SetMaxRetries(10),
		elastic.SetSniff(false),
		elastic.SetScheme(`https`))

	if err != nil {
		panic(err)
	}

	//client.DeleteIndex(`jobs`).Do()
	//client.CreateIndex(`jobs`).Do()

	storage := service.NewStorage(client)
	skillParser := service.NewSkillParser(storage)
	//salaryParser := service.SalaryParser{}
	//dateParser := service.DateParser{}
	/*crawlAction := action.NewCrawlJobsAction(
		&salaryParser,
		&skillParser,
		&dateParser,
		config,
		&storage)
	*/

	crawlAction := action.NewCrawlTagsAction(&skillParser, config, storage)

	crawlAction.Run()
}
