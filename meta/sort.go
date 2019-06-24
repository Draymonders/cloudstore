package meta

import (
	"cloudstore/config"
	"time"
)

type ByCreateTime []FileMeta

func (a ByCreateTime) Len() int {
	return len(a)
}

func (a ByCreateTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByCreateTime) Less(i, j int) bool {
	iTime, _ := time.Parse(config.BaseFormat, a[i].CreateTime)
	jTime, _ := time.Parse(config.BaseFormat, a[j].CreateTime)
	return iTime.UnixNano() > jTime.UnixNano()
}
