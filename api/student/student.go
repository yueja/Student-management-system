package student

import (
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
	"xiangmu/Student/structure_type"
	"xiangmu/Student/data_conn"
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
	stuNu := r.Form["stuNu"][0]
	stuName := r.Form["stuName"][0]
	stuSex := r.Form["stuSex"][0]
	stuBirth := r.Form["stuBirth"][0]
	className := r.Form["className"][0]

	var classId,stuId int
	if stuNu == "" || stuName == "" || stuSex == "" || stuBirth == "" || className == "" {
		s := structure_type.Things{"请将信息输入完整",false}
		render.JSON(w, r, s)
		return
	}

	//查询对应班级id
	rows, err := student.db.Model(&data_conn.ClassInfo{}).Where("ClassName=?", className).Select("ClassId ").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&classId)
		if err != nil {
			return
		}
	}
	if classId == 0 {
		s := structure_type.Things{"班级信息输入错误，班级不存在",false}
		render.JSON(w, r, s)
		return
	}
	//判断学生是否已经存在
	rows, err = student.db.Model(&data_conn.StudentInfo{}).Where("StuNu=?", stuNu).Select("StuId").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&stuId)
		if err != nil {
			return
		}
	}
	if stuId != 0 {
		s := structure_type.Things{"本学生信息已存在",false}
		render.JSON(w, r, s)
		return
	}
	t, _ := time.Parse("2006-01-02", stuBirth) //字符串转时间戳
	err = student.db.Create(&data_conn.StudentInfo{StuNu: stuNu, StuName: stuName, StuSex: stuSex, StuBirth: t, ClaId: classId}).Error
	if err != nil {
		return
	}
	s := structure_type.Things{"班级信息添加成功",true}
	render.JSON(w, r, s)
}

