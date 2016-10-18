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

	fmt.Println("Initialized")


	mux := http.NewServeMux()
	mux.HandleFunc(`/skills`, func (rs http.ResponseWriter, rq *http.Request) {
		rs1 := service.SearchGetSkills(&client, 0, 100)

		fmt.Println(rs1)


		t := template.New(`list.html`)
		t, _ = t.ParseFiles(templatePath + `/skills/list.html`)
		err := t.Execute(rs, nil)
		if err != nil {
			fmt.Println(err)
		}
	})

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

