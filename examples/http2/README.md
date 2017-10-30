## How to generate RSA private key and digital certificate

1. Install Openssl

Please visit https://github.com/openssl/openssl to get pkg and install.

2. Generate RSA private key

```sh
$ mkdir testdata
$ openssl genrsa -out ./testdata/server.key 2048
```

3. Generate digital certificate

```sh
$ openssl req -new -x509 -key ./testdata/server.key -out ./testdata/server.pem -days 365
```
