package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type SysProject struct {
	Id          int64  `orm:"column(id);auto" description:"主键ID"`
	ProjectName string `orm:"column(project_name);size(32)" description:"项目名称"`
	Detail      string `orm:"column(detail);size(255)" description:"备注"`
	CreateId    int64  `orm:"column(create_id)" description:"创建者ID"`
	UpdateId    int64  `orm:"column(update_id)" description:"修改这ID"`
	Status      uint8  `orm:"column(status)" description:"状态1-正常，0-删除"`
	CreateTime  int64  `orm:"column(create_time)" description:"添加时间"`
	UpdateTime  int64  `orm:"column(update_time)" description:"修改时间"`
}

func (t *SysProject) TableName() string {
	return "sys_project"
}

func init() {
	orm.RegisterModel(new(SysProject))
}

// AddSysProject insert a new SysProject into database and returns
// last inserted Id on success.
func AddSysProject(m *SysProject) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetSysProjectById retrieves SysProject by Id. Returns error if
// Id doesn't exist
func GetSysProjectById(id int64) (v *SysProject, err error) {
	o := orm.NewOrm()
	v = &SysProject{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllSysProject retrieves all SysProject matches certain condition. Returns empty list if
// no records exist
func GetAllSysProject(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(SysProject))
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

	var l []SysProject
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

// UpdateSysProject updates SysProject by Id and returns error if
// the record to be updated doesn't exist
func UpdateSysProjectById(m *SysProject) (err error) {
	o := orm.NewOrm()
	v := SysProject{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSysProject deletes SysProject by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSysProject(id int64) (err error) {
	o := orm.NewOrm()
	v := SysProject{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SysProject{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
