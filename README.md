# About

Nhite stands for _Nhite is hip-terraform_

TODO: A lot of documentation, code and bug-fixing....

## Getting started

`go get -v -u github.com/nhite/nhite`

## Setting up a certificate

nhite is a grpc service and it requires a TLS/SSL connexion.
Therefore is requires ssl certificates.

You can generate tests certificates with `openssl` or you can use the tool [certstrap](https://github.com/square/certstrap) which does not rely on openssl.

### Quick start

```shell
certstrap init --common-name "test" 
certstrap request-cert -ip 127.0.0.1
certstrap sign 127.0.0.1 --CA test
```

Then points the env variables to the correct files:

```shell
export NHITE_CERT_FILE="out/127.0.0.1.crt"
export NHITE_KEY_FILE="out/127.0.0.1.key"
```

# FAQ

## I have an error `panic: http: multiple registrations for /debug/requests` in runtime

This is related to this [issue](https://github.com/grpc/grpc-go/issues/566).
Please remove the directory `$GOPATH/src/github.com/hashicorp/terraform/vendor/golang.org/x/net/trace` and build the tool again
