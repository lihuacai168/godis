# 什么是godis
- 使用`golang`实现的`redis cli`,二进制文件直接运行，无需任何依赖
- 支持安全模式，限制只能运行只读命令
- 支持常用`cluster`和`alone`模式的常用命令

# godis加入环境变量
- 下载release中对应平台的压缩包
- 解压，得到godis可执行文件
- 把`godis`可执行文件移动到`/usr/local/bin/`目录中（适用于Mac和Linux）
- `Windows`需要把`godis`所在的路径加入到系统的环境变量中

配置好之后，在命令行终端输入`godis`就能直接使用啦~
```
$ godis
A utility redis command line

Usage:
  godis [command]

Available Commands:
  config      Handle godis configuration
  del         Delete a key
  exists      Assure a key is exists
  hcopy       Copy a hash key
  hdel        Hash key hdel
  help        Help about any command
  hget        Hash key hget
  hgetall     Hash key hgetall
  hmset       Add a hash key, auto unpack jsonValue
  renamenx    Rename key, if new_key is exist return fail, else success
  sadd        Add a set key
  smembers    Get set key values
  ttl         Get key ttl
  type        Get key type
```



# 常用命令
## config
- 新增配置
```bash
$ godis config add local -a=127.0.0.1:6379 --desc="local redis client"
2021/03/12 10:29:39 not configured, please use "godis config add" to set a cluster configuration
Added cluster.
```
- 激活配置
```bash
$ godis config use local
2021/03/12 10:30:41 not configured, please use "godis config add" to set a cluster configuration
Switched to cluster "local".
```

## hash
- hmset
```
$ godis hmset myhash '{"a":1,"b":2}'
2021/03/12 10:35:23 connect success, using alone mode, conf: local, addr is [127.0.0.1:6379]
hmset success, hash key is myhash , value is
{
  "a": "1",
  "b": "2"
}
```

- hgetall  
hga是hgetall的aliase
```

$ godis hga myhash
2021/03/12 10:36:23 connect success, using alone mode, conf: local, addr is [127.0.0.1:6379]
{
  "a": "1",
  "b": "2"
}

$ godis hgetall myhash
2021/03/12 10:37:55 connect success, using alone mode, conf: local, addr is [127.0.0.1:6379]
{
  "a": "1",
  "b": "2"
}
```

- hget  
hg是hge的aliase
```
$ godis hget myhash a
2021/03/12 10:40:12 connect success, using alone mode, conf: local, addr is [127.0.0.1:6379]
1

$ godis hg myhash a
2021/03/12 10:40:27 connect success, using alone mode, conf: local, addr is [127.0.0.1:6379]
1
```

# set
- sadd
```
$ godis sadd myset a b c d
2021/03/12 10:41:35 connect success, using alone mode, conf: local, addr is [127.0.0.1:6379]
sadd success, set key is myset , value is
[
  "d",
  "c",
  "a",
  "b"
]


$ godis sadd myjsonset '{"a":1}' '{"b":2}'
2021/03/12 10:42:33 connect success, using alone mode, conf: local, addr is [127.0.0.1:6379]
sadd success, set key is myjsonset , value is
[
  {
    "a": 1
  },
  {
    "b": 2
  }
]
```
- smembers  
sget是smembers的aliase
```
$ godis smembers myset
2021/03/12 10:43:20 connect success, using alone mode, conf: local, addr is [127.0.0.1:6379]
[
  "d",
  "c",
  "a",
  "b"
]

$ godis sget myjsonset
2021/03/12 10:43:57 connect success, using alone mode, conf: local, addr is [127.0.0.1:6379]
[
  {
    "a": 1
  },
  {
    "b": 2
  }
]
```

# 安全模式
## 新增安全模式的配置
```
$ godis config add local_safe -a=127.0.0.1:6379 --desc="local safe redis client" --isSafeMode=true
2021/03/12 10:47:47 connect success, using alone mode, conf: local, addr is [127.0.0.1:6379]
Added cluster.
```

## 激活安全模式的配置
```
$ godis config use local_safe
2021/03/12 10:49:18 connect success, using alone mode, conf: local, addr is [127.0.0.1:6379]
Switched to cluster "local_safe".
```

## 安全模式正常操作  
支持的所有命令：`["hg", "hget", "hgetall", "hga", "sget", "smembers", "type", "config"]`
```
$ godis hga myhash
2021/03/12 10:50:41 connect success, using alone mode, conf: local_safe, addr is [127.0.0.1:6379]
{
  "a": "1",
  "b": "2"
}

$ godis sget myjsonset
2021/03/12 10:51:10 connect success, using alone mode, conf: local_safe, addr is [127.0.0.1:6379]
[
  {
    "a": 1
  },
  {
    "b": 2
  }
]

...
```

## 临时切换配置
```
$ godis del myhash -c local
2021/03/12 10:59:35 connect success, using alone mode, conf: local, addr is [127.0.0.1:6379]
delete key success
```
`-c`设置当前的操作使用`local`这个配置，只对当前这个命令有效，不会修改默认的配置。
通过`godis config ls`查看，默认配置还是`local_safe`
```
$ godis config ls
2021/03/12 11:01:57 connect success, using alone mode, conf: local_safe, addr is [127.0.0.1:6379]
[
  {
    "Addrs": [
      "127.0.0.1:6379"
    ],
    "Description": "local redis client",
    "IsSafeMode": false,
    "Name": "local",
    "Password": ""
  },
  {
    "Addrs": [
      "127.0.0.1:6379"
    ],
    "Description": "local safe redis client",
    "IsSafeMode": false,
    "Name": "local_safe",
    "Password": ""
  }
]
CurrentCluster: local_safe
```