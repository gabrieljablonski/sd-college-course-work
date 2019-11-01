package db

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"spidServer/entities"
	eh "spidServer/errorHandling"
)

func (m *Manager) GetRemoteUsersFromFile() *entities.Users {
	log.Print("Reading remote users.")
	usersFromFile, err := m.FileManager.ReadFile(DefaultRemoteUsersLocation)
	if err != nil {
		log.Printf("Failed to read remote users from file: %s", err)
		return nil
	}
	users, err := entities.UnmarshalUsers(usersFromFile)
	if err != nil {
		log.Printf("Failed to parse remote users from file: %s", err)
		return nil
	}
	if users.Users == nil {
		users = entities.NewUsers()
	}
	log.Printf("Read remote users: %s.", users)
	return users
}

func (m *Manager) WriteRemoteUsersToFile() {
	marshaledUsers, err := entities.MarshalUsers(m.RemoteUsers)
	if err != nil {
		log.Printf("Failed to write remote users to file: %s", err)
		return
	}
	log.Printf("Making backup of remote users file...")
	src := m.FileManager.BasePath + string(os.PathSeparator) + DefaultRemoteUsersLocation
	dst := src + ".bk"
	err = os.Rename(src, dst)
	eh.HandleFatal(err)

	log.Printf("Writing remote users: %s", m.RemoteUsers)
	err = m.FileManager.WriteToFile(DefaultRemoteUsersLocation, marshaledUsers)
	if err != nil {
		log.Printf("Failed to write remote users to file: %s", err)
		return
	}
	log.Print("Finished writing remote users.")
}

func (m *Manager) GetRemoteUsers() *entities.Users {
	log.Print("Querying remote users.")
	return m.RemoteUsers
}

func (m *Manager) QueryRemoteUser(userID uuid.UUID) (*entities.User, error) {
	log.Printf("Querying remote user with ID %s.", userID)
	s, ok := m.RemoteUsers.Users[userID]
	if !ok {
		err := fmt.Errorf("remote user with ID %s not found", userID)
		return nil, err
	}
	log.Printf("Remote user found: %s", s)
	return s, nil
}

func (m *Manager) AddRemoteUser(user *entities.User) error {
	log.Printf("Adding remote user: %s.", user)
	m.RemoteUsers.Users[user.ID] = user
	log.Print("Remote user added.")
	return nil
}

func (m *Manager) UpdateRemoteUser(user *entities.User) error {
	_, err := m.QueryRemoteUser(user.ID)
	if err != nil {
		return err
	}
	log.Printf("Updating remote user: %s.", user)
	m.RemoteUsers.Users[user.ID] = user
	log.Print("Remote user updated.")
	return nil
}

func (m *Manager) RemoveRemoteUser(user *entities.User) error {
	_, err := m.QueryUser(user.ID)
	if err != nil {
		return err
	}
	log.Printf("Removing remote user: %s.", user)
	delete(m.Users.Users, user.ID)
	log.Print("Remote user removed.")
	return nil
}

