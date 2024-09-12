# File Storage Service
## Prerequisites
1.  **Docker**: Ensure Docker is installed and running.
2.  **Docker-compose**: Ensure Docker Compose is installed.
## How To Run
### Start Server
At the current repository, run the following command:

`make setup`
### API Documentation
Access the API documentation by visiting the link: 

http://localhost:8080/docs/index.html
### Test
After starting the server as described above, you can test the endpoints in two ways:

**1. Using Postman**

Import the `File Storage Service.postman_collection.json` file from the `internal/testdata` directory into Postman.
Make sure to switch the test file in your directory to the appropriate form.

**2. Using CURL**

You can test the API with the following curl commands:

* Ping the server

`curl --location 'http://localhost:8080/api/v1/ping'`

* Upload a file

`curl --location 'http://localhost:8080/api/v1/upload' --form 'file=@"{file_location}"'`

* Get file data

`curl --location 'http://localhost:8080/api/v1/files-data'`

* Download a file

`curl --location 'http://localhost:8080/api/v1/download?file_id={file_id}'`
## Clean Up
Remove Docker Containers and Images

To stop and remove Docker containers and images, run:

`make down`

`make clean-image`

To remove everything, including Docker containers, images, and volumes:

`make clean-all`
## Restart and Rebuild Docker Images and Containers
To rebuild and restart the server:
`make rebuild`
