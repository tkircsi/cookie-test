# Transfer cookie with HTTP redirect

## Goals

1. The goal is making a RESTful API which writes http cookie to the client, redirect the client to a different site which reads the http cookie.
2. Store and encrypt the custom claim in a jwt token and decript and read the other site
3. Make a Docker container to run the application

## Demo data

`UserData{ Name: "Kiss Csaba", Age: 45, Address: "1222 Budapest, Szent János utca 7.", Agree: true, FulkID: "FLK-001122", }`

## Docker Environment variables

**JWT_SECRET**: this will be used for encoding the jwt token
**ADDR**: the HTTP Server port. (Default: ":5000")

**REMOTE**: the remote URL to redirect (Default: "http://localhost:5000/redirpage")

## Usage

Run the container on site1 and site2. Site1 will redirect the client to site2.

### Site1

`docker run -it -p 5000:5000 --env REMOTE="http://[site2_url]:5002/redirpage" tkircsi/cookie-test`

### Site2

`docker run -it -p 5002:5002 --env ADDR=":5002" tkircsi/cookie-test`
