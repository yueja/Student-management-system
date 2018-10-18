package major

import (
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"xiangmu/Student/data_conn"
	"xiangmu/Student/structure_type"
)

type MajorAPi struct {
	db *gorm.DB
}

func MakeDb(db *gorm.DB) *MajorAPi {
	DB := &MajorAPi{db}
	return DB
}
func (m *MajorAPi) AddMajor(w http.ResponseWriter, r *http.Request) {
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

	rows, err := m.db.Model(&data_conn.Major{}).Where("MajName=?", majName).Select("MajId").Rows()
	if err != nil {
		log.Printf("err:%s", err)
	}
	for rows.Next() {
		err = rows.Scan(&majId)
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	if majId != 0 {
		s := structure_type.Things{"本专业已存在", false}
		render.JSON(w, r, s)
		return
	}
	err = m.db.Create(&data_conn.Major{MajName: majName, MajPrin: majPrin, MajLink: majLink, MajPhone: majPhone}).Error
	if err != nil {
		log.Printf("err:%s", err)
	}
	s := structure_type.Things{"专业信息添加成功", true}
	render.JSON(w, r, s)
}

func (m *MajorAPi) BrowMajor(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	s := structure_type.MajorTotal{}
	tem := structure_type.Major{}
	majName := r.Form["majName"][0]

	if majName == "" {
		rows, err := m.db.Model(&data_conn.Major{}).Select("MajName,majPrin,majLink,majPhone").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&tem.MajName, &tem.MajPrin, &tem.MajLink, &tem.MajPhone)
			if err != nil {
				log.Printf("err:%s", err)
			}
			s.MajorList = append(s.MajorList, tem)
		}
		s.IsSuccess = true
		render.JSON(w, r, s)
	}

	if majName != "" {
		rows, err := m.db.Model(&data_conn.Major{}).Where("Maj_name=?", majName).Select("majName,majPrin,majLink,majPhone").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&tem.MajName, &tem.MajPrin, &tem.MajLink, &tem.MajPhone)
			if err != nil {
				log.Printf("err:%s", err)
			}
			s.MajorList = append(s.MajorList, tem)
		}
		s.IsSuccess = true
		render.JSON(w, r,s)
	}
}

func (m *MajorAPi) UpMajor(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	majId := r.Form["majId"][0]
	majName := r.Form["majName"][0]
	majPrin := r.Form["majPrin"][0]
	majLink := r.Form["majLink"][0]
	majPhone := r.Form["majPhone"][0]

	if majName != "" {
		err := m.db.Model(&data_conn.Major{}).Where("MajId=?", majId).Update(&data_conn.Major{MajName: majName}).Error
		if err != nil {
			log.Printf("err:%s", err)
		}
	}

	if majPrin != "" {
		err := m.db.Model(&data_conn.Major{}).Where("MajId=?", majId).Update(&data_conn.Major{MajPrin: majPrin}).Error
		if err != nil {
			log.Printf("err:%s", err)
		}
	}

	if majLink != "" {
		err := m.db.Model(&data_conn.Major{}).Where("MajId=?", majId).Update(&data_conn.Major{MajLink: majLink}).Error
		if err != nil {
			log.Printf("err:%s", err)
		}
	}

	if majPhone != "" {
		err := m.db.Model(&data_conn.Major{}).Where("MajId=?", majId).Update(&data_conn.Major{MajPhone: majPhone}).Error
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	s := structure_type.Things{"更新专业信息成功", true}
	render.JSON(w, r, s)
}

func (m *MajorAPi) DelMajor(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	majId := r.Form["majId"][0]

	err := m.db.Model(&data_conn.Major{}).Where("MajId=?", majId).Delete(&data_conn.Major{}).Error
	if err != nil {
		log.Printf("err:%s", err)
	}
	s := structure_type.Things{"删除专业信息成功", true}
	render.JSON(w, r, s)
}
