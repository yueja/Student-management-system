package subject

import (
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"xiangmu/Student/data_conn"
	"xiangmu/Student/structure_type"
)

type SubjectAPi struct {
	db *gorm.DB
}

func Make_db(db *gorm.DB) *SubjectAPi {
	DB := &SubjectAPi{db}
	return DB
}
func (subject *SubjectAPi) AddSubject(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	subName := r.Form["subName"][0]
	subType := r.Form["subType"][0]
	subTimes := r.Form["subTimes"][0]
	var subId int

	if subName == "" || subType == "" || subTimes == "" {
		s := structure_type.Things{"请将信息输入完整", false}
		render.JSON(w, r, s)
		return
	}

	//判断课程是否已经存在
	rows, err := subject.db.Model(&data_conn.Subject{}).Where("SubName=?", subName).Select("SubId").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&subId)
		if err != nil {
			return
		}
	}
	if subId != 0 {
		s := structure_type.Things{"该课程已经存在", false}
		render.JSON(w, r, s)
		return
	}
	subTime, err := strconv.Atoi(subTimes)
	if err != nil {
		return
	}
	err = subject.db.Create(&data_conn.Subject{SubName: subName, SubType: subType, SubTimes: subTime}).Error
	if err != nil {
		return
	}
	s := structure_type.Things{"课程信息添加成功", true}
	render.JSON(w, r, s)
}

func (subject *SubjectAPi) BrowSubject(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	subName := r.Form["subName"][0]
	subType := r.Form["subType"][0]

	s := structure_type.SubjectTotal{}
	tem := structure_type.Subject{}

	//按课程名称搜索
	if subName != "" && subType == "" {
		rows, err := subject.db.Model(&data_conn.Subject{}).Where("SubName=?", subName).Select("SubId,SubName,SubType,SubTimes").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.SubId, &tem.SubName, &tem.SubType, &tem.SubTimes)
			if err != nil {
				return
			}
			s.SubjectList = append(s.SubjectList, tem)
		}
		render.JSON(w, r, s)
	}
	//按课程类型搜索
	if subName == "" && subType != "" {
		rows, err := subject.db.Model(&data_conn.Subject{}).Where("SubType=?", subType).Select("SubId,SubName,SubType,SubTimes").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.SubId, &tem.SubName, &tem.SubType, &tem.SubTimes)
			if err != nil {
				return
			}
			s.SubjectList = append(s.SubjectList, tem)
		}
		render.JSON(w, r, s)
	}
}

func (subject *SubjectAPi) UpSubject(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	subId := r.Form["subId"][0]
	subName := r.Form["subName"][0]
	subType := r.Form["subType"][0]
	subTimes := r.Form["subTimes"][0]

	if subName != "" {
		err := subject.db.Model(&data_conn.Subject{}).Where("SubId=?", subId).Update(&data_conn.Subject{SubName: subName}).Error
		if err != nil {
			return
		}
	}

	if subType != "" {
		err := subject.db.Model(&data_conn.Subject{}).Where("SubId=?", subId).Update(&data_conn.Subject{SubType: subType}).Error
		if err != nil {
			return
		}
	}

	if subTimes != "" {
		subTime, err := strconv.Atoi(subTimes)
		if err != nil {
			return
		}
		err = subject.db.Model(&data_conn.Subject{}).Where("SubId=?", subId).Update(&data_conn.Subject{SubTimes: subTime}).Error
		if err != nil {
			return
		}
	}
	s := structure_type.Things{"更新课程信息成功", true}
	render.JSON(w, r, s)
}

func (subject *SubjectAPi) DelSubject(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	subId := r.Form["subId"][0]

	err := subject.db.Model(&data_conn.Subject{}).Where("SubId=?", subId).Delete(&data_conn.Subject{}).Error
	if err != nil {
		return
	}
	s := structure_type.Things{"删除课程信息成功", true}
	render.JSON(w, r, s)
}
