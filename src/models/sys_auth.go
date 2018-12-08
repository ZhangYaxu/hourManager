/**
* @Project: hourManager
* @Package models
* @Description: 权限管理
* @author : wj
* @date Date : 2018/12/06/ 16:00
* @version V1.0
 */
package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type SysAuth struct {
	Id         int64  `orm:"column(id);auto" description:"自增ID"`
	Pid        int64  `orm:"column(pid)" description:"上级ID，0为顶级"`
	AuthName   string `orm:"column(auth_name);size(64)" description:"权限名称"`
	AuthUrl    string `orm:"column(auth_url);size(255)" description:"URL地址"`
	Sort       int    `orm:"column(sort)" description:"排序，越小越前"`
	Icon       string `orm:"column(icon);size(255)"`
	IsShow     int    `orm:"column(is_show)" description:"是否显示，0-隐藏，1-显示"`
	UserId     int64  `orm:"column(user_id)" description:"操作者ID"`
	CreateId   int64  `orm:"column(create_id)" description:"创建者ID"`
	UpdateId   int64  `orm:"column(update_id)" description:"修改者ID"`
	Status     int    `orm:"column(status)" description:"状态，1-正常，0-删除"`
	CreateTime int64  `orm:"column(create_time)" description:"创建时间"`
	UpdateTime int64  `orm:"column(update_time)" description:"更新时间"`
}

func (t *SysAuth) TableName() string {
	return "sys_auth"
}

func init() {
	orm.RegisterModel(new(SysAuth))
}

// AddSysAuth insert a new SysAuth into database and returns
// last inserted Id on success.
func AddSysAuth(m *SysAuth) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetSysAuthById retrieves SysAuth by Id. Returns error if
// Id doesn't exist
func GetSysAuthById(id int64) (v *SysAuth, err error) {
	o := orm.NewOrm()
	v = &SysAuth{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetSysAuthList(page, pageSize int, filters ...interface{}) ([]*SysAuth, int64) {
	offset := (page - 1) * pageSize
	list := make([]*SysAuth, 0)
	query := orm.NewOrm().QueryTable("sys_auth")
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	query.OrderBy("pid", "sort").Limit(pageSize, offset).All(&list)

	return list, total
}

// GetAllSysAuth retrieves all SysAuth matches certain condition. Returns empty list if
// no records exist
func GetAllSysAuth(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(SysAuth))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []SysAuth
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateSysAuth updates SysAuth by Id and returns error if
// the record to be updated doesn't exist
func UpdateSysAuthById(m *SysAuth) (err error) {
	o := orm.NewOrm()
	v := SysAuth{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func (a *SysAuth) UpdateSysAuth(fields ...string) error {
	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
		return err
	}
	return nil
}

// DeleteSysAuth deletes SysAuth by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSysAuth(id int64) (err error) {
	o := orm.NewOrm()
	v := SysAuth{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SysAuth{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
