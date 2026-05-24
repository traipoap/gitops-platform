package main

import (
	"exporter/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routers.SetupRoutes(r)
	r.Run()
}
