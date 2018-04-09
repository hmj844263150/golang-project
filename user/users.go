package user

import (
	"espressif.com/chip/factory/dal"
	"espressif.com/chip/factory/rpc"
)

func Users(req *rpc.Request, resp *rpc.Response) {
	switch req.Method {
	case rpc.Get:
		usersGet(req, resp)
	default:
		resp.Err = rpc.MethodNotAllowed
	}
}

func usersGet(req *rpc.Request, resp *rpc.Response) {
	users := listUser(req)
	resp.Body["users"] = users
}

func listUser(req *rpc.Request) []*dal.User {
	if req.Factory != nil {
		if req.Factory.IsStaff {
			return dal.ListUserAll(req.Ctx, 0, 1000)
		}
		return dal.ListUserByFactorySid(req.Ctx, req.Factory.Sid, 0, 1000)
	}
	return nil
}
