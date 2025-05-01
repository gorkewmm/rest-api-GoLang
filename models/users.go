package models

import (
	"errors"
	"example/db"
	"example/utils"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `binding:"required"`
	Password string `binding:"required"`
	Role     string `json:"role"`
}

func (user User) Save() error {
	query := `INSERT INTO users(email, password,role)
	values(?,?,?)`
	stmt, err := db.DB.Prepare(query) // query kullanılmaya hazır mı?kontrol et
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(user.Email, hashedPassword, user.Role) // struct'taki email ve hashlenmiş şifre, veritabanına bu satırda gönderiliyor.
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	user.ID = id
	return err
}

// login yaparken haslenmıs sıfreyle, girilen şifreyi karsılaştırıyor
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

func GetAllUsers() ([]User, error) {
	query := `
	SELECT id, email, role FROM users
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var user []User
	for rows.Next() {
		var users User
		err = rows.Scan(&users.ID, &users.Email, &users.Role) // rows.scan struc veya değişken ister
		if err != nil {
			return nil, err
		}
		user = append(user, users)
	}
	return user, nil
}

func GetUserById(id int64) (*User, error) {
	query := `
	SELECT * FROM users WHERE ID =?
	`
	row := db.DB.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (user User) UserUpdate(id int64) error {
	query := `
	UPDATE users SET email =?,
	password =?, role=? WHERE id =?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Email, user.Password, user.Role, id)
	if err != nil {
		return err
	}
	return err
}

func DeleteUser(id int64) error {
	query := `
	DELETE FROM users WHERE id =?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return err
}

func (user User) ChangePassword(id int64) error {
	query := `
	UPDATE users SET password = ? WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(hashedPassword, id)
	return err
}
