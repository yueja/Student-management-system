package ClassSubjectMation

import (
		"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"xiangmu/Student/DataConn"
	"xiangmu/Student/StructureType"
	)

type Cla2subAPi struct {
	db *gorm.DB
}

func Make_db(db *gorm.DB) *Cla2subAPi {
	DB := &Cla2subAPi{db}
	return DB
}
func (cla_sub *Cla2subAPi) AddCla_Sub(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	class_name := r.Form["Class_name"][0]
	sub_name := r.Form["Sub_name"][0]
	tec_name := r.Form["Tec_name"][0]

	var class_id, sub_id, tec_id, cla2sub_id int

	if class_name == "" || sub_name == "" || tec_name == "" {
		s := StructureType.Things{"请将信息输入完整"}
		render.JSON(w, r, s)
		return
	}
	//判断班级是否存在
	rows, err := cla_sub.db.Table("Classinfos").Where("Class_name=?", class_name).Select("Class_id").Rows()
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
		s := StructureType.Things{"本班级信息不存在"}
		render.JSON(w, r, s)
		return
	}
	//判断课程是否存在
	rows, err = cla_sub.db.Table("Subjects").Where("Sub_name=?", sub_name).Select("Sub_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&sub_id)
		if err != nil {
			return
		}
	}
	if sub_id == 0 {
		s := StructureType.Things{"本课程信息不存在"}
		render.JSON(w, r, s)
		return
	}

	//判断教师是否存在
	rows, err = cla_sub.db.Table("Teacherinfos").Where("Tec_name=?", tec_name).Select("Tec_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&tec_id)
		if err != nil {
			return
		}
	}
	if tec_id == 0 {
		s := StructureType.Things{"本教师信息不存在"}
		render.JSON(w, r, s)
		return
	}
	//判断本班级是否已存在该课程
	rows, err = cla_sub.db.Table("Cla2subs").Where("Cla_id=? and Sub_id=?", class_id, sub_id).Select("Cla2sub_id").Rows()
	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&cla2sub_id)
		if err != nil {
			return
		}
	}
	if cla2sub_id != 0 {
		s := StructureType.Things{"本班级该课程信息已存在"}
		render.JSON(w, r, s)
		return
	}
	//添加课程到班级
	err = cla_sub.db.Create(&DataConn.Cla2sub{Cla_id: class_id, Sub_id: sub_id, Tec_id: tec_id}).Error
	if err != nil {
		s := StructureType.Things{"添加班级课程信息失败"}
		render.JSON(w, r, s)
		return
	}
	s := StructureType.Things{"添加班级课程信息成功"}
	render.JSON(w, r, s)
}

func (cla_sub *Cla2subAPi) BrowCla_Sub(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	class_name := r.Form["Class_name"][0]
	sub_name := r.Form["Sub_name"][0]
	tec_name := r.Form["Tec_name"][0]

	s := StructureType.Cla2subtotal{}
	tem := StructureType.Cla2sub{}

	var class_id, sub_id,tec_id int

	//按班级查询
	if class_name != "" && sub_name == "" && tec_name == "" {
		rows, err := cla_sub.db.Table("Classinfos").Where("Class_name=?", class_name).Select("Class_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&class_id)
			if err != nil {
				return
			}
		}

		rows, err = cla_sub.db.Table("Cla2subs").Where("Cla_id=?", class_id).Select("Cla2sub_id,Sub_id,Tec_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Cla2sub_id, &tem.Sub_name, &tem.Tec_name)
			if err != nil {
				return
			}
			tem.Cla_name = class_name
			s.Cla2subList = append(s.Cla2subList, tem)
		}
		for i := 0; i <len(s.Cla2subList); i++ {
			//查询科目名称
			rows, err = cla_sub.db.Table("Subjects").Where("Sub_id=?",s.Cla2subList[i].Sub_name).Select("Sub_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].Sub_name)
				if err != nil {
					return
				}
			}
			//查询教师名字
			rows, err = cla_sub.db.Table("Teacherinfos").Where("Tec_id=?",s.Cla2subList[i].Tec_name).Select("Tec_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].Tec_name)
				if err != nil {
					return
				}
			}
		}
		render.JSON(w, r, s)
	}

	//按课程查询
	if sub_name !=""&& class_name == "" && tec_name == "" {
		rows, err := cla_sub.db.Table("Subjects").Where("Sub_name=?", sub_name).Select("Sub_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&sub_id)
			if err != nil {
				return
			}
		}

		rows, err = cla_sub.db.Table("Cla2subs").Where("Sub_id=?", sub_id).Select("Cla2sub_id,Cla_id,Tec_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Cla2sub_id, &tem.Cla_name, &tem.Tec_name)
			if err != nil {
				return
			}
			tem.Sub_name = sub_name
			s.Cla2subList = append(s.Cla2subList, tem)
		}
		for i := 0; i <len(s.Cla2subList); i++ {
			//查询班级名称
			rows, err = cla_sub.db.Table("Classinfos").Where("Class_id=?",s.Cla2subList[i].Cla_name).Select("Class_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].Cla_name)
				if err != nil {
					return
				}
			}
			//查询教师名字
			rows, err = cla_sub.db.Table("Teacherinfos").Where("Tec_id=?",s.Cla2subList[i].Tec_name).Select("Tec_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].Tec_name)
				if err != nil {
					return
				}
			}
		}
		render.JSON(w, r, s)
	}

	//按老师查询
	if tec_name !=""&& class_name == "" && sub_name == "" {
		rows, err := cla_sub.db.Table("Teacherinfos").Where("Tec_name=?", tec_name).Select("Tec_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tec_id)
			if err != nil {
				return
			}
		}

		rows, err = cla_sub.db.Table("Cla2subs").Where("Tec_id=?", tec_id).Select("Cla2sub_id,Cla_id,Sub_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Cla2sub_id, &tem.Cla_name, &tem.Sub_name)
			if err != nil {
				return
			}
			tem.Tec_name = tec_name
			s.Cla2subList = append(s.Cla2subList, tem)
		}
		for i := 0; i <len(s.Cla2subList); i++ {
			//查询班级名称
			rows, err = cla_sub.db.Table("Classinfos").Where("Class_id=?",s.Cla2subList[i].Cla_name).Select("Class_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].Cla_name)
				if err != nil {
					return
				}
			}
			//查询教师名字
			rows, err = cla_sub.db.Table("Subjects").Where("Sub_id=?",s.Cla2subList[i].Sub_name).Select("Sub_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].Sub_name)
				if err != nil {
					return
				}
			}
		}
		render.JSON(w, r, s)
	}

	//所有班级所有课程
	if class_name == "" || sub_name == "" || tec_name == "" {

		rows, err := cla_sub.db.Table("Cla2subs").Select("Cla2sub_id,Cla_id,Sub_id,Tec_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Cla2sub_id, &tem.Cla_name, &tem.Sub_name, &tem.Tec_name)
			if err != nil {
				return
			}
			s.Cla2subList=append(s.Cla2subList,tem)
		}

		for i := 0; i < len(s.Cla2subList); i++{
			//查班级名字
			rows, err = cla_sub.db.Table("Classinfos").Where("Class_id=?",s.Cla2subList[i].Cla_name).Select("Class_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].Cla_name)
				if err != nil {
					return
				}
			}
//查课程名字
			rows, err = cla_sub.db.Table("Subjects").Where("Sub_id=?",s.Cla2subList[i].Sub_name).Select("Sub_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].Sub_name)
				if err != nil {
					return
				}
			}
			//查老师的名字
			rows, err := cla_sub.db.Table("Teacherinfos").Where("Tec_id=?",s.Cla2subList[i].Tec_name).Select("Tec_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].Tec_name)
				if err != nil {
					return
				}
			}
		}
		render.JSON(w, r, s)
	}
}


