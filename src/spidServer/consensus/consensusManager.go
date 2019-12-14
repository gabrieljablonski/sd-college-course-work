package consensus

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"spidServer/db"
	"strings"
	"time"
)

const (
	RequestTimeout = 10 * time.Second
	DialTimeout = 2 * time.Second
	KeyPrefix = "cmd"
)

type Manager struct {
	DBManager   db.Manager
	Client      *clientv3.Client
	KV          clientv3.KV
	LastCommand int
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
		DBManager: db.NewManager(basePath),
		Client:      client,
		KV:          kv,
		LastCommand: 0,
	}
}

func (m *Manager) PutCommand(cmd db.WriteAction) error {
	ctx, _ := context.WithTimeout(context.Background(), RequestTimeout)
	key := fmt.Sprintf("%s_%03d", KeyPrefix, m.LastCommand)
	value, _ := cmd.Json()
	_, err := m.KV.Put(ctx, key, value)
	if err == nil {
		m.LastCommand++
	}
	return err
}

func (m *Manager) ProcessCommand(cmdString string) error {
	var cmd db.WriteAction
	err := json.Unmarshal([]byte(cmdString), &cmd)
	if err != nil {
		log.Fatalf("Failed to unmarshal command string `%s`: %v", cmdString, err)
	}
	return m.DBManager.ProcessWriteAction(cmd)
}

func (m *Manager) WatchChanges() {
	watchChannel := m.Client.Watch(context.Background(), KeyPrefix, clientv3.WithPrefix())
	for response := range watchChannel {
		for _, event := range response.Events {
			command := string(event.Kv.Value)
			log.Printf("Detected new change: %s", command)
			err := m.ProcessCommand(command)
			if err != nil {
				log.Fatalf("Failed to process command  `%s`:`%s`: %v", event.Kv.Key, command, err)
			}
		}
	}
}
