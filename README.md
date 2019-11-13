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

## Planned Testing
- Stress testing for several vehicles operating simultaneously
- Server recovery after crash
- Client recovery after crash

## TCP Implementation

The original TCP implementation (including readme) is available with:

```git checkout d0f4eca```

## gRPC Implementation

The gRCP implementation used the original network protocol implemented in TCP as a baseline, applying it to the services/remote procedures model, described by [`.proto` files](src/spidServer/proto_files).

### Multiple Servers

The SPID network is no longer limited to one server. In the server bootstrap stage, a request is sent to a [server mapper](src/spidServer/spid_mapper.py), which is responsible for building the routing table of each one of the servers that will compose the network (number is set when creating the server instances).

#### Client to Server Connections

It is assumed that a client will somehow know the address to a random server in the network. In the client's point of view, it does not (and should not) need to know about the existence of the other servers.

#### Server to Server Connections

The server network can be thought of as a 2D square matrix with order `k`, with each server numbered from `0` to `n-1`, with ![n=k^2](http://latex.codecogs.com/gif.latex?n=k^2), from the top-left to the bottom-right. It should be noted the amount of servers can only be a perfect square (i.e. `1, 4, 9, 16, ...`).

One of the requisites, outlined in the previous section, is that a client connecting to the network must be able to access data maintained by **any** server by making requests to **any** other server. To accomplish this specification, a server must be connected to a reasonable amount of other servers (not all, and not too few), the following logic is being used: each server is directly connected to the servers with distances equal to powers of two (i.e. `1, 2, 4, 8...` in all six main cardinal directions (N, S, W, E, NW, NE, SW, SE). The following diagram represents this logic applied to some servers in a network with 2500 servers in total (50 by 50).

![SPID network connections](server_connections_2500.gif)

Thus, it can be shown that the maximum amount of connections `m` a single server will make is upper-bounded by the case in which we consider the server in the dead center of the matrix. That value is given by:

![Max server connections](http://latex.codecogs.com/gif.latex?m_{max}=6*\left\lfloor{log_2(n-1)}\right\rfloor)

This way, when a client makes a request to a random server, this server will be able to redirect the request to the one closest to the server responsible for the data in the request, process which can be repeated until reaching the target server.

#### Data Mapping

The distribution of the data related to the system entities is done in two ways across the network. 

The first one considers the entity unique identifier, taking into account the need for an uniform distribution, avoiding the overload of a small set of servers. To do that, it is assumed the algorithm used for the generation of the entity's [`uuid`](https://en.wikipedia.org/wiki/Universally_unique_identifier) is reasonably random. 

That way, it is possible to use a simple relation, as ![Compute server number with modulo](http://latex.codecogs.com/gif.latex?s=uuid\mod{n}), to compute the server number for the entity's "home agent", that is, the server which will always be responsible for that entity's data throughout all of it's existence in the network.

The second method of mapping offers both data redundancy and convenience for the network clients.
