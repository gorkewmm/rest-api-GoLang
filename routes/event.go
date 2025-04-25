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
	// 1. URL'den gelen `id` parametresini string olarak al
	idStr := context.Param("id")

	// 2. String olan `id` parametresini int64'e çevir
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	// 3. Bu ID'ye sahip bir event var mı kontrol et (veritabanından çekmeye çalış)
	_, err = models.GetById(id)
	if err != nil {
		// Eğer kayıt yoksa veya hata oluşursa 500 hatası döndür
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the event"})
		return
	}

	// 4. İstek gövdesinden (request body) gelen JSON veriyi `updatedEvent` adlı struct'a aktar
	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		// JSON veri eksik veya hatalıysa 400 Bad Request döndür
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}

	// 5. `updatedEvent` struct'ının ID'sini güncelle (önceki adımda gelen path parametresinden alınmıştı)
	updatedEvent.ID = id

	// 6. Veritabanındaki kaydı güncelle
	err = updatedEvent.Update()
	if err != nil {
		// Güncelleme sırasında hata olursa 500 Internal Server Error döndür
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update event"})
		return
	}

	// 7. Her şey başarılıysa 200 OK ile başarı mesajı dön
	context.JSON(http.StatusOK, gin.H{"message": "Event updated successfuly"})
}

func deleteEvent(context *gin.Context) {
	idStr := context.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	_, err = models.GetById(id)
	if err != nil {
		// Eğer kayıt yoksa veya hata oluşursa 500 hatası döndür
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the event"})
		return
	}

	var event models.Event
	event.ID = id

	err = event.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete the event"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Event successfully deleted."})

}
