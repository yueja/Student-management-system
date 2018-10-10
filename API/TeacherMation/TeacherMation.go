package TeacherMation

import (
		"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
	"xiangmu/Student/DataConn"
	"xiangmu/Student/StructureType"
)

type TeacherAPi struct {
	db *gorm.DB
}

func Make_db(db *gorm.DB) *TeacherAPi {
	DB := &TeacherAPi{db}
	return DB
}

func (teacher *TeacherAPi) AddTeacher(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tec_name := r.Form["Tec_name"][0]
	tec_sex := r.Form["Tec_sex"][0]
	tec_birth := r.Form["Tec_birth"][0]
	tec_major := r.Form["Tec_major"][0]
	tec_phone := r.Form["Tec_phone"][0]

	var tec_id, maj_id int
	if tec_name == "" || tec_sex == "" || tec_birth == "" || tec_major == "" || tec_phone == "" {
		s := StructureType.Things{"请将信息输入完整"}
		render.JSON(w, r, s)
		return
	}

	//判断教师信息是否已经存在
	rows, err := teacher.db.Table("Teacherinfos").Where("Tec_name=?", tec_name).Select("Tec_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&tec_id)
		if err != nil {
			return
		}
	}
	if tec_id != 0 {
		s := StructureType.Things{"本教师信息已存在"}
		render.JSON(w, r, s)
		return
	}
	t, _ := time.Parse("2006-01-02", tec_birth) //字符串转时间戳
	//查看该专业是否存在
	rows, err = teacher.db.Table("Majors").Where("Maj_name=?", tec_major).Select("Maj_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&maj_id)
		if err != nil {
			return
		}
	}
	if maj_id == 0 {
		s := StructureType.Things{"该专业不存在，请重新输入"}
		render.JSON(w, r, s)
		return
	}
	err = teacher.db.Create(&DataConn.Teacherinfo{Tec_name: tec_name, Tec_sex: tec_sex, Tec_birth: t, Tec_major: tec_major, Tec_phone: tec_phone}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"添加教师信息成功"}
	render.JSON(w, r, s)
}

func (teacher *TeacherAPi) BrowTeacher(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tec_major := r.Form["Tec_major"][0]
	tec_name := r.Form["Tec_name"][0]

	s := StructureType.Teachertotal{}
	tem:=StructureType.Teacherinfo{}

	//按照专业查找
	a:="Tec_id,Ope_id,Tec_name,Tec_sex,Tec_birth,Tec_major,Tec_phone"
	if tec_major != "" {
		rows, err := teacher.db.Table("Teacherinfos").Where("Tec_major=?", tec_major).Select(a).Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Tec_id, &tem.Ope_id, &tem.Tec_name, &tem.Tec_sex, &tem.Tec_birth, &tem.Tec_major, &tem.Tec_phone)
			if err != nil {
				return
			}
			s.TeacherList = append(s.TeacherList,tem)
		}
		render.JSON(w, r, s)
	}

	//按照个人查询
	if tec_name != "" {
		rows, err := teacher.db.Table("Teacherinfos").Where("Tec_name=?", tec_name).Select(a).Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Tec_id, &tem.Ope_id, &tem.Tec_name, &tem.Tec_sex, &tem.Tec_birth, &tem.Tec_major, &tem.Tec_phone)
			if err != nil {
				return
			}
			s.TeacherList = append(s.TeacherList,tem)
		}
		render.JSON(w, r, s)
	}

	//查询全部教师
	if tec_major == "" && tec_name == "" {
		rows,err := teacher.db.Table("Teacherinfos").Select(a).Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Tec_id, &tem.Ope_id, &tem.Tec_name, &tem.Tec_sex, &tem.Tec_birth, &tem.Tec_major, &tem.Tec_phone)
			if err != nil {
				return
			}
			s.TeacherList = append(s.TeacherList,tem)
		}
		render.JSON(w, r, s)
	}
}

func (teacher *TeacherAPi) UpTeacher(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tec_id := r.Form["Tec_id"][0]
	tec_name := r.Form["Tec_name"][0]
	tec_sex := r.Form["Tec_sex"][0]
	tec_birth := r.Form["Tec_birth"][0]
	tec_major := r.Form["Tec_major"][0]
	tec_phone := r.Form["Tec_phone"][0]

	if tec_name != "" {
		err := teacher.db.Model(&DataConn.Teacherinfo{}).Where("Tec_id=?", tec_id).Update(&DataConn.Teacherinfo{Tec_name: tec_name}).Error
		if err != nil {
			return
		}
	}

	if tec_sex != "" {
		err := teacher.db.Model(&DataConn.Teacherinfo{}).Where("Tec_id=?", tec_id).Update(&DataConn.Teacherinfo{Tec_sex: tec_sex}).Error
		if err != nil {
			return
		}
	}

	if tec_birth != "" {
		t, _ := time.Parse("2006-01-02", tec_birth) //字符串转时间戳
		err := teacher.db.Model(&DataConn.Teacherinfo{}).Where("Tec_id=?", tec_id).Update(&DataConn.Teacherinfo{Tec_birth: t}).Error
		if err != nil {
			return
		}
	}

	if tec_major != "" {
		err := teacher.db.Model(&DataConn.Teacherinfo{}).Where("Tec_id=?", tec_id).Update(&DataConn.Teacherinfo{Tec_major: tec_major}).Error
		if err != nil {
			return
		}
	}

	if tec_phone != "" {
		err := teacher.db.Model(&DataConn.Teacherinfo{}).Where("Tec_id=?", tec_id).Update(&DataConn.Teacherinfo{Tec_phone: tec_phone}).Error
		if err != nil {
			return
		}
	}
	s := StructureType.Things{"更新教师信息成功"}
	render.JSON(w, r, s)
}

func (teacher *TeacherAPi) SubTeacher(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tec_id := r.Form["Tec_id"][0]

	err := teacher.db.Model(&DataConn.Teacherinfo{}).Where("Tec_id=?", tec_id).Delete(&DataConn.Teacherinfo{}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"删除教师信息成功"}
	render.JSON(w, r, s)
}
