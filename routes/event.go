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

	// 3. Başarılıysa JSON olarak döndür
	context.JSON(http.StatusOK, event)
}

func createEvent(context *gin.Context) {
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
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create events. Try again later."})
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

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	userId := context.GetInt64("userid")
	if evnt.UserID != userId {
		context.JSON(http.StatusForbidden, gin.H{"message": "Not authorized to update the event"})
		return
	}

	updatedEvent.ID = id
	err = updatedEvent.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update event"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Event updated successfuly"})
}

func deleteEvent(context *gin.Context) {
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

	userId := context.GetInt64("userid")

	if evnt.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized tu update the event"})
		return
	}

	// Kullanıcı event'in sahibiyse, silmesine izin ver
	evnt.ID = id

	err = evnt.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete the event"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event successfully deleted."})
}
