# Using Temporal.io for Workflow Orchestration

This is source code for the article [Evaluating Temporal for Workflow Orchestration](google.com).

## Prerequisites/Setup 

*  MySQL 5.7 or above
    * if running within Docker then,
        * Run `docker run -p 3307:3306  --name=mysqldb -e MYSQL_ROOT_PASSWORD=passwd -d mysql:5.7`
        * Access db using `docker exec -it mysqldb bash` && `mysql -u root -p`
        * Create database using `create database temporaldb` on mysql CLI
    
* Go 1.15 or above 
    * 

