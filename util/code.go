package util

// ErrorCode : 错误类型
type ErrorCode int32

const (
	_ int32 = iota + (-1)
	// StatusOK : 0 正常
	StatusOK

	// StatusParamInvalid :  1 请求参数无效
	StatusParamInvalid

	// StatusServerError : 2 服务出错
	StatusServerError

	// StatusRegisterFailed : 3 注册失败
	StatusRegisterFailed

	// StatusLoginFailed : 4 登录失败
	StatusLoginFailed

	// StatusInvalidToken : 5 token无效
	StatusInvalidToken

	// StatusFormReadError : 6 form中读取信息失败
	StatusFormReadError

	// StatusCreateFileError : 7 创建文件失败
	StatusCreateFileError

	// StatusCopyFileError : 8 复制文件失败
	StatusCopyFileError

	// StatusFastUploadError : 9 秒传失败，请检查数据库数据
	StatusFastUploadError

	// StatusCephUploadError : 10 ceph 上传失败
	StatusCephUploadError

	// StatusKodoUploadError : 11 kdo 上传失败
	StatusKodoUploadError

	// StatusStoreToFileError : 12 更新tb_file表失败
	StatusStoreToFileError

	// StatusStoreToUserFileError : 13 更新tb_user_file表失败
	StatusStoreToUserFileError

	// StatusQueryUserFilesError : 14 查询 user file表失败
	StatusQueryUserFilesError

	// StatusFileMetaUpdateError : 15 文件元数据更新失败
	StatusFileMetaUpdateError

	// StatusQueryFileError 16 通过hash查询文件失败
	StatusQueryFileError

	// StatusDownloadError : 17 文件下载失败
	StatusDownloadError

	// StatusFileOpenError : 18 文件打开失败
	StatusFileOpenError

	// StatusFileReadError : 19 文件读取失败
	StatusFileReadError

	// StatusRedisGetError : 20 redis获取数据失败
	StatusRedisGetError

	// StatusPartNotFull : 21 分块没有上传完毕
	StatusPartNotFull

	// StatusMergeError : 22 分块归并失败
	StatusMergeError
)
