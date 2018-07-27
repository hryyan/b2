# Go client and library for B2 Cloud Storage

[![Build Status](https://travis-ci.org/hryyan/b2.svg)](https://travis-ci.org/hryyan/b2)
[![codecov](https://codecov.io/gh/hryyan/b2/branch/master/graph/badge.svg)](https://codecov.io/gh/hryyan/b2)
[![Go Report Card](https://goreportcard.com/badge/github.com/hryyan/b2)](https://goreportcard.com/report/github.com/hryyan/b2)

[README](README.md) | [中文文档](README_zh.md)

## What is b2?

![Alt text](./doc/intro.svg)

## What can I do with b2?

* Transmit file across machines.
* OSS for your web server.
* Share file with your partner and friend.
* Intergrate b2 library in your program.

## How to use b2?

1. Download the latest programs from [Release](https://github.com/hryyan/b2/releases) page according to your os and arch.

2. Register a b2 account by
```shell
b2 register
```

3. Set environment `B2_ACCOUNT_ID` and `B2_APPLICATION_KEY` which you got from b2. You can put these in your ~/.bashrc
```shell
export B2_ACCOUNT_ID="XXXX"
export B2_APPLICATION_KEY="XXXX"
```

## Dependencies

b2 uses:

* [cobra](https://github.com/spf13/cobra)
* [viper](https://github.com/spf13/viper)
* [mpb](https://github.com/vbauerster/mpb)
* [color](https://github.com/fatih/color)