func (student *StudentAPi) BrowStudent(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	majName := r.Form["majName"][0]
	className := r.Form["className"][0]
	stuName := r.Form["stuName"][0]

	s := structure_type.StudentTotal{}
	tem := structure_type.StudentInfo{}
	var majId,classId int
	//按专业浏览学生信息
	if majName != "" {
		rows, err := student.db.Model(&data_conn.Major{}).Where("MajName=?", majName).Select("MajId").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&majId)
			if err != nil {
				return
			}
		}
		//查同专业的班级id
		rows, err = student.db.Model(&data_conn.ClassInfo{}).Where("MajId=?", majId).Select("ClassId").Rows()
		if err != nil {
			return
		}
		var classId_1 []int
		for rows.Next() {
			err = rows.Scan(&classId)
			if err != nil {
				return
			}
			classId_1 = append(classId_1, classId)
		}

		for i := 0; i < len(classId_1); i++ {
			a := "StuId,StuNu,StuName,StuSex,StuBirth,StuPic"
			var className_1 string
			rows, err = student.db.Model(&data_conn.ClassInfo{}).Where("ClassId=?", classId_1[i]).Select("ClassName").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&className_1)
				if err != nil {
					return
				}
			}
			rows, err = student.db.Model(&data_conn.StudentInfo{}).Where("ClaId=?", classId_1[i]).Select(a).Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&tem.StuId, &tem.StuNu, &tem.StuName, &tem.StuSex, &tem.StuBirth, &tem.StuPic)
				if err != nil {
					return
				}
				tem.MajName = majName
				tem.ClaName = className_1
				s.StudentList = append(s.StudentList, tem)
			}
		}
		render.JSON(w, r, s)
	}

	//按照班级查询学生信息
	if className != "" {
		var classId, majId int
		var majName_1 string

		rows, err := student.db.Model(&data_conn.ClassInfo{}).Where("ClassName=?", className).Select("ClassId,MajId").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&classId, &majId)
			if err != nil {
				return
			}
		}
		//查询专业名字
		rows, err = student.db.Model(&data_conn.Major{}).Where("MajId=?", majId).Select("MajName").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&majName_1)
			if err != nil {
				return
			}
		}
		//查询学生信息
		a := "StuId,StuNu,StuName,StuSex,StuBirth,StuPic"
		rows, err = student.db.Model(&data_conn.StudentInfo{}).Where("ClaId=?", classId).Select(a).Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.StuId,  &tem.StuNu, &tem.StuName, &tem.StuSex, &tem.StuBirth, &tem.StuPic)
			if err != nil {
				return
			}
			tem.ClaName = className
			tem.MajName = majName_1
			s.StudentList = append(s.StudentList, tem)
		}
		render.JSON(w, r, s)
	}

	//按照个人查询学生信息
	if stuName != "" {
		var classId, majId int
		a := "Stu_id,Stu_nu,Stu_name,Stu_sex,Stu_birth,Stu_pic,Cla_id"

		rows, err := student.db.Model(&data_conn.StudentInfo{}).Where("StuName=?", stuName).Select(a).Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.StuId, &tem.StuNu, &tem.StuName, &tem.StuSex, &tem.StuBirth, &tem.StuPic,&classId)
			if err != nil {
				return
			}
		}

		//查询班级名和专业id
		rows, err = student.db.Model(&data_conn.ClassInfo{}).Where("ClassId=?", classId).Select("ClassName,MajId").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.ClaName,&majId)
			if err != nil {
				return
			}
		}
		//查询专业
		rows, err = student.db.Model(&data_conn.Major{}).Where("MajId=?", majId).Select("MajName").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.MajName)
			if err != nil {
				return
			}
		}
		s.StudentList = append(s.StudentList, tem)
		render.JSON(w, r, s)
	}

	//查询所有学生信息
	if majName == "" && className == "" && stuName == "" {
		rows, err := student.db.Model(&data_conn.StudentInfo{}).Select("StuId,StuNu,StuName,StuSex,StuBirth,StuPic,ClaId").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.StuId, &tem.StuNu, &tem.StuName, &tem.StuSex, &tem.StuBirth, &tem.StuPic, &tem.ClaName)
			if err != nil {
				return
			}
			s.StudentList = append(s.StudentList, tem)
		}
		for i := 0; i < len(s.StudentList); i++ {
			rows, err := student.db.Model(&data_conn.ClassInfo{}).Where("ClassId=?", s.StudentList[i].ClaName).Select("ClassName,MajId").Rows()
			if err != nil {
				return
			}
			var majId int
			for rows.Next() {
				err = rows.Scan(&s.StudentList[i].ClaName, &majId)
				if err != nil {
					return
				}
			}
			rows, err = student.db.Model(&data_conn.Major{}).Where("MajId=?", majId).Select("MajName").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.StudentList[i].MajName)
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
	stuId := r.Form["stuId"][0]
	stuNu := r.Form["stuNu"][0]
	stuName := r.Form["stuName"][0]
	stuSex := r.Form["stuSex"][0]
	stuBirth := r.Form["stuBirth"][0]
	className := r.Form["className"][0]

	if stuNu != "" {
		err := student.db.Model(&data_conn.StudentInfo{}).Where("StuId=?", stuId).Update(&data_conn.StudentInfo{StuNu: stuNu}).Error
		if err != nil {
			return
		}
	}

	if stuName != "" {
		err := student.db.Model(&data_conn.StudentInfo{}).Where("StuId=?", stuId).Update(&data_conn.StudentInfo{StuName: stuName}).Error
		if err != nil {
			return
		}
	}

	if stuSex != "" {
		err := student.db.Model(&data_conn.StudentInfo{}).Where("StuId=?", stuId).Update(&data_conn.StudentInfo{StuSex: stuSex}).Error
		if err != nil {
			return
		}
	}

	if stuBirth != "" {
		t, _ := time.Parse("2006-01-02", stuBirth) //字符串转时间戳
		err := student.db.Model(&data_conn.StudentInfo{}).Where("StuId=?", stuId).Update(&data_conn.StudentInfo{StuBirth: t}).Error
		if err != nil {
			return
		}
	}

	if className != "" {
		rows, err := student.db.Model(&data_conn.ClassInfo{}).Where("ClassName=?", className).Select("ClassId").Rows()
		if err != nil {
			return
		}
		var classId int
		for rows.Next() {
			err = rows.Scan(&classId)
			if err != nil {
				return
			}
		}
		err = student.db.Model(&data_conn.StudentInfo{}).Where("StuId=?", stuId).Update(&data_conn.StudentInfo{ClaId: classId}).Error
		if err != nil {
			return
		}
	}
	s := structure_type.Things{"更新学生信息成功",true}
	render.JSON(w, r, s)
}

func (student *StudentAPi) DelStudent(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	stuId := r.Form["stuId"][0]

	err := student.db.Model(&data_conn.StudentInfo{}).Where("StuId=?", stuId).Delete(&data_conn.StudentInfo{}).Error
	if err != nil {
		return
	}
	s := structure_type.Things{"删除学生信息成功",true}
	render.JSON(w, r, s)
}
