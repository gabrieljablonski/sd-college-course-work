package consensus

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"spidServer/db"
	"strconv"
	"strings"
	"time"
)

const (
	RequestTimeout = 10 * time.Second
	DialTimeout = 2 * time.Second
	KeyPrefix = "cmd_"
)

type Manager struct {
	DBManager   db.Manager
	Client      *clientv3.Client
	KV          clientv3.KV
	LastCommand int
	Recovered   bool
}

func NewManager(endpoints []string, basePath string) Manager {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: DialTimeout,
	})
	log.Printf("Setting up etcd client with endpoints: %s", strings.Join(endpoints, " "))
	if err != nil {
		log.Fatalf("Failed to setup etcd client: %v", err)
	}
	log.Print("Finished setting up etcd client.")
	kv := clientv3.NewKV(client)
	return Manager{
		DBManager:   db.NewManager(basePath),
		Client:      client,
		KV:          kv,
		LastCommand: -1,
		Recovered:   false,
	}
}

func (m *Manager) Recover() error {
	ctx, _ := context.WithTimeout(context.Background(), RequestTimeout)
	options := []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
	}
	response, err := m.KV.Get(ctx, KeyPrefix, options...)
	if err != nil {
		log.Fatalf("Failed to recover commands from cluster: %v", err)
	}
	for _, item := range response.Kvs {
		k, v := string(item.Key), string(item.Value)
		log.Printf("Processing command %s: %s", k, v)
		err = m.ProcessCommand(k, v)
		if err != nil {
			log.Printf("Failed to recover command %s", k)
		}
	}
	log.Print("Finished recovering commands from cluster")
	m.Recovered = true
	return nil
}

func (m *Manager) PutCommand(cmd db.WriteAction) error {
	ctx, _ := context.WithTimeout(context.Background(), RequestTimeout)
	key := fmt.Sprintf("%s%03d", KeyPrefix, m.LastCommand+1)
	value, _ := cmd.Json()
	_, err := m.KV.Put(ctx, key, value)
	if err == nil {
		m.LastCommand++
	}
	return err
}

func (m *Manager) ProcessCommand(key, cmdString string) error {
	v := key[len(KeyPrefix):]
	cmdValue, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("failed to parse command number: %v", err)
	}
	if cmdValue <= m.LastCommand {
		return fmt.Errorf("command with value %d already processed. last was %d", cmdValue, m.LastCommand)
	}
	var cmd db.WriteAction
	err = json.Unmarshal([]byte(cmdString), &cmd)
	if err != nil {
		log.Fatalf("Failed to unmarshal command string `%s`: %v", cmdString, err)
	}
	err = m.DBManager.ProcessWriteAction(cmd)
	if err != nil {
		return err
	}
	m.LastCommand = cmdValue
	return nil
}

func (m *Manager) WatchChanges() {
	watchChannel := m.Client.Watch(context.Background(), KeyPrefix, clientv3.WithPrefix())
	for response := range watchChannel {
		for !m.Recovered {
			// wait until the previous commands are processed before processing new ones
			time.Sleep(time.Second) // diminish cpu hogging
		}
		for _, event := range response.Events {
			key, command := string(event.Kv.Key), string(event.Kv.Value)
			log.Printf("Detected new change: %s", command)
			err := m.ProcessCommand(key, command)
			if err != nil {
				log.Printf("Failed to process command `%s`:`%s`: %v", key, command, err)
			}
		}
	}
}
