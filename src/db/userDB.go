package db

import (
	mydb "db/mysql"
	"fmt"
)

type User struct {
	Username     string
	Email        string
	Phone        string
	CreateTime   string
	LastEditTime string
	Status       int
}

// UserSignUp : sign up a user
func UserSignUp(username string, passwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tb_user(`username`, `password`, `status`) values(?,?,1)")
	if err != nil {
		fmt.Println("Failded to prepare statement, err: ", err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}
	rows, err := ret.RowsAffected()
	fmt.Printf("rows:%d\n", rows)
	if nil == err && rows > 0 {
		fmt.Printf("user %s creat OK\n", username)
		return true
	}
	return false
}

// UserSignin : sign in a user
func UserSignin(username string, password string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"select * from tb_user where username= ? and password = ? limit 1")
	if err != nil {
		fmt.Println("stmt err:", err.Error())
		return false
	}
	defer stmt.Close()

	// query username
	rows, err := stmt.Query(username, password)
	if err != nil {
		fmt.Println("query user err: ", err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found: " + username)
		return false
	}
	fmt.Println("username:", username, " log in OK")
	return true
}

// UpdateToken : create token or update token
func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"replace into tb_user_token(username, user_token) values(?,?)")
	if err != nil {
		fmt.Println("stmt err:", err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println("update token err:", err.Error())
		return false
	}
	return true
}

// GetUserInfo : return user info of username
func GetUserInfo(username string) (User, error) {
	// create a user
	user := User{}

	stmt, err := mydb.DBConn().Prepare(
		"select username, create_time from tb_user where username = ?")
	if err != nil {
		fmt.Println("stmt err:", err.Error())
		return user, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.CreateTime)
	// fmt.Printf("user info: %s\n",)
	if err != nil {
		fmt.Println("user info query  err:", err.Error())
		return user, err
	}
	return user, nil
}
