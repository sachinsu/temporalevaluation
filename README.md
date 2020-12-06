## Setup 
-  Start Mysql in Docker, 
* Run `docker run -p 3307:3306  --name=mysqldb -e MYSQL_ROOT_PASSWORD=passwd -d mysql:5.7`
* Access db using `docker exec -it mysqldb bash`
