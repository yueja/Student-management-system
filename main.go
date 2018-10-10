package main

import (
	"github.com/go-chi/chi"
	"net/http"
	"xiangmu/Student/API/ClassMation"
	"xiangmu/Student/API/ClassSubjectMation"
	"xiangmu/Student/API/OperatorMation"
	"xiangmu/Student/API/ScoreMation"
	"xiangmu/Student/API/StudentMation"
	"xiangmu/Student/API/SubjectMation"
	"xiangmu/Student/API/TeacherMation"
	"xiangmu/Student/DataConn"
	"xiangmu/Student/StudentManagement"
)

func main() {
	DB := DataConn.DB_Mysql()
	defer DB.Close()
	r := chi.NewRouter()

	//用户管理
	User := UserMation.Make_db(DB)
	r.Route("/user", func(r chi.Router) {
		r.Post("/registerstuuser", User.RegisterStuUser) //注册学生用户
		r.Post("/registerteauser", User.RegisterTeaUser) //注册教师用户
		r.Post("/loginstuuser", User.LoginStuUser)       //学生用户登录
		r.Post("/loginteauser", User.LoginTeaUser)       //教师用户登录
		r.Post("/userpwdmodify", User.UserPwdModify)     //用户信息修改
		r.Get("/browuser", User.BrowUser)                //查看所有用户（按角色、全部）
		r.Delete("/subuser", User.SubUser)               //删除用户信息
		//r.Post("/RegisterAdmini", User.RegisterAdmini) //申请为管理员
	})

	//专业信息管理
	Major := ProInformation.Make_db(DB)
	r.Route("/major", func(r chi.Router) {
		r.Post("/addmajor", Major.AddMajor)   //增加专业
		r.Get("/browmajor", Major.BrowMajor)  //浏览所有专业
		r.Post("/upmajor", Major.UpMajor)     //编辑专业
		r.Delete("/submajor", Major.SubMajor) //删除专业
	})

	//班级信息管理
	Class := ClassMation.Make_db(DB)
	r.Route("/class", func(r chi.Router) {
		r.Post("/addclass", Class.AddClass)   //增加班级
		r.Get("/browclass", Class.BrowClass)  //浏览班级
		r.Post("/upclass", Class.UpClass)     //编辑班级
		r.Delete("/subclass", Class.SubClass) //删除班级
	})

	//学生信息管理
	Student := StudentMation.Make_db(DB)
	r.Route("/student", func(r chi.Router) {
		r.Post("/addstudent", Student.AddStudent)   //增加学生信息
		r.Get("/browstudent", Student.BrowStudent)  //浏览学生信息（按专业、班级、个人、全部学生）
		r.Post("/upstudent", Student.UpClass)       //编辑学生信息
		r.Delete("/substudent", Student.SubStudent) //删除学生信息
	})

	//课程信息管理
	Subject := SubjectMation.Make_db(DB)
	r.Route("/subject", func(r chi.Router){
		r.Post("/addsubject", Subject.AddSubject)   //增加课程信息
		r.Get("/browsubject", Subject.BrowSubject)  //浏览课程信息（按名字、类型、全部）
		r.Post("/upsubject", Subject.UpSubject)     //编辑课程信息
		r.Delete("/subsubject", Subject.SubSubject) //删除课程信息
	})

	//教师信息管理
	Teacher := TeacherMation.Make_db(DB)
	r.Route("/teacher", func(r chi.Router) {
		r.Post("/addteacher", Teacher.AddTeacher)   //增加教师信息
		r.Get("/browteacher", Teacher.BrowTeacher)  //浏览教师信息（按专业、个人、全部）
		r.Post("/upteacher", Teacher.UpTeacher)     //编辑教师信息
		r.Delete("/subteacher", Teacher.SubTeacher) //删除教师信息
	})

	//班级课程管理
	ClassSubject := ClassSubjectMation.Make_db(DB)
	r.Route("/class_subject", func(r chi.Router) {
		r.Post("/addcla_sub", ClassSubject.AddCla_Sub)   //增加班级课程信息
		r.Get("/browcla_sub", ClassSubject.BrowCla_Sub)  //浏览班级课程信息（按班级、课程、老师、所有）
		r.Post("/upcla_sub", ClassSubject.UpCla_Sub)     //编辑班级课程信息
		r.Delete("/subcla_sub", ClassSubject.SubCla_Sub) //删除班级课程信息
	})

	//成绩管理
	Score := ScoreMation.Make_db(DB)
	r.Route("/score", func(r chi.Router) {
		r.Post("/addscore", Score.AddScore)   //增加学生成绩信息
		r.Get("/browscore", Score.BrowScore)  //浏览学生成绩信息（按班级、课程、学生查看）
		r.Post("/upscore", Score.UpScore)     //编辑学生成绩信息
		r.Delete("/subscore", Score.SubScore) //删除学生成绩信息
	})

	//查询报表

	// 绑定端口是8888
	http.ListenAndServe(":8888", r)
}

/*
	// 注册一个默认的路由器
	router := gin.Default()

	//专业信息管理
	Major:= ProInformation.Make_db(DB)
	router.Group("/major")
	{
		router.POST("/addmajor",Major.AddMajor)   //增加专业
		router.GET("/browmajor",Major.BrowMajor)   //浏览所有专业
		//router.POST("/upmajor",Major.UpMajor)   //编辑专业
		//router.POST("/submajor",Major.SubMajor)   //删除专业
	}
	// 绑定端口是8888
	router.Run(":6666")
}*/
