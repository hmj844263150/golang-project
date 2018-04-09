package dal

import (
	"bytes"
	"context"
	"espressif.com/chip/factory/db"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var UserTid = 1
var _ db.Doer = (*User)(nil)
var usercols = []db.Col{1, 2, 3, 5, 6, 7, 8, 9, 10, 11, 12}
var userfields = []string{"id", "created", "updated", "visibly", "account", "password", "name", "factory_sid", "group_id", "email", "description"}

var UserCol = struct {
	Id, Created, Updated, Visibly, Account, Password, Name, FactorySid, GroupId, Email, Description, _ db.Col
}{1, 2, 3, 5, 6, 7, 8, 9, 10, 11, 12, 0}

type User struct {
	Id          int
	Created     time.Time
	Updated     time.Time
	Visibly     bool
	Account     string
	Password    string
	Name        string
	FactorySid  string
	GroupId     int
	Email       string
	Description string

	// ext, not persistent field
	ext      *Ext
	paddings map[string]interface{}
}

func NewUser(ctx context.Context) *User {
	now := time.Now()
	u := &User{Created: now, Updated: now, Visibly: true}
	u.ext = GetExtFromContext(ctx)
	defaultUser(ctx, u)
	return u
}

func FindUser(ctx context.Context, id int) *User {
	dos, err := db.Open("User").Query(newUserDest, true, userSqls[6], id)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		user, _ := do.(*User)
		if ext != nil {
			user.ext = ext
		}
		return user
	}
	return nil
}

func ListUser(ctx context.Context, ids ...int) []*User {
	holders := make([]string, len(ids))
	generic := make([]interface{}, len(ids))
	for ii, id := range ids {
		holders[ii] = "?"
		generic[ii] = id
	}
	sql := fmt.Sprintf(userSqls[7], strings.Join(holders, ", "))
	dos, err := db.Open("User").Query(newUserDest, true, sql, generic...)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	users := make([]*User, len(dos))
	for ii, do := range dos {
		user, _ := do.(*User)
		if ext != nil {
			user.ext = ext
		}
		users[ii] = user
	}
	return users
}

func ListUserAll(ctx context.Context, offset int, rowCount int) []*User {
	dos, err := db.Open("User").Query(newUserDest, true, userSqls[8], offset, rowCount)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	users := make([]*User, len(dos))
	for ii, do := range dos {
		user, _ := do.(*User)
		if ext != nil {
			user.ext = ext
		}
		users[ii] = user
	}
	return users
}

func FindUserByAccount(ctx context.Context, account string) *User {
	dos, err := db.Open("User").Query(newUserDest, true, userSqls[9], account)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		user, _ := do.(*User)
		if ext != nil {
			user.ext = ext
		}
		return user
	}
	return nil
}

func ListUserByGroupId(ctx context.Context, groupId int, offset int, rowCount int) []*User {
	dos, err := db.Open("User").Query(newUserDest, true, userSqls[10], groupId, offset, rowCount)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	users := make([]*User, len(dos))
	for ii, do := range dos {
		user, _ := do.(*User)
		if ext != nil {
			user.ext = ext
		}
		users[ii] = user
	}
	return users
}

func ListUserByFactorySid(ctx context.Context, factorySid string, offset int, rowCount int) []*User {
	dos, err := db.Open("User").Query(newUserDest, true, userSqls[11], factorySid, offset, rowCount)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	users := make([]*User, len(dos))
	for ii, do := range dos {
		user, _ := do.(*User)
		if ext != nil {
			user.ext = ext
		}
		users[ii] = user
	}
	return users
}

func (u *User) Save() error {
	now := time.Now()
	u.Created, u.Updated, u.Visibly = now, now, true
	var id int64
	var err error
	if u.Id == 0 {
		id, _, err = db.Open("User").Exec(userSqls[0], u.Account, u.Password, u.Name, u.FactorySid, u.GroupId, u.Email, u.Description)
	} else {
		id, _, err = db.Open("User").Exec(userSqls[1], u.Id, u.Account, u.Password, u.Name, u.FactorySid, u.GroupId, u.Email, u.Description)
	}
	if err != nil {
		return err
	}
	u.Id = int(id)
	return nil
}

