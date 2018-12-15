package mongoInstances

import (
	"core"
	"core/libs/mongo"
)

func Global() *mongo.Client {
	return core.Service.GetMongoClient("global")
}
