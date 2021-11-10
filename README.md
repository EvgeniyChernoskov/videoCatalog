#CRUD application providing Web API data

##Description 

Storing records in the database

## Install
1. https://github.com/EvgeniyChernoskov/videoCatalog
2. Create table (DATABASE: PostgreSQL)
```sql
CREATE TABLE videos
(
id          SERIAL       NOT NULL UNIQUE,
title       VARCHAR(255) NOT NULL,
description VARCHAR(255),
url         VARCHAR(255) NOT NULL
);
```
3. edit settings ./configs/config.yml and .env (DB_PASSWORD) 

### To start the application:

```
go run main.go
```
####Remark:
Logging to file: video.log

## API


###POST /videos
Creates new record

##### Example Input:
```json
{
"title": "SQL на примере Postgres",
"description": "Теория и практика SQL",
"url": "https://www.youtube.com/watch?v=i5-1HNf3W_Y"
}
```

### GET /videos/id

Get record by id

##### Example Response:
```json
{
"id": 7,
"title": "SQL на примере Postgres",
"description": "Теория и практика SQL",
"url": "https://www.youtube.com/watch?v=i5-1HNf3W_Y"
}
```

### PUT /videos/id

Edit record by id

##### Example Request :
```json
{   
    "id": 7,
    "title": "NEW",
    "description": "NEW",
    "url": "NEW"
}
```
### GET /videos
Get all records

### DELETE /videos/id

Remove record

