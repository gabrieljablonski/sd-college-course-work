package db

import (
	"github.com/google/uuid"
	"log"
	"main/entities"
)

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

func (d Manager) QueryUser(userID uuid.UUID) entities.User {
	users := d.GetUsersFromFile()
	log.Printf("Querying user with ID %s.", userID)
	_, ok := users.Users[userID]
	if ok {
		log.Printf("User found: %s", users.Users[userID])
		return users.Users[userID]
	}
	log.Printf("User with ID %s not found.", users.Users[userID])
	return entities.User{}
}

func (d Manager) RegisterUser(user entities.User) {
	if d.QueryUser(user.ID) != (entities.User{}) {
		log.Fatalf("User with ID %s already exists.", user.ID)
	}
	log.Printf("Registering user: %s.", user)
	d.writeUser(user)
	log.Print("User registered.")
}

func (d Manager) UpdateUser(user entities.User) {
	if d.QueryUser(user.ID) == (entities.User{}) {
		log.Fatalf("User with ID %s doesn't exist.", user.ID)
	}
	log.Printf("Updating user: %s.", user)
	d.writeUser(user)
	log.Print("User updated.")
}

func (d Manager) DeleteUser(user entities.User) {
	if d.QueryUser(user.ID) == (entities.User{}) {
		log.Fatalf("User with ID %s doesn't exist.", user.ID)
	}
	log.Printf("Deleting user: %s.", user)
	users := d.GetUsersFromFile()
	delete(users.Users, user.ID)
	d.WriteUsersToFile(entities.MarshalUsers(users))
	log.Print("User updated.")
}
