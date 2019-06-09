package meta

import (
	mydb "db"
	"fmt"
)

// FileMeta : 文件元信息结构
type FileMeta struct {
	FileName   string
	FilePath   string
	FileSize   int64
	Hash       string
	CreateTime string
}

var fileMetas map[string]FileMeta

// init : 初始化 fileMetas key为FileName,value为FileMeta
func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta : 新增 文件元数据信息
// func UpdateFileMeta(fmeta FileMeta) {
// 	fileMetas[fmeta.FileName] = fmeta
// }

// GetFileMeta : 获取 文件元数据信息
// func GetFileMeta(filename string) FileMeta {
// 	return fileMetas[filename]
// }

// RemoveFileMeta : 删除文件元信息
// func RemoveFileMeta(filename string) {
// 	delete(fileMetas, filename)
// }

// GetLastFileMetas : 返回最近上传的count个元数据信息
// func GetLastFileMetas(count int) []FileMeta {
// 	ll := len(fileMetas)
// 	if count > ll {
// 		count = ll
// 	}
// 	fMetaArray := make([]FileMeta, ll)
// 	for _, v := range fileMetas {
// 		fMetaArray = append(fMetaArray, v)
// 	}
// 	sort.Sort(ByCreateTime(fMetaArray))
// 	return fMetaArray[:count]
// }

// CreateFileMetaDB : store file meta to DB
func CreateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(fmeta.FileName, fmeta.FileSize, fmeta.FilePath, fmeta.Hash)
}

// GetFileMetaDB : get file meta from db
func GetFileMetaDB(filename string) FileMeta {
	tfile, err := mydb.GetFileMeta(filename)
	if err != nil {
		fmt.Println("GetFileMeta from db err: ", err.Error())
		return FileMeta{}
	}
	fileMeta := FileMeta{
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		FilePath: tfile.FilePath.String,
		Hash:     tfile.Hash}
	return fileMeta
}

// GetFileMetaListsDB : get file meta lists from db
func GetFileMetaListsDB(count int) []FileMeta {
	tfiles, err := mydb.GetFileMetaLists(count)
	if err != nil {
		fmt.Println("GetFileMetaListsDB err: ", err.Error())
		return nil
	}
	fileMetas := make([]FileMeta, count)
	for _, tfile := range tfiles {
		fileMeta := FileMeta{
			FileName: tfile.FileName.String,
			FileSize: tfile.FileSize.Int64,
			FilePath: tfile.FilePath.String,
			Hash:     tfile.Hash}
		fileMetas = append(fileMetas, fileMeta)
	}
	return fileMetas
}

// RemoveFileMeta : set the file meta col's status = 2 (unvaild)
func RemoveFileMetaDB(filename string) bool {
	return mydb.OnFileRemoved(filename)
}

// UpdateFileMetaFromfilename : update oldFilename to newFilename
func UpdateFileMetaFromfilenameDB(oldFilename, newFilename string) bool {
	return mydb.OnFileMetaUpdate(oldFilename, newFilename)
}
