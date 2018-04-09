package rpc

import (
	"espressif.com/chip/factory/db"
)

func MultiPut(req *Request, name string, doerFunc func(req *Request, x map[string]interface{}) (db.Doer, error), cols ...db.Col) ([]db.Doer, error) {
	xs, err := req.GetSliceMap(name, req.Body)
	if err != nil {
		return nil, err
	}
	doers := make([]db.Doer, len(xs))
	for i, x := range xs {
		doer, err := doerFunc(req, x)
		if err != nil {
			return nil, err
		}
		updatedCols, err := doer.UnmarshalMap(req.Ctx, x, cols...)
		if err != nil {
			return nil, err
		}
		err = doer.Update(updatedCols...)
		if err != nil {
			return nil, err
		}
		doers[i] = doer
	}
	return doers, nil
}

func MultiDelete(req *Request, name string, doerFunc func(req *Request, x map[string]interface{}) (db.Doer, error)) ([]db.Doer, error) {
	xs, err := req.GetSliceMap(name, req.Body)
	if err != nil {
		return nil, Forbidden
	}
	doers := make([]db.Doer, len(xs))
	for i, x := range xs {
		doer, err := doerFunc(req, x)
		if err != nil {
			return nil, err
		}
		if doer != nil {
			doer.Invisibly()
		}
		doers[i] = doer
	}
	return doers, err
}
