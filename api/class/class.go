package class

import (
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"xiangmu/Student/data_conn"
	"xiangmu/Student/structure_type"
)

type ClassAPi struct {
	db *gorm.DB
}

func MakeDb(db *gorm.DB) *ClassAPi {
	DB := &ClassAPi{db}
	return DB
}

func (class *ClassAPi) AddClass(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	className := r.Form["className"][0]
	classTec := r.Form["classTec"][0]
	majName := r.Form["majName"][0]

	var majId, classId int
	if className == "" || classTec == "" || majName == "" {
		s := structure_type.Things{"请将信息输入完整", false}
		render.JSON(w, r, s)
		return
	}
	//查询对应专业id
	rows, err := class.db.Model(&data_conn.Major{}).Where("MajName=?", majName).Select("MajId").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&majId)
	}

	//判断班级是否已经存在
	rows, err = class.db.Model(&data_conn.ClassInfo{}).Where("ClassName=?", className).Select("ClassId").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&classId)
	}
	if classId != 0 {
		s := structure_type.Things{"本班级已存在", false}
		render.JSON(w, r, s)
		return
	}
	err = class.db.Create(&data_conn.ClassInfo{ClassName: className, ClassTec: classTec, MajId: majId}).Error
	if err != nil {
		return
	}
	s := structure_type.Things{"班级信息添加成功", true}
	render.JSON(w, r, s)
}

func (class *ClassAPi) BrowClass(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	m := structure_type.ClassTotal{}
	tem := structure_type.Class{}

	className := r.Form["className"][0]
	var majId int

	if className == "" {
		rows, err := class.db.Model(&data_conn.ClassInfo{}).Select("ClassName,ClassTec,MajId").Rows()
		if err != nil {
			return
		}

		for rows.Next() {
			err = rows.Scan(&tem.ClassName, &tem.ClassTec, &majId)
			tem.MajName = strconv.Itoa(majId)
			m.ClassList = append(m.ClassList, tem)
		}
		//查询班级对应的专业
		for i := 0; i < len(m.ClassList); i++ {
			rows, err = class.db.Model(&data_conn.Major{}).Where("MajId=?", m.ClassList[i].MajName).Select("MajName").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&m.ClassList[i].MajName)
				if err != nil {
					return
				}
			}
		}
		render.JSON(w, r, m)
	}

	if className != "" {
		rows, err := class.db.Model(&data_conn.ClassInfo{}).Where("ClassName=?", className).Select("ClassTec,MajId").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.ClassTec, &majId)
			if err != nil {
				return
			}
			//查询班级对应的专业
			rows, err = class.db.Model(&data_conn.Major{}).Where("MajId=?", majId).Select("MajName").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&tem.MajName)
				if err != nil {
					return
				}
			}
		}
		tem.ClassName = className
		m.ClassList = append(m.ClassList, tem)
		render.JSON(w, r, m)
	}
}

func (class *ClassAPi) UpClass(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	classId := r.Form["classId"][0]
	className := r.Form["className"][0]
	classTec := r.Form["classTec"][0]
	majName := r.Form["majName"][0]

	if className != "" {
		err := class.db.Model(&data_conn.ClassInfo{}).Where("ClassId=?", classId).Update(&data_conn.ClassInfo{ClassName: className}).Error
		if err != nil {
			return
		}
	}

	if classTec != "" {
		err := class.db.Model(&data_conn.ClassInfo{}).Where("ClassId=?", classId).Update(&data_conn.ClassInfo{ClassTec: classTec}).Error
		if err != nil {
			return
		}
	}

	if majName != "" {
		var majId int
		rows, err := class.db.Model(&data_conn.Major{}).Where("MajName=?", majName).Select("MajId").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&majId)
			if err != nil {
				return
			}
		}
		err = class.db.Model(&data_conn.ClassInfo{}).Where("ClassId=?", classId).Update(&data_conn.ClassInfo{MajId: majId}).Error
		if err != nil {
			return
		}
	}
	s := structure_type.Things{"更新专业信息成功", true}
	render.JSON(w, r, s)
}

func (class *ClassAPi) DelClass(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	classId := r.Form["classId"][0]

	err := class.db.Model(&data_conn.ClassInfo{}).Where("ClassId=?", classId).Delete(&data_conn.ClassInfo{}).Error
	if err != nil {
		return
	}
	s := structure_type.Things{"删除班级信息成功", true}
	render.JSON(w, r, s)
}
