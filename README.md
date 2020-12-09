# Using Temporal.io for Workflow Orchestration

This is source code for the article [Evaluating Temporal for Workflow Orchestration](https://sachinsu.github.io/posts/temporalworkflow/).

## Prerequisites/Setup 

*  MySQL 5.7 or above
    * if running within Docker then,
        * Run `docker run -p 3307:3306  --name=mysqldb -e MYSQL_ROOT_PASSWORD=passwd -d mysql:5.7`
        * Access db using `docker exec -it mysqldb bash` && `mysql -u root -p`
        * Create database using `create database temporaldb` on mysql CLI    
* Temporal  
    * Follow the instructions provided [here](https://docs.temporal.io/docs/install-temporal-server) to start Temporal server in Docker. 
* Install Go 1.15 or above 
* From root folder of repository,
    * Start Worker using `go run worker\main.go`
    * Start workflow using `go run start\main.go`
    * Start HTTP Service using `go run service\server.go`
    * As such workflow imports dummy set of users defined in `shared.go`. 
        * Check for list of unapproved users at http://localhost:8080/Users 
        * Simulate user approval by post Ids using CURL or any other http tool. Using [Curl](https://linuxize.com/post/curl-post-request/) `curl -X POST -d 'id\n1\2\n' https://localhost:8080/Users`
        * Revisit Users URL to verify that users are approved (i.e. Users with ID 1 and 2 are approved.)
    * Check the Temporal Web UI for details on workflow execution at `http://localhost:8088' (Replace localhost with Docker IP if needed)
