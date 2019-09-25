package utils

import (
	"github.com/google/uuid"
	"log"
	"main/entities"
)

const UsersDefaultLocation = "/db/users.spd"
const SpidsDefaultLocation = "/db/spids.spd"

type DBManager struct {
	FM FileManager
}

func NewManager(basePath string) DBManager {
	return DBManager{FM: FileManager{BasePath: basePath}}
}

func (d DBManager) GetUsersFromFile() entities.Users {
	log.Print("Reading users.")
	users := entities.UnmarshalUsers(d.FM.ReadFile(UsersDefaultLocation))
	log.Printf("Read users: %s.", users)
	return users
}

func (d DBManager) WriteUsersToFile(users []byte) {
	log.Printf("Writing users: %s", string(users))
	d.FM.WriteToFile(UsersDefaultLocation, users)
	log.Print("Finished writing users.")
}

func (d DBManager) writeUser(user entities.User) {
	users := d.GetUsersFromFile()
	users.Users[user.ID] = user
	log.Printf("Writing user: %s.", user)
	d.WriteUsersToFile(entities.MarshalUsers(users))
	log.Print("User written.")
}

func (d DBManager) QueryUser(userID uuid.UUID) entities.User {
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

func (d DBManager) RegisterUser(user entities.User) {
	if d.QueryUser(user.ID) != (entities.User{}) {
		log.Fatalf("User with ID %s already exists.", user.ID)
	}
	log.Printf("Registering user: %s.", user)
	d.writeUser(user)
	log.Print("User registered.")
}

func (d DBManager) UpdateUser(user entities.User) {
	if d.QueryUser(user.ID) == (entities.User{}) {
		log.Fatalf("User with ID %s doesn't exist.", user.ID)
	}
	log.Printf("Updating user: %s.", user)
	d.writeUser(user)
	log.Print("User updated.")
}

func (d DBManager) DeleteUser(user entities.User) {
	if d.QueryUser(user.ID) == (entities.User{}) {
		log.Fatalf("User with ID %s doesn't exist.", user.ID)
	}
	log.Printf("Deleting user: %s.", user)
	users := d.GetUsersFromFile()
	delete(users.Users, user.ID)
	d.WriteUsersToFile(entities.MarshalUsers(users))
	log.Print("User updated.")
}

func (d DBManager) GetSpidsFromFile() entities.Spids {
	log.Print("Reading spids.")
	spids := entities.UnmarshalSpids(d.FM.ReadFile(SpidsDefaultLocation))
	log.Printf("Read spids: %s.", spids)
	return spids
}

func (d DBManager) WriteSpidsToFile(spids []byte) {
	log.Printf("Writing spids: %s", string(spids))
	d.FM.WriteToFile(SpidsDefaultLocation, spids)
	log.Print("Finished writing spids.")
}

func (d DBManager) writeSpid(spid entities.Spid) {
	spids := d.GetSpidsFromFile()
	spids.Spids[spid.ID] = spid
	log.Printf("Writing spid: %s.", spid)
	d.WriteUsersToFile(entities.MarshalSpids(spids))
	log.Print("Spid written.")
}

func (d DBManager) QuerySpid(spidID uuid.UUID) entities.Spid {
	spids := d.GetSpidsFromFile()
	log.Printf("Querying spid with ID %s.", spidID)
	_, ok := spids.Spids[spidID]
	if ok {
		log.Printf("Spid found: %s", spids.Spids[spidID])
		return spids.Spids[spidID]
	}
	log.Printf("Spid with ID %s not found.", spids.Spids[spidID])
	return entities.Spid{}
}

func (d DBManager) RegisterSpid(spid entities.Spid) {
	if d.QuerySpid(spid.ID) != (entities.Spid{}) {
		log.Fatalf("Spid with ID %s already exists.", spid.ID)
	}
	log.Printf("Registering spid: %s.", spid)
	d.writeSpid(spid)
	log.Print("Spid registered.")
}

func (d DBManager) UpdateSpid(spid entities.Spid) {
	if d.QuerySpid(spid.ID) == (entities.Spid{}) {
		log.Fatalf("Spid with ID %s doesn't exist.", spid.ID)
	}
	log.Printf("Updating spid: %s.", spid)
	d.writeSpid(spid)
	log.Print("Spid updated.")
}

func (d DBManager) DeleteSpid(spid entities.Spid) {
	if d.QueryUser(spid.ID) == (entities.User{}) {
		log.Fatalf("Spid with ID %s doesn't exist.", spid.ID)
	}
	log.Printf("Deleting spid: %s.", spid)
	spids := d.GetSpidsFromFile()
	delete(spids.Spids, spid.ID)
	d.WriteSpidsToFile(entities.MarshalSpids(spids))
	log.Print("Spid updated.")
}
