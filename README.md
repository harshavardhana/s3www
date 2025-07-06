# s3www: Serve Static Files from S3-Compatible Storage

s3www is a lightweight, open-source tool to serve static files from any S3-compatible object storage service, such as AWS S3, MinIO, or DigitalOcean Spaces. It supports HTTPS with Let's Encrypt for secure hosting and is ideal for hosting static websites, single-page applications (SPAs), or file servers.

## Features

- **S3 Compatibility**: Works with any S3-compatible storage provider.
- **HTTPS Support**: Automatic TLS certificates via Let's Encrypt.
- **Lightweight**: Built in Go for performance and minimal resource usage.
- **SPA Support**: Configurable single-page application routing.
- **Customizable**: Supports custom error pages and CORS configuration.
- **Cross-Platform**: Runs on Linux, macOS, FreeBSD, and more.

## Prerequisites

- An S3-compatible storage bucket with static files (e.g., `index.html`).
- S3 credentials with `s3:GetObject` and `s3:ListBucket` permissions.
- Go 1.21+ (for building from source) or Docker/Podman (for containerized deployment).
- Network access to your S3 endpoint and port 8080 (or your chosen port).

### Required S3 Permissions

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowGetObject",
      "Effect": "Allow",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::<Bucket Name>/*"
    },
    {
      "Sid": "AllowListBucket",
      "Effect": "Allow",
      "Action": "s3:ListBucket",
      "Resource": "arn:aws:s3:::<Bucket Name>"
    }
  ]
}
```

## Installation

### Option 1: Download Prebuilt Binaries

1. Visit the [Releases page](https://github.com/harshavardhana/s3www/releases) and download the binary for your platform (e.g., `s3www_0.9.0_linux_amd64.tar.gz`).
2. Extract the archive:
   ```bash
   tar -xzf s3www_0.9.0_linux_amd64.tar.gz
   ```
3. Move the binary to a directory in your PATH:
   ```bash
   sudo mv s3www /usr/local/bin/
   ```

### Option 2: Run with Docker

```bash
docker run --rm -p 8080:8080 y4m4/s3www:latest \
  -endpoint "https://s3.amazonaws.com" \
  -accessKey "accessKey" \
  -secretKey "secretKey" \
  -bucket "mysite" \
  -address "0.0.0.0:8080"
```

### Option 3: Build from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/harshavardhana/s3www.git
   cd s3www
   ```
2. Build the binary:
   ```bash
   go build
   ```
3. Move the binary to a directory in your PATH:
   ```bash
   sudo mv s3www /usr/local/bin/
   ```

## Usage

1. **Basic Command**:
   Serve files from an S3 bucket over HTTP:
   ```bash
   s3www -endpoint "https://s3.amazonaws.com" \
         -accessKey "accessKey" \
         -secretKey "secretKey" \
         -bucket "mysite"
   ```
   Output:
   ```
   s3www: Started listening on http://127.0.0.1:8080
   ```
   Open `http://127.0.0.1:8080` in your browser to verify.

2. **With Let's Encrypt**:
   Serve over HTTPS with automatic TLS certificates:
   ```bash
   s3www -endpoint "https://s3.amazonaws.com" \
         -accessKey "accessKey" \
         -secretKey "secretKey" \
         -bucket "mysite" \
         -lets-encrypt \
         -address "example.com"
   ```
   Output:
   ```
   s3www: Started listening on https://example.com
   ```
   Open `https://example.com` in your browser.

3. **Single-Page Application (SPA)**:
   Serve an SPA by specifying a fallback file:
   ```bash
   s3www -endpoint "https://s3.amazonaws.com" \
         -accessKey "accessKey" \
         -secretKey "secretKey" \
         -bucket "mysite" \
         -spa "index.html"
   ```

## Configuration

You can configure s3www via command-line flags or environment variables:

| Flag             | Environment Variable    | Description                              | Default              |
|------------------|------------------------|------------------------------------------|----------------------|
| `-endpoint`      | `S3WWW_ENDPOINT`       | S3 endpoint URL                          |                      |
| `-accessKey`     | `S3WWW_ACCESS_KEY`     | S3 access key                            |                      |
| `-secretKey`     | `S3WWW_SECRET_KEY`     | S3 secret key                            |                      |
| `-bucket`        | `S3WWW_BUCKET`         | S3 bucket name                           |                      |
| `-address`       | `S3WWW_ADDRESS`        | Host and port to listen on               | `127.0.0.1:8080`     |
| `-lets-encrypt`  |                        | Enable Let's Encrypt TLS                  | `false`              |
| `-spa`           | `S3WWW_SPA`            | Fallback file for SPA routing            |                      |
| `-tls-cert`      | `S3WWW_TLS_CERT`       | Path to TLS certificate file             |                      |
| `-tls-key`       | `S3WWW_TLS_KEY`        | Path to TLS key file                     |                      |

Example using environment variables:
```bash
export S3WWW_ENDPOINT="https://s3.amazonaws.com"
export S3WWW_ACCESS_KEY="accessKey"
export S3WWW_SECRET_KEY="secretKey"
export S3WWW_BUCKET="mysite"
export S3WWW_ADDRESS="0.0.0.0:8080"
s3www
```

## Use Cases

- **Static Website Hosting**: Serve HTML, CSS, and JavaScript files for blogs, portfolios, or documentation.
- **Single-Page Applications**: Host React, Vue, or Angular apps with proper routing.
- **File Sharing**: Share files securely from an S3 bucket over HTTPS.
- **Development Testing**: Quickly test static assets from S3 during development.

## Troubleshooting

- **404 Errors**: Ensure your bucket contains an `index.html` file or specify a custom SPA file with `-spa`.
- **Permission Denied**: Verify your S3 credentials have `s3:GetObject` and `s3:ListBucket` permissions.
- **Let's Encrypt Issues**: Ensure your domain is publicly accessible and port 80 is open for certificate issuance.
- **Connection Errors**: Check the S3 endpoint URL and your network connectivity.

## Contributing

We welcome contributions! To get started:

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature-name`).
3. Commit your changes (`git commit -m "Add feature"`).
4. Push to the branch (`git push origin feature-name`).
5. Open a pull request.

Please read the [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for more information.