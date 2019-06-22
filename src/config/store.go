package config

// StoreType : 存储类型
type StoreType int

const (
	_ StoreType = iota
	// StoreLocal : 节点本地
	StoreLocal
	// StoreCeph : Ceph集群
	StoreCeph
	// StoreOSS : ali OSS
	StoreOSS
	// StoreMix : Ceph and OSS
	StoreMix
	// StoreAll : all store
	StoreAll
)

const (
	// TempLocalRootDir : 本地临时存储地址的路径
	TempLocalRootDir = "/d/tmp/"
	// CurrentStoreType : 设置当前文件的存储类型
	CurrentStoreType = StoreCeph
)
