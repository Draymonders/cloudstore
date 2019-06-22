## ceph集群docker部署
通过docker可以快速部署小规模Ceph集群的流程，可用于开发测试。
以下的安装流程是通过linux shell来执行的；假设你只有一台机器，装了linux(如Ubuntu)系统和docker环境，那么可以参考以下步骤安装Ceph:


```bash
# 要用root用户创建, 或有sudo权限
# 注: 建议使用这个docker镜像源:https://registry.docker-cn.com
# 1. 修改docker镜像源
cat > /etc/docker/daemon.json << EOF
{
  "registry-mirrors": [
    "https://registry.docker-cn.com"
  ]
}
EOF
# 关于镜像， 这里需要用到三个: ceph/mon, ceph/osd, ceph/radosgw
# 如果下载不了的话，可以尝试下载我打包的: moxiaomomo/ceph-mon, moxiaomomo/ceph-osd, moxiaomomo/ceph-radosgw
# 下载完之后，可以重命名成官方镜像名，比如: docker tag moxiaomomo/ceph-mon:latest ceph/mon:latest
# 重启docker
systemctl restart docker
# 2. 创建Ceph专用网络
docker network create --driver bridge --subnet 172.20.0.0/16 ceph-network
docker network inspect ceph-network
# 3. 删除旧的ceph相关容器
docker rm -f $(docker ps -a | grep ceph | awk '{print $1}')
# 4. 清理旧的ceph相关目录文件，加入有的话
rm -rf /www/ceph /var/lib/ceph/  /www/osd/
# 5. 创建相关目录及修改权限，用于挂载volume
mkdir -p /www/ceph /var/lib/ceph/osd /www/osd/
chown -R 64045:64045 /var/lib/ceph/osd/
chown -R 64045:64045 /www/osd/
# 6. 创建monitor节点
docker run -itd --name monnode --network ceph-network --ip 172.20.0.10 -e MON_NAME=monnode -e MON_IP=172.20.0.10 -v /www/ceph:/etc/ceph ceph/mon
# 7. 在monitor节点上标识3个osd节点
docker exec monnode ceph osd create
docker exec monnode ceph osd create
docker exec monnode ceph osd create
# 8. 创建OSD节点
docker run -itd --name osdnode0 --network ceph-network -e CLUSTER=ceph -e WEIGHT=1.0 -e MON_NAME=monnode -e MON_IP=172.20.0.10 -v /www/ceph:/etc/ceph -v /www/osd/0:/var/lib/ceph/osd/ceph-0 ceph/osd 
docker run -itd --name osdnode1 --network ceph-network -e CLUSTER=ceph -e WEIGHT=1.0 -e MON_NAME=monnode -e MON_IP=172.20.0.10 -v /www/ceph:/etc/ceph -v /www/osd/1:/var/lib/ceph/osd/ceph-1 ceph/osd
docker run -itd --name osdnode2 --network ceph-network -e CLUSTER=ceph -e WEIGHT=1.0 -e MON_NAME=monnode -e MON_IP=172.20.0.10 -v /www/ceph:/etc/ceph -v /www/osd/2:/var/lib/ceph/osd/ceph-2 ceph/osd
# 9. 增加monitor节点，组件成集群
docker run -itd --name monnode_1 --network ceph-network --ip 172.20.0.11 -e MON_NAME=monnode_1 -e MON_IP=172.20.0.11 -v /www/ceph:/etc/ceph ceph/mon
docker run -itd --name monnode_2 --network ceph-network --ip 172.20.0.12 -e MON_NAME=monnode_2 -e MON_IP=172.20.0.12 -v /www/ceph:/etc/ceph ceph/mon
# 10. 创建gateway节点
docker run -itd --name gwnode --network ceph-network --ip 172.20.0.9 -p 9080:80 -e RGW_NAME=gwnode -v /www/ceph:/etc/ceph ceph/radosgw
# 11. 查看ceph集群状态
sleep 10 && docker exec monnode ceph -s
```

## 创建用户
```
docker exec -it gwnode radosgw-admin user create --uid=draymonder --display-name=draymonder
# 然后就得到如下信息
{
    "user_id": "draymonder",
    "display_name": "draymonder",
    "email": "",
    "suspended": 0,
    "max_buckets": 1000,
    "auid": 0,
    "subusers": [],
    "keys": [
        {
            "user": "draymonder",
            "access_key": "xxxxxxxxx",
            "secret_key": "xxxxxxxxxxxxxxxxxxxxxxxxxxx"
        }
    ],
    "swift_keys": [],
    "caps": [],
    "op_mask": "read, write, delete",
    "default_placement": "",
    "placement_tags": [],
    "bucket_quota": {
        "enabled": false,
        "max_size_kb": -1,
        "max_objects": -1
    },
    "user_quota": {
        "enabled": false,
        "max_size_kb": -1,
        "max_objects": -1
    },
    "temp_url_keys": []
}
```

## 参考链接
Docker简单部署Ceph测试集群 https://www.imooc.com/article/282861