#!/bin/bash

# generate memory management for float64
go run *.go -pkg="mms" | gofmt > ../mms_float64.go
