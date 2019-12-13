# Welcome

<img align="right" width="350px" height="400px" src="./assets/logo/gasperlogo.svg">

Gasper is an intelligent Platform as a Service (PaaS) used for deploying and managing 
applications and databases in any cloud topology.

## The Dilemma
Imagine you have a couple of *Bare Metal Servers* and/or *Virtual Machines* (collectively called nodes) at your disposal. Now you want to deploy a couple of applications/services to these nodes in such a manner so as to not put too much load on a single node.

## Naive Approach
Your 1st option is to manually decide which application goes to which node, then use ssh/telnet to manually
setup all of your applications in each node one by one.

## A Wise Choice
But you are smarter than that, hence you go for the 2nd option which is [Kubernetes](https://kubernetes.io/). You setup Kubernetes in all of your nodes which forms a cluster, and now you can deploy your applications without worrying about load distribution. But Kubernetes requires a lot of configuration for each application(deployments, services, stateful-sets etc) not to mention pipelines for creating the corresponding docker image.<br>

## The Ultimatum
Here comes (🥁drumroll please 🥁) **Gasper**, your 3rd option!<br>
Gasper builds and runs applications in docker containers **directly from source code** instead of docker images.
It requires minimal parameters for deploying an application, so minimal that you can count them on fingers in one hand 🤚. Same goes for Gasper provisioned databases. Gone are the days of hard labour (writing configurations).

!!!question "What is Gasper in a nutshell ?"
    Your Cloud in a Binary :)
