package StudentMation

import (
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
	"xiangmu/Student/DataConn"
	"xiangmu/Student/StructureType"
	)

type StudentAPi struct {
	db *gorm.DB
}

func Make_db(db *gorm.DB) *StudentAPi {
	DB := &StudentAPi{db}
	return DB
}

func (student *StudentAPi) AddStudent(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	stu_nu := r.Form["Stu_nu"][0]
	stu_name := r.Form["Stu_name"][0]
	stu_sex := r.Form["Stu_sex"][0]
	stu_birth := r.Form["Stu_birth"][0]
	class_name := r.Form["Class_name"][0]

	var class_id,stu_id int
	if stu_nu == "" || stu_name == "" || stu_sex == "" || stu_birth == "" || class_name == "" {
		s := StructureType.Things{"请将信息输入完整"}
		render.JSON(w, r, s)
		return
	}

	//查询对应班级id
	rows, err := student.db.Table("Classinfos").Where("Class_name=?", class_name).Select("Class_id ").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&class_id)
		if err != nil {
			return
		}
	}
	if class_id == 0 {
		s := StructureType.Things{"班级信息输入错误，班级不存在"}
		render.JSON(w, r, s)
		return
	}
	//判断学生是否已经存在
	rows, err = student.db.Table("Studentinfos").Where("Stu_nu=?", stu_nu).Select("Stu_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&stu_id)
		if err != nil {
			return
		}
	}
	if stu_id != 0 {
		s := StructureType.Things{"本学生信息已存在"}
		render.JSON(w, r, s)
		return
	}
	t, _ := time.Parse("2006-01-02", stu_birth) //字符串转时间戳
	err = student.db.Create(&DataConn.Studentinfo{Stu_nu: stu_nu, Stu_name: stu_name, Stu_sex: stu_sex, Stu_birth: t, Cla_id: class_id}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"班级信息添加成功"}
	render.JSON(w, r, s)
}

func (student *StudentAPi) BrowStudent(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	maj_name := r.Form["Maj_name"][0]
	class_name := r.Form["Class_name"][0]
	stu_name := r.Form["Stu_name"][0]

	s := StructureType.Studenttotal{}
	tem := StructureType.Studentinfo{}
	var maj_id,class_id int
	//按专业浏览学生信息
	if maj_name != "" {
		rows, err := student.db.Table("Majors").Where("Maj_name=?", maj_name).Select("Maj_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&maj_id)
			if err != nil {
				return
			}
		}
		//查同专业的班级id
		rows, err = student.db.Table("Classinfos").Where("Maj_id=?", maj_id).Select("Class_id").Rows()
		if err != nil {
			return
		}
		var class_id_1 []int
		for rows.Next() {
			err = rows.Scan(&class_id)
			if err != nil {
				return
			}
			class_id_1 = append(class_id_1, class_id)
		}

		for i := 0; i < len(class_id_1); i++ {
			a := "Stu_id,Ope_id,Stu_nu,Stu_name,Stu_sex,Stu_birth,Stu_pic"
			var class_name string
			rows, err = student.db.Table("Classinfos").Where("Class_id=?", class_id_1[i]).Select("Class_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&class_name)
				if err != nil {
					return
				}
			}
			rows, err = student.db.Table("Studentinfos").Where("Cla_id=?", class_id_1[i]).Select(a).Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&tem.Stu_id, &tem.Ope_id, &tem.Stu_nu, &tem.Stu_name, &tem.Stu_sex, &tem.Stu_birth, &tem.Stu_pic)
				if err != nil {
					return
				}
				tem.Maj_name = maj_name
				tem.Cla_name = class_name
				s.StudentList = append(s.StudentList, tem)
			}
		}
		render.JSON(w, r, s)
	}

	//按照班级查询学生信息
	if class_name != "" {
		var class_id, maj_id int
		var maj_name string

		rows, err := student.db.Table("Classinfos").Where("Class_name=?", class_name).Select("Class_id,Maj_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&class_id, &maj_id)
			if err != nil {
				return
			}
		}
		//查询专业名字
		rows, err = student.db.Table("Majors").Where("Maj_id=?", maj_id).Select("Maj_name").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&maj_name)
			if err != nil {
				return
			}
		}
		//查询学生信息
		a := "Stu_id,Ope_id,Stu_nu,Stu_name,Stu_sex,Stu_birth,Stu_pic"
		rows, err = student.db.Table("Studentinfos").Where("Cla_id=?", class_id).Select(a).Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Stu_id, &tem.Ope_id, &tem.Stu_nu, &tem.Stu_name, &tem.Stu_sex, &tem.Stu_birth, &tem.Stu_pic)
			if err != nil {
				return
			}
			tem.Cla_name = class_name
			tem.Maj_name = maj_name
			s.StudentList = append(s.StudentList, tem)
		}
		render.JSON(w, r, s)
	}

	//按照个人查询学生信息
	if stu_name != "" {
		var class_id, maj_id int
		a := "Stu_id,Ope_id,Stu_nu,Stu_name,Stu_sex,Stu_birth,Stu_pic,Cla_id"

		rows, err := student.db.Table("Studentinfos").Where("Stu_name=?", stu_name).Select(a).Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Stu_id, &tem.Ope_id, &tem.Stu_nu, &tem.Stu_name, &tem.Stu_sex, &tem.Stu_birth, &tem.Stu_pic,&class_id)
			if err != nil {
				return
			}
		}

		//查询班级名和专业id
		rows, err = student.db.Table("Classinfos").Where("Class_id=?", class_id).Select("Class_name,Maj_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Cla_name,&maj_id)
			if err != nil {
				return
			}
		}
		//查询专业
		rows, err = student.db.Table("Majors").Where("Maj_id=?", maj_id).Select("Maj_name").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Maj_name)
			if err != nil {
				return
			}
		}
		s.StudentList = append(s.StudentList, tem)
		render.JSON(w, r, s)
	}

	//查询所有学生信息
	if maj_name == "" && class_name == "" && stu_name == "" {
		rows, err := student.db.Table("Studentinfos").Select("Stu_id,Ope_id,Stu_nu,Stu_name,Stu_sex,Stu_birth,Stu_pic,Cla_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Stu_id, &tem.Ope_id, &tem.Stu_nu, &tem.Stu_name, &tem.Stu_sex, &tem.Stu_birth, &tem.Stu_pic, &tem.Cla_name)
			if err != nil {
				return
			}
			s.StudentList = append(s.StudentList, tem)
		}
		for i := 0; i < len(s.StudentList); i++ {
			rows, err := student.db.Table("Classinfos").Where("Class_id=?", s.StudentList[i].Cla_name).Select("Class_name,Maj_id").Rows()
			if err != nil {
				return
			}
			var maj_id int
			for rows.Next() {
				err = rows.Scan(&s.StudentList[i].Cla_name, &maj_id)
				if err != nil {
					return
				}
			}
			rows, err = student.db.Table("Majors").Where("Maj_id=?", maj_id).Select("Maj_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.StudentList[i].Maj_name)
				if err != nil {
					return
				}
			}
		}
		render.JSON(w, r, s)
	}
}