func (u *User) Update(cs ...db.Col) error {
	if u.Id == 0 {
		return logError("dal.User Error: can not update row while id is zero")
	}
	u.Updated = time.Now()
	if len(cs) == 0 {
		_, _, err := db.Open("User").Exec(userSqls[2], u.Visibly, u.Account, u.Password, u.Name, u.FactorySid, u.GroupId, u.Email, u.Description, u.Id)
		return err
	}
	cols, args, err := colsAndArgsUser(u, cs...)
	if err != nil {
		return err
	}
	args = append(args, u.Id)
	sqlstr := fmt.Sprintf(userSqls[3], strings.Join(cols, ", "))
	_, _, err = db.Open("User").Exec(sqlstr, args...)
	return err
}

func (u *User) Invisibly() error {
	if u.Id == 0 {
		return logError("dal.User Error: can not invisibly row while id is zero")
	}
	u.Updated = time.Now()
	u.Visibly = false
	_, _, err := db.Open("User").Exec(userSqls[4], u.Id)
	return err
}

func (u *User) Delete() error {
	if u.Id == 0 {
		return logError("dal.User Error: can not delete row while id is zero")
	}
	u.Updated = time.Now()
	_, _, err := db.Open("User").Exec(userSqls[5], u.Id)
	return err
}

func (u *User) Valid() error {
	return u.valid()
}

func (u *User) SetExt(ext *Ext) {
	u.ext = ext
}

func (u *User) Padding(pkey string, pvalue interface{}) {
	if u.ext == nil {
		u.ext = &Ext{Loc: DefaultLoc}
	}
	if u.paddings == nil {
		u.paddings = make(map[string]interface{})
	}
	u.paddings[pkey] = pvalue
	u.ext.IsComplex = true
}

func (u *User) AsMap(isColumnName bool, cs ...db.Col) map[string]interface{} {
	mm := make(map[string]interface{})
	for _, cc := range cs {
		switch cc {
		case UserCol.Id:
			if isColumnName {
				mm["id"] = u.Id
			} else {
				mm["Id"] = u.Id
			}
		case UserCol.Created:
			if isColumnName {
				mm["created"] = u.Created
			} else {
				mm["Created"] = u.Created
			}
		case UserCol.Updated:
			if isColumnName {
				mm["updated"] = u.Updated
			} else {
				mm["Updated"] = u.Updated
			}
		case UserCol.Visibly:
			if isColumnName {
				mm["visibly"] = u.Visibly
			} else {
				mm["Visibly"] = u.Visibly
			}
		case UserCol.Account:
			if isColumnName {
				mm["account"] = u.Account
			} else {
				mm["Account"] = u.Account
			}
		case UserCol.Password:
			if isColumnName {
				mm["password"] = u.Password
			} else {
				mm["Password"] = u.Password
			}
		case UserCol.Name:
			if isColumnName {
				mm["name"] = u.Name
			} else {
				mm["Name"] = u.Name
			}
		case UserCol.FactorySid:
			if isColumnName {
				mm["factory_sid"] = u.FactorySid
			} else {
				mm["FactorySid"] = u.FactorySid
			}
		case UserCol.GroupId:
			if isColumnName {
				mm["group_id"] = u.GroupId
			} else {
				mm["GroupId"] = u.GroupId
			}
		case UserCol.Email:
			if isColumnName {
				mm["email"] = u.Email
			} else {
				mm["Email"] = u.Email
			}
		case UserCol.Description:
			if isColumnName {
				mm["description"] = u.Description
			} else {
				mm["Description"] = u.Description
			}
		default:
			logError(fmt.Sprintf("dal.User Error: unknow column num %d in talbe user", cc))
		}
	}
	return mm
}

