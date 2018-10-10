package ScoreMation

import (
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"xiangmu/Student/DataConn"
	"xiangmu/Student/StructureType"
	)

type ScoreAPi struct {
	db *gorm.DB
}

func Make_db(db *gorm.DB) *ScoreAPi {
	DB := &ScoreAPi{db}
	return DB
}
func (Score *ScoreAPi) AddScore(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var stu_id, sco_id, sub_id, cla_id, cla2sub_id int
	sco_daily := r.Form["Sco_daily"][0]
	sub_exam := r.Form["Sub_exam"][0]
	stu_name := r.Form["Stu_name"][0]
	sub_name := r.Form["Sub_name"][0]

	if sco_daily == "" || sub_exam == "" || stu_name == "" || sub_name == "" {
		s := StructureType.Things{"请将信息输入完整"}
		render.JSON(w, r, s)
		return
	}
	//判断学生成绩是否已经存在
	rows, err := Score.db.Table("Studentinfos").Where("Stu_name=?", stu_name).Select("Stu_id,Cla_id").Rows()

	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&stu_id, &cla_id)
		if err != nil {
			return
		}
	}
	rows, err = Score.db.Table("Scores").Where("Stu_id=?", stu_id).Select("Sco_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&sco_id)
		if err != nil {
			return
		}
	}
	if sco_id != 0 {
		s := StructureType.Things{"本学生成绩信息已存在"}
		render.JSON(w, r, s)
		return
	}
	//查询课程id
	rows, err = Score.db.Table("Subjects").Where("Sub_name=?", sub_name).Select("Sub_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&sub_id)
		if err != nil {
			return
		}
	}
	//查询课程表id
	rows, err = Score.db.Table("Cla2subs").Where("Cla_id=? and Sub_id=?", cla_id, sub_id).Select("Cla2sub_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&cla2sub_id)
		if err != nil {
			return
		}
	}
	//增加学生成绩信息
	sco, err := strconv.ParseFloat(sco_daily, 64)
	if err != nil {
		return
	}
	sub, err := strconv.ParseFloat(sub_exam, 64)
	if err != nil {
		return
	}
	wco_count := strconv.Itoa(int(sco*0.4 + sub*0.6))
	err = Score.db.Create(&DataConn.Score{Sco_daily: sco_daily, Sub_exam: sub_exam, Wco_count: wco_count, Stu_id: stu_id, Sub_id: sub_id, Cla2sub_id: cla2sub_id, Cla_id: cla_id}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"添加学生成绩信息成功"}
	render.JSON(w, r, s)
}

func (Score *ScoreAPi) BrowScore(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	class_name := r.Form["Class_name"][0]
	sub_name := r.Form["Sub_name"][0]
	stu_name := r.Form["Stu_name"][0]

	 s:=StructureType.Scoretotal{}
	 tem:=StructureType.Score{}

	var stu_id, sub_id, cla_id int

//按班级查询
	if class_name != ""{
		rows, err:= Score.db.Table("Classinfos").Where("Class_name=?", class_name).Select("Class_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&cla_id)
			if err != nil {
				return
			}
		}

		a:="Sco_id,Sco_daily,Sub_exam,Wco_count,Stu_id,Sub_id"
		rows, err = Score.db.Table("Scores").Where("Cla_id=?", cla_id).Select(a).Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Sco_id,&tem.Sco_daily,&tem.Sub_exam,&tem.Wco_count,&tem.Stu_name,&tem.Sub_name)
			if err != nil {
				return
			}
			tem.Cla_name=class_name
			s.ScoreList=append(s.ScoreList,tem)
		}

		for i:=0;i<len(s.ScoreList);i++{
			//查询科目名称
			rows, err = Score.db.Table("Subjects").Where("Sub_id=?",s.ScoreList[i].Sub_name).Select("Sub_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.ScoreList[i].Sub_name)
				if err != nil {
					return
				}
			}
			//查询学生名字
			rows, err = Score.db.Table("Studentinfos").Where("Stu_id=?",s.ScoreList[i].Stu_name).Select("Stu_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.ScoreList[i].Stu_name)
				if err != nil {
					return
				}
			}
		}
		render.JSON(w, r, s)
	}

//按科目查询
	if sub_name != ""{
		rows, err:= Score.db.Table("Subjects").Where("Sub_name=?", sub_name).Select("Sub_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&sub_id)
			if err != nil {
				return
			}
		}

		a:="Sco_id,Sco_daily,Sub_exam,Wco_count,Stu_id,Cla_id"
		rows, err = Score.db.Table("Scores").Where("Sub_id=?", sub_id).Select(a).Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Sco_id,&tem.Sco_daily,&tem.Sub_exam,&tem.Wco_count,&tem.Stu_name,&tem.Cla_name)
			if err != nil {
				return
			}
			tem.Sub_name=sub_name
			s.ScoreList=append(s.ScoreList,tem)
		}

		for i:=0;i<len(s.ScoreList);i++{
			//查询班级名称
			rows, err = Score.db.Table("Classinfos").Where("Class_id=?",s.ScoreList[i].Cla_name).Select("Class_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.ScoreList[i].Cla_name)
				if err != nil {
					return
				}
			}
			//查询学生名字
			rows, err = Score.db.Table("Studentinfos").Where("Stu_id=?",s.ScoreList[i].Stu_name).Select("Stu_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.ScoreList[i].Stu_name)
				if err != nil {
					return
				}
			}
		}
		render.JSON(w, r, s)
	}

