package meta

import (
	"sort"
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

// UpdateFileMeta : 新增/更新 文件元数据信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileName] = fmeta
}

// GetFileMeta : 获取 文件元数据信息
func GetFileMeta(filename string) FileMeta {
	return fileMetas[filename]
}

// RemoveFileMeta : 删除文件元信息
func RemoveFileMeta(filename string) {
	delete(fileMetas, filename)
}

// GetLastFileMetas : 返回最近上传的count个元数据信息
func GetLastFileMetas(count int) []FileMeta {
	ll := len(fileMetas)
	if count > ll {
		count = ll
	}
	fMetaArray := make([]FileMeta, ll)
	for _, v := range fileMetas {
		fMetaArray = append(fMetaArray, v)
	}
	sort.Sort(ByCreateTime(fMetaArray))
	return fMetaArray[:count]
}
