# Distributed Systems - College Course Work

This project seeks to implement a distributed system for geolocalization and access control of means of transportion, typically electric (scooters, bikes, etc.).

## Project Specifications

### Embedded System
To be embedded on the vehicle. 
Responsible for keeping track of geoposition and making sure the vehicle is only usable once the server authorizes access.

### Server Side
Responsible for keeping track of general status of users, vehicles, and other equipment. 
Should be able to detect improper use of the vehicles (exiting a specified zone - geofencing -, for example).

### User Interface
Allows the user to request access to the server and unlock a vehicle.

## Planned testing
- Stress testing for several vehicles operating simultaneously
- Server recovery after crash
- Client recovery after crash

## Running on Docker

For simplicity's sake, the docker images `jablonski/spid-server`, `jablonski/spid-embedded`, and `jablonski/spid-user` can be used to run the whole system.
With docker installed, follow the instruction below.

### Setup SPID Server

Run the following command to pull and start running the server container:

`docker run -t -e SERVER_PORT="8001" --name spid_server jablonski/spid-server`

The `SERVER_PORT` environment variable can be set to change to a different port from the default one (8001). The `-t` flag is used to set output to the terminal.

### Setup SPID Clients

To setup the clients, you must obtain the server container IP address using 

`docker inspect spid_server` 

and looking for the `IPAddress` key. Alternatively, you can run the following command to output the address directly. Make sure the `spid_server` container is up before running the command.

`docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' spid_server`

With the address in hands, you can start instantiating the other two images with the following commands.

`docker run -it -e SERVER_HOST="<spid_server ip>" -e SERVER_PORT="8001" --name spid_embedded jablonski/spid-embedded`

`docker run -it -e SERVER_HOST="<spid_server ip>" -e SERVER_PORT="8001" --name spid_user jablonski/spid-user`

The `SERVER_HOST` must be set to the IP address obtained previously, as well as `SERVER_PORT`, if you changed the port from the default one when instantiating the server. The additional flag `-i` is used to direct input from the terminal to the container.
At this point, any number of client instances can be created.

The following command can be used to restart both server and client containers after stopping:

`docker start -ia <container name>`

The available commands for the SPID user client are:

| Command | Use |
|:-------:|:---:|
| `view user` | view local user info |
| `load user` | load user id from file |
| `register user` | register new user |
| `query user` | query user info from the server.<br>if no user is loaded, must type id manually |
| `save user` | save user id to file, which can be loaded after |
| `update location` | update current user location |
| `delete user` | delete current user from the server |
| `associate spid` | request association to spid.<br>it can be loaded from file or typed manually.<br>to get the spid id, run `view` from the desired spid client |
| `save spid` | save current associated spid id to file |
| `dissociate` | dissociate from current spid |
| `query spid` | query info for current associated spid |
| `exit` | close the process |

For the SPID embedded client, the commands are:


| Command | Use |
|:-------:|:---:|
| `view` | view local spid info |
| `load` | load spid id to be used from file |
| `register` | register new spid |
| `query` | query spid info from the server.<br>if no spid is loaded, must type id manually |
| `save` | save spid id to file, which can be loaded after |
| `update location` | update current spid location |
| `delete` | delete current spid from server |
| `exit` | close the process |

The save and load funcionality are useful for closing/reopening and not having to type the ids manually.