//按学生查询
	if stu_name != "" {
		rows, err := Score.db.Table("Studentinfos").Where("Stu_name=?", stu_name).Select("Stu_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&stu_id)
			if err != nil {
				return
			}
		}

		a := "Sco_id,Sco_daily,Sub_exam,Wco_count,Sub_id,Cla_id"
		rows, err = Score.db.Table("Scores").Where("Stu_id=?", stu_id).Select(a).Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Sco_id, &tem.Sco_daily, &tem.Sub_exam, &tem.Wco_count, &tem.Sub_name, &tem.Cla_name)
			if err != nil {
				return
			}
			tem.Stu_name = stu_name
			s.ScoreList = append(s.ScoreList, tem)
		}

		for i := 0; i < len(s.ScoreList); i++ {
			//查询班级名称
			rows, err = Score.db.Table("Classinfos").Where("Class_id=?", s.ScoreList[i].Cla_name).Select("Class_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.ScoreList[i].Cla_name)
				if err != nil {
					return
				}
			}
			//查询科目名字
			rows, err = Score.db.Table("Subjects").Where("Sub_id=?", s.ScoreList[i].Sub_name).Select("Sub_name").Rows()
			if err != nil {
				return
			}
			for rows.Next() {
				err = rows.Scan(&s.ScoreList[i].Sub_name)
				if err != nil {
					return
				}
			}
		}
		render.JSON(w, r, s)
	}
}

func (Score *ScoreAPi) UpScore(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	sco_id := r.Form["Sco_id"][0]
	sco_daily := r.Form["Sco_daily"][0]
	sub_exam := r.Form["Sub_exam"][0]
	stu_name := r.Form["Stu_name"][0]
	sub_name := r.Form["Sub_name"][0]

	var sco_daily_1,sub_exam_1 string
	var cla2sub_id,stu_id,cla_id,sub_id int
	//更新平时成绩
	if sco_daily!=""{
		//查考试成绩
		rows, err := Score.db.Table("Scores").Where("Sco_id=?",sco_id).Select("Sub_exam").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&sub_exam_1)
			if err != nil {
				return
			}
		}
		//计算总成绩
		sco, err := strconv.ParseFloat(sco_daily, 64)
		if err != nil {
			return
		}
		sub, err := strconv.ParseFloat(sub_exam_1, 64)
		if err != nil {
			return
		}
		wco_count := strconv.Itoa(int(sco*0.4 + sub*0.6))
		//更新平时成绩和总成绩
		err = Score.db.Model(&DataConn.Score{}).Where("Sco_id=?",sco_id).Update(&DataConn.Score{Sco_daily: sco_daily,Wco_count:wco_count}).Error
		if err != nil {
			return
		}
	}
	//更新考试成绩
	if sub_exam!=""{
		//查平时成绩
		rows, err := Score.db.Table("Scores").Where("Sco_id=?",sco_id).Select("Sco_daily").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&sco_daily_1)
			if err != nil {
				return
			}
		}
		//计算总成绩
		sco, err := strconv.ParseFloat(sco_daily_1, 64)
		if err != nil {
			return
		}
		sub, err := strconv.ParseFloat(sub_exam, 64)
		if err != nil {
			return
		}
		wco_count := strconv.Itoa(int(sco*0.4 + sub*0.6))
		//更新平时成绩和总成绩
		err = Score.db.Model(&DataConn.Score{}).Where("Sco_id=?",sco_id).Update(&DataConn.Score{Sub_exam: sub_exam,Wco_count:wco_count}).Error
		if err != nil {
			return
		}
	}
	//更新学生
	if stu_name!=""{
		rows, err := Score.db.Table("Studentinfos").Where("Stu_name=?",stu_name).Select("Stu_id,Cla_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&stu_id,&cla_id)
			if err != nil {
				return
			}
		}

		rows, err = Score.db.Table("Scores").Where("Sco_id=?",sco_id).Select("Sub_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&sub_id)
			if err != nil {
				return
			}
		}

		rows, err = Score.db.Table("Cla2sub").Where("Sub_id=? and Cla_id ",sub_id,cla_id).Select("Cla2sub_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&cla2sub_id)
			if err != nil {
				return
			}
		}

		err = Score.db.Model(&DataConn.Score{}).Where("Sco_id=?",sco_id).Update(&DataConn.Score{Stu_id: stu_id,Cla_id:cla_id,Cla2sub_id:cla2sub_id}).Error
		if err != nil {
			return
		}
	}
	//更新科目
	if sub_name!=""{
		rows, err := Score.db.Table("Subjects").Where("Sub_name=?",sub_name).Select("Sub_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&sub_id)
			if err != nil {
				return
			}
		}

		rows, err = Score.db.Table("Scores").Where("Sco_id=?",sco_id).Select("Cla_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&cla_id)
			if err != nil {
				return
			}
		}

		rows, err = Score.db.Table("Cla2sub").Where("Sub_id=? and Cla_id ",sub_id,cla_id).Select("Cla2sub_id").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&cla2sub_id)
			if err != nil {
				return
			}
		}

		err = Score.db.Model(&DataConn.Score{}).Where("Sco_id=?",sco_id).Update(&DataConn.Score{Sub_id: sub_id,Cla2sub_id:cla2sub_id}).Error
		if err != nil {
			return
		}
	}
	s := StructureType.Things{"更新成绩信息成功"}
	render.JSON(w, r, s)
}

func (score *ScoreAPi) SubScore(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sco_id := r.Form["Sco_id"][0]
	err := score.db.Model(&DataConn.Score{}).Where("Sco_id=?", sco_id).Delete(&DataConn.Score{}).Error
	if err != nil {
		s := StructureType.Things{"删除学生成绩信息失败"}
		render.JSON(w, r, s)
		return
	}
	s := StructureType.Things{"删除学生成绩信息成功"}
	render.JSON(w, r, s)
}
