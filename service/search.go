package service

import (
	"gopkg.in/olivere/elastic.v3"
	"strings"
	"sync"
	"github.com/moazzamk/moz-tech/structures"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
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

func SearchHasJobWithUrl(client **elastic.Client, url string) bool {
	searchClient := *client

	hasher := md5.New()
	hasher.Write([]byte(url))
	md5hash := hex.EncodeToString(hasher.Sum(nil))

	esMutex.Lock()
	rs, err := searchClient.Get().
							Index(`jobs`).
							Type(`job`).
							Id(md5hash).
							Do()
	esMutex.Unlock()

	if err != nil {
		return false
	}

	return rs.Found
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

func SearchGetJobs(client **elastic.Client, search string, start int, end int) []structures.JobDetail {
	var ret []structures.JobDetail
	var tmp structures.JobDetail

	searchClient := *client

	esMutex.Lock()
	query := elastic.NewTermQuery(`skill`, search)
	searchResult, _ := searchClient.Search().
								Index(`jobs`).
								Type(`skills`).
								Query(query).
								Pretty(true).
								Do()
	esMutex.Unlock()

	for _, item := range searchResult.Each(reflect.TypeOf(tmp)) {
		ret = append(ret, item.(structures.JobDetail))
	}

	return ret
}

func SearchGetSkills(client **elastic.Client, start int, end int) []map[string]string {
	ret := []map[string]string{}
	searchClient := *client

	esMutex.Lock()
	searchResult, _ := searchClient.Search().
		Index(`jobs`).
		Type(`skills`).
//		Query(searchQuery).
		Pretty(true).
		Do()
	esMutex.Unlock()

	var tt map[string]string
	for _, item := range searchResult.Each(reflect.TypeOf(tt)) {
		fmt.Println(item, "DDDDD")
		break
	}

/*
	fmt.Println(searchResult.TotalHits(), "ITS")
	t := make(map[string]interface{})
	for _, hit := range searchResult.Hits.Hits {
		_ = json.Unmarshal(*hit.Source, tt)
		fmt.Println(tt, string(*hit.Source))
		ret = append(ret, t[`skill`].(string))

	}


*/

	fmt.Println(ret)
	return ret

}
