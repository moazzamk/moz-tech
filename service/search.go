package service

import (
	"gopkg.in/olivere/elastic.v3"
	"strings"
	"sync"
)

var esMutex sync.Mutex

func SearchHasSkill(client **elastic.Client, skill string) bool {
	searchClient := *client

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

	//fmt.Println(searchResult.Hits.TotalHits , " RESULTS FOUND FOR " + skill)
	if searchResult.Hits.TotalHits > 0 {
		return true
	}

	return false
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