func (cla_sub *Cla2subAPi) UpCla_Sub(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cla2sub_id_1 := r.Form["Cla2sub_id"][0]
	class_name := r.Form["Class_name"][0]
	sub_name := r.Form["Sub_name"][0]
	tec_name := r.Form["Tec_name"][0]

	var class_id,sub_id,tec_id int

	if class_name != ""{
		rows, err := cla_sub.db.Table("Classinfos").Where("Class_name=?", class_name).Select("Class_id").Rows()
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
			s := StructureType.Things{"本班级信息不存在"}
			render.JSON(w, r, s)
			return
		}
		err = cla_sub.db.Model(&DataConn.Cla2sub{}).Where("Cla2sub_id=?", cla2sub_id_1).Update(&DataConn.Cla2sub{Cla_id: class_id}).Error
		if err != nil {
			return
		}
	}

	if class_name != "" {
		rows, err := cla_sub.db.Table("Subjects").Where("Sub_name=?", sub_name).Select("Sub_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&sub_id)
			if err != nil {
				return
			}
		}
		if sub_id == 0 {
			s := StructureType.Things{"本课程信息不存在"}
			render.JSON(w, r, s)
			return
		}
		err = cla_sub.db.Model(&DataConn.Cla2sub{}).Where("Cla2sub_id=?", cla2sub_id_1).Update(&DataConn.Cla2sub{Sub_id: sub_id}).Error
		if err != nil {
			return
		}
	}

	if tec_name != "" {
		rows, err := cla_sub.db.Table("Teacherinfos").Where("Tec_name=?", tec_name).Select("Tec_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tec_id)
			if err != nil {
				return
			}
		}
		if tec_id == 0 {
			s := StructureType.Things{"本教师信息不存在"}
			render.JSON(w, r, s)
			return
		}
		err = cla_sub.db.Model(&DataConn.Cla2sub{}).Where("Cla2sub_id=?", cla2sub_id_1).Update(&DataConn.Cla2sub{Tec_id: tec_id}).Error
		if err != nil {
			return
		}
	}
	s := StructureType.Things{"编辑班级课程信息成功"}
	render.JSON(w, r, s)
}

func (cla_sub *Cla2subAPi) SubCla_Sub(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cla2sub_id_1 := r.Form["Cla2sub_id"][0]

	err := cla_sub.db.Model(&DataConn.Cla2sub{}).Where("Cla2sub_id=?", cla2sub_id_1).Delete(&DataConn.Cla2sub{}).Error
	if err != nil {
		s := StructureType.Things{"删除班级课程信息失败"}
		render.JSON(w, r, s)
		return
	}
	s := StructureType.Things{"删除班级课程信息成功"}
	render.JSON(w, r, s)
}