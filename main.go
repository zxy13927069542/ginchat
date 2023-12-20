package main

import (
	"ginchat/config"
	"ginchat/models"
	"ginchat/redisc"
	"ginchat/router"
	"ginchat/utils"
)

func main() {
	utils.Init()
	c := config.Init("./etc")
	_ = models.Init(c)
	redisc.Init(c)

	r := router.Router()
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run(":8080")
}
