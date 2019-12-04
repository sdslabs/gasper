# Concepts

## Terminology

* **Node**: A Virtual Machine or a Bare-Metal Server
* **Worker Node**: A node which handles creation/management of applications and databases
* **Miscellaneous Node**: A node which handles Reverse-Proxy to other nodes or DNS entries or SSH access to docker containers
* **Master Node**: A node which manages all other nodes and keeps them in check

## Components

Gasper is divided into number of components each playing a specific role in the ecosystem. Lets have a look.

### Mizu 💧 

Mizu service deals with creating and managing applications and their life-cycles.

It currently supports applications of the following types

* Static web pages
* PHP
* Python 2
* Python 3
* Node.js
* Golang
* Ruby

It ain't much but it's honest work 🥳

!!!info
    A node with **Mizu** deployed is a Worker Node

### Kaen 🔥

Kaen service deals with creating and managing databases and their life-cycles

It currently supports databases of the following types

* MySQL
* MongoDB

It ain't.... (complete the rest yourself)

!!!info
    A node with **Kaen** deployed is a Worker Node

### Enrai ⚡

Enrai service deals with reverse-proxying HTTP, HTTPS, HTTP/2, Websocket and gRPC requests to the desired application's IPv4 address and port based on the hostname.

!!!info
    A node with **Enrai** deployed is a Miscellaneous Node

### Hikari 💡

Hikari service deals with creating and managing DNS records of all deployed applications. All DNS records point to the IPv4 addresses of Enrai ⚡ instances which in turn reverse-proxies the request to the desired application's IPv4 address and port.

!!!info
    **Hikari** stores DNS records in such a manner that all requests are equally distributed among all available **Enrai** instances. The records dynamically change with the addition/deletion of **Enrai** instances.

!!!info
    A node with **Hikari** deployed is a Miscellaneous Node

### Iwa 🗿 

Iwa service provides [SSH](https://www.ssh.com/ssh/protocol/) access directly to an application's docker container to the end user.
The SSH command will be automatically returned to the user on application creation provided the node where the application is deployed has the Iwa service deployed.

!!!info
    A node with **Iwa** deployed is a Miscellaneous Node

### Kaze 🌪 

Kaze is the master of the entire Gasper ecosystem which performs the following tasks

* Equal distribution of applications and databases among worker nodes
* User Authentication based on JWT (JSON Web Token)
* User API for performing operations on any application/database in any node (Identity Access Management is handled with JWT)
* Admin API for fetching and managing information of all nodes, applications, databases and users
* Removal of inactive nodes from the cloud ecosystem
* Re-scheduling of applications in case of node failure

Kaze API docs are available [here](/api)

!!!info
    You can interact with the entire Gasper ecosystem (example:- create/manage applications or databases) only through the REST API provided by **Kaze**

!!!info
    A node with **Kaze** deployed is a Master Node