func (u *User) MarshalJSON() ([]byte, error) {
	if u == nil {
		return []byte("null"), nil
	}
	loc := DefaultLoc
	var numericEnum bool
	if u.ext != nil {
		if u.ext.IsComplex {
			return u.marshalJSONComplex()
		}
		loc = u.ext.Loc
		numericEnum = u.ext.NumericEnum
	}
	var buf bytes.Buffer
	buf.WriteString(`{"id":`)
	buf.WriteString(strconv.FormatInt(int64(u.Id), 10))
	buf.WriteString(`, "created":`)
	if loc == Epoch {
		buf.WriteString(strconv.FormatInt(u.Created.Unix(), 10))
	} else {
		u.Created = u.Created.In(loc)
		buf.WriteString(`"` + u.Created.Format(DefaultTimeFormat) + `"`)
	}
	buf.WriteString(`, "updated":`)
	if loc == Epoch {
		buf.WriteString(strconv.FormatInt(u.Updated.Unix(), 10))
	} else {
		u.Updated = u.Updated.In(loc)
		buf.WriteString(`"` + u.Updated.Format(DefaultTimeFormat) + `"`)
	}
	buf.WriteString(`, "visibly":`)
	if numericEnum {
		if u.Visibly {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
	} else {
		buf.WriteString(strconv.FormatBool(u.Visibly))
	}
	buf.WriteString(`, "account":`)
	WriteJsonString(&buf, u.Account)
	buf.WriteString(`, "password":`)
	WriteJsonString(&buf, u.Password)
	buf.WriteString(`, "name":`)
	WriteJsonString(&buf, u.Name)
	buf.WriteString(`, "factory_sid":`)
	WriteJsonString(&buf, u.FactorySid)
	buf.WriteString(`, "group_id":`)
	buf.WriteString(strconv.FormatInt(int64(u.GroupId), 10))
	buf.WriteString(`, "email":`)
	WriteJsonString(&buf, u.Email)
	buf.WriteString(`, "description":`)
	WriteJsonString(&buf, u.Description)
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (u *User) marshalJSONComplex() ([]byte, error) {
	if u == nil {
		return []byte("null"), nil
	}
	if u.ext == nil {
		return nil, logError("dal.User Error: can not marshalJSONComplex with .ext == nil")
	}
	loc := u.ext.Loc
	numericEnum := u.ext.NumericEnum
	var includes, excludes map[db.Col]interface{}
	if vv, ok := dalVerboses[UserTid]; ok {
		if vvv, ok := vv[u.ext.Verbose]; ok {
			includes, excludes = vvv[0], vvv[1]
		}
	}
	paddings := u.paddings
	var buf bytes.Buffer
	var isRender bool
	isRender = isRenderField(UserCol.Id, "id", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "id":`)
		buf.WriteString(strconv.FormatInt(int64(u.Id), 10))
	}
	isRender = isRenderField(UserCol.Created, "created", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "created":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(u.Created.Unix(), 10))
		} else {
			u.Created = u.Created.In(loc)
			buf.WriteString(`"` + u.Created.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(UserCol.Updated, "updated", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "updated":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(u.Updated.Unix(), 10))
		} else {
			u.Updated = u.Updated.In(loc)
			buf.WriteString(`"` + u.Updated.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(UserCol.Visibly, "visibly", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "visibly":`)
		if numericEnum {
			if u.Visibly {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}
		} else {
			buf.WriteString(strconv.FormatBool(u.Visibly))
		}
	}
	isRender = isRenderField(UserCol.Account, "account", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "account":`)
		WriteJsonString(&buf, u.Account)
	}
	isRender = isRenderField(UserCol.Password, "password", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "password":`)
		WriteJsonString(&buf, u.Password)
	}
	isRender = isRenderField(UserCol.Name, "name", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "name":`)
		WriteJsonString(&buf, u.Name)
	}
	isRender = isRenderField(UserCol.FactorySid, "factory_sid", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "factory_sid":`)
		WriteJsonString(&buf, u.FactorySid)
	}
	isRender = isRenderField(UserCol.GroupId, "group_id", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "group_id":`)
		buf.WriteString(strconv.FormatInt(int64(u.GroupId), 10))
	}
	isRender = isRenderField(UserCol.Email, "email", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "email":`)
		WriteJsonString(&buf, u.Email)
	}
	isRender = isRenderField(UserCol.Description, "description", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "description":`)
		WriteJsonString(&buf, u.Description)
	}
	if paddings != nil {
		var kk string
		var vv interface{}
		var str string
		var err error
		for kk, vv = range paddings {
			buf.WriteString(`, "` + kk + `":`)
			str, err = General(vv)
			if err != nil {
				return nil, err
			}
			buf.WriteString(str)
		}
	}
	buf.WriteString("}")
	bs := buf.Bytes()
	bs[0] = '{'
	return bs, nil
}

