package service

import (
	"gopkg.in/olivere/elastic.v3"
	"strings"
)

func SearchHasSkill(client **elastic.Client, skill string) bool {

	searchClient := *client
	searchQuery := elastic.NewTermQuery(`skill`, skill)
	searchResult, err := searchClient.Search().
		Index(`jobs`).
		Type(`skills`).
		Query(searchQuery).
		Pretty(true).
		Do()

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
	_, err := searchClient.Index().
		Index(`jobs`).
		Type(`skills`).
		BodyString(`{"skill":"` + strings.Replace(skill, `"`, `\"`, -1) + `"}`).
		Refresh(true).
		Do()

	if err != nil {
		return false, err
	}

	return true, nil
}
