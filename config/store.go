package config

// StoreType : 存储类型
type StoreType int

const (
	_ StoreType = iota
	// StoreLocal : 节点本地
	StoreLocal
	// StoreCeph : Ceph集群
	StoreCeph
	// StoreKodo : qiniu kodo
	StoreKodo
	// StoreMix : Ceph and OSS
	StoreMix
	// StoreAll : all store
	StoreAll
)

// CurrentStoreType : 当前存储类型
const CurrentStoreType = StoreKodo
