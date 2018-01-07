# s3www
Serve static files from any S3 compatible object storage endpoints.

## Install
```
go get github.com/harshavardhana/s3www
```

## Run
```
s3www -h

Usage of s3www:
  -accessKey string
    	Access key of S3 storage.
  -address string
    	Bind to a specific ADDRESS:PORT, ADDRESS can be an IP or hostname. (default ":8080")
  -bucket string
    	Bucket name which hosts static files.
  -endpoint string
    	S3 server endpoint. (default "https://s3.amazonaws.com")
  -secretKey string
    	Secret key of S3 storage.
```

```
s3www -endpoint "https://s3.amazonaws.com" -accessKey "accessKey" -secretKey "secretKey" -bucket "website-bucket"
```

## License
This project is licensed under Apache License version 2.0