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

## Single File Mode

Single file mode redirects all requests to a single file, which is useful for quickly transferring a single file from
one device to another.

To use single file mode, pass the path to a file with the `-file` flag

```shell
# Uses single file mode to redirect all requests to file.txt
quickserve -file file.txt
```
