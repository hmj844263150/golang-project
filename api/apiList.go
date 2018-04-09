package api

import (
	"espressif.com/chip/factory/factory"
	//"espressif.com/chip/factory/module"
	"espressif.com/chip/factory/batch"
	"espressif.com/chip/factory/hengha"
	"espressif.com/chip/factory/mpn"
	"espressif.com/chip/factory/testdata"
	"espressif.com/chip/factory/user"
)

var auth = hengha.Auth
var authStaff = hengha.AuthStaff

var apiList = []*Api{

	// static TODO
	// {"/static/filepath:*", hengha.Static, nil},

	// db
	{"/hengha/db/create", hengha.Createdb, authStaff},
	{"/hengha/db/drop", hengha.Dropdb, authStaff},

	// factory
	{"/factory", factory.Factory, authStaff},
	{"/factorys", factory.Factorys, auth},

	// module

	// testdata
	{"/testdata", testdata.Testdata, nil},
	{"/testdata/print", testdata.Print, nil},
	{"/testdata/logs", testdata.Logs, nil},
	{"/testdata/stats", testdata.Stats, nil},
	{"/testdata/dump", testdata.Dump, nil},

	// batch
	{"/batch", batch.Batch, auth},
	{"/batch/file", batch.File, authStaff},
	{"/batch/stats", batch.Stats, auth},
	{"/batchs", batch.Batchs, auth},

	// mpn
	{"/mpn", mpn.Mpn, nil},

	// user
	{"/user", user.User, authStaff},
	{"/user/login", user.Login, nil},
	{"/user/modify", user.Modify, authStaff},
	{"/users", user.Users, auth},
}
