package MajorMation

import (
		"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"xiangmu/Student/DataConn"
	"xiangmu/Student/StructureType"
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
	maj_name := r.Form["Maj_name"][0]
	maj_prin := r.Form["Maj_prin"][0]
	maj_link := r.Form["Maj_link"][0]
	maj_phone := r.Form["Maj_phone"][0]

	var maj_id int
	if maj_name == "" || maj_prin == "" || maj_link == "" || maj_phone == "" {
		s := StructureType.Things{"请将信息输入完整"}
		render.JSON(w, r, s)
		return
	}

	rows, err := major.db.Model(&DataConn.Major{}).Where("Maj_name=?", maj_name).Select("Maj_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&maj_id)
		if err != nil {
			return
		}
	}
	if maj_id != 0 {
		s := StructureType.Things{"本专业已存在"}
		render.JSON(w, r, s)
		return
	}
	err = major.db.Create(&DataConn.Major{Maj_name: maj_name, Maj_prin: maj_prin, Maj_link: maj_link, Maj_phone: maj_phone}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"专业信息添加成功"}
	render.JSON(w, r, s)
}

func (major *MajorAPi) BrowMajor(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	m := StructureType.Majortotal{}
	tem:=StructureType.Major{}
	maj_name := r.Form["Maj_name"][0]

	if maj_name == ""{
		rows, err := major.db.Table("Majors").Select("maj_name,maj_prin,maj_link,maj_phone").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Maj_name, &tem.Maj_prin, &tem.Maj_link, &tem.Maj_phone)
			if err != nil {
				return
			}
			m.MajorList = append(m.MajorList,tem)
		}
		render.JSON(w, r, m)
	}

	if maj_name != "" {
		rows, err := major.db.Table("Majors").Where("Maj_name=?", maj_name).Select("maj_name,maj_prin,maj_link,maj_phone").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Maj_name, &tem.Maj_prin, &tem.Maj_link, &tem.Maj_phone)
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
	maj_id := r.Form["Maj_id"][0]
	maj_name := r.Form["Maj_name"][0]
	maj_prin := r.Form["Maj_prin"][0]
	maj_link := r.Form["Maj_link"][0]
	maj_phone := r.Form["Maj_phone"][0]

	if maj_name != "" {
		err := major.db.Model(&DataConn.Major{}).Where("Maj_id=?",maj_id).Update(&DataConn.Major{Maj_name: maj_name}).Error
		if err != nil {
			return
		}
	}

	if maj_prin != "" {
		err := major.db.Model(&DataConn.Major{}).Where("Maj_id=?",maj_id).Update(&DataConn.Major{Maj_prin: maj_prin}).Error
		if err != nil {
			return
		}
	}

	if maj_prin != "" {
		err := major.db.Model(&DataConn.Major{}).Where("Maj_id=?",maj_id).Update(&DataConn.Major{Maj_link: maj_link}).Error
		if err != nil {
			return
		}
	}

	if maj_prin != "" {
		err := major.db.Model(&DataConn.Major{}).Where("Maj_id=?",maj_id).Update(&DataConn.Major{Maj_phone: maj_phone}).Error
		if err != nil {
			return
		}
	}
	s := StructureType.Things{"更新专业信息成功"}
	render.JSON(w, r, s)
}

func (major *MajorAPi) SubMajor(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	maj_id := r.Form["Maj_id"][0]

	err := major.db.Model(&DataConn.Major{}).Where("Maj_id=?", maj_id).Delete(&DataConn.Major{}).Error
	if err!= nil {
		return
	}
	s := StructureType.Things{"删除专业信息成功"}
	render.JSON(w, r, s)
}
