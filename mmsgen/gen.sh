#!/bin/bash

# generate memory management for float64
go run *.go 					\
	-pkg="mms"					\
	-type="[]float64"			\
	-new="make([]float64,size)"	\
	-name="Float64sCache"		\
	| gofmt > ../mms_float64.go

# generate memory management for float32
go run *.go 					\
	-pkg="mms"					\
	-type="[]float32"			\
	-new="make([]float32,size)"	\
	-name="Float32sCache"		\
	| gofmt > ../mms_float32.go

# generate memory management for int
go run *.go 					\
	-pkg="mms"					\
	-type="[]int"				\
	-new="make([]int,size)"		\
	-name="IntsCache"			\
	| gofmt > ../mms_int.go

# generate memory management for int64
go run *.go 					\
	-pkg="mms"					\
	-type="[]int64"				\
	-new="make([]int64,size)"	\
	-name="Int64sCache"			\
	| gofmt > ../mms_int64.go

# generate memory management for int32
go run *.go 					\
	-pkg="mms"					\
	-type="[]int32"				\
	-new="make([]int32,size)"	\
	-name="Int32sCache"			\
	| gofmt > ../mms_int32.go
