package class_subject

import (
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"xiangmu/Student/data_conn"
	"xiangmu/Student/structure_type"
)

type Cla2subAPi struct {
	db *gorm.DB
}

func MakeDb(db *gorm.DB) *Cla2subAPi {
	DB := &Cla2subAPi{db}
	return DB
}

func (c *Cla2subAPi) AddCla_Sub(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	className := r.Form["className"][0]
	subName := r.Form["subName"][0]
	tecName := r.Form["tecName"][0]

	var classId, subId, tecId, cla2subId int

	if className == "" || subName == "" || tecName == "" {
		s := structure_type.Things{"请将信息输入完整", false}
		render.JSON(w, r, s)
		return
	}
	//判断班级是否存在
	rows, err := c.db.Model(&data_conn.ClassInfo{}).Where("ClassName=?", className).Select("ClassId").Rows()
	if err != nil {
		log.Printf("err:%s", err)
	}
	for rows.Next() {
		err = rows.Scan(&classId)
		if err != nil {
			log.Printf("err:%s", err)
		}
	}

	if classId == 0 {
		s := structure_type.Things{"本班级信息不存在", false}
		render.JSON(w, r, s)
		return
	}
	//判断课程是否存在
	rows, err = c.db.Model(&data_conn.Subject{}).Where("SubName=?", subName).Select("SubId").Rows()
	if err != nil {
		log.Printf("err:%s", err)
	}
	for rows.Next() {
		err = rows.Scan(&subId)
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	if subId == 0 {
		s := structure_type.Things{"本课程信息不存在", false}
		render.JSON(w, r, s)
		return
	}

	//判断教师是否存在
	rows, err = c.db.Model(&data_conn.TeacherInfo{}).Where("TecName=?", tecName).Select("TecId").Rows()
	if err != nil {
		log.Printf("err:%s", err)
	}
	for rows.Next() {
		err = rows.Scan(&tecId)
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	if tecId == 0 {
		s := structure_type.Things{"本教师信息不存在", false}
		render.JSON(w, r, s)
		return
	}
	//判断本班级是否已存在该课程
	rows, err = c.db.Model(&data_conn.Cla2sub{}).Where("ClaId=? and SubId=?",
		classId, subId).Select("Cla2subId").Rows()
	if err != nil {
		log.Printf("err:%s", err)
	}

	for rows.Next() {
		err = rows.Scan(&cla2subId)
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	if cla2subId != 0 {
		s := structure_type.Things{"本班级该课程信息已存在", false}
		render.JSON(w, r, s)
		return
	}
	//添加课程到班级
	err = c.db.Create(&data_conn.Cla2sub{ClaId: classId, SubId: subId, TecId: tecId}).Error
	if err != nil {
		log.Printf("err:%s", err)
	}
	s := structure_type.Things{"添加班级课程信息成功", true}
	render.JSON(w, r, s)
}

func (c *Cla2subAPi) BrowCla_Sub(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	className := r.Form["className"][0]
	subName := r.Form["subName"][0]
	tecName := r.Form["tecName"][0]

	s := structure_type.Cla2subTotal{}
	tem := structure_type.Cla2sub{}

	var classId, subId, tecId int

	//按班级查询
	if className != "" && subName == "" && tecName == "" {
		rows, err := c.db.Model(&data_conn.ClassInfo{}).Where("ClassName=?", className).Select("ClassId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&classId)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}

		rows, err = c.db.Model(&data_conn.Cla2sub{}).Where("ClaId=?", classId).Select("Cla2subId,SubId,TecId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&tem.Cla2subId, &tem.SubName, &tem.TecName)
			if err != nil {
				log.Printf("err:%s", err)
			}
			tem.ClaName = className
			s.Cla2subList = append(s.Cla2subList, tem)
		}
		for i := 0; i < len(s.Cla2subList); i++ {
			//查询科目名称
			rows, err = c.db.Model(&data_conn.Subject{}).Where("SubId=?", s.Cla2subList[i].SubName).Select("SubName").Rows()
			if err != nil {
				log.Printf("err:%s", err)
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].SubName)
				if err != nil {
					log.Printf("err:%s", err)
				}
			}
			//查询教师名字
			rows, err = c.db.Model(&data_conn.TeacherInfo{}).Where("TecId=?", s.Cla2subList[i].TecName).Select("TecName").Rows()
			if err != nil {
				log.Printf("err:%s", err)
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].TecName)
				if err != nil {
					log.Printf("err:%s", err)
				}
			}
		}
		s.IsSuccess = true
		render.JSON(w, r, s)
	}

	//按课程查询
	if subName != "" && className == "" && tecName == "" {
		rows, err := c.db.Model(&data_conn.Subject{}).Where("SubName=?", subName).Select("SubId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&subId)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}

		rows, err = c.db.Model(&data_conn.Cla2sub{}).Where("SubId=?", subId).Select("Cla2subId,ClaId,TecId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&tem.Cla2subId, &tem.ClaName, &tem.TecName)
			if err != nil {
				log.Printf("err:%s", err)
			}
			tem.SubName = subName
			s.Cla2subList = append(s.Cla2subList, tem)
		}
		for i := 0; i < len(s.Cla2subList); i++ {
			//查询班级名称
			rows, err = c.db.Model(&data_conn.ClassInfo{}).Where("ClassId=?", s.Cla2subList[i].ClaName).Select("ClassName").Rows()
			if err != nil {
				log.Printf("err:%s", err)
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].ClaName)
				if err != nil {
					log.Printf("err:%s", err)
				}
			}
			//查询教师名字
			rows, err = c.db.Model(&data_conn.TeacherInfo{}).Where("TecId=?", s.Cla2subList[i].TecName).Select("TecName").Rows()
			if err != nil {
				log.Printf("err:%s", err)
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].TecName)
				if err != nil {
					log.Printf("err:%s", err)
				}
			}
		}
		s.IsSuccess = true
		render.JSON(w, r, s)
	}

	//按老师查询
	if tecName != "" && className == "" && subName == "" {
		rows, err := c.db.Model(&data_conn.TeacherInfo{}).Where("TecName=?", tecName).Select("TecId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&tecId)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}

		rows, err = c.db.Model(&data_conn.Cla2sub{}).Where("TecId=?", tecId).Select("Cla2subId,ClaId,SubId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&tem.Cla2subId, &tem.ClaName, &tem.SubName)
			if err != nil {
				log.Printf("err:%s", err)
			}
			tem.TecName = tecName
			s.Cla2subList = append(s.Cla2subList, tem)
		}
		for i := 0; i < len(s.Cla2subList); i++ {
			//查询班级名称
			rows, err = c.db.Model(&data_conn.ClassInfo{}).Where("ClassId=?", s.Cla2subList[i].ClaName).Select("ClassName").Rows()
			if err != nil {
				log.Printf("err:%s", err)
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].ClaName)
				if err != nil {
					log.Printf("err:%s", err)
				}
			}
			//查询教师名字
			rows, err = c.db.Model(&data_conn.Subject{}).Where("SubId=?", s.Cla2subList[i].SubName).Select("SubName").Rows()
			if err != nil {
				log.Printf("err:%s", err)
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].SubName)
				if err != nil {
					log.Printf("err:%s", err)
				}
			}
		}
		s.IsSuccess = true
		render.JSON(w, r, s)
	}

	//所有班级所有课程
	if className == "" || subName == "" || tecName == "" {

		rows, err := c.db.Model(&data_conn.Cla2sub{}).Select("Cla2subId,ClaId,SubId,TecId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&tem.Cla2subId, &tem.ClaName, &tem.SubName, &tem.TecName)
			if err != nil {
				log.Printf("err:%s", err)
			}
			s.Cla2subList = append(s.Cla2subList, tem)
		}

		for i := 0; i < len(s.Cla2subList); i++ {
			//查班级名字
			rows, err = c.db.Model(&data_conn.ClassInfo{}).Where("ClassId=?", s.Cla2subList[i].ClaName).Select("ClassName").Rows()
			if err != nil {
				log.Printf("err:%s", err)
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].ClaName)
				if err != nil {
					log.Printf("err:%s", err)
				}
			}
			//查课程名字
			rows, err = c.db.Model(&data_conn.Subject{}).Where("SubId=?", s.Cla2subList[i].SubName).Select("SubName").Rows()
			if err != nil {
				log.Printf("err:%s", err)
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].SubName)
				if err != nil {
					log.Printf("err:%s", err)
				}
			}
			//查老师的名字
			rows, err := c.db.Model(&data_conn.TeacherInfo{}).Where("TecId=?", s.Cla2subList[i].TecName).Select("TecName").Rows()
			if err != nil {
				log.Printf("err:%s", err)
			}
			for rows.Next() {
				err = rows.Scan(&s.Cla2subList[i].TecName)
				if err != nil {
					log.Printf("err:%s", err)
				}
			}
		}
		s.IsSuccess = true
		render.JSON(w, r, s)
	}
}

