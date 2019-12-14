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
//	requestTimeout = 10 * time.Hour
//)
//
//func main() {
//	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
//	cli, _ := clientv3.New(clientv3.Config{
//		Endpoints: []string{"127.0.0.1:12340", "127.0.0.1:12350", "127.0.0.1:12360"},
//		DialTimeout: dialTimeout,
//	})
//	defer cli.Close()
//	kv := clientv3.NewKV(cli)
//
//	go func () {
//		time.Sleep(time.Second)
//		for i := 0; i < 100; i++ {
//			k := fmt.Sprintf("key_%02d", i)
//			kv.Put(ctx, k, strconv.Itoa(i))
//		}
//	}()
//
//	wch := cli.Watch(ctx, "key", clientv3.WithPrefix())
//	for wresp := range wch {
//		for _, ev := range wresp.Events {
//			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
//		}
//	}
//
//	//gr, _ := kv.Get(ctx, "key", clientv3.WithPrefix())
//	//fmt.Print(gr.Kvs)
//	//GetSingleValueDemo(ctx, kv)
//}