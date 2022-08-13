![s3www](https://raw.githubusercontent.com/harshavardhana/s3www/master/s3www.png)

Serve static files from any S3 compatible object storage endpoints.

> Simplifies and secure [AWS S3 Static Website Hosting](https://docs.aws.amazon.com/AmazonS3/latest/userguide/WebsiteHosting.html) by keeping your bucket private, secure, run-anywhere you like and with Automatic TLS based on Let's Encrypt.

## Features
- Automatic credentials rotation when deployed on AWS EC2, ECS or EKS services for your AWS S3 buckets - yay! 🔒😍
- Automatic certs renewal for your DOMAIN along with OCSP stapling, full suite of ACME features, HTTP->HTTPS redirection (all thanks to [certmagic](github.com/caddyserver/certmagic)).

## Install
Released binaries are available [here](https://github.com/harshavardhana/s3www/releases), or you can compile yourself from source.

> NOTE: minimum Go version needed is v1.18

```
go install github.com/harshavardhana/s3www@latest
```

## Binary
Make sure you have `index.html` under `mysite`
```
s3www -endpoint "https://s3.amazonaws.com" -accessKey "accessKey" \
	  -secretKey "secretKey" -bucket "mysite"

s3www: Started listening on http://127.0.0.1:8080
```

Point your web browser to http://127.0.0.1:8080 ensure your `s3www` is serving your `index.html` successfully.

## Container
Make sure you have `index.html` under `mysite`

```
podman run --rm -p 8080:8080 y4m4/s3www:latest \
	  -endpoint "https://s3.amazonaws.com" \
	  -accessKey "accessKey" \
	  -secretKey "secretKey" \
	  -bucket "mysite" \
	  -address "0.0.0.0:8080"

s3www: Started listening on http://0.0.0.0:8080
```

Point your web browser to http://127.0.0.1:8080 ensure your `s3www` is serving your `index.html` successfully.

## Auto TLS
Make sure you have `index.html` under `mysite`
```
s3www -endpoint "https://s3.amazonaws.com" -accessKey "accessKey" \
	  -secretKey "secretKey" -bucket "mysite" \
	  -lets-encrypt -address "example.com"

s3www: Started listening on https://example.com
```

Point your web browser to https://example.com ensure your `s3www` is serving your `index.html` successfully.

## Permissions

### AWS IAM Policy
s3www requires access to view and list all files in the bucket.

```
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "",
			"Effect": "Allow",
			"Action": "s3:GetObject",
			"Resource": "arn:aws:s3:::<Bucket Name>/*"
		},
		{
			"Sid": "",
			"Effect": "Allow",
			"Action": "s3:ListBucket",
			"Resource": "arn:aws:s3:::<Bucket Name>"
		}
	]
}
```
# License
This project is distributed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0), see [LICENSE](./LICENSE) for more information.