func (student *StudentAPi) UpClass(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	stu_id := r.Form["Stu_id"][0]
	stu_nu := r.Form["Stu_nu"][0]
	stu_name := r.Form["Stu_name"][0]
	stu_sex := r.Form["Stu_sex"][0]
	stu_birth := r.Form["Stu_birth"][0]
	class_name := r.Form["Class_name"][0]

	if stu_nu != "" {
		err := student.db.Model(&DataConn.Studentinfo{}).Where("Stu_id=?", stu_id).Update(&DataConn.Studentinfo{Stu_nu: stu_nu}).Error
		if err != nil {
			return
		}
	}

	if stu_name != "" {
		err := student.db.Model(&DataConn.Studentinfo{}).Where("Stu_id=?", stu_id).Update(&DataConn.Studentinfo{Stu_name: stu_name}).Error
		if err != nil {
			return
		}
	}

	if stu_sex != "" {
		err := student.db.Model(&DataConn.Studentinfo{}).Where("Stu_id=?", stu_id).Update(&DataConn.Studentinfo{Stu_sex: stu_sex}).Error
		if err != nil {
			return
		}
	}

	if stu_birth != "" {
		t, _ := time.Parse("2006-01-02", stu_birth) //字符串转时间戳
		err := student.db.Model(&DataConn.Studentinfo{}).Where("Stu_id=?", stu_id).Update(&DataConn.Studentinfo{Stu_birth: t}).Error
		if err != nil {
			return
		}
	}

	if class_name != "" {
		rows, err := student.db.Table("Classinfos").Where("Class_name=?", class_name).Select("Class_id").Rows()
		if err != nil {
			return
		}
		var class_id int
		for rows.Next() {
			err = rows.Scan(&class_id)
			if err != nil {
				return
			}
		}
		err = student.db.Model(&DataConn.Studentinfo{}).Where("Stu_id=?", stu_id).Update(&DataConn.Studentinfo{Cla_id: class_id}).Error
		if err != nil {
			return
		}
	}
	s := StructureType.Things{"更新学生信息成功"}
	render.JSON(w, r, s)
}

func (student *StudentAPi) SubStudent(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	stu_id := r.Form["Stu_id"][0]

	err := student.db.Model(&DataConn.Studentinfo{}).Where("Stu_id=?", stu_id).Delete(&DataConn.Studentinfo{}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"删除学生信息成功"}
	render.JSON(w, r, s)
}
