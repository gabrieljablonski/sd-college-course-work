package main
//
//import (
//	"context"
//	"fmt"
//	"go.etcd.io/etcd/clientv3"
//	"strconv"
//	"time"
//)
//
//var (
//	dialTimeout    = 2 * time.Second
//	requestTimeout = 30 * time.Second
//)
//
//func main() {
//	ctx := context.Background()
//	cli, _ := clientv3.New(clientv3.Config{
//		Endpoints: []string{"127.0.0.1:12340", "127.0.0.1:12350", "127.0.0.1:12360"},
//		DialTimeout: dialTimeout,
//	})
//	defer cli.Close()
//	kv := clientv3.NewKV(cli)
//
//	for i := 0; i < 100; i++ {
//		k := fmt.Sprintf("key_%02d", i)
//		kv.Put(ctx, k, strconv.Itoa(i))
//	}
//
//	gr, _ := kv.Get(ctx, "key", clientv3.WithPrefix())
//
//	for _, item := range gr.Kvs {
//		fmt.Println(string(item.Key), string(item.Value))
//	}
//}