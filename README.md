![GitHub](https://img.shields.io/github/license/AdrianOrlow/files-api)
# Files API

My personal file sharing service API. Made with Go, GORM, Gorilla Mux and MySQL.

[Files frontend](https://github.com/AdrianOrlow/files)

![thumbnail](https://user-images.githubusercontent.com/10941338/71479248-d0b0b800-27f3-11ea-96dd-2c98a82453d2.png)

## Getting started

Firstly, rename `.env.sample` to `.env` and fill all the fields with your data.
It should me mentioned that `ADMIN_GMAIL_ADDRESSES` is array of Google accounts email addresses which
can login to the system and perform CUD operations, separated with a comma.

Once you filled the config you can run the server via

```
go run main.go
```

If you want to build the package, run

```
go build
```

## Deployment (Dokku)

Create the app container

```
dokku apps:create app_name
```

Create the app container

```
dokku apps:create app_name
```

Add all the env variables

```
dokku config:set HASH_ID_SALT= ...
```

## License

[MIT](https://choosealicense.com/licenses/mit/)
