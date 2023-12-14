# GO MongoDB

## Deskripsi

Dalam repo ini, belajar pemrograman Golang menggunakan basis data MongoDB. mempelajari cara terhubung ke MongoDn, mengirim perintah MongoDb. Menggunakan database driver yang disediakan oleh Third-party dari [Go Mongo Developer](go.mongodb.org/mongo-driver).

## Fitur

- Sistem CRUD
- Akses MongoDB
- Driver GO
- Rest API

## Cara Membuat Aplikasi

### Langkah 1: 

Pastikan sudah menginstal go.msi atau download menggunakan link ini:

[Download Go for Windows](https://go.dev/dl/go1.21.5.windows-amd64.msi).

### Langkah 2: 

Jalan kan perintah ini untuk membuat package main:
```
go mod init main
```

### Langkah 3: 

Jalan kan perintah ini untuk mendapatkan package driver:
```go
go get -u 
github.com/golang/snappy 
github.com/klauspost/compress 
github.com/montanaflynn/stats 
github.com/xdg-go/pbkdf2 
github.com/xdg-go/scram 
github.com/xdg-go/stringprep 
github.com/youmark/pkcs8 
go.mongodb.org/mongo-driver 
golang.org/x/crypto 
golang.org/x/sync 
golang.org/x/text
```

### Langkah 4: 

Buat Collection diMongoDB dengan Nama ``` Users ``` 

atau Import dari [Collection Users](https://github.com/panntod/Go-Mongo/tree/main/mongo).
