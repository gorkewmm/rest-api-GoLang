package models

import "example/db"

type User struct {
	ID       int64
	Email    string `binding :"required"`
	Password string `binding :"required"`
}

func (user User) Save() error {
	query := `INSERT INTO users(email, password)
	values(?,?)`
	stmt, err := db.DB.Prepare(query) // query kullanılmaya hazır mı?kontrol et
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Email, user.Password) //Bu satır, az önce hazırlanan stmt objesini çalıştırır.
	if err != nil {
		return err
	}

	user.ID, err = result.LastInsertId()
	return err
}
