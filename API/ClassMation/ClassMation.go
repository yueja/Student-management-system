package ClassMation

import (
		"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"xiangmu/Student/DataConn"
	"xiangmu/Student/StructureType"
)

type ClassAPi struct {
	db *gorm.DB
}

func Make_db(db *gorm.DB) *ClassAPi {
	DB := &ClassAPi{db}
	return DB
}

func (class *ClassAPi) AddClass(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	class_name := r.Form["Class_name"][0]
	class_tec := r.Form["Class_tec"][0]
	maj_name := r.Form["Maj_name"][0]

	var maj_id,class_id int
	if class_name == "" || class_tec == "" || maj_name == "" {
		s := StructureType.Things{"请将信息输入完整"}
		render.JSON(w, r, s)
		return
	}
	//查询对应专业id
	rows, err := class.db.Table("Majors").Where("Maj_name=?", maj_name).Select("Maj_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&maj_id)
	}

	//判断班级是否已经存在
	rows, err = class.db.Table("Classinfos").Where("Class_name=?", class_name).Select("Class_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&class_id)
	}
	if class_id != 0 {
		s := StructureType.Things{"本班级已存在"}
		render.JSON(w, r, s)
		return
	}
	err = class.db.Create(&DataConn.Classinfo{Class_name: class_name, Class_tec: class_tec, Maj_id: maj_id}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"班级信息添加成功"}
	render.JSON(w, r, s)
}

func (class *ClassAPi) BrowClass(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	m := StructureType.Classtotal{}
	tem := StructureType.Class{}

	class_name := r.Form["Class_name"][0]
	var maj_id int

	if class_name == "" {
		rows, err := class.db.Table("Classinfos").Select("Class_name,Class_tec,Maj_id").Rows()
		if err != nil {
			return
		}

		for rows.Next() {
			err = rows.Scan(&tem.Class_name, &tem.Class_tec, &maj_id)
			tem.Maj_name = strconv.Itoa(maj_id)
			m.ClassList = append(m.ClassList, tem)
		}
		//查询班级对应的专业
		for i := 0; i < len(m.ClassList); i++ {
			rows, err = class.db.Table("Majors").Where("Maj_id=?", m.ClassList[i].Maj_name).Select("Maj_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&m.ClassList[i].Maj_name)
				if err != nil {
					return
				}
			}
		}
		render.JSON(w, r, m)
	}

	if class_name != "" {
		rows, err := class.db.Table("Classinfos").Where("Class_name=?", class_name).Select("Class_tec,Maj_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Class_tec, &maj_id)
			if err != nil {
				return
			}
			//查询班级对应的专业
			rows, err = class.db.Table("Majors").Where("Maj_id=?", maj_id).Select("Maj_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&tem.Maj_name)
				if err != nil {
					return
				}
			}
		}
		tem.Class_name = class_name
		m.ClassList = append(m.ClassList, tem)
		render.JSON(w, r, m)
	}
}

func (class *ClassAPi) UpClass(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	class_id := r.Form["Class_id"][0]
	class_name := r.Form["Class_name"][0]
	class_tec := r.Form["Class_tec"][0]
	maj_name := r.Form["Maj_name"][0]

	if class_name != "" {
		err := class.db.Model(&DataConn.Classinfo{}).Where("Class_id=?", class_id).Update(&DataConn.Classinfo{Class_name: class_name}).Error
		if err != nil {
			return
		}
	}

	if class_tec != "" {
		err := class.db.Model(&DataConn.Classinfo{}).Where("Class_id=?", class_id).Update(&DataConn.Classinfo{Class_tec: class_tec}).Error
		if err != nil {
			return
		}
	}

	if maj_name != "" {
		var maj_id int
		rows, err := class.db.Table("Majors").Where("Maj_name=?", maj_name).Select("Maj_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&maj_id)
			if err != nil {
				return
			}
		}
		err = class.db.Model(&DataConn.Classinfo{}).Where("Class_id=?", class_id).Update(&DataConn.Classinfo{Maj_id: maj_id}).Error
		if err != nil {
			return
		}
	}
	s := StructureType.Things{"更新专业信息成功"}
	render.JSON(w, r, s)
}

func (class *ClassAPi) SubClass(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	class_id := r.Form["Class_id"][0]

	err := class.db.Model(&DataConn.Classinfo{}).Where("Class_id=?", class_id).Delete(&DataConn.Classinfo{}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"删除班级信息成功"}
	render.JSON(w, r, s)
}
