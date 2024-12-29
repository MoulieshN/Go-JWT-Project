package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/MoulieshN/Go-JWT-Project.git/models"
)

type UserRepository interface {
	GetUser(userid string) (models.User, error)
	GetUsers(limit, offset int) ([]models.User, error)
	CreateTable() error
	CreateUser(user models.User) (int64, error)
	UpdateUserToken(token string, refreshToken string, userId string) error
	GetUserByEmail(email string) (models.User, error)
}

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) UserRepository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) CreateTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			user_id binary(16) NOT NULL,
			first_name varchar(32) NOT NULL,
			last_name varchar(32) DEFAULT NULL,
			user_type enum('ADMIN','USER') NOT NULL DEFAULT 'USER',
			email varchar(64) NOT NULL,
			phone varchar(10) NOT NULL,
			password VARCHAR(100) DEFAULT NULL,
			token VARCHAR(100) DEFAULT NULL,
			refresh_token VARCHAR(100) DEFAULT NULL,
			created_on datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_on datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id),
			UNIQUE KEY username (email)
		);
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.DB.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating product table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}
	log.Printf("Rows affected when creating table: %d", rows)
	return nil
}

func (r *Repository) GetUser(userid string) (models.User, error) {
	query := `SELECT user_id, username, first_name, last_name, user_type, email, phone, password, created_on, updated_on FROM users WHERE user_id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.DB.QueryRowContext(ctx, query, userid)

	var user models.User
	err := row.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.UserType, &user.Email, &user.Phone, &user.Password)
	if err != nil {
		log.Printf("Error %s when getting user", err)
		return user, err
	}
	return user, nil
}

func (r *Repository) GetUsers(limit, offset int) ([]models.User, error) {
	query := `SELECT user_id, username, first_name, last_name, user_type, email, phone, password, created_on, updated_on FROM users;`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		log.Printf("Error %s when getting users", err)
		return nil, err
	}

	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.UserType, &user.Email, &user.Phone, &user.Password); err != nil {
			log.Printf("Error %s when scanning user", err)
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error %s when closing rows", err)
		return nil, err
	}

	return users, nil
}

func (r *Repository) CreateUser(user models.User) (int64, error) {
	query := `INSERT INTO users (first_name, last_name, user_type, email, phone) VALUES (?, ?, ?, ?, ?)`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, query, user.FirstName, user.LastName, user.UserType, user.Email, user.Phone)
	if err != nil {
		log.Printf("Error %s when inserting user", err)
		return -1, err
	}
	num, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error %s when getting last insert ID", err)
		return -1, err
	}
	return num, nil
}

func (r *Repository) UpdateUserToken(token string, refreshToken string, userId string) error {
	query := `UPDATE users SET token = ?, refresh_token = ? WHERE user_id = ?`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.DB.ExecContext(ctx, query, token, refreshToken, userId)
	if err != nil {
		log.Printf("Error %s when updating user token", err)
		return err
	}
	return nil
}

func (r *Repository) GetUserByEmail(email string) (models.User, error) {
	query := `SELECT user_id, username, first_name, last_name, user_type, email, phone, password, created_on, updated_on FROM users WHERE email = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.DB.QueryRowContext(ctx, query, email)

	var user models.User
	err := row.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.UserType, &user.Email, &user.Phone, &user.Password)
	if err != nil {
		log.Printf("Error %s when getting user", err)
		return user, err
	}
	return user, nil

}
