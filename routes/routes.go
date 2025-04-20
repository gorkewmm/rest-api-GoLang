package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/events", getEvents) //GET,POST,PUT,PATCH,DELETE
	server.GET("/events/:id", getEventById)
	server.POST("/events", createEvent)      // veri ekleme
	server.PUT("/events/:id", updateEvent)   //guncelleme
	server.DELETE("events/:id", deleteEvent) // delete ,silme
	server.POST("/signup", signup)           //the goal here is the create new users
}
