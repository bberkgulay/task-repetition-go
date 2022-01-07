# Task Repetition

It helps to remember task after done using spaced repetition principles.

### User Scenario

The user first adds the task that interests for user to the application and takes notes about it. After finishing the studies, user completes the task. Then the application reminds the user to remember this task by using the "spaced repetition" principle.

## Documentation

[Documentation](https://documenter.getpostman.com/view/18749435/UVXeqx7E#intro) on Postman

## Tech Stack

* Go for backend  
* MongoDB for database 
* Heroku for deployment
* Postman for api test and documentation

## Run Locally

Clone the project

```bash
  git clone https://github.com/bberkgulay/task-repetition-go
```

Go to the project directory

```bash
  cd task-repetition-go
```

Set the environments (.env)

```bash
  CONNECTION_URL=
  DATABASE_NAME=
  PORT=
  VERSIONING=

```

Start the server

```bash
  go run main.go
```
