package SubjectMation

import (
		"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"xiangmu/Student/DataConn"
	"xiangmu/Student/StructureType"
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
	sub_name := r.Form["Sub_name"][0]
	sub_type := r.Form["Sub_type"][0]
	sub_times := r.Form["Sub_times"][0]

	if sub_name == "" || sub_type == "" || sub_times == "" {
		s := StructureType.Things{"请将信息输入完整"}
		render.JSON(w, r, s)
		return
	}

	//判断课程是否已经存在
	rows, err := subject.db.Table("Subjects").Where("Sub_name=?", sub_name).Select("Sub_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		var sub_id int
		err = rows.Scan(&sub_id)
		if sub_id != 0 {
			return
		}
	}
	sub_time, err := strconv.Atoi(sub_times)
	if err != nil {
		return
	}
	err = subject.db.Create(&DataConn.Subject{Sub_name: sub_name, Sub_type: sub_type, Sub_times: sub_time}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"课程信息添加成功"}
	render.JSON(w, r, s)
}

func (subject *SubjectAPi) BrowSubject(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sub_name := r.Form["Sub_name"][0]
	sub_type := r.Form["Sub_type"][0]

	s := StructureType.Subjecttotal{}
	tem := StructureType.Subject{}

	//按课程名称搜索
	if sub_name != "" && sub_type == "" {
		rows, err := subject.db.Table("Subjects").Where("Sub_name=?", sub_name).Select("Sub_id,Sub_name,Sub_type,Sub_times").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Sub_id, &tem.Sub_name, &tem.Sub_type, &tem.Sub_times)
			if err != nil {
				return
			}
			s.SubjectList = append(s.SubjectList, tem)
		}
		render.JSON(w, r, s)
	}
	//按课程类型搜索
	if sub_name == "" && sub_type != "" {
		rows, err := subject.db.Table("Subjects").Where("Sub_type=?", sub_type).Select("Sub_id,Sub_name,Sub_type,Sub_times").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Sub_id, &tem.Sub_name, &tem.Sub_type, &tem.Sub_times)
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
	sub_id := r.Form["Sub_id"][0]
	sub_name := r.Form["Sub_name"][0]
	sub_type := r.Form["Sub_type"][0]
	sub_times := r.Form["Sub_times"][0]

	if sub_name != "" {
		err := subject.db.Model(&DataConn.Subject{}).Where("Sub_id=?", sub_id).Update(&DataConn.Subject{Sub_name: sub_name}).Error
		if err != nil {
			return
		}
	}

	if sub_type != "" {
		err := subject.db.Model(&DataConn.Subject{}).Where("Sub_id=?", sub_id).Update(&DataConn.Subject{Sub_type: sub_type}).Error
		if err != nil {
			return
		}
	}

	if sub_times != "" {
		sub_time, err := strconv.Atoi(sub_times)
		if err != nil {
			return
		}
		err = subject.db.Model(&DataConn.Subject{}).Where("Sub_id=?", sub_id).Update(&DataConn.Subject{Sub_times: sub_time}).Error
		if err != nil {
			return
		}
	}
	s := StructureType.Things{"更新课程信息成功"}
	render.JSON(w, r, s)
}

func (subject *SubjectAPi) SubSubject(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sub_id := r.Form["Sub_id"][0]

	err := subject.db.Model(&DataConn.Subject{}).Where("Sub_id=?", sub_id).Delete(&DataConn.Subject{}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"删除课程信息成功"}
	render.JSON(w, r, s)
}
