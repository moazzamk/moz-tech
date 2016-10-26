package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/moazzamk/moz-tech/structures"
	"gopkg.in/olivere/elastic.v3"
	"reflect"
	"strings"
	"sync"
)

var esMutex = &sync.Mutex{}


type Storage interface {
	HasSkill(string) bool
	AddSkill(string) (bool, error)
	HasJobWithUrl(string) bool
	AddJob(structures.JobDetail)
	GetJobs(string, int, int) []structures.JobDetail
	GetSkills(int, int) []map[string]string
}

type ElasticSearchStorage struct {
	esMutex      sync.Mutex
	skills       map[string]bool
	searchClient *elastic.Client
}

func NewStorage(client *elastic.Client) *ElasticSearchStorage {
	ret := new(ElasticSearchStorage)
	ret.skills =  make(map[string]bool)
	ret.searchClient = client

	return ret
}

func (r *ElasticSearchStorage) HasSkill(skill string) bool {
	searchClient := r.searchClient
	id := r.getHash(skill)

	esMutex.Lock()
	ret, found := r.skills[skill]
	esMutex.Unlock()

	fmt.Println(skill + ` slice `, found, " " + id)

	if found == false {
		esMutex.Lock()
		_, err := searchClient.Get().
								Index(`jobs`).
								Type(`skills`).
								Id(id).
								Do()

		if err != nil {
			r.skills[skill] = false
			fmt.Println(skill, `skill not in cachey `, err )
		} else {
			r.skills[skill] = true
			fmt.Println(skill, `skill in cachey`)
		}

		ret = r.skills[skill]

		esMutex.Unlock()
	}

	fmt.Println(`YOYOMA`, skill, ret)

	return ret
}

func (r *ElasticSearchStorage) AddSkill(skill string) (bool, error) {
	id := r.getHash(strings.ToLower(skill))
	esMutex.Lock()
	_, err := r.searchClient.Index().
							Index(`jobs`).
							Type(`skills`).
							BodyString(`{"skill":"` + strings.Replace(skill, `"`, `\"`, -1) + `"}`).
							Id(id).
							Refresh(true).
							Do()
	esMutex.Unlock()

	fmt.Println(`added`, skill)

	if err != nil {
		return false, err
	}

	esMutex.Lock()
	r.skills[skill] = true
	esMutex.Unlock()

	return true, nil
}

func (r *ElasticSearchStorage) HasJobWithUrl(url string) bool {
	md5hash := r.getHash(url)
	esMutex.Lock()
	rs, err := r.searchClient.Get().
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

func (r *ElasticSearchStorage) getHash(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	md5hash := hex.EncodeToString(hasher.Sum(nil))

	return md5hash
}

func (r *ElasticSearchStorage) AddJob(job structures.JobDetail) {
	md5hash := r.getHash(job.Link)
	esMutex.Lock()
	_, err := r.searchClient.
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

func (r *ElasticSearchStorage) GetJobs(search string, start int, end int) []structures.JobDetail {
	var ret []structures.JobDetail
	var tmp structures.JobDetail

	esMutex.Lock()
	query := elastic.NewTermQuery(`skill`, search)
	searchResult, _ := r.searchClient.Search().
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

func (r *ElasticSearchStorage) GetSkills(start int, end int) []map[string]string {
	ret := []map[string]string{}

	esMutex.Lock()
	searchResult, _ := r.searchClient.Search().
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

	return ret

}
