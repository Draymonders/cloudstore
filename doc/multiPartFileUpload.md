## 分块上传
在这里使用到了`redis`的`hash object`来存储分块信息

[![Zpz3xx.png](https://s2.ax1x.com/2019/06/22/Zpz3xx.png)](https://imgchr.com/i/Zpz3xx)

## 分块合并
由于项目中对大文件进行`5MB`为一个分块上传(多线程，提升上传效率)，因此所有分块上传完是要求合并的。

filepath: 分块所在目录,分块按照数字标号来有序存储
filestore: 生成一个新文件的目录位置

### 脚本文件
```bash
#!/bin/bash

filepath=$1
filestore=$2

echo "filepath: " $filepath
echo "filestorepath: "  $filestore

if [ ! -f $filestore ]; then
        echo "$filestore not exist"
else
        rm -f $filestore
fi

for item in `ls $filepath | sort -n`
do
        `cat $filepath/${item} >> ${filestore}`
        echo "merge ${filepath/${item}}  to $filestore ok"
done

echo "file store ok"
```

### Go进行脚本控制
```go
package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	// dirPath     = "/data/tmp/"
	dirPath = "d:\\tmp\\"
)

func main() {
	var cmd *exec.Cmd
	filepath := dirPath + "/root5d0dd1ec/"
	filestore := dirPath + "1111111.pdf"

	cmd = exec.Command(dirPath+"mergeAll.sh", filepath, filestore)
	// cmd.Run()
	if _, err := cmd.Output(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println(filestore, " has been merge complete")
	}
}
```