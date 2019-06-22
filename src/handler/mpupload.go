package handler

import (
	rPool "cache/redis"
	. "config"
	mydb "db"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
	"util"

	"github.com/garyburd/redigo/redis"
)

const chunkSize = 5 * 1024 * 1024

// MultipartUploadInfo : multi part upload info
type MultipartUploadInfo struct {
	UploadID   string
	FileSize   int
	FileHash   string
	ChunkSize  int
	ChunkCount int
}

// InitialMultipartUploadHandler : init of multi part
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// 1. get param
	username := getUserName(r)
	filehash := getHash(r)
	filesize := getFileSize(r)

	// 2. get a redis connection
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3. create multi parts
	upInfo := MultipartUploadInfo{
		UploadID:   username + fmt.Sprintf("%x", time.Now().Unix()),
		FileSize:   filesize,
		FileHash:   filehash,
		ChunkSize:  chunkSize,
		ChunkCount: int(math.Ceil(float64(filesize) / float64(chunkSize))),
	}

	// 4. store to redis
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)

	// 5. return client ok
	w.Write(util.NewRespMsg(0, "OK", upInfo).JSONByte())
}

// UploadPartHandler  : upload part chunk file
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	fpath := DirPath + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Upload path create file failed", nil).JSONByte())
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)
	w.Write(util.NewRespMsg(0, "OK", nil).JSONByte())
}

// CompleteUploadHandler : complete upload then merge it
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := getUserName(r)
	filename := getFileName(r)
	filehash := getHash(r)
	filesize := getFileSize(r)
	uploadId := r.Form.Get("uploadid")

	// get a connection of redis
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadId))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "redis get all err:"+err.Error(), nil).JSONByte())
		return
	}

	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}

	if totalCount != chunkCount {
		w.Write(util.NewRespMsg(-1, "has not upload all, please wait...", nil).JSONByte())
		return
	}

	// merge all datas to one data
	// use shell to merge
	// cat `ls | sort -n` > /tmp/filename

	filepath := DirPath + "/uploadid"
	filestore := DirPath + filename
	if _, err := mergeAllPartFile(filepath, filestore); err != nil {
		w.Write(util.NewRespMsg(-2, "merge datas failed, please commit again...", nil).JSONByte())
		return
	}

	// store file to DB
	mydb.OnFileUploadFinished(filename, int64(filesize), filestore, filehash)
	mydb.OnUserFileUploadFinished(username, filename, filehash, int64(filesize))

	w.Write(util.NewRespMsg(0, "OK", nil).JSONByte())
}

// mergeAllPartFile: filepath： 分块存储的路径 filestore： 文件最终地址
func mergeAllPartFile(filepath, filestore string) (bool, error) {
	var cmd *exec.Cmd
	cmd = exec.Command(MergeAllShell, filepath, filestore)

	if _, err := cmd.Output(); err != nil {
		fmt.Println(err)
		return false, err
	} else {
		fmt.Println(filestore, " has been merge complete")
		return true, nil
	}
}
