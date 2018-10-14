package major

import (
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"xiangmu/Student/data_conn"
	"xiangmu/Student/structure_type"
)

type MajorAPi struct {
	db *gorm.DB
}

func Make_db(db *gorm.DB) *MajorAPi {
	DB := &MajorAPi{db}
	return DB
}
func (major *MajorAPi) AddMajor(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	majName := r.Form["majName"][0]
	majPrin := r.Form["majPrin"][0]
	majLink := r.Form["majLink"][0]
	majPhone := r.Form["majPhone"][0]

	var majId int
	if majName == "" || majPrin == "" || majLink == "" || majPhone == "" {
		s := structure_type.Things{"请将信息输入完整", false}
		render.JSON(w, r, s)
		return
	}

	rows, err := major.db.Model(&data_conn.Major{}).Where("MajName=?", majName).Select("MajId").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&majId)
		if err != nil {
			return
		}
	}
	if majId != 0 {
		s := structure_type.Things{"本专业已存在", false}
		render.JSON(w, r, s)
		return
	}
	err = major.db.Create(&data_conn.Major{MajName: majName, MajPrin: majPrin, MajLink: majLink, MajPhone: majPhone}).Error
	if err != nil {
		return
	}
	s := structure_type.Things{"专业信息添加成功", true}
	render.JSON(w, r, s)
}

func (major *MajorAPi) BrowMajor(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	m := structure_type.MajorTotal{}
	tem := structure_type.Major{}
	majName := r.Form["majName"][0]

	if majName == "" {
		rows, err := major.db.Model(&data_conn.Major{}).Select("MajName,majPrin,majLink,majPhone").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.MajName, &tem.MajPrin, &tem.MajLink, &tem.MajPhone)
			if err != nil {
				return
			}
			m.MajorList = append(m.MajorList, tem)
		}
		render.JSON(w, r, m)
	}

	if majName != "" {
		rows, err := major.db.Model(&data_conn.Major{}).Where("Maj_name=?", majName).Select("majName,majPrin,majLink,majPhone").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.MajName, &tem.MajPrin, &tem.MajLink, &tem.MajPhone)
			if err != nil {
				return
			}
			m.MajorList = append(m.MajorList, tem)
		}
		render.JSON(w, r, m)
	}
}

func (major *MajorAPi) UpMajor(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	majId := r.Form["majId"][0]
	majName := r.Form["majName"][0]
	majPrin := r.Form["majPrin"][0]
	majLink := r.Form["majLink"][0]
	majPhone := r.Form["majPhone"][0]

	if majName != "" {
		err := major.db.Model(&data_conn.Major{}).Where("MajId=?", majId).Update(&data_conn.Major{MajName: majName}).Error
		if err != nil {
			return
		}
	}

	if majPrin != "" {
		err := major.db.Model(&data_conn.Major{}).Where("MajId=?", majId).Update(&data_conn.Major{MajPrin: majPrin}).Error
		if err != nil {
			return
		}
	}

	if majLink != "" {
		err := major.db.Model(&data_conn.Major{}).Where("MajId=?", majId).Update(&data_conn.Major{MajLink: majLink}).Error
		if err != nil {
			return
		}
	}

	if majPhone != "" {
		err := major.db.Model(&data_conn.Major{}).Where("MajId=?", majId).Update(&data_conn.Major{MajPhone: majPhone}).Error
		if err != nil {
			return
		}
	}
	s := structure_type.Things{"更新专业信息成功", true}
	render.JSON(w, r, s)
}

func (major *MajorAPi) DelMajor(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	majId := r.Form["majId"][0]

	err := major.db.Model(&data_conn.Major{}).Where("MajId=?", majId).Delete(&data_conn.Major{}).Error
	if err != nil {
		return
	}
	s := structure_type.Things{"删除专业信息成功", true}
	render.JSON(w, r, s)
}