func (c *Cla2subAPi) UpCla_Sub(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cla2subId := r.Form["cla2subId"][0]
	className := r.Form["className"][0]
	subName := r.Form["subName"][0]
	tecName := r.Form["tecName"][0]

	var classId, subId, tecId int

	if className != "" {
		rows, err := c.db.Model(&data_conn.ClassInfo{}).Where("ClassName=?", className).Select("ClassId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&classId)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}
		if classId == 0 {
			s := structure_type.Things{"本班级信息不存在", false}
			render.JSON(w, r, s)
			return
		}
		err = c.db.Model(&data_conn.Cla2sub{}).Where("Cla2subId=?", cla2subId).Update(&data_conn.Cla2sub{ClaId: classId}).Error
		if err != nil {
			log.Printf("err:%s", err)
		}
	}

	if className != "" {
		rows, err := c.db.Model(&data_conn.Subject{}).Where("SubName=?", subName).Select("SubId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&subId)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}
		if subId == 0 {
			s := structure_type.Things{"本课程信息不存在", false}
			render.JSON(w, r, s)
			return
		}
		err = c.db.Model(&data_conn.Cla2sub{}).Where("Cla2subId=?", cla2subId).Update(&data_conn.Cla2sub{SubId: subId}).Error
		if err != nil {
			log.Printf("err:%s", err)
		}
	}

	if tecName != "" {
		rows, err := c.db.Model(&data_conn.TeacherInfo{}).Where("TecName=?", tecName).Select("TecId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&tecId)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}
		if tecId == 0 {
			s := structure_type.Things{"本教师信息不存在", false}
			render.JSON(w, r, s)
			return
		}
		err = c.db.Model(&data_conn.Cla2sub{}).Where("Cla2sub_id=?", cla2subId).Update(&data_conn.Cla2sub{TecId: tecId}).Error
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	s := structure_type.Things{"编辑班级课程信息成功", true}
	render.JSON(w, r, s)
}

func (c *Cla2subAPi) DelCla_Sub(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cla2subId := r.Form["cla2subId"][0]

	err := c.db.Model(&data_conn.Cla2sub{}).Where("Cla2subId=?", cla2subId).Delete(&data_conn.Cla2sub{}).Error
	if err != nil {
		log.Printf("err:%s", err)
	}
	s := structure_type.Things{"删除班级课程信息成功", true}
	render.JSON(w, r, s)
}
