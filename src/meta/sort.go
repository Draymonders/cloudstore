package meta

import "time"

const baseFormat = "2006-01-02 15:04:05"

type ByCreateTime []FileMeta

func (a ByCreateTime) Len() int {
	return len(a)
}

func (a ByCreateTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByCreateTime) Less(i, j int) bool {
	iTime, _ := time.Parse(baseFormat, a[i].CreateTime)
	jTime, _ := time.Parse(baseFormat, a[j].CreateTime)
	return iTime.UnixNano() > jTime.UnixNano()
}
