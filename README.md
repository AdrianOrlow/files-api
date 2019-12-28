![License](https://img.shields.io/github/license/AdrianOrlow/files-api)
![Go](https://img.shields.io/github/go-mod/go-version/AdrianOrlow/files-api)
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

create the mysql database container

```
dokku mysql:create app_name-db
```

link database to the container

```
dokku mysql:link app_name-db app_name
```

set all the environment variables
   
```
dokku config:set PORT=5000 HASH_ID_SALT= ...
```

create storage symlink

```
dokku storage:mount app_name /var/lib/dokku/data/storage/app_name:/storage
```

add Dokku remote repository

```
git remote add dokku dokku@server_ip:app_name
```

and finally push code to the repo

```
git push dokku master
```

## License

[MIT](https://choosealicense.com/licenses/mit/)
