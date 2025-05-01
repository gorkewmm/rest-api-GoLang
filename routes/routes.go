package routes

import (
	"example/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/events", getEvents)
	server.GET("/events/:id", getEventById)
	server.GET("/admin/users", middlewares.Authenticate, getUsers)
	server.GET("/admin/users/:id", middlewares.Authenticate, getUserById)

	server.POST("/admin/events", middlewares.Authenticate, createEvent) // veri ekleme
	server.PUT("/events/:id", middlewares.Authenticate, updateEvent)    //guncelleme
	server.DELETE("/events/:id", middlewares.Authenticate, deleteEvent) // delete ,silme

	server.POST("/signup", signup) //the goal here is the create new users
	server.POST("/login", login)

	server.POST("/events/:id/register", middlewares.Authenticate, registerForEvent)
	server.DELETE("/events/:id/register", middlewares.Authenticate, cancelRegistration)

	server.PUT("/users/:id", middlewares.Authenticate, updateUser) //kullanıcı kendi bilgilerini günceller
	server.DELETE("/users/:id", middlewares.Authenticate, deleteUser)
	server.PUT("/users/:id/password", middlewares.Authenticate, changePassword)

}
