# s3www
Serve static files from any S3 compatible object storage endpoints.

## Install
If you do not have a working Golang environment, please follow [How to install Golang](https://golang.org/doc/install). Minimum version required is [go-latest](https://golang.org/dl/#stable)

```
go get github.com/harshavardhana/s3www
```

## Run with Let's Encrypt
Make sure you have `index.html` under `website-bucket`
```
s3www -endpoint "https://s3.amazonaws.com" -accessKey "accessKey" \
      -secretKey "secretKey" -bucket "website-bucket" \
      -lets-encrypt -address "example.com"

s3www: Started listening on https://example.com
```

### Test
Point your web browser to https://example.com ensure your `s3www` is serving your `index.html` successfully.


## Run locally
Make sure you have `index.html` under `website-bucket`
```
s3www -endpoint "https://s3.amazonaws.com" -accessKey "accessKey" \
      -secretKey "secretKey" -bucket "website-bucket"

s3www: Started listening on http://127.0.0.1:8080
```

### Test
Point your web browser to http://127.0.0.1:8080 ensure your `s3www` is serving your `index.html` successfully.

## License
This project is distributed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0), see [LICENSE](./LICENSE) for more information.

