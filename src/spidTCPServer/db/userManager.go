package db

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"main/entities"
)

// TODO: error handling for file stuff

func (d Manager) GetUsersFromFile() entities.Users {
	log.Print("Reading users.")
	users := entities.UnmarshalUsers(d.FM.ReadFile(UsersDefaultLocation))
	log.Printf("Read users: %s.", users)
	return users
}

func (d Manager) WriteUsersToFile(users []byte) {
	log.Printf("Writing users: %s", string(users))
	d.FM.WriteToFile(UsersDefaultLocation, users)
	log.Print("Finished writing users.")
}

func (d Manager) writeUser(user entities.User) {
	users := d.GetUsersFromFile()
	users.Users[user.ID] = user
	log.Printf("Writing user: %s.", user)
	d.WriteUsersToFile(entities.MarshalUsers(users))
	log.Print("User written.")
}

func (d Manager) QueryUser(userID uuid.UUID) (entities.User, error) {
	users := d.GetUsersFromFile()
	log.Printf("Querying user with ID %s.", userID)
	_, ok := users.Users[userID]
	if !ok {
		err := fmt.Errorf("user with ID %s not found", userID)
		log.Print(err)
		return entities.User{}, err
	}
	log.Printf("User found: %s", users.Users[userID])
	return users.Users[userID], nil
}

func (d Manager) RegisterUser(user entities.User) error {
	_, err := d.QueryUser(user.ID)
	if err != nil {
		return err
	}
	log.Printf("Registering user: %s.", user)
	d.writeUser(user)
	log.Print("User registered.")
	return nil
}

func (d Manager) UpdateUser(user entities.User) error {
	_, err := d.QueryUser(user.ID)
	if err != nil {
		return err
	}
	log.Printf("Updating user: %s.", user)
	d.writeUser(user)
	log.Print("User updated.")
	return nil
}

func (d Manager) DeleteUser(user entities.User) error {
	_, err := d.QueryUser(user.ID)
	if err != nil {
		return err
	}
	log.Printf("Deleting user: %s.", user)
	users := d.GetUsersFromFile()
	delete(users.Users, user.ID)
	d.WriteUsersToFile(entities.MarshalUsers(users))
	log.Print("User deleted.")
	return nil
}
