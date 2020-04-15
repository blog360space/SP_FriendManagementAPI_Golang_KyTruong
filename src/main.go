package main

import (
	"friendMgmt/configs"
	"friendMgmt/data"
	"friendMgmt/docs"
)

func main() {
	docs.SwaggerInfo.Title = "Friend Management APIs"
	docs.SwaggerInfo.Description = ""
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8081"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}

	db, _ := data.InitDB()
	defer db.Close()

	configs.ConfigRoutes(db)
}
