package models

import (
	"example/db"
	"fmt"
	"time"
)

type Event struct {
	ID              int64
	Name            string    `binding :"required"`
	Description     string    `binding :"required"`
	Location        string    `binding :"required"`
	DateTime        time.Time `binding :"required"`
	UserID          int64
	Price           float64 `binding :"required"`
	RegisteredUsers int64   `json:"registered_users"`
}

var events = []Event{} //Event struct'larının tutulduğu bir slice (dinamik dizi).
// = ile events değişkenine boş bir slice atanıyor. (Henüz içinde hiç Event yok.)

func (e *Event) Save() error { //database kaydetme işlemi
	// 0 kayıtla başlasın
	e.RegisteredUsers = 0

	query := `
	INSERT INTO events(name,description,location,datetime,user_id,price,registered_users)
	VALUES (?,?,?,?,?,?,?)`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close() //stmt nesnesi başarıyla oluşturulmuşsa, kullanılmasından bağımsız olarak kapanacağı garanti edilir.

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID, e.Price, e.RegisteredUsers)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	e.ID = id
	return err
}

func GetAllEvents() ([]Event, error) {
	query := `
	SELECT * FROM events`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event // birden fazla Event structu tutmak için başta boş bir slice
	for rows.Next() {
		var event Event // her döngüde yeni bir event structu oluşturuluyor
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID, &event.Price, &event.RegisteredUsers)
		//Başlangıçta boş olan event structını o satırda buldugumuz tüm verilerle doldurduk
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func GetById(id int64) (*Event, error) {
	var event Event
	query := `
	SELECT * FROM events where id =?
	`
	row := db.DB.QueryRow(query, id) //SQL sorgusunu çalıştırır ve tek satır sonuç bekler.

	err := row.Scan(
		&event.ID,
		&event.Name,
		&event.Description,
		&event.Location,
		&event.DateTime,
		&event.UserID,
		&event.Price,
		&event.RegisteredUsers,
	)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (event Event) Update() error {
	query := `
	UPDATE events
	SET name =?, description =?, location =?,datetime =?,price =?
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close() //stmt nesnesi başarıyla oluşturulmuşsa, kullanılmasından bağımsız olarak kapanacağı garanti edilir.

	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.Price, event.ID)
	return err
}

func (e Event) Delete() error {
	fmt.Println("Delete fonksiyonu çağrıldı.")
	query := `DELETE FROM events
	WHERE id =?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close() //stmt nesnesi başarıyla oluşturulmuşsa, kullanılmasından bağımsız olarak kapanacağı garanti edilir.

	result, err := stmt.Exec(e.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println(rowsAffected)

	if rowsAffected == 0 {
		return fmt.Errorf("kayıt bulunamadı, silinemedi")
	}

	return nil

}
func (e Event) Register(userId int64) error {
	query := `
		INSERT INTO registrations
		(event_id ,user_id)
		VALUES(?,?)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)
	if err != nil {
		fmt.Println("SQL INSERT error:", err)
		return err
	}

	_, err = db.DB.Exec("UPDATE events SET registered_users = registered_users + 1 WHERE id = ?", e.ID)
	if err != nil {
		fmt.Println("SQL UPDATE error:", err)
		return err
	}
	return err
}

func (e Event) DeleteRegistration(userId int64) error {
	query := `
	DELETE FROM registrations WHERE event_id=? and user_id =?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(e.ID, userId)
	if err != nil {
		fmt.Println("SQL DELETE error:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("kayıt bulunamadı, silinemedi")
	}

	if rowsAffected > 0 {
		_, err = db.DB.Exec("UPDATE events SET registered_users = registered_users - 1 WHERE id = ?", e.ID)
		if err != nil {
			fmt.Println("SQL UPDATE error:", err)
			return err
		}
	}
	return nil
}

func (e *Event) GetRegistrationCount() (int64, error) {
	query := `
	SELECT COUNT(*) FROM registrations WHERE event_id = ?
	`
	row := db.DB.QueryRow(query, e.ID)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
