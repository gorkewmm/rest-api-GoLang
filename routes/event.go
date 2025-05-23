package routes

import (
	"example/models"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch events. Try again later."})
		return
	}
	context.JSON(http.StatusOK, events)
}

func getEventById(context *gin.Context) {
	// 1. URL'den gelen id'yi al
	idStr := context.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	// 2. Veritabanından ID'ye göre eventi çek
	event, err := models.GetById(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event."})
		return
	}

	// Burada kayıtlı kişi sayısını alıyoruz
	registrationCount, err := event.GetRegistrationCount()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch registration count"})
		return
	}

	// Event bilgisi + kayıtlı kişi sayısını dönüyoruz
	context.JSON(http.StatusOK, gin.H{
		"event":            event,
		"registered_users": registrationCount,
	})
}

func createEvent(context *gin.Context) {
	role := context.GetString("role")
	if role != "admin" {
		context.JSON(http.StatusForbidden, gin.H{"message": "Only admins can create events"})
		return
	}

	userId := context.GetInt64("userid")

	var event models.Event
	event.UserID = userId

	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request"})
		return
	}

	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Event created!", "event": event})
}

func updateEvent(context *gin.Context) {
	idStr := context.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	evnt, err := models.GetById(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the event"})
		return
	}

	var newEvent models.Event
	err = context.ShouldBindJSON(&newEvent) //newEvent kullanıcının gönderdiği yeni verileri taşıyor
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}

	userid := context.GetInt64("userid")
	role := context.GetString("role")

	if userid != evnt.UserID && role != "admin" {
		context.JSON(http.StatusForbidden, gin.H{"message": "You are not authorized to delete this event"})
		return
	}

	err = newEvent.Update()
	newEvent.ID = evnt.ID
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update the event"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Event successfully updateed"})
}

func deleteEvent(context *gin.Context) {
	idStr := context.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}
	event, err := models.GetById(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id.", "error": err.Error()})
		return
	}
	userId := context.GetInt64("userid")
	role := context.GetString("role")
	if role != "admin" && userId != event.UserID {
		context.JSON(http.StatusForbidden, gin.H{"message": "You are not authorized to delete this event"})
		return
	}

	var events models.Event
	events.ID = id
	err = events.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete the event"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Event successfully deleted"})
}

func getUserEvents(context *gin.Context) {
	idStr := context.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse user id."})
		return
	}
	events, err := models.GetEventsByUserId(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch events."})
		return
	}
	context.JSON(http.StatusOK, events)
}
