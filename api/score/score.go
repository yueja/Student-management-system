package score

import (
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
	"xiangmu/Student/data_conn"
	"xiangmu/Student/structure_type"
)

type ScoreAPi struct {
	db *gorm.DB
}

func MakeDb(db *gorm.DB) *ScoreAPi {
	DB := &ScoreAPi{db}
	return DB
}
func (Score *ScoreAPi) AddScore(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	scoDaily := r.Form["scoDaily"][0]
	subExam := r.Form["subExam"][0]
	stuName := r.Form["stuName"][0]
	subName := r.Form["subName"][0]

	var stuId, scoId, subId, claId, cla2subId int
	if scoDaily == "" || subExam == "" || stuName == "" || subName == "" {
		s := structure_type.Things{"请将信息输入完整", false}
		render.JSON(w, r, s)
		return
	}
	//判断学生成绩是否已经存在
	rows, err := Score.db.Model(&data_conn.StudentInfo{}).Where("StuName=?", stuName).Select("StuId,ClaId").Rows()

	if err != nil {
		log.Printf("err:%s", err)
	}
	for rows.Next() {
		err = rows.Scan(&stuId, &claId)
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	rows, err = Score.db.Model(&data_conn.Score{}).Where("StuId=?", stuId).Select("ScoId").Rows()
	if err != nil {
		log.Printf("err:%s", err)
	}
	for rows.Next() {
		err = rows.Scan(&scoId)
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	if scoId != 0 {
		s := structure_type.Things{"本学生成绩信息已存在", false}
		render.JSON(w, r, s)
		return
	}
	//查询课程id
	rows, err = Score.db.Model(&data_conn.Subject{}).Where("SubName=?", subName).Select("SubId").Rows()
	if err != nil {
		log.Printf("err:%s", err)
	}
	for rows.Next() {
		err = rows.Scan(&subId)
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	//查询课程表id
	rows, err = Score.db.Model(&data_conn.Cla2sub{}).Where("ClaId=? and SubId=?", claId, subId).Select("Cla2subId").Rows()
	if err != nil {
		log.Printf("err:%s", err)
	}
	for rows.Next() {
		err = rows.Scan(&cla2subId)
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	//增加学生成绩信息
	sco, err := strconv.ParseFloat(scoDaily, 64)
	if err != nil {
		log.Printf("err:%s", err)
	}
	sub, err := strconv.ParseFloat(subExam, 64)
	if err != nil {
		log.Printf("err:%s", err)
	}
	wcoCount := strconv.Itoa(int(sco*0.4 + sub*0.6))
	err = Score.db.Create(&data_conn.Score{ScoDaily: scoDaily, SubExam: subExam, WcoCount: wcoCount,
		StuId: stuId, SubId: subId, Cla2subId: cla2subId, ClaId: claId}).Error
	if err != nil {
		log.Printf("err:%s", err)
	}
	s := structure_type.Things{"添加学生成绩信息成功", true}
	render.JSON(w, r, s)
}

func (Score *ScoreAPi) BrowScore(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	className := r.Form["className"][0]
	subName := r.Form["subName"][0]
	stuName := r.Form["stuName"][0]

	s := structure_type.ScoreTotal{}
	tem := structure_type.Score{}

	var stuId, subId, claId int

	//按班级查询
	if className != "" {
		rows, err := Score.db.Model(&data_conn.ClassInfo{}).Where("ClassName=?", className).Select("ClassId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&claId)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}

		a := "ScoId,ScoDaily,SubExam,WcoCount,StuId,SubId"
		rows, err = Score.db.Model(&data_conn.Score{}).Where("ClaId=?", claId).Select(a).Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&tem.ScoId, &tem.ScoDaily, &tem.SubExam, &tem.WcoCount, &tem.StuName, &tem.SubName)
			if err != nil {
				log.Printf("err:%s", err)
			}
			tem.ClaName = className
			s.ScoreList = append(s.ScoreList, tem)
		}

		for i := 0; i < len(s.ScoreList); i++ {
			//查询科目名称
			rows, err = Score.db.Model(&data_conn.Subject{}).Where("SubId=?", s.ScoreList[i].SubName).Select("SubName").Rows()
			if err != nil {
				log.Printf("err:%s", err)
			}
			for rows.Next() {
				err = rows.Scan(&s.ScoreList[i].SubName)
				if err != nil {
					log.Printf("err:%s", err)
				}
			}
			//查询学生名字
			rows, err = Score.db.Model(&data_conn.StudentInfo{}).Where("StuId=?", s.ScoreList[i].StuName).Select("StuName").Rows()
			if err != nil {
				log.Printf("err:%s", err)
			}
			for rows.Next() {
				err = rows.Scan(&s.ScoreList[i].StuName)
				if err != nil {
					log.Printf("err:%s", err)
				}
			}
		}
		s.IsSuccess = true
		render.JSON(w, r, s)
	}

	//按科目查询
	if subName != "" {
		rows, err := Score.db.Model(&data_conn.Subject{}).Where("SubName=?", subName).Select("SubId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&subId)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}

		a := "ScoId,ScoDaily,SubExam,WcoCount,StuId,ClaId"
		rows, err = Score.db.Model(&data_conn.Score{}).Where("SubId=?", subId).Select(a).Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&tem.ScoId, &tem.ScoDaily, &tem.SubExam, &tem.WcoCount, &tem.StuName, &tem.ClaName)
			if err != nil {
				log.Printf("err:%s", err)
			}
			tem.SubName = subName
			s.ScoreList = append(s.ScoreList, tem)
		}

		for i := 0; i < len(s.ScoreList); i++ {
			//查询班级名称
			rows, err = Score.db.Model(&data_conn.ClassInfo{}).Where("ClassId=?", s.ScoreList[i].ClaName).Select("ClassName").Rows()
			if err != nil {
				log.Printf("err:%s", err)
			}
			for rows.Next() {
				err = rows.Scan(&s.ScoreList[i].ClaName)
				if err != nil {
					log.Printf("err:%s", err)
				}
			}
			//查询学生名字
			rows, err = Score.db.Model(&data_conn.StudentInfo{}).Where("StuId=?", s.ScoreList[i].StuName).Select("StuName").Rows()
			if err != nil {
				log.Printf("err:%s", err)
			}
			for rows.Next() {
				err = rows.Scan(&s.ScoreList[i].StuName)
				if err != nil {
					log.Printf("err:%s", err)
				}
			}
		}
		s.IsSuccess = true
		render.JSON(w, r, s)
	}

	//按学生查询
	if stuName != "" {
		rows, err := Score.db.Model(&data_conn.StudentInfo{}).Where("StuName=?", stuName).Select("StuId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&stuId)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}

		a := "ScoId,ScoDaily,SubExam,WcoCount,SubId,ClaId"
		rows, err = Score.db.Model(&data_conn.Score{}).Where("StuId=?", stuId).Select(a).Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&tem.ScoId, &tem.ScoDaily, &tem.SubExam, &tem.WcoCount, &tem.SubName, &tem.ClaName)
			if err != nil {
				log.Printf("err:%s", err)
			}
			tem.StuName = stuName
			s.ScoreList = append(s.ScoreList, tem)
		}

		for i := 0; i < len(s.ScoreList); i++ {
			//查询班级名称
			rows, err = Score.db.Model(&data_conn.ClassInfo{}).Where("ClassId=?", s.ScoreList[i].ClaName).Select("ClassName").Rows()
			if err != nil {
				log.Printf("err:%s", err)
			}
			for rows.Next() {
				err = rows.Scan(&s.ScoreList[i].ClaName)
				if err != nil {
					log.Printf("err:%s", err)
				}
			}
			//查询科目名字
			rows, err = Score.db.Model(&data_conn.Subject{}).Where("SubId=?", s.ScoreList[i].SubName).Select("SubName").Rows()
			if err != nil {
				log.Printf("err:%s", err)
			}
			for rows.Next() {
				err = rows.Scan(&s.ScoreList[i].SubName)
				if err != nil {
					log.Printf("err:%s", err)
				}
			}
		}
		s.IsSuccess = true
		render.JSON(w, r, s)
	}
}

func (Score *ScoreAPi) UpScore(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	scoId := r.Form["scoId"][0]
	scoDaily := r.Form["scoDaily"][0]
	subExam := r.Form["subExam"][0]
	stuName := r.Form["stuName"][0]
	subName := r.Form["subName"][0]

	var scoDaily_1, subExam_1 string
	var cla2subId, stuId, claId, subId int

	//更新平时成绩
	if scoDaily != "" {
		//查考试成绩
		rows, err := Score.db.Model(&data_conn.Score{}).Where("ScoId=?", scoId).Select("SubExam").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&subExam_1)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}
		//计算总成绩
		sco, err := strconv.ParseFloat(scoDaily, 64)
		if err != nil {
			log.Printf("err:%s", err)
		}
		sub, err := strconv.ParseFloat(subExam_1, 64)
		if err != nil {
			log.Printf("err:%s", err)
		}
		wcoCount := strconv.Itoa(int(sco*0.4 + sub*0.6))
		//更新平时成绩和总成绩
		err = Score.db.Model(&data_conn.Score{}).Where("ScoId=?", scoId).Update(&data_conn.Score{ScoDaily: scoDaily, WcoCount: wcoCount}).Error
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	//更新考试成绩
	if subExam != "" {
		//查平时成绩
		rows, err := Score.db.Model(&data_conn.Score{}).Where("ScoId=?", scoId).Select("ScoDaily").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&scoDaily_1)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}
		//计算总成绩
		sco, err := strconv.ParseFloat(scoDaily_1, 64)
		if err != nil {
			log.Printf("err:%s", err)
		}
		sub, err := strconv.ParseFloat(subExam, 64)
		if err != nil {
			log.Printf("err:%s", err)
		}
		wcoCount := strconv.Itoa(int(sco*0.4 + sub*0.6))
		//更新平时成绩和总成绩
		err = Score.db.Model(&data_conn.Score{}).Where("ScoId=?", scoId).Update(&data_conn.Score{SubExam: subExam, WcoCount: wcoCount}).Error
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	//更新学生
	if stuName != "" {
		rows, err := Score.db.Model(&data_conn.StudentInfo{}).Where("StuName=?", stuName).Select("StuId,ClaId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&stuId, &claId)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}

		rows, err = Score.db.Model(&data_conn.Score{}).Where("ScoId=?", scoId).Select("SubId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&subId)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}

		rows, err = Score.db.Model(&data_conn.Cla2sub{}).Where("SubId=? and ClaId ", subId, claId).Select("Cla2subId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&cla2subId)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}

		err = Score.db.Model(&data_conn.Score{}).Where("ScoId=?", scoId).Update(&data_conn.Score{StuId: stuId, ClaId: claId, Cla2subId: cla2subId}).Error
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	//更新科目
	if subName != "" {
		rows, err := Score.db.Model(&data_conn.Subject{}).Where("SubName=?", subName).Select("SubId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&subId)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}

		rows, err = Score.db.Model(&data_conn.Score{}).Where("ScoId=?", scoId).Select("ClaId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&claId)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}

		rows, err = Score.db.Model(&data_conn.Cla2sub{}).Where("SubId=? and ClaId ", subId, claId).Select("Cla2subId").Rows()
		if err != nil {
			log.Printf("err:%s", err)
		}
		for rows.Next() {
			err = rows.Scan(&cla2subId)
			if err != nil {
				log.Printf("err:%s", err)
			}
		}

		err = Score.db.Model(&data_conn.Score{}).Where("ScoId=?", scoId).Update(&data_conn.Score{SubId: subId, Cla2subId: cla2subId}).Error
		if err != nil {
			log.Printf("err:%s", err)
		}
	}
	s := structure_type.Things{"更新成绩信息成功", true}
	render.JSON(w, r, s)
}

func (score *ScoreAPi) DelScore(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	scoId := r.Form["scoId"][0]
	err := score.db.Model(&data_conn.Score{}).Where("ScoId=?", scoId).Delete(&data_conn.Score{}).Error
	if err != nil {
		log.Printf("err:%s", err)
	}
	s := structure_type.Things{"删除学生成绩信息成功", true}
	render.JSON(w, r, s)
}
