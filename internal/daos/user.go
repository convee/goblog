package daos

import (
	"database/sql"
	"errors"

	"github.com/convee/goblog/internal/model"
)

func GetUser(email string) (user model.User, err error) {
	row := db.QueryRow("select email,password from user where email=?", email)
	if err := row.Scan(&user.Email, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, sql.ErrNoRows
		}
		return user, err
	}
	return user, nil
}

func AddUser(user model.User) (id int, err error) {
	rs, err := db.Exec("insert into user (email, password) values (?, ?)", user.Email, user.Password)
	if err != nil {
		return
	}
	id64, err := rs.LastInsertId()
	return int(id64), err
}
