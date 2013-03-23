package penguin

import (
    "labix.org/v2/mgo/bson"
)

type Queue struct { 
        Id           bson.ObjectId  `json:"_id" bson:"_id"`
        Name         string         `json:"name"`
        Stories      []Story        `json:"stories"`
} 

type Story struct { 
        Id           bson.ObjectId  `json:"_id" bson:"_id"`
        Author       string         `json:"author"`       
        Merged       bool           `json:"merged"`          
        Reference    string         `json:"reference"`
        Title        string         `json:"title"`
} 

