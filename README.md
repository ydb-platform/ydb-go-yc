# ydb-go-yc

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/ydb-platform/ydb/blob/main/LICENSE)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ydb-platform/ydb-go-yc)](https://pkg.go.dev/github.com/ydb-platform/ydb-go-yc)
![tests](https://github.com/ydb-platform/ydb-go-yc/workflows/tests/badge.svg?branch=master)
![lint](https://github.com/ydb-platform/ydb-go-yc/workflows/lint/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/ydb-platform/ydb-go-yc)](https://goreportcard.com/report/github.com/ydb-platform/ydb-go-yc)
[![codecov](https://codecov.io/gh/ydb-platform/ydb-go-yc/badge.svg?precision=2)](https://app.codecov.io/gh/ydb-platform/ydb-go-yc)
![Code lines](https://sloc.xyz/github/ydb-platform/ydb-go-yc/?category=code)
[![WebSite](https://img.shields.io/badge/website-ydb.tech-blue.svg)](https://ydb.tech)

Helpers to connect to YDB inside yandex-cloud.

## Table of contents
1. [Overview](#Overview)
2. [About semantic versioning](#SemVer)
3. [Prerequisites](#Prerequisites)
4. [Installation](#Install)
5. [Usage](#Usage)

## Overview <a name="Overview"></a>

Currently package provides helpers to connect to YDB inside yandex-cloud.

## About semantic versioning <a name="SemVer"></a>

We follow the **[SemVer 2.0.0](https://semver.org)**. In particular, we provide backward compatibility in the `MAJOR` releases. New features without loss of backward compatibility appear on the `MINOR` release. In the minor version, the patch number starts from `0`. Bug fixes and internal changes are released with the third digit (`PATCH`) in the version.

There are, however, some changes with the loss of backward compatibility that we consider to be `MINOR`:
* extension or modification of internal `ydb-go-yc` interfaces. We understand that this will break the compatibility of custom implementations of the `ydb-go-yc` internal interfaces. But we believe that the internal interfaces of `ydb-go-yc` are implemented well enough that they do not require custom implementation. We are working to ensure that all internal interfaces have limited access only inside `ydb-go-yc`.
* major changes to (including removal of) the public interfaces and types that have been previously exported by `ydb-go-yc`. We understand that these changes will break the backward compatibility of early adopters of these interfaces. However, these changes are generally coordinated with early adopters and have the concise interfacing with `ydb-go-yc` as a goal.

Internal interfaces outside from `internal` directory are marked with comment such as
```
// Warning: only for internal usage inside ydb-go-yc
```

We publish the planned breaking `MAJOR` changes:
* via the comment `Deprecated` in the code indicating what should be used instead
* through the file [`NEXT_MAJOR_RELEASE.md`](#NEXT_MAJOR_RELEASE.md)

## Prerequisites <a name="Prerequisites"></a>

Requires Go 1.13 or later.

## Installation <a name="Installation"></a>

```bash
go get -u github.com/ydb-platform/ydb-go-yc
```

## Usage <a name="Usage"></a>

```go
import (
    yc "github.com/ydb-platform/ydb-go-yc"
)
...
    db, err := ydb.Open(ctx, os.Getenv("YDB_CONNECTION_STRING"),
        yc.WithInternalCA(),
        yc.WithServiceAccountKeyFileCredentials("~/.ydb/sa.json"), // auth from service account key file
        // yc.WithMetadataCredentials(), // auth inside cloud (virual machine or yandex function)
    )
    
```
