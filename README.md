# xk6-output-error

This is a [k6](https://go.k6.io/k6) extension using the [xk6](https://github.com/grafana/xk6) system.
`xk6-output-error` is a k6 extension to add more information into StdErr k6. Moreover, you can send logs with errors in your Elasticsearch or other log storage system used for this vector.dev, filebeat, etc. 

## Contents
* [Build](#build)
* [Usage](#usage)

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. [Install](https://github.com/grafana/xk6) `xk6`:

  ```shell
go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. Build the binary:

  ```shell
CGO_ENABLED=1
xk6 build --with github.com/tinkoff/xk6-output-error@latest
  ```

If you use Windows:

```shell
set CGO_ENABLED=1
xk6 build master --with github.com/tinkoff/xk6-output-error
```

## Usage

You can use default params:

```shell
k6 run -o xk6-output-error example.js
```

Also, you can add special fields for output:

```shell
k6 run -o xk6-output-error=fields="proto,tls_version" example.js
```

Or set an environment variable:

```shell
export K6_OUTPUTERROR_FIELDS="proto,tls_version"
k6 run -o xk6-output-error example.js
```

Moreover, you can set output into a file:

```shell
k6 run -o xk6-output-error=fields="proto,tls_version" example.js 2>error.log
```
