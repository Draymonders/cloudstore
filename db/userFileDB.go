package db

import (
	mydb "cloudstore/db/mysql"
	"fmt"
)

// UserFile : user file property
type UserFile struct {
	Username     string
	FileName     string
	FileSize     int64
	Hash         string
	CreateTime   string
	LastEditTime string
	DownLoadUrl  string
}

// OnUserFileUploadFinished : user have finished upload file
func OnUserFileUploadFinished(username, filename, hash string, filesize int64) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tb_user_file(`username`,`filename`,`filesize`,`hash`,`status`) " +
			" values(?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failded to prepare statement, err: ", err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, filename, filesize, hash)
	if err != nil {
		fmt.Println("OnUserFileUploadFinished err :", err.Error())
		return false
	}
	return true
}

// QueryUserFileMetas : search file metas from username and limit
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select filename, filesize, hash, create_time, last_edit_time " +
			" from tb_user_file where username = ? and status = 1 limit ?")
	if err != nil {
		fmt.Println("Failded to prepare statement, err: ", err.Error())
		return nil, nil
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		fmt.Println("QueryUserFileMetas, err: ", err.Error())
		return nil, nil
	}
	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(&ufile.FileName, &ufile.FileSize, &ufile.Hash, &ufile.CreateTime, &ufile.LastEditTime)
		if err != nil {
			fmt.Println("QueryUserFileMetas err: ", err.Error())
			break
		}
		userFiles = append(userFiles, ufile)
	}
	return userFiles, nil
}

// RenameFilename : rename filename
func RenameFilename(username, hash, filename string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"update tb_user_file set filename = ? where username = ? and hash = ? limit 1")
	if err != nil {
		fmt.Println("Failded to prepare statement, err: ", err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(filename, username, hash)
	if err != nil {
		fmt.Println("RenameFilename err: ", err.Error())
		return false
	}
	return true
}

// DeleteUserFile : delete file
func DeleteUserFile(username, filehash string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"update tb_user_file set status=2 where username=? and hash=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

// QueryUserFileMeta :
func QueryUserFileMeta(username, hash string) (*UserFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select filename,filesize,hash,create_time," +
			"last_edit_time from tb_user_file where username=? and hash=?  limit 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, hash)
	if err != nil {
		return nil, err
	}

	ufile := UserFile{}
	if rows.Next() {
		err = rows.Scan(&ufile.FileName, &ufile.FileSize, &ufile.Hash,
			&ufile.CreateTime, &ufile.LastEditTime)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return &ufile, nil
}
