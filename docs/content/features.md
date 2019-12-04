# Features

The following functionalities are provided by the Gasper Ecosystem

* Worker services for creating/managing databases and applications
* Master service for:-
    * Checking the status of worker services
    * Intelligently distributing applications/databases among them
    * Transferring applications from one worker node to another in case of node failure
    * Removing dead worker nodes from the cloud
* REST API interface for the entire ecosystem
* Reverse-proxy service with HTTPS, HTTP/2 and Websocket support for creating bridge connections to an application
* DNS service which automatically creates DNS entries for all applications which in turn are resolved inside containers
* SSH service for providing ssh access directly to an application's docker container
* Dynamic addition/deletion of nodes and services without configuration changes or restarts
* All of the above packaged with ❤️ in a **single binary**
