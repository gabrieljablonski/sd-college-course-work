package db

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"spidServer/entities"
	eh "spidServer/errorHandling"
)

func (m *Manager) GetUsersFromFile() *entities.Users {
	log.Print("Reading users.")
	usersFromFile, err := m.FileManager.ReadFile(DefaultUsersLocation)
	if err != nil {
		log.Fatalf("Failed to read users from file: %s", err)
		return nil
	}
	users, err := entities.UnmarshalUsers(usersFromFile)
	if err != nil {
		log.Fatalf("Failed to parse users from file: %s", err)
		return nil
	}
	if users.Users == nil {
		users = entities.NewUsers()
	}
	log.Printf("Read users: %s.", users)
	return users
}

func (m *Manager) WriteUsersToFile() {
	marshaledUsers, err := entities.MarshalUsers(m.Users)
	if err != nil {
		log.Printf("Failed to write users to file: %s", err)
		return
	}

	log.Printf("Making backup of users file...")
	src := m.FileManager.GetAbsolutePath(DefaultUsersLocation)
	dst := src + ".bk"
	err = os.Rename(src, dst)
	eh.HandleFatal(err)

	log.Printf("Writing users: %s", m.Users)
	err = m.FileManager.WriteToFile(DefaultUsersLocation, marshaledUsers)
	if err != nil {
		log.Printf("Failed to write users to file: %s", err)
		return
	}
	log.Print("Finished writing users.")
}

func (m *Manager) QueryUser(userID uuid.UUID) (*entities.User, error) {
	log.Printf("Querying user with ID %s.", userID)
	u, ok := m.Users.Users[userID]
	if !ok {
		err := fmt.Errorf("user with ID %s not found", userID)
		log.Print(err)
		return nil, err
	}
	log.Printf("User found: %s", u)
	return u, nil
}

func (m *Manager) RegisterUser(user *entities.User) error {
	log.Printf("Registering user: %s.", user)
	m.Users.Users[user.ID] = user
	log.Print("User registered.")
	return m.logWriteAction(WriteAction{
		Location:   Local,
		EntityType: User,
		Type:       Register,
		Entity:     user,
	})
}

func (m *Manager) UpdateUser(user *entities.User) error {
	_, err := m.QueryUser(user.ID)
	if err != nil {
		return err
	}
	log.Printf("Updating user: %s.", user)
	m.Users.Users[user.ID] = user
	log.Print("User updated.")
	return m.logWriteAction(WriteAction{
		Location:   Local,
		EntityType: User,
		Type:       Update,
		Entity:     user,
	})
}

func (m *Manager) DeleteUser(user *entities.User) error {
	_, err := m.QueryUser(user.ID)
	if err != nil {
		return err
	}
	log.Printf("Deleting user: %s.", user)
	delete(m.Users.Users, user.ID)
	log.Print("User deleted.")
	return m.logWriteAction(WriteAction{
		Location:   Local,
		EntityType: User,
		Type:       Delete,
		Entity:     user,
	})
}
