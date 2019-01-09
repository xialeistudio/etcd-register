# etcd-register

对etcd注册微服务的简单封装，通过一行代码向etcd注册微服务，并且会定时注册，保持微服务存活

```golang
service, err := NewService([]string{"http://192.168.1.144:2379"}, "/article/spider", "192.168.144:10000", nil)
if err != nil {
    log.Fatal("discovery initialize failed.", err.Error());
}

if err = service.Start(); err != nil {
    log.Fatal("discovery start failed.", err.Error())
}
```