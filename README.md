# Gasper

> Your Cloud in a Binary

<img align="right" width="350px" height="400px" src="./docs/assets/gasperlogo.svg">

[![Build Status](https://api.travis-ci.org/sdslabs/gasper.svg)](https://travis-ci.org/sdslabs/gasper)
[![codecov](https://codecov.io/gh/sdslabs/gasper/branch/develop/graph/badge.svg)](https://codecov.io/gh/sdslabs/gasper)
[![Go Report Card](https://goreportcard.com/badge/github.com/sdslabs/gasper)](https://goreportcard.com/report/github.com/sdslabs/gasper)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/sdslabs/gasper/blob/develop/LICENSE.md)

Gasper is an intelligent Platform as a Service (PaaS) used for deploying and managing 
applications and databases in any cloud topology.

## Contents

* [Overview](#overview)
* [Features](#features)
* [Supported Languages](#supported-languages)
* [Supported Databases](#supported-databases)
* [Dependencies](#dependencies)
* [Quickstart](#quickstart)
* [Contributing](#contributing)
* [Meet the A-Team](#meet-the-a-team)
* [Contact](#contact)

## Overview

Imagine you have a couple of *Bare Metal Servers* and/or *Virtual Machines* (collectively called nodes) at your disposal. Now you want to deploy a couple of applications/services to these nodes in such a manner so as to not put too much load on a single node.<br>

Your 1st option is to manually decide which application goes to which node, then use ssh/telnet to manually
setup all of your applications in each node one by one.<br>

But you are smarter than that, hence you go for the 2nd option which is Kubernetes. You setup Kubernetes in all of your
nodes which forms a cluster, and now you can deploy your applications without worrying about load distribution. But
Kubernetes requires a lot of configuration for each application(deployments, services, stateful-sets etc) not to mention
pipelines for creating the corresponding docker image.<br>

Here comes (🥁drumroll please 🥁) **Gasper**, your 3rd option!<br>
Gasper builds and runs applications in docker containers **directly from source code** instead of docker images.
It requires minimal parameters for deploying an application, so minimal that you can count them on fingers in one hand 🤚. Same goes for Gasper provisioned databases. Gone are the days of hard labour (writing configurations).

## Features

Fear not because the reduction in complexity doesn't imply the reduction in features. You can rest assured because
Gasper has:-

* Worker services for creating/managing databases and applications
* Master service for:-
    * Checking the status of worker services
    * Intelligently distributing applications/databases among them
    * Transferring applications from one worker node to another in case of node failure
    * Removing dead worker nodes from the cloud
* REST APIs for master and worker services 
* Reverse-proxy service with HTTPS, HTTP/2 and Websocket support for creating bridge connections to an application
* DNS service which automatically creates DNS entries for all applications which in turn are resolved inside containers
* SSH service for providing ssh access directly to an application's docker container
* Dynamic addition/deletion of nodes and services without configuration changes or restarts
* All of the above packaged with ❤️ in a **single binary**

## Supported Languages

Gasper currently supports applications of the following types:-

* Static web pages
* PHP
* Python 2
* Python 3
* Node.js
* Golang

It ain't much but it's honest work 🥳

## Supported Databases

The following databases are supported by Gasper:-

* MySQL
* MongoDB

It ain't.... (complete the rest yourself)

## Dependencies

The following softwares are required for running Gasper:-

* [Golang 1.12.x](https://golang.org/dl/)
* [Docker](https://www.docker.com/)
    * [For Linux](https://runnable.com/docker/install-docker-on-linux)
    * [For MacOS](https://docs.docker.com/docker-for-mac/install/)
    * [For Windows](https://docs.docker.com/docker-for-windows/install/)
* [MongoDB](https://www.mongodb.com/download-center/community)
* [Redis](https://redis.io/download)

## Quickstart

Open your favourite terminal and perform the following tasks:-

1. Cross-check your golang version.

    ```bash
    $ go version
    go version go1.12.7 darwin/amd64
    ```

2. Clone this repository.

    ```bash
    $ git clone https://github.com/sdslabs/gasper
    ```

3. Go inside the cloned directory and list available *makefile* commands.

    ```bash
    $ cd gasper && make help

    Gasper: Your cloud in a binary

    install   Install missing dependencies
    build     Build the project binary
    tools     Install development tools
    start     Start in development mode with hot-reload enabled
    clean     Clean build files
    fmt       Format entire codebase
    vet       Vet entire codebase
    lint      Check codebase for style mistakes
    test      Run tests
    help      Display this help
    ```

4. Setup project configuration and make changes if required. The configuration file is well-documented so you
won't have a hard time looking around.

    ```bash
    $ cp config.sample.toml config.toml
    $ vim config.toml
    ```

5. Start the development server.

    ```bash
    $ make start
    ```

## Contributing

If you'd like to contribute to the project, refer to the [contributing documentation](./CONTRIBUTING.md).

## Meet the A-Team

* Anish Mukherjee [@alphadose](https://github.com/alphadose)
* Vaibhav [@vrongmeal](https://github.com/vrongmeal)
* Supratik Das [@supra08](https://github.com/supra08)
* Karanpreet Singh [@karan0299](https://github.com/karan0299)

## Contact

If you have a query regarding the product or just want to say hello then feel free to visit
[chat.sdslabs.co](http://chat.sdslabs.co/) or drop a mail at [contact@sdslabs.co.in](mailto:contact@sdslabs.co.in)
