
![Logo](https://rpsoftech.net/assets/svgs/logo-rp.svg)


# Multi Tenant Bullion Server

Multiple Microservice to manage multiple Bullion Application on minimal servers.


## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`APP_ENV` can `"LOCAL" || "PRODUCTION" || "CI"` 

`PORT` Like `"5000"`

`DB_URL` Mongo Database Url `"mongodb+srv://username:password@url.com/databasename"`

`DB_NAME` Database Name

`ACCESS_TOKEN_KEY` 128 char long string

`REFRESH_TOKEN_KEY` 128 char long string

`FIREBASE_JSON_STRING` Firebase Admin Serialized JSON

`FIREBASE_DATABASE_URL` Firebase Realtime Database URL

`REDIS_DB_URL` Redis Database URL `"host:port" "127.0.0.1:7777"`

`REDIS_DB_PASSWORD` (Optional)Redis Database Password

`REDIS_DB_DATABASE` (Optional)Redis Database Number

## Authors

- [Keyur Shah](https://www.github.com/keyurboss)

