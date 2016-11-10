package main

import (
	"net/http"
	"fmt"
	"github.com/urfave/negroni"
	"log"
	"html/template"

	"github.com/moazzamk/moz-tech/service"
	"github.com/moazzamk/moz-tech/app"
	"gopkg.in/olivere/elastic.v3"
	"github.com/bgentry/que-go"
	"github.com/jackc/pgx"
	"os"
	"encoding/json"
	"strconv"
	"io/ioutil"
)

var (
	templatePath = `/Users/moz/gosites/src/github.com/moazzamk/moz-tech/cmd/web/views`
	qc      *que.Client
	pgxpool *pgx.ConnPool
	esClient *elastic.Client
)

// queueIndexRequest into the que as an encoded JSON object
func queueScanSkillsRequest(ir moz_tech.ScanSkillsRequest) error {
	j := que.Job{
		Type: moz_tech.ScanSkillsJob,
		Args: []byte("{}"),
	}

	return qc.Enqueue(&j)
}

func queueScanJobsRequest(ir moz_tech.ScanJobsRequest) error {
	j := que.Job{
		Type: moz_tech.ScanJobsJob,
	}

	return qc.Enqueue(&j)
}


func main() {
	var err error

	esUrl, pgUrl := moz_tech.GetConfigVars()
	templatePath = os.Getenv(`PWD`) + `/cmd/web/views`
	pgxpool, qc, err =	 moz_tech.SetupDb(pgUrl)
	if err != nil {
		fmt.Println(err, "ERRRRRRRRRR")
	}
	defer pgxpool.Close()

	fmt.Println(esUrl)

	// Initialize ES
	// Keep sniffing to false. It causes ES library to fail
	client, err := elastic.NewClient(
		elastic.SetURL(esUrl),
		elastic.SetMaxRetries(10),
		elastic.SetSniff(false),
		elastic.SetScheme(`https`))

	if err != nil {
		panic(err)
	}

	// Make sure ES works
	_, _, err = client.Ping(esUrl).Do()
	if err != nil {
		panic(err)
	}

	esClient = client

	fmt.Println("Webserver Initialized")
	fmt.Println("Template path:" + templatePath)



	mux := http.NewServeMux()

		mux.Handle("/static", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))

	// Routes
	mux.HandleFunc(`/index/delete`, func (rs http.ResponseWriter, rq *http.Request) {
		esClient.DeleteIndex(`jobs`).Do()
		esClient.CreateIndex(`jobs`).Do()

		t := template.New(`index.html`)
		t, _ = t.ParseFiles(templatePath + `/admin/index.html`)
		err := t.Execute(rs, make(map[string]string))
		if err != nil {
			fmt.Println(err)
		}
	})

	mux.HandleFunc(`/jobs/delete`, func (rs http.ResponseWriter, rq *http.Request) {
		service := elastic.NewDeleteByQueryService(esClient)
		service.Index(`jobs`)
		service.Type(`job`)
		service.QueryString(`stack`)
		r1, err := service.Do()
		if err != nil {
			fmt.Println("hi", err)
		}
		fmt.Println(r1)

		t := template.New(`index.html`)
		t, _ = t.ParseFiles(templatePath + `/admin/index.html`)
		err = t.Execute(rs, make(map[string]string))
		if err != nil {
			fmt.Println(err)
		}
	})

	mux.HandleFunc(`/admin`, func (rs http.ResponseWriter, rq *http.Request) {

		t := template.New(`index.html`)
		t, _ = t.ParseFiles(templatePath + `/admin/index.html`)
		err := t.Execute(rs, make(map[string]string))
		if err != nil {
			fmt.Println(err)
		}

	})

	mux.HandleFunc(`/scan/skills`, func (rs http.ResponseWriter, rq *http.Request) {
		queueScanSkillsRequest(moz_tech.ScanSkillsRequest{})
	})

	mux.HandleFunc(`/scan/jobs`, func (rs http.ResponseWriter, rq *http.Request) {
		queueScanJobsRequest(moz_tech.ScanJobsRequest{})
	})


	// Search jobs
	mux.HandleFunc(`/search`, func (rs http.ResponseWriter, rq *http.Request) {
		requestData := rq.URL.Query()

		if query, ok := requestData[`q`]; ok {
			start := 0
			if val, ok := requestData[`start`]; ok {
				start, _ = strconv.Atoi(val[0])
			}
			jobs, totalCount := service.NewStorage(esClient).GetJobs(query[0], start, 10)

			rs1 := make(map[string]interface{})
			rs1[`total`] = totalCount
			rs1[`data`] = jobs
			rs1[`success`] = 1

			val, _ := json.Marshal(rs1)

			rs.Write(val)
			return
		}

		f, _ := ioutil.ReadFile(templatePath + `/jobs/search.html`)
		rs.Write(f)
	})

	/**

	 * Delete a skill by it's ID
	 */
	mux.HandleFunc(`/skills/del`, func (rs http.ResponseWriter, rq *http.Request) {
		_, err := client.Delete().Index(`jobs`).Type(`skills`).Id(rq.URL.Query()[`id`][0]).Do()
		if err != nil {
			fmt.Println(`ERRRRR`, err)
		}
		http.Redirect(rs, rq, `/skills`, 302)
	})

	/**
	 * List all skills
	 */
	mux.HandleFunc(`/skills`, func (rs http.ResponseWriter, rq *http.Request) {
		rs1 := service.NewStorage(esClient).GetSkills(0, 100)

		t := template.New(`list.html`)
		t, _ = t.ParseFiles(templatePath + `/skills/list.html`)
		err := t.Execute(rs, rs1)
		if err != nil {
			fmt.Println(err)
		}
	})

	/**
	 * Home page
	 */
	mux.HandleFunc(`/`, func (rs http.ResponseWriter, rq *http.Request) {
		t := template.New(`list.html`)
		t, _ = t.ParseFiles(templatePath + `/skills/list.html`)
		err := t.Execute(rs, nil)
		if err != nil {
			fmt.Println(err)
		}
	})


	n := negroni.Classic()
	n.UseHandler(mux)

	portString := os.Getenv(`PORT`)
	if portString == `` {
		portString = `7000`
	}
	s := &http.Server{
		Addr: ":" + portString,
		Handler: n,
		MaxHeaderBytes: 1 <<20,
	}

	log.Fatal(s.ListenAndServe())
}
