package main

/*
 * Entry point for running the penguin service.
 * This package only knows how to process command line args, 
 * and call the main service to start. 
 */
import (
    "flag"
    "github.com/simonajones/penguin-go"
)

var dbUrl string
var port string
var swaggerHost string
var swaggerFilePath string

func init() {
    flag.StringVar(&dbUrl, "dbUrl", "localhost:27017/penguin", "The host:port/db of the Mongo database to connect to.")
    flag.StringVar(&port, "port", "9091", "The port number that the Rest API will listen on.")
    flag.StringVar(&swaggerHost, "swaggerHost", "localhost", "The hostname that swagger will use to server the Swagger Web UI.")
    flag.StringVar(&swaggerFilePath, "swaggerFilePath", "/home/simon/Downloads/swagger-ui-1.1.7", 
        "The local path the unzipped swagger files.")
}

func main() {
    flag.Parse()
    penguin.StartService(penguin.Config{dbUrl, port, swaggerHost, swaggerFilePath})
}

