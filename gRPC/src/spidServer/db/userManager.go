package db

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"spidServer/entities"
	eh "spidServer/errorHandling"
)

func (m *Manager) GetUsersFromFile() entities.Users {
	log.Print("Reading users.")
	usersFromFile, err := m.FileManager.ReadFile(UsersDefaultLocation)
	if err != nil {
		log.Fatalf("Failed to read users from file: %s", err)
		return entities.Users{}
	}
	users, err := entities.UnmarshalUsers(usersFromFile)
	if err != nil {
		log.Fatalf("Failed to parse users from file: %s", err)
		return entities.Users{}
	}
	if users.Users == nil {
		users.Users = map[uuid.UUID]entities.User{}
	}
	log.Printf("Read users: %s.", users.ToString())
	return users
}

func (m *Manager) WriteUsersToFile() {
	marshaledUsers, err := entities.MarshalUsers(m.Users)
	if err != nil {
		log.Printf("Failed to write users to file: %s", err)
		return
	}

	log.Printf("Making copy of users file...")
	src := m.FileManager.BasePath + string(os.PathSeparator) + UsersDefaultLocation
	dst := src + ".bk"
	source, err := os.Open(src)
	if err != nil {
		eh.HandleFatal(err)
	}
	defer source.Close()
	_ = os.Remove(dst)
	destination, err := os.Create(dst)
	if err != nil {
		eh.HandleFatal(err)
	}
	defer destination.Close()
	_, err := io.Copy(destination, source)
	eh.HandleFatal(err)

	log.Printf("Writing users: %s", m.Users.ToString())
	err = m.FileManager.WriteToFile(UsersDefaultLocation, marshaledUsers)
	if err != nil {
		log.Printf("Failed to write users to file: %s", err)
		return
	}
	log.Print("Finished writing users.")
}

func (m *Manager) QueryUser(userID uuid.UUID) (user entities.User, err error) {
	log.Printf("Querying user with ID %s.", userID)
	_, ok := m.Users.Users[userID]
	if !ok {
		err = fmt.Errorf("user with ID %s not found", userID)
		log.Print(err)
		return user, err
	}
	log.Printf("User found: %s", m.Users.Users[userID].ToString())
	return m.Users.Users[userID], nil
}

func (m *Manager) RegisterUser(userName string) (user entities.User, err error) {
	user = entities.NewUser(userName)
	log.Printf("Registering user: %s.", user.ToString())
	m.Users.Users[user.ID] = user
	log.Print("User registered.")
	return user, nil
}

func (m *Manager) UpdateUser(user entities.User) error {
	_, err := m.QueryUser(user.ID)
	if err != nil {
		return err
	}
	log.Printf("Updating user: %s.", user.ToString())
	m.Users.Users[user.ID] = user
	log.Print("User updated.")
	return nil
}

func (m *Manager) DeleteUser(user entities.User) error {
	_, err := m.QueryUser(user.ID)
	if err != nil {
		return err
	}
	log.Printf("Deleting user: %s.", user.ToString())
	delete(m.Users.Users, user.ID)
	log.Print("User deleted.")
	return nil
}
