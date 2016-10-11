package service

import (
	"gopkg.in/olivere/elastic.v3"
	"strings"
	"sync"
	"github.com/moazzamk/moz-tech/structures"
)

var esMutex sync.Mutex
var hasSkill = make(map[string]bool)

func SearchHasSkill(client **elastic.Client, skill string) bool {
	searchClient := *client

	ret, err := hasSkill[skill]
	if err == false {
		esMutex.Lock()
		searchQuery := elastic.NewTermQuery(`skill`, skill)
		searchResult, err := searchClient.Search().
			Index(`jobs`).
			Type(`skills`).
			Query(searchQuery).
			Pretty(true).
			Do()

		esMutex.Unlock()

		if err != nil {
			panic(err)
		}

		if searchResult.Hits.TotalHits > 0 {
			hasSkill[skill] = true
		} else {
			hasSkill[skill] = false
		}

		ret = hasSkill[skill]
	}

	return ret
}

func SearchAddSkill(client **elastic.Client, skill string) (bool, error) {
	//return true, nil
	searchClient := *client

	esMutex.Lock()
	_, err := searchClient.Index().
		Index(`jobs`).
		Type(`skills`).
		BodyString(`{"skill":"` + strings.Replace(skill, `"`, `\"`, -1) + `"}`).
		Refresh(true).
		Do()
	esMutex.Unlock()

	if err != nil {
		return false, err
	}

	return true, nil
}

func SearchAddJob(client **elastic.Client, job structures.JobDetail) {
	searchClient := *client

	esMutex.Lock()
	_, err := searchClient.
						Index().
						Index("jobs").
						Type("job").
						Id(job.Link).
						BodyJson(job).
						Refresh(true).
						Do()
	esMutex.Unlock()

	if err != nil {
		panic(err)
	}
}