package rpc

import (
	"errors"
)

var Skip = errors.New("skip")

var BadRequest = errors.New("400 Bad Request")
var Forbidden = errors.New("403 Forbidden")
var NotFound = errors.New("404 Not Found")
var MethodNotAllowed = errors.New("405 Method Not Allowed")
var Conflict = errors.New("409 Conflict")

var MustDeliver = errors.New("must deliver request")
var MustHaveAction = errors.New("must have action")
var ActionUnknow = errors.New("action unknow")

var DispatchFailed = errors.New("dispatch failed")
