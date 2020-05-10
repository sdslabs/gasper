# Concepts

## Terminology

* **Node**: A Virtual Machine or a Bare-Metal Server
* **Worker Node**: A node which handles creation/management of applications and databases
* **Support Node**: A node which supports the functioning of applications and databases in worker nodes
* **Master Node**: A node which manages all other nodes and keeps them in check

## Components

Gasper is divided into number of components each playing a specific role in the ecosystem. Lets have a look.

### AppMaker 💧 

AppMaker service deals with creating and managing applications and their life-cycles.

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
    A node with **AppMaker** deployed is a Worker Node

### DbMaker 🔥

DbMaker service deals with creating and managing databases and their life-cycles

It currently supports databases of the following types

* MySQL
* MongoDB
* PostgreSQL
* Redis

It ain't.... (complete the rest yourself)

!!!info
    A node with **DbMaker** deployed is a Worker Node

### GenProxy ⚡

GenProxy service deals with reverse-proxying HTTP, HTTPS, HTTP/2, Websocket and gRPC requests to the desired application's IPv4 address and port based on the hostname.

!!!info
    A node with **GenProxy** deployed is a Support Node

### GenDNS 💡

GenDNS service deals with creating and managing DNS records of all deployed applications. All DNS records point to the IPv4 addresses of GenProxy ⚡ instances which in turn reverse-proxies the request to the desired application's IPv4 address and port.

!!!note
    **GenDNS** stores DNS records in such a manner that all requests are equally distributed among all available **GenProxy** instances. The records dynamically change with the addition/deletion of **GenProxy** instances.

!!!info
    A node with **GenDNS** deployed is a Support Node

### GenSSH 🗿 

GenSSH service provides [SSH](https://www.ssh.com/ssh/protocol/) access directly to an application's docker container to the end user.
The SSH command will be automatically returned to the user on application creation provided the node where the application is deployed has the GenSSH service deployed.

!!!info
    A node with **GenSSH** deployed is a Support Node

### Master 🌪 

Master is the master of the entire Gasper ecosystem which performs the following tasks

* Equal distribution of applications and databases among worker nodes
* User Authentication based on JWT (JSON Web Token)
* User API for performing operations on any application/database in any node (Identity Access Management is handled with JWT)
* Admin API for fetching and managing information of all nodes, applications, databases and users
* Removal of inactive nodes from the cloud ecosystem
* Re-scheduling of applications in case of node failure

Master API docs are available [here](/api)

!!!info
    You can interact with the entire Gasper ecosystem (example:- create/manage applications or databases) only through the REST API provided by **Master**

!!!info
    A node with **Master** deployed is a Master Node
