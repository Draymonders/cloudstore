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
