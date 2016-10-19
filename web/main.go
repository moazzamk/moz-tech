package main

import (
	"net/http"
	"fmt"
	"github.com/urfave/negroni"
	"log"
	"html/template"

	"github.com/moazzamk/moz-tech/app"
	"gopkg.in/olivere/elastic.v3"
	"github.com/moazzamk/moz-tech/service"
)

var templatePath = `/Users/mkhan/gosites/src/github.com/moazzamk/moz-tech/web/views`

func main() {

	config := moz_tech.NewAppConfig(`config/config.txt`)
	esUrl, _ := config.Get(`es_url`)

	fmt.Println(config)

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

	fmt.Println("Webserver Initialized")


	mux := http.NewServeMux()

	// Routes

	// Search jobs
	mux.HandleFunc(`/search`, func (rs http.ResponseWriter, rq *http.Request) {
		t := template.New(`search.html`)
		t, _ = t.ParseFiles(templatePath + `/jobs/search.html`)
		err := t.Execute(rs, make(map[string]string))
		if err != nil {
			fmt.Println(err)
		}
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
		rs1 := service.SearchGetSkills(&client, 0, 100)

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

	s := &http.Server{
		Addr: ":7000",
		Handler: n,
		MaxHeaderBytes: 1 <<20,
	}

	log.Fatal(s.ListenAndServe())
}
