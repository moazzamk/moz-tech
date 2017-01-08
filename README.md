# job-tech
=======

A project to crawl job listings on the web and calculate which technologies are generally asked for together.
For example: PHP jobs almost always require MySQL knowledge.


## For cli:

go run cmd/cli/main.go

## For web:

go run cmd/web/main.go

For any errors of missing dependencies, you can do:

go get &lt;dependency>

##Generate mocks:
mockgen --source=somefile.go --destination=mock/somefile.go

## To run tests:
ginkgo
