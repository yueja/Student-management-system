package DataConn

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

//学生信息
type Studentinfo struct {
	Stu_id    int       `gorm:"auto_increment;primary_key"` //学生id
	Ope_id    int       `gorm:"not null"`                   //操作员id
	Stu_nu    string    `gorm:"not null"`                   //学生学号
	Stu_name  string    `gorm:"not null"`                   //学生名字
	Stu_sex   string    `gorm:"not null"`                   //学生性别
	Stu_birth time.Time `gorm:"not null"`                   //学生生日
	Stu_pic   string    `gorm:"not null"`                   //学生照片
	Cla_id    int       `gorm:"not null"`                   //班级id
}

//教师信息
type Teacherinfo struct {
	Tec_id    int       `gorm:"auto_increment;primary_key"` //教师id
	Ope_id    int       `gorm:"not null"`                   //操作员id
	Tec_name  string    `gorm:"not null"`                   //教师名字
	Tec_sex   string    `gorm:"not null"`                   //教师性别
	Tec_birth time.Time `gorm:"not null"`                   //教师生日
	Tec_major string    `gorm:"not null"`                   //教师专业
	Tec_phone string    `gorm:"not null"`                   //教师电话
}

//班级表
type Classinfo struct {
	Class_id   int    `gorm:"auto_increment;primary_key"` //班级id
	Class_name string `gorm:"not null"`                   //班级名字
	Class_tec  string `gorm:"not null"`                   //班级教师
	Maj_id     int    `gorm:"not null"`                   //专业id
}

//专业表
type Major struct {
	Maj_id    int    `gorm:"auto_increment;primary_key"` //专业id
	Maj_name  string `gorm:"not null"`                   //专业名称
	Maj_prin  string `gorm:"not null"`                   //专业负责人
	Maj_link  string `gorm:"not null"`                   //专业联系人
	Maj_phone string `gorm:"not null"`                   //专业联系人电话
}

//课程信息表
type Subject struct {
	Sub_id    int    `gorm:"auto_increment;primary_key"` //科目id
	Sub_name  string `gorm:"not null"`                   //科目名称
	Sub_type  string `gorm:"not null"`                   //课程类型
	Sub_times int    `gorm:"not null"`                   //课时
}

//成绩表
type Score struct {
	Sco_id     int    `gorm:"auto_increment;primary_key"` //成绩id
	Sco_daily  string `gorm:"not null"`                   //平时成绩
	Sub_exam   string `gorm:"not null"`                   //考试成绩
	Wco_count  string `gorm:"not null"`                   //总成绩
	Stu_id     int    `gorm:"not null"`                   //学生id
	Sub_id     int    `gorm:"not null"`                   //科目id
	Cla2sub_id int    `gorm:"not null"`                   //课程表id
	Cla_id     int    `gorm:"not null"`                   //班级id
}

//课程表
type Cla2sub struct {
	Cla2sub_id int `gorm:"auto_increment;primary_key"` //课程表id
	Cla_id     int `gorm:"not null"`                   //班级id
	Sub_id     int `gorm:"not null"`                   //科目id
	Tec_id     int `gorm:"not null"`                   //主讲老师id
}

//功能表
type Privilege struct {
	Pri_id    int    `gorm:"auto_increment;primary_key"` //功能id
	Pri_name  string `gorm:"not null"`                   //模块名称
	Pri_url   string `gorm:"not null"`                   //模块链接
	Menu_name string `gorm:"not null"`                   //菜单名称
	Rol_id    int    `gorm:"not null"`                   //角色id
}

//操作员表
type Userinfo struct {
	User_id   int    `gorm:"auto_increment;primary_key"` //用户id
	User_name string `gorm:"not null"`                   //用户名字
	User_pwd  string `gorm:"not null"`                   //用户密码密码
	Role_name string `gorm:"not null"`                   //角色名字
}

func DB_Mysql() *gorm.DB {
	db, err := gorm.Open("mysql", "root:root123@(127.0.0.1:3306)/stu?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("连接数据库失败")
	}
	// 自动迁移模式
	db.AutoMigrate(&Studentinfo{}, &Teacherinfo{}, &Classinfo{}, &Major{}, &Subject{},
		&Score{}, &Cla2sub{}, &Privilege{}, &Userinfo{})
	return db
}
