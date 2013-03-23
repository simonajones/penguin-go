package penguin

import (
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
)

func NewQueueService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path("/api").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/queues").To(queueList).
		// docs	
		Doc("get all queues"))

	ws.Route(ws.GET("/queue/{id}").To(queueGet).
		// docs	
		Doc("get a queue by id").
		Param(ws.PathParameter("id", "identifier of the queue")))
	
	ws.Route(ws.POST("/queues").To(queueCreate).
		// docs	
		Doc("create a queue"))
	
    ws.Route(ws.PUT("/queue/{id}").To(queueUpdate).
		// docs	
		Doc("Update a queue").
		Param(ws.PathParameter("name", "name of the queue")))
	
    ws.Route(ws.POST("/queue/{queue-id}/stories").To(storyCreate).
		// docs	
		Doc("create a story"))
		
    ws.Route(ws.GET("/queue/{queue-id}/story/{id}").To(storyGet).
		// docs	
		Doc("get a story"))

    ws.Route(ws.POST("/queue/{queue-id}/story/{id}").To(storyDelete).
		// docs	
		Doc("delete a story"))

	return ws
}

func queueList(request *restful.Request, response *restful.Response) {

    session, collection := getDB()
    defer session.Close()

    result := []Queue{}
    collection.Find(bson.M{}).All(&result)
    log.Printf("Found: ", result)
	response.WriteEntity(&result)
}

func queueGet(request *restful.Request, response *restful.Response) {
    idParam := request.PathParameter("id")
	id := bson.ObjectIdHex(idParam)

    session, collection := getDB()
    defer session.Close()

    result := Queue{}
    collection.Find(bson.M{"_id": id}).One(&result)
    log.Printf("Found: ", result)

	if len(result.Id) == 0 {
		response.WriteError(http.StatusNotFound, nil)
	} else {
		response.WriteEntity(&result)
	}
}

func queueCreate(request *restful.Request, response *restful.Response) {
	queue := new(Queue)
	err := request.ReadEntity(&queue)
    queue.Id = bson.NewObjectId()

    session, collection := getDB()
    defer session.Close()

    err = collection.Insert(&queue)
    log.Printf("Inserted: ", queue, " err ", err)

	if err == nil {
		response.WriteHeader(http.StatusCreated)
    	response.WriteEntity(&queue)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
    }
}

func queueUpdate(request *restful.Request, response *restful.Response) {
    idParam := request.PathParameter("id")
	id := bson.ObjectIdHex(idParam)

	queue := new(Queue)
	err := request.ReadEntity(&queue)

    session, collection := getDB()
    defer session.Close()

    err = collection.UpdateId(id, &queue)
    log.Printf("Updated: ", queue, " err ", err)

	if err == nil {
    	response.WriteEntity(&queue)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
    }
}


func storyCreate(request *restful.Request, response *restful.Response) {
    idParam := request.PathParameter("queue-id")
	queueOid := bson.ObjectIdHex(idParam)

	story := new(Story)
	err := request.ReadEntity(&story)
    story.Id = bson.NewObjectId()

    session, collection := getDB()
    defer session.Close()

    err = collection.UpdateId(queueOid, bson.M{"$push": bson.M{"stories": story}})
    log.Printf("Created story: ", story, " err ", err)

	if err == nil {
		response.WriteHeader(http.StatusCreated)
    	response.WriteEntity(&story)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
    }
}

func storyGet(request *restful.Request, response *restful.Response) {
    queueId := request.PathParameter("queue-id")
	queueOid := bson.ObjectIdHex(queueId)

    id := request.PathParameter("id")
	storyOid := bson.ObjectIdHex(id)

    session, collection := getDB()
    defer session.Close()

    story := Story{}
    err := collection.Find(bson.M{"_id": queueOid, "stories._id": storyOid}).One(&story)
    log.Printf("Found story: ", story, " err ", err)

	if err == nil {
		response.WriteHeader(http.StatusCreated)
    	response.WriteEntity(&story)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
    }
}

func storyDelete(request *restful.Request, response *restful.Response) {
    queueId := request.PathParameter("queue-id")
	queueOid := bson.ObjectIdHex(queueId)

    id := request.PathParameter("id")
	story := Story{Id: bson.ObjectIdHex(id)}

    session, collection := getDB()
    defer session.Close()

    err := collection.UpdateId(queueOid, bson.M{"$pull": bson.M{"stories": story}})
    log.Printf("Created story: ", story, " err ", err)

	if err == nil {
		response.WriteHeader(http.StatusCreated)
    	response.WriteEntity(&story)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
    }
}

func getDB() (session *mgo.Session, c *mgo.Collection) {
    session, err := mgo.Dial(config.DbUrl)
    if err != nil {
            panic(err)
    }
    return session, session.DB("penguin").C("queues")
}

var config Config

func StartService(configuration Config) {
    config = configuration
	restful.Add(NewQueueService())
	
	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// Open http://localhost:8080/apidocs and enter http://localhost:8080/apidocs.json in the api input field.
	swaggerConfig := restful.SwaggerConfig{ 
		WebServicesUrl: "http://"+config.SwaggerHost+ ":" + config.Port,
		ApiPath: "/apidocs.json",
		SwaggerPath: "/apidocs/",
		SwaggerFilePath: config.SwaggerFilePath }	
	restful.InstallSwaggerService(swaggerConfig)
	
	log.Printf("start listening on localhost:"+config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}


