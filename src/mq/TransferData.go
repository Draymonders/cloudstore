package mq

import (
	"config"
)

// TransferData : 传送的数据
type TransferData struct {
	FileHash      string
	CurPath       string
	DestPath      string
	DestStoreType config.StoreType
}
