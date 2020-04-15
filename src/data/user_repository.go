package data

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type IUserRepository interface {
	FindAll() []string
	Create(email string) bool
	CheckUserExist(email string) int64
	CheckUsersExist(emails []string) []int64
}

type UserRepository struct {
	DB *sql.DB
}

func (repo UserRepository) FindAll() []string {
	query := `SELECT email FROM user ORDER BY id;`

	rows, _ := repo.DB.Query(query)

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil
		}
		emails = append(emails, email)
	}

	return emails
}

func (repo UserRepository) Create(email string) bool {
	query := `INSERT INTO user (email) VALUES (?)`

	rows, err := repo.DB.Prepare(query)
	if err != nil {
		return false
	}

	rows.Exec(email)

	return true
}

func (repo UserRepository) CheckUserExist(email string) int64 {

	query := `SELECT id FROM user WHERE email =? limit 1;`

	var id int64
	row := repo.DB.QueryRow(query, email)
	err := row.Scan(&id)

	if err != nil {
		return -1
	}

	return id
}

func (repo UserRepository) CheckUsersExist(emails []string) []int64 {

	stmt := `
		select id from user where email in ('%s')
	`

	query := fmt.Sprintf(stmt, strings.Join(emails, "','"))

	rows, err := repo.DB.Query(query)
	if err != nil {
		fmt.Println(err)
	}

	var ids []int64
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids
}
