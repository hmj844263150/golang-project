package rpc

import (
	"errors"
)

type Version uint8
type Method uint8
type Mode uint8
type Proto uint8
type Format uint8

const (
	V1 Version = 1
	V2 Version = 2
	V3 Version = 3

	Options Method = 1
	Get     Method = 2
	Head    Method = 3
	Post    Method = 4
	Put     Method = 5
	Delete  Method = 6
	Trace   Method = 7
	Connect Method = 8
	Patch   Method = 9

	None Mode = 0
	V    Mode = 1
	Vv   Mode = 2
	Vvv  Mode = 3

	Http  Proto = 1
	Http2 Proto = 2
	Json  Proto = 3
	Ws    Proto = 4
	Mqtt  Proto = 5

	Plain Format = 1
	Msgp  Format = 2
	Pb    Format = 3
)

var nilMap = errors.New("map is nil")

var methodMap map[string]Method
var modeMap map[string]Mode

func init() {
	methodMap = map[string]Method{"": Get, "Options": Options, "Get": Get, "Head": Head, "Post": Post, "Put": Put, "Delete": Delete, "Trace": Trace, "Connect": Connect, "Patch": Patch, "OPTIONS": Options, "GET": Get, "HEAD": Head, "POST": Post, "PUT": Put, "DELETE": Delete, "TRACE": Trace, "CONNECT": Connect, "PATCH": Patch}
	modeMap = map[string]Mode{"": None, "v": V, "vv": Vv, "vvv": Vvv, "V": V, "Vv": Vv, "Vvv": Vvv}
}
