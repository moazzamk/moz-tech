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

type Storage struct {
	esMutex sync.Mutex
	skills map[string]bool
	searchClient *elastic.Client
}

func NewStorage(client **elastic.Client) Storage {
	return Storage{
		skills: make(map[string]bool),
		searchClient: *client,
	}
}

func (r *Storage) HasSkill(skill string) bool {
	searchClient := r.searchClient
	ret, err := r.skills[skill]
	if err == false {
		r.esMutex.Lock()
		searchQuery := elastic.NewTermQuery(`skill`, skill)
		searchResult, err := searchClient.Search().
											Index(`jobs`).
											Type(`skills`).
											Query(searchQuery).
											Pretty(true).
											Do()

		r.esMutex.Unlock()

		if err != nil {
			panic(err)
		}

		//fmt.Println(searchResult.TotalHits(), `SEARCHY`)

		if searchResult.Hits.TotalHits > 0 {
			r.skills[skill] = true
		} else {
			r.skills[skill] = false
		}

		ret = r.skills[skill]
	}

	return ret
}

func (r *Storage) HasJobWithUrl(url string) bool {

	hasher := md5.New()
	hasher.Write([]byte(url))
	md5hash := hex.EncodeToString(hasher.Sum(nil))


	r.esMutex.Lock()
	rs, err := r.searchClient.Get().
							Index(`jobs`).
							Type(`job`).
							Id(md5hash).
							Do()
	r.esMutex.Unlock()

	if err != nil {
		return false
	}

	return rs.Found
}

func (r *Storage) AddSkill(skill string) (bool, error) {
	hasher := md5.New()
	hasher.Write([]byte(skill))
	md5hash := hex.EncodeToString(hasher.Sum(nil))

	r.esMutex.Lock()
	_, err := r.searchClient.Index().
		Index(`jobs`).
		Type(`skills`).
		BodyString(`{"skill":"` + strings.Replace(skill, `"`, `\"`, -1) + `"}`).
		Id(md5hash).
		Refresh(true).
		Do()
	r.esMutex.Unlock()

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *Storage) AddJob(job structures.JobDetail) {
	hasher := md5.New()
	hasher.Write([]byte(job.Link))
	md5hash := hex.EncodeToString(hasher.Sum(nil))

	//fmt.Println(md5hash, "HASHY")

	r.esMutex.Lock()
	_, err := r.searchClient.
						Index().
						Index("jobs").
						Type("job").
						Id(md5hash).
						BodyJson(job).
						Refresh(true).
						Do()
	r.esMutex.Unlock()

	if err != nil {
		panic(err)
	}
}

func (r *Storage) GetJobs(search string, start int, end int) []structures.JobDetail {
	var ret []structures.JobDetail
	var tmp structures.JobDetail

	r.esMutex.Lock()
	query := elastic.NewTermQuery(`skill`, search)
	searchResult, _ := r.searchClient.Search().
								Index(`jobs`).
								Type(`skills`).
								Query(query).
								Pretty(true).
								Do()
	r.esMutex.Unlock()

	for _, item := range searchResult.Each(reflect.TypeOf(tmp)) {
		ret = append(ret, item.(structures.JobDetail))
	}

	return ret
}

func (r *Storage) GetSkills(start int, end int) []map[string]string {
	ret := []map[string]string{}

	r.esMutex.Lock()
	searchResult, _ := r.searchClient.Search().
									Index(`jobs`).
									Type(`skills`).
							//		Query(searchQuery).
									Pretty(true).
									Do()
	r.esMutex.Unlock()

	var tt map[string]string
	for _, item := range searchResult.Each(reflect.TypeOf(tt)) {
		fmt.Println(item, "DDDDD")
		break
	}

	return ret

}
