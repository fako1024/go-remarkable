# Tooling for the ReMarkable (2) Paper Tablet (in Go)

[![Github Release](https://img.shields.io/github/release/fako1024/go-remarkable.svg)](https://github.com/fako1024/go-remarkable/releases)
[![GoDoc](https://godoc.org/github.com/fako1024/go-remarkable?status.svg)](https://godoc.org/github.com/fako1024/go-remarkable/)
[![Go Report Card](https://goreportcard.com/badge/github.com/fako1024/go-remarkable)](https://goreportcard.com/report/github.com/fako1024/go-remarkable)
[![Build/Test Status](https://github.com/fako1024/go-remarkable/workflows/Go/badge.svg)](https://github.com/fako1024/go-remarkable/actions?query=workflow%3AGo)

This package allows to interact with a ReMarkable (2) device on various levels, trying to provide all functionality currently available in various projects within one package. A

**NOTE: This package is currently work in progress. Interfaces and implementation are subject to change.**

## Features
- Access to framebuffer / screen data
  - Screenshot
  - Low latency, live stream (via client application) with low power consumption (using input detection)
  - Low-overhead "broadcast" functionality supporting multiple clients at the same time

## Installation
```bash
go get -u github.com/fako1024/go-remarkable
cd cmd/rm-agent && GOOS=linux GOARCH=arm GOARM=7 go build
```
