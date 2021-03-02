package mongoInstances

import (
	"GoGameServer/core"
	"GoGameServer/core/libs/mongo"
)

func Global() *mongo.Client {
	return core.Service.GetMongoClient("global")
}
