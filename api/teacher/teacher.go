package teacher

import (
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"time"
	"xiangmu/Student/data_conn"
	"xiangmu/Student/structure_type"
)

type TeacherAPi struct {
	db *gorm.DB
}

func MakeDb(db *gorm.DB) *TeacherAPi {
	DB := &TeacherAPi{db}
	return DB
}

func (teacher *TeacherAPi) AddTeacher(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tecName := r.Form["tecName"][0]
	tecSex := r.Form["tecSex"][0]
	tecBirth := r.Form["tecBirth"][0]
	tecMajor := r.Form["tecMajor"][0]
	tecPhone := r.Form["tecPhone"][0]

	var tecId, majId int
	if tecName == "" || tecSex == "" || tecBirth == "" || tecMajor == "" || tecPhone == "" {
		s := structure_type.Things{"请将信息输入完整", false}
		render.JSON(w, r, s)
		return
	}

	//判断教师信息是否已经存在
	rows, err := teacher.db.Model(&data_conn.TeacherInfo{}).Where("TecName=?", tecName).Select("TecId").Rows()
	if err != nil {
		log.Printf("err:%s", err)
	}
	for rows.Next() {
		err = rows.Scan(&tecId)
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	if tecId != 0 {
		s := structure_type.Things{"本教师信息已存在", false}
		render.JSON(w, r, s)
		return
	}
	t, _ := time.Parse("2006-01-02", tecBirth) //字符串转时间戳
	//查看该专业是否存在
	rows, err = teacher.db.Model(&data_conn.Major{}).Where("MajName=?", tecMajor).Select("MajId").Rows()
	if err != nil {
		log.Printf("err:%s", err)
	}
	for rows.Next() {
		err = rows.Scan(&majId)
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	if majId == 0 {
		s := structure_type.Things{"该专业不存在，请重新输入", false}
		render.JSON(w, r, s)
		return
	}
	err = teacher.db.Create(&data_conn.TeacherInfo{TecName: tecName, TecSex: tecSex, TecBirth: t, TecMajor: tecMajor, TecPhone: tecPhone}).Error
	if err != nil {
		log.Printf("err:%s", err)
	}
	s := structure_type.Things{"添加教师信息成功", true}
	render.JSON(w, r, s)
}

func (teacher *TeacherAPi) BrowTeacher(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tecMajor := r.Form["tecMajor"][0]
	tecName := r.Form["tecName"][0]

	s := structure_type.TeacherTotal{}
	tem := structure_type.TeacherInfo{}

	//按照专业查找
	a := "TecId,TecName,TecSex,TecBirth,TecMajor,TecPhone"
	if tecMajor != "" {
		rows, err := teacher.db.Model(&data_conn.TeacherInfo{}).Where("TecMajor=?", tecMajor).Select(a).Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&tem.TecId, &tem.TecName, &tem.TecSex, &tem.TecBirth, &tem.TecMajor, &tem.TecPhone)
			if err != nil {
				log.Printf("err:%s", err)
			}
			s.TeacherList = append(s.TeacherList, tem)
		}
		s.IsSuccess = true
		render.JSON(w, r, s)
	}

	//按照个人查询
	if tecName != "" {
		rows, err := teacher.db.Model(&data_conn.TeacherInfo{}).Where("TecName=?", tecName).Select(a).Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&tem.TecId, &tem.TecName, &tem.TecSex, &tem.TecBirth, &tem.TecMajor, &tem.TecPhone)
			if err != nil {
				log.Printf("err:%s", err)
			}
			s.TeacherList = append(s.TeacherList, tem)
		}
		s.IsSuccess = true
		render.JSON(w, r, s)
	}

	//查询全部教师
	if tecMajor == "" && tecName == "" {
		rows, err := teacher.db.Model(&data_conn.TeacherInfo{}).Select(a).Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&tem.TecId, &tem.TecName, &tem.TecSex, &tem.TecBirth, &tem.TecMajor, &tem.TecPhone)
			if err != nil {
				log.Printf("err:%s", err)
			}
			s.TeacherList = append(s.TeacherList, tem)
		}
		s.IsSuccess = true
		render.JSON(w, r, s)
	}
}

func (teacher *TeacherAPi) UpTeacher(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tecId := r.Form["tecId"][0]
	tecName := r.Form["tecName"][0]
	tecSex := r.Form["tecSex"][0]
	tecBirth := r.Form["tecBirth"][0]
	tecMajor := r.Form["tecMajor"][0]
	tecPhone := r.Form["tecPhone"][0]

	if tecName != "" {
		err := teacher.db.Model(&data_conn.TeacherInfo{}).Where("TecId=?", tecId).Update(&data_conn.TeacherInfo{TecName: tecName}).Error
		if err != nil {
			log.Printf("err:%s", err)
		}
	}

	if tecSex != "" {
		err := teacher.db.Model(&data_conn.TeacherInfo{}).Where("TecId=?", tecId).Update(&data_conn.TeacherInfo{TecSex: tecSex}).Error
		if err != nil {
			log.Printf("err:%s", err)
		}
	}

	if tecBirth != "" {
		t, _ := time.Parse("2006-01-02", tecBirth) //字符串转时间戳
		err := teacher.db.Model(&data_conn.TeacherInfo{}).Where("TecId=?", tecId).Update(&data_conn.TeacherInfo{TecBirth: t}).Error
		if err != nil {
			log.Printf("err:%s", err)
		}
	}

	if tecMajor != "" {
		err := teacher.db.Model(&data_conn.TeacherInfo{}).Where("TecId=?", tecId).Update(&data_conn.TeacherInfo{TecMajor: tecMajor}).Error
		if err != nil {
			log.Printf("err:%s", err)
		}
	}

	if tecPhone != "" {
		err := teacher.db.Model(&data_conn.TeacherInfo{}).Where("TecId=?", tecId).Update(&data_conn.TeacherInfo{TecPhone: tecPhone}).Error
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	s := structure_type.Things{"更新教师信息成功", true}
	render.JSON(w, r, s)
}

func (teacher *TeacherAPi) DelTeacher(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tecId := r.Form["tecId"][0]

	err := teacher.db.Model(&data_conn.TeacherInfo{}).Where("TecId=?", tecId).Delete(&data_conn.TeacherInfo{}).Error
	if err != nil {
		log.Printf("err:%s", err)
	}
	s := structure_type.Things{"删除教师信息成功", true}
	render.JSON(w, r, s)
}
