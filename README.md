# Gasper

> Your Cloud in a Binary

<img align="right" width="350px" height="400px" src="./docs/content/assets/logo/gasperlogo.svg">

[![Build Status](https://travis-ci.org/sdslabs/gasper.svg?branch=develop)](https://travis-ci.org/sdslabs/gasper)
[![Docs](https://img.shields.io/badge/docs-current-brightgreen.svg)](https://gasper-docs.netlify.com/)
[![Go Report Card](https://goreportcard.com/badge/github.com/sdslabs/gasper)](https://goreportcard.com/report/github.com/sdslabs/gasper)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/sdslabs/gasper/blob/develop/LICENSE.md)

Gasper is an intelligent Platform as a Service (PaaS) used for deploying and managing applications and databases in any cloud topology.

## Contents

* [Overview](#overview)
* [Features](#features)
* [Supported Languages](#supported-languages)
* [Supported Databases](#supported-databases)
* [Documentation](#documentation)
* [Dependencies](#dependencies)
* [Download](#download)
* [Development](#development)
* [Contributing](#contributing)
* [Meet the A-Team](#meet-the-a-team)
* [Contact](#contact)

## Overview

### The Dilemma
Imagine you have a couple of *Bare Metal Servers* and/or *Virtual Machines* (collectively called nodes) at your disposal. Now you want to deploy a couple of applications/services to these nodes in such a manner so as to not put too much load on a single node.

### Naive Approach
Your 1st option is to manually decide which application goes to which node, then use ssh/telnet to manually
setup all of your applications in each node one by one.

### A Wise Choice
But you are smarter than that, hence you go for the 2nd option which is [Kubernetes](https://kubernetes.io/). You setup Kubernetes in all of your
nodes which forms a cluster, and now you can deploy your applications without worrying about load distribution. But
Kubernetes requires a lot of configuration for each application(deployments, services, stateful-sets etc) not to mention
pipelines for creating the corresponding docker image.

### The Ultimatum
Here comes (🥁drumroll please 🥁) **Gasper**, your 3rd option!<br>
Gasper builds and runs applications in docker containers **directly from source code**. You no longer need to create application specific docker images and build pipelines, let Gasper do the heavylifting for you 😊.
Gasper requires minimal parameters for deploying an application, so minimal that you can count them on fingers in one hand 🤚. Same goes for Gasper provisioned databases. Gone are the days of hard labour (writing configurations).

## Features

Fear not because the reduction in complexity doesn't imply the reduction in features. You can rest assured because Gasper has:-

* Worker services for creating/managing databases and applications
* Master service for:-
    * Checking the status of worker services
    * Intelligently distributing applications/databases among them
    * Transferring applications from one worker node to another in case of node failure
    * Removing dead worker nodes from the cloud
* REST API interface for the entire ecosystem
* Reverse-proxy service with HTTPS, HTTP/2, Websocket and gRPC support for accessing deployed applications
* DNS service which automatically creates DNS entries for all applications which in turn are resolved inside containers
* SSH service for providing ssh access directly to an application's docker container
* Virtual terminal for interacting with your application's docker container from your browser
* Dynamic addition/removal of nodes and services without configuration changes or restarts
* Compatibility with Linux, Windows, MacOS, FreeBSD and OpenBSD
* All of the above packaged with ❤️ in a **single binary**

## Supported Languages

Gasper currently supports applications of the following types:-

* Static web pages
* PHP
* Python 2
* Python 3
* Node.js
* Golang
* Ruby
* Rust

It ain't much but it's honest work 🥳

## Supported Databases

The following databases are supported by Gasper:-

* MySQL
* MongoDB
* PostgreSQL
* Redis

It ain't.... (complete the rest yourself)

## Documentation

You can find the complete documentation of Gasper at [https://gasper-docs.netlify.app/](https://gasper-docs.netlify.app/)

## Dependencies

The only thing you need for running Gasper is [Docker](https://www.docker.com/). Here are the installation guides for:-

* [Linux](https://runnable.com/docker/install-docker-on-linux)
* [MacOS](https://docs.docker.com/docker-for-mac/install/)
* [Windows](https://docs.docker.com/docker-for-windows/install/)

If you perhaps need a higher degree of control over your entire cloud then you may setup [MongoDB](https://www.mongodb.com/download-center/community) and [Redis](https://redis.io/download) separately within your infrastructure and make the necessary changes in the `mongo` and `redis` sections of `config.toml`.

## Download

Assuming you have the [dependencies](#dependencies) installed, head over to Gasper's [releases](https://github.com/sdslabs/gasper/releases) page and grab the latest binary according to your operating system and system architecture

Run the downloaded binary with the [sample configuration file](./config.sample.toml)

```bash
$ ./gasper --conf ./config.toml
```

## Development

You need to have [Golang 1.13.x](https://golang.org/dl/) or higher installed along with the mentioned [dependencies](#dependencies)

Open your favourite terminal and perform the following tasks:-

1. Cross-check your golang version.

    ```bash
    $ go version
    go version go1.13.5 darwin/amd64
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
    release   Build release binaries
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
    ```

5. Start the development server.

    ```bash
    $ make start
    ```

## Contributing

If you'd like to contribute to this project, refer to the [contributing documentation](./CONTRIBUTING.md).

## Meet the A-Team

* Anish Mukherjee [@alphadose](https://github.com/alphadose)
* Vaibhav [@vrongmeal](https://github.com/vrongmeal)
* Supratik Das [@supra08](https://github.com/supra08)
* Karanpreet Singh [@karan0299](https://github.com/karan0299)
* Mohit Sharma [@Scar26](https://github.com/Scar26)

Gasper Logo: Leshna Balara [@leshnabalara](https://github.com/leshnabalara)

You can find the entire list of contributors [here](https://github.com/sdslabs/gasper/graphs/contributors)

Created with 💖 by [SDSLabs](https://github.com/sdslabs)

## Contact

If you have a query regarding the product or just want to say hello then feel free to visit
[chat.sdslabs.co](http://chat.sdslabs.co/) or drop a mail at [contact@sdslabs.co.in](mailto:contact@sdslabs.co.in)
