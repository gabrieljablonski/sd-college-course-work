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
- ...
