package mysql

import "blog/pkg/model"

func GetUser(email string) (user model.User) {
	row := db.QueryRow("select email,password from user where email=?", email)
	row.Scan(&user.Email, &user.Password)
	return
}

func AddUser(user model.User) (id int, err error) {
	rs, err := db.Exec("insert into user (email, password) values (?, ?)", user.Email, user.Password)
	if err != nil {
		return
	}
	id64, err := rs.LastInsertId()
	return int(id64), err
}
