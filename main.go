package main

import (
	"flag"
	"github.com/go-chi/chi"
	"net/http"
	"xiangmu/Student/api/class"
	"xiangmu/Student/api/class_subject"
	"xiangmu/Student/api/major"
	"xiangmu/Student/api/score"
	"xiangmu/Student/api/student"
	"xiangmu/Student/api/subject"
	"xiangmu/Student/api/teacher"
	"xiangmu/Student/api/user"
	"xiangmu/Student/data_conn"
)

func main() {
	DB := data_conn.DB_Mysql()
	defer DB.Close()
	r := chi.NewRouter()

	//用户管理
	users := user.Make_db(DB)
	r.Route("/user", func(r chi.Router) {
		r.Post("/register_stu", users.RegisterStuUser) //注册学生用户
		r.Post("/register_tea", users.RegisterTeaUser) //注册教师用户
		r.Post("/login_stu", users.LoginStuUser)       //学生用户登录
		r.Post("/login_tea", users.LoginTeaUser)       //教师用户登录
		r.Post("/pwd_modify", users.UserPwdModify)     //用户信息修改
		r.Get("/browse", users.BrowUser)               //查看所有用户（按角色、全部）
		r.Delete("/delete", users.DelUser)             //删除用户信息
	})

	//专业信息管理
	majors := major.Make_db(DB)
	r.Route("/major", func(r chi.Router) {
		r.Post("/addition", majors.AddMajor) //增加专业
		r.Get("/browse", majors.BrowMajor)   //浏览所有专业
		r.Post("/updata", majors.UpMajor)    //编辑专业
		r.Delete("/delete", majors.DelMajor) //删除专业
	})

	//班级信息管理
	classes := class.Make_db(DB)
	r.Route("/class", func(r chi.Router) {
		r.Post("/addition", classes.AddClass) //增加班级
		r.Get("/browse", classes.BrowClass)   //浏览班级
		r.Post("/updata", classes.UpClass)    //编辑班级
		r.Delete("/delete", classes.DelClass) //删除班级
	})

	//学生信息管理
	Students := student.Make_db(DB)
	r.Route("/student", func(r chi.Router) {
		r.Post("/addition", Students.AddStudent) //增加学生信息
		r.Get("/browse", Students.BrowStudent)   //浏览学生信息（按专业、班级、个人、全部学生）
		r.Post("/updata", Students.UpClass)      //编辑学生信息
		r.Delete("/delete", Students.DelStudent) //删除学生信息
	})

	//课程信息管理
	Subject := subject.Make_db(DB)
	r.Route("/subject", func(r chi.Router) {
		r.Post("/addition", Subject.AddSubject) //增加课程信息
		r.Get("/browse", Subject.BrowSubject)   //浏览课程信息（按名字、类型、全部）
		r.Post("/updata", Subject.UpSubject)    //编辑课程信息
		r.Delete("/delete", Subject.DelSubject) //删除课程信息
	})

	//教师信息管理
	teachers := teacher.Make_db(DB)
	r.Route("/teacher", func(r chi.Router) {
		r.Post("/addition", teachers.AddTeacher) //增加教师信息
		r.Get("/browse", teachers.BrowTeacher)   //浏览教师信息（按专业、个人、全部）
		r.Post("/updata", teachers.UpTeacher)    //编辑教师信息
		r.Delete("/delete", teachers.DelTeacher) //删除教师信息
	})

	//班级课程管理
	class_subjects := class_subject.Make_db(DB)
	r.Route("/class_subject", func(r chi.Router) {
		r.Post("/addition", class_subjects.AddCla_Sub) //增加班级课程信息
		r.Get("/browse", class_subjects.BrowCla_Sub)   //浏览班级课程信息（按班级、课程、老师、所有）
		r.Post("/updata", class_subjects.UpCla_Sub)    //编辑班级课程信息
		r.Delete("/delete", class_subjects.DelCla_Sub) //删除班级课程信息
	})

	//成绩管理
	scores := score.Make_db(DB)
	r.Route("/score", func(r chi.Router) {
		r.Post("/addition", scores.AddScore) //增加学生成绩信息
		r.Get("/browse", scores.BrowScore)   //浏览学生成绩信息（按班级、课程、学生查看）
		r.Post("/updata", scores.UpScore)    //编辑学生成绩信息
		r.Delete("/delete", scores.DelScore) //删除学生成绩信息
	})

	// 默认端口是8888
	address := flag.String("address", ":8888", "")
	flag.Parse()
	http.ListenAndServe(*address, r)
}
