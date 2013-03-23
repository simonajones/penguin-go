package main

import (
	"github.com/emicklei/go-restful"
	"log"
    "flag"
	"net/http"
    "github.com/simonajones/penguin-go"
)

var dbUrl string
var port string
var swaggerHost string

func init() {
    flag.StringVar(&dbUrl, "dbUrl", "localhost:27017/penguin", "The host:port/db of the Mongo database to connect to.")
    flag.StringVar(&port, "port", "9091", "The port number that the Rest API will listen on.")
    flag.StringVar(&swaggerHost, "swaggerHost", "localhost", "The hostname that swagger will use to server the Swagger Web UI.")
}

func main() {
    flag.Parse()
	restful.Add(penguin.NewQueueService(dbUrl))
	
	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// Open http://localhost:8080/apidocs and enter http://localhost:8080/apidocs.json in the api input field.
	config := restful.SwaggerConfig{ 
		WebServicesUrl: "http://"+swaggerHost+ ":" +port,
		ApiPath: "/apidocs.json",
		SwaggerPath: "/apidocs/",
		SwaggerFilePath: "/home/simon/Downloads/swagger-ui-1.1.7" }	
	restful.InstallSwaggerService(config)
	
	log.Printf("start listening on localhost:"+port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

