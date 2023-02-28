package repository

import (
	logger "go_rest_api_with_mysql/pkg/log"

	"database/sql"
	"go_rest_api_with_mysql/entity"
	"time"

	"go.uber.org/zap"
)

var log *zap.SugaredLogger = logger.GetLogger().Sugar()

// UserMySQL mysql repo
type UserMySQL struct {
	db *sql.DB
}

// NewUserMySQL create new repository
func NewUserMySQL(db *sql.DB) *UserMySQL {
	return &UserMySQL{
		db: db,
	}
}

// create an user
func (r *UserMySQL) CreateUser(e *entity.User) (entity.ID, error) {
	stmt, err := r.db.Prepare(`insert into tbl_users(id, email, password, first_name, last_name, created_at, updated_at) values(?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Errorf("Error prepare user statement. %v", err)
		return e.ID, err
	}

	_, err = stmt.Exec(
		e.ID,
		e.Email,
		e.Password,
		e.FirstName,
		e.LastName,
		time.Now().Format("2006-01-02"),
		time.Now().Format("2006-01-02"),
	)
	if err != nil {
		log.Errorf("unable to save user details. %v", err)
		return e.ID, err
	}
	err = stmt.Close()
	if err != nil {
		log.Errorf("error closing statement. %v", err)
		return e.ID, err
	}
	return e.ID, nil
}

// get an user
func (r *UserMySQL) GetUser(id entity.ID) (*entity.User, error) {
	return getUserDetails(id, r.db)
}

// get user details by id
func getUserDetails(id entity.ID, db *sql.DB) (*entity.User, error) {
	stmt, err := db.Prepare(`select id, email, first_name, last_name, created_at from tbl_users where id = ?`)
	if err != nil {
		log.Errorf("Error prepare select user statement. %v", err)
		return nil, err
	}
	var u entity.User
	rows, err := stmt.Query(id)
	if err != nil {
		log.Errorf("Error executing query. %v", err)
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt)
		if err != nil {
			log.Errorf("Error getting user details. %v", err)
			return nil, err
		}
	}
	return &u, nil
}

// get an user list
func (r *UserMySQL) GetUserList() ([]*entity.User, error) {
	stmt, err := r.db.Prepare(`select id, email, first_name, last_name, created_at from tbl_users`)
	if err != nil {
		log.Errorf("Unable to prepare user statement. %v", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	var users []*entity.User
	for rows.Next() {
		var email, first_name, last_name string
		var id entity.ID
		var created_at time.Time
		err = rows.Scan(&id, &email, &first_name, &last_name, &created_at)
		if err != nil {
			log.Errorf("Unable to query user details. %v", err)
			return nil, err
		}
		user := &entity.User{
			ID:        id,
			Email:     email,
			FirstName: first_name,
			LastName:  last_name,
			CreatedAt: created_at,
		}
		users = append(users, user)
	}
	return users, nil
}

// Update user
func (r *UserMySQL) UpdateUser(userConfigs *entity.User, userId entity.ID) error {
	userConfigs.UpdatedAt = time.Now()
	stmt, err := r.db.Exec(`update tbl_users set email = ?, first_name = ?, last_name = ?, updated_at = ? where id = ?`, userConfigs.Email, userConfigs.FirstName, userConfigs.LastName, userConfigs.UpdatedAt, userId)
	if err != nil {
		log.Errorf("Unable to update user details. %v", err)
		return err
	}
	log.Debug("User details updated.", stmt)
	return nil
}

// delete user
func (r *UserMySQL) DeleteUser(userId entity.ID) error {
	_, err := r.db.Exec(`delete from tbl_users where id = $1`, userId)
	if err != nil {
		log.Errorf("Unable to delete user details. %v", err)
		return err
	}
	return nil
}