func (u *User) UnmarshalMap(ctx context.Context, vi interface{}, cols ...db.Col) ([]db.Col, error) {
	var err error
	if vi == nil {
		return nil, logError("can not UnmarshalUser with null value")
	}
	vv, ok := vi.(map[string]interface{})
	if !ok {
		return nil, logError("type asserted failed when UnmarshalUser")
	}
	updatedCols := []db.Col{}
	if len(cols) == 0 {
		cols = usercols
	}
	loc := DefaultLoc
	for _, col := range cols {
		switch col {
		case UserCol.Id:
			vvv, ok := vv["id"]
			if !ok {
				continue
			}
			u.Id, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case UserCol.Created:
			vvv, ok := vv["created"]
			if !ok {
				continue
			}
			u.Created, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case UserCol.Updated:
			vvv, ok := vv["updated"]
			if !ok {
				continue
			}
			u.Updated, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case UserCol.Visibly:
			vvv, ok := vv["visibly"]
			if !ok {
				continue
			}
			u.Visibly, err = Bool(vvv)
			updatedCols = append(updatedCols, col)
		case UserCol.Account:
			vvv, ok := vv["account"]
			if !ok {
				continue
			}
			u.Account, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case UserCol.Password:
			vvv, ok := vv["password"]
			if !ok {
				continue
			}
			u.Password, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case UserCol.Name:
			vvv, ok := vv["name"]
			if !ok {
				continue
			}
			u.Name, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case UserCol.FactorySid:
			vvv, ok := vv["factory_sid"]
			if !ok {
				continue
			}
			u.FactorySid, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case UserCol.GroupId:
			vvv, ok := vv["group_id"]
			if !ok {
				continue
			}
			u.GroupId, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case UserCol.Email:
			vvv, ok := vv["email"]
			if !ok {
				continue
			}
			u.Email, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case UserCol.Description:
			vvv, ok := vv["description"]
			if !ok {
				continue
			}
			u.Description, err = String(vvv)
			updatedCols = append(updatedCols, col)
		}
		if err != nil {
			return nil, err
		}
	}
	return cols, nil
}

func UnmarshalUser(ctx context.Context, vi interface{}, cols ...db.Col) (*User, error) {
	u := NewUser(ctx)
	_, err := u.UnmarshalMap(ctx, vi, cols...)
	if err != nil {
		return nil, err
	}
	return u, err
}

func UnmarshalUsers(ctx context.Context, vi interface{}, cols ...db.Col) ([]*User, error) {
	var err error
	if vi == nil {
		return nil, logError("can not UnmarshalUsers with null value")
	}
	vv, ok := vi.([]interface{})
	if !ok {
		return nil, logError("type asserted failed when UnmarshalUsers")
	}
	users := make([]*User, len(vv))
	for ii, vvv := range vv {
		var u *User
		u, err = UnmarshalUser(ctx, vvv, cols...)
		if err != nil {
			return nil, err
		}
		users[ii] = u
	}
	return users, nil
}

func newUserDest(cols ...string) (db.Doer, []interface{}, error) {
	u := &User{}
	if cols == nil || len(cols) == 0 {
		return u, []interface{}{&u.Id, &u.Created, &u.Updated, &u.Visibly, &u.Account, &u.Password, &u.Name, &u.FactorySid, &u.GroupId, &u.Email, &u.Description}, nil
	}
	dest := make([]interface{}, len(cols))
	for ii, col := range cols {
		switch col {
		case "id":
			dest[ii] = &u.Id
		case "created":
			dest[ii] = &u.Created
		case "updated":
			dest[ii] = &u.Updated
		case "visibly":
			dest[ii] = &u.Visibly
		case "account":
			dest[ii] = &u.Account
		case "password":
			dest[ii] = &u.Password
		case "name":
			dest[ii] = &u.Name
		case "factory_sid":
			dest[ii] = &u.FactorySid
		case "group_id":
			dest[ii] = &u.GroupId
		case "email":
			dest[ii] = &u.Email
		case "description":
			dest[ii] = &u.Description
		default:
			return nil, nil, logError("dal.User Error: unknow column " + col + " in talbe user")
		}
	}
	return u, dest, nil
}

