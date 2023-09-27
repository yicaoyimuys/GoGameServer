package mongoInstances

import (
	"github.com/yicaoyimuys/GoGameServer/core"
	"github.com/yicaoyimuys/GoGameServer/core/libs/mongo"
)

func Global() *mongo.Client {
	return core.Service.GetMongoClient("global")
}
