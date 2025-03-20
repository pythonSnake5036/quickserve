# Quickserve

A quick and simple HTTP file server written in Go.

# Usage

```shell
# Serves files in local directory on port 8080
quickserve

# Serves files in local directory on port 8000
quickserve -port 8000

# Serves files in the directory public
quickserve -dir public
```

## HTTPS

```shell
# Serves with HTTPS using the certificate cert.pem and key key.pem
quickserve -cert cert.pem -key key.pem
```