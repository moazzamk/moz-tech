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

type Storage struct {
	esMutex      sync.Mutex
	skills       map[string]bool
	searchClient *elastic.Client
}

func NewStorage(client *elastic.Client) *Storage {
	ret := new(Storage)
	ret.skills =  make(map[string]bool)
	ret.searchClient = client

	return ret
}

func (r *Storage) HasSkill(skill string) bool {
	searchClient := r.searchClient

	esMutex.Lock()
	ret, found := r.skills[skill]
	esMutex.Unlock()

	if found == false {
		esMutex.Lock()
		_, err := searchClient.Get().
			Index(`jobs`).
			Type(`skills`).
			Id(r.getHash(skill)).
			Do()

		if err != nil {
			r.skills[skill] = false
			//fmt.Println(skill, `skill not in cachey`)
		} else {
			r.skills[skill] = true
			//fmt.Println(skill, `skill in cachey`)
		}

		ret = r.skills[skill]

		esMutex.Unlock()
	}

	return ret
}

func (r *Storage) HasJobWithUrl(url string) bool {
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

func (r *Storage) getHash(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	md5hash := hex.EncodeToString(hasher.Sum(nil))

	return md5hash
}

func (r *Storage) AddSkill(skill string) (bool, error) {
	esMutex.Lock()
	_, err := r.searchClient.Index().
		Index(`jobs`).
		Type(`skills`).
		BodyString(`{"skill":"` + strings.Replace(skill, `"`, `\"`, -1) + `"}`).
		Id(r.getHash(skill)).
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

func (r *Storage) AddJob(job structures.JobDetail) {
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

func (r *Storage) GetJobs(search string, start int, end int) []structures.JobDetail {
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

func (r *Storage) GetSkills(start int, end int) []map[string]string {
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
