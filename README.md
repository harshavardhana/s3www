# s3www
Serve static files from any S3 compatible object storage endpoints.

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [s3www](#s3www)
    - [Install](#install)
    - [Binary](#binary)
    - [Container](#container)
    - [Auto TLS](#auto-tls)
- [License](#license)

<!-- markdown-toc end -->

## Install
Released binaries are available [here](https://github.com/harshavardhana/s3www/releases), or you can compile yourself from source.

> NOTE: minimum Go version needed is v1.17

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

# License
This project is distributed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0), see [LICENSE](./LICENSE) for more information.

<a class="bmc-button" target="_blank" href="https://www.buymeacoffee.com/y4m4"><img src="https://cdn.buymeacoffee.com/buttons/bmc-new-btn-logo.svg" alt="Buy me a coffee"><span style="margin-left:5px;font-size:24px !important;">Buy me a coffee</span></a>
