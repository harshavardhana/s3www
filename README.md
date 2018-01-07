# s3www
Serve static files from any S3 compatible object storage endpoints.

## Install
```
go get github.com/harshavardhana/s3www
```

## Run
```
s3www -endpoint "https://s3.amazonaws.com" -accessKey "accessKey" -secretKey "secretKey" -bucket "website-bucket"
```

For more options look at `s3www -h`

## License
This project is licensed under Apache License version 2.0
