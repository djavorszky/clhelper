#!/bin/bash

cp $1 $1.bak
cp $1.go $1.go.bak

export GOOS="darwin"
echo "buliding mac version"
go build $1.go
mv $1 $1.mac

export GOOS="windows"
echo "buliding windows version"
go build $1.go

export GOOS="linux"
echo "buliding linux version"
go build $1.go
mv $1 $1.linux
