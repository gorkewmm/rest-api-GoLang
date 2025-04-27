package models

import (
	"errors"
	"example/db"
	"example/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
	Role     string `binding:"required"`
}

func (user User) Save() error {
	query := `INSERT INTO users(email, password)
	values(?,?)`
	stmt, err := db.DB.Prepare(query) // query kullanılmaya hazır mı?kontrol et
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(user.Email, hashedPassword) // struct'taki email ve hashlenmiş şifre, veritabanına bu satırda gönderiliyor.
	if err != nil {
		return err
	}

	user.ID, err = result.LastInsertId()
	return err
}

func (u *User) ValidateCredentials() error {
	query := `SELECT id,password FROM users WHERE email = ?`
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)

	if err != nil { // If no user was found for that email
		return errors.New("Credentials invalid")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)
	if !passwordIsValid {
		return errors.New("Credentials invalid")
	}

	return nil

}

func FindUserByEmail(email string) (User, error) {
	query := `
	SELECT * FROM users where email =?
	`
	row := db.DB.QueryRow(query, email)
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return user, err
	}

	return user, nil
}