func colsAndArgsUser(u *User, cs ...db.Col) ([]string, []interface{}, error) {
	len := len(cs)
	if len == 0 {
		return nil, nil, logError("dal.User Error: at least one column to colsAndArgsUser")
	}
	cols := make([]string, len)
	args := make([]interface{}, len)
	for ii, cc := range cs {
		switch cc {
		case UserCol.Id:
			cols[ii] = "`id` = ?"
			args[ii] = u.Id
		case UserCol.Created:
			cols[ii] = "`created` = ?"
			args[ii] = u.Created
		case UserCol.Updated:
			cols[ii] = "`updated` = ?"
			args[ii] = u.Updated
		case UserCol.Visibly:
			cols[ii] = "`visibly` = ?"
			args[ii] = u.Visibly
		case UserCol.Account:
			cols[ii] = "`account` = ?"
			args[ii] = u.Account
		case UserCol.Password:
			cols[ii] = "`password` = ?"
			args[ii] = u.Password
		case UserCol.Name:
			cols[ii] = "`name` = ?"
			args[ii] = u.Name
		case UserCol.FactorySid:
			cols[ii] = "`factory_sid` = ?"
			args[ii] = u.FactorySid
		case UserCol.GroupId:
			cols[ii] = "`group_id` = ?"
			args[ii] = u.GroupId
		case UserCol.Email:
			cols[ii] = "`email` = ?"
			args[ii] = u.Email
		case UserCol.Description:
			cols[ii] = "`description` = ?"
			args[ii] = u.Description
		default:
			return nil, nil, logError(fmt.Sprintf("dal.User Error: unknow column num %d in talbe user", cc))
		}
	}
	return cols, args, nil
}

var UserEnum = struct {
}{}

var userSqls = []string{
	/*
		CREATE TABLE `user` (
		  `id` int(11) NOT NULL AUTO_INCREMENT,
		  `created` datetime NOT NULL,
		  `updated` datetime NOT NULL,
		  `visibly` bool NOT NULL,
		  `account` varchar(64) NOT NULL,
		  `password` varchar(64) NOT NULL,
		  `name` varchar(64) NOT NULL,
		  `factory_sid` varchar(64) NOT NULL,
		  `group_id` int(11) NOT NULL,
		  `email` varchar(128) NOT NULL,
		  `description` varchar(256) NOT NULL,
		  PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;

	*/
	/*0*/ "insert into user(`created`, `updated`, `visibly`, `account`, `password`, `name`, `factory_sid`, `group_id`, `email`, `description`) values(now(), now(), 1, ?, ?, ?, ?, ?, ?, ?)",
	/*1*/ "insert into user(`id`, `created`, `updated`, `visibly`, `account`, `password`, `name`, `factory_sid`, `group_id`, `email`, `description`) values(?, now(), now(), 1, ?, ?, ?, ?, ?, ?, ?)",
	/*2*/ "update user set updated = now(), `visibly` = ?, `account` = ?, `password` = ?, `name` = ?, `factory_sid` = ?, `group_id` = ?, `email` = ?, `description` = ? where id = ?",
	/*3*/ "update user set updated = now(), %s where id = ?",
	/*4*/ "update user set visibly = 0, updated = now() where id = ?",
	/*5*/ "delete from user where id = ?",
	/*6*/ "select `id`, `created`, `updated`, `visibly`, `account`, `password`, `name`, `factory_sid`, `group_id`, `email`, `description` from user where id = ? and visibly = 1",
	/*7*/ "select `id`, `created`, `updated`, `visibly`, `account`, `password`, `name`, `factory_sid`, `group_id`, `email`, `description` from user where id in (%s) and visibly = 1",

	/*8*/ "select `id`, `created`, `updated`, `visibly`, `account`, `password`, `name`, `factory_sid`, `group_id`, `email`, `description` from user where visibly = 1 order by id desc limit ?, ?",
	/*9*/ "select `id`, `created`, `updated`, `visibly`, `account`, `password`, `name`, `factory_sid`, `group_id`, `email`, `description` from user where visibly = 1 and account = ? limit 0, 1",
	/*10*/ "select `id`, `created`, `updated`, `visibly`, `account`, `password`, `name`, `factory_sid`, `group_id`, `email`, `description` from user where visibly = 1 and group_id = ? limit ?, ?",
	/*11*/ "select `id`, `created`, `updated`, `visibly`, `account`, `password`, `name`, `factory_sid`, `group_id`, `email`, `description` from user where visibly = 1 and factory_sid = ? limit ?, ?",
}
