package user

import (
	"database/sql"
	"fmt"

	"github.com/TenacityLabs/retrospect-backend/types"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func scanRowIntoUser(row *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userStore *UserStore) GetUserByEmail(email string) (*types.User, error) {
	rows, err := userStore.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (userStore *UserStore) GetUserById(userId uint) (*types.User, error) {
	rows, err := userStore.db.Query("SELECT * FROM users WHERE id = ?", userId)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.ID != userId {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (userStore *UserStore) CreateUser(firstName string, lastName string, email string, password string) error {
	_, err := userStore.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)", firstName, lastName, email, password)
	if err != nil {
		return err
	}
	return nil
}

// delete user leaves memory leaks (eg. capsules where the owner is deleted)
// but this feature is only for testing, so it's fine
func (userStore *UserStore) DeleteUser(userId uint) error {
	_, err := userStore.db.Exec("DELETE FROM users WHERE id = ?", userId)
	if err != nil {
		return err
	}
	return nil
}

func (userStore *UserStore) UpdateUser(userId uint, firstName string, lastName string, email string) error {
	_, err := userStore.db.Exec("UPDATE users SET firstName = ?, lastName = ?, email = ? WHERE id = ?", firstName, lastName, email, userId)
	if err != nil {
		return err
	}
	return nil
}

func (userStore *UserStore) UpdateUserPassword(userId uint, password string) error {
	_, err := userStore.db.Exec("UPDATE users SET password = ? WHERE id = ?", password, userId)
	if err != nil {
		return err
	}
	return nil
}
