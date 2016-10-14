package service

import (
	"gopkg.in/olivere/elastic.v3"
	"strings"
	"sync"
	"github.com/moazzamk/moz-tech/structures"
	"crypto/md5"
	"encoding/hex"
	"fmt"
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

		//fmt.Println(searchResult.TotalHits(), `SEARCHY`)

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

	hasher := md5.New()
	hasher.Write([]byte(job.Link))
	md5hash := hex.EncodeToString(hasher.Sum(nil))

	fmt.Println(md5hash, "SSSS")

	esMutex.Lock()
	_, err := searchClient.
						Index().
						Index("jobs").
						Type("job").
						Id(md5hash).
						BodyJson(job).
						Refresh(true).
						Do()
	esMutex.Unlock()

	if err != nil {
		panic(err)
	}
}
