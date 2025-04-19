package main

import (
	"example/db"
	"example/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default() // perde arkasında bir HTTP sunucusu yapılandırıyor

	routes.RegisterRoutes(server)
	server.Run(":8080") // Sunucuyu başlatır, gelen istekleri dinlemeye başlar, localhost:8080
}
