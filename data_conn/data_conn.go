package data_conn

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

//学生信息
type StudentInfo struct {
	StuId    int       `gorm:"auto_increment;primary_key"` //学生id
	StuNu    string    `gorm:"not null"`                   //学生学号
	StuName  string    `gorm:"not null"`                   //学生名字
	StuSex   string    `gorm:"not null"`                   //学生性别
	StuBirth time.Time `gorm:"not null"`                   //学生生日
	StuPic   string    `gorm:"not null"`                   //学生照片
	ClaId    int       `gorm:"not null"`                   //班级id
}

//教师信息
type TeacherInfo struct {
	TecId    int       `gorm:"auto_increment;primary_key"` //教师id
	TecName  string    `gorm:"not null"`                   //教师名字
	TecSex   string    `gorm:"not null"`                   //教师性别
	TecBirth time.Time `gorm:"not null"`                   //教师生日
	TecMajor string    `gorm:"not null"`                   //教师专业
	TecPhone string    `gorm:"not null"`                   //教师电话
}

//班级表
type ClassInfo struct {
	ClassId   int    `gorm:"auto_increment;primary_key"` //班级id
	ClassName string `gorm:"not null"`                   //班级名字
	ClassTec  string `gorm:"not null"`                   //班级教师
	MajId     int    `gorm:"not null"`                   //专业id
}

//专业表
type Major struct {
	MajId    int    `gorm:"auto_increment;primary_key"` //专业id
	MajName  string `gorm:"not null"`                   //专业名称
	MajPrin  string `gorm:"not null"`                   //专业负责人
	MajLink  string `gorm:"not null"`                   //专业联系人
	MajPhone string `gorm:"not null"`                   //专业联系人电话
}

//课程信息表
type Subject struct {
	SubId    int    `gorm:"auto_increment;primary_key"` //科目id
	SubName  string `gorm:"not null"`                   //科目名称
	SubType  string `gorm:"not null"`                   //课程类型
	SubTimes int    `gorm:"not null"`                   //课时
}

//成绩表
type Score struct {
	ScoId     int    `gorm:"auto_increment;primary_key"` //成绩id
	ScoDaily  string `gorm:"not null"`                   //平时成绩
	SubExam   string `gorm:"not null"`                   //考试成绩
	WcoCount  string `gorm:"not null"`                   //总成绩
	StuId     int    `gorm:"not null"`                   //学生id
	SubId     int    `gorm:"not null"`                   //科目id
	Cla2subId int    `gorm:"not null"`                   //课程表id
	ClaId     int    `gorm:"not null"`                   //班级id
}

//课程表
type Cla2sub struct {
	Cla2subId int `gorm:"auto_increment;primary_key"` //课程表id
	ClaId     int `gorm:"not null"`                   //班级id
	SubId     int `gorm:"not null"`                   //科目id
	TecId     int `gorm:"not null"`                   //主讲老师id
}

//操作员表
type UserInfo struct {
	UserId   int    `gorm:"auto_increment;primary_key"` //用户id
	UserName string `gorm:"not null"`                   //用户名字
	UserPwd  string `gorm:"not null"`                   //用户密码密码
	RoleName string `gorm:"not null"`                   //角色名字
}

func DB_Mysql() *gorm.DB {
	db, err := gorm.Open("mysql", "root:root123@(127.0.0.1:3306)/stu?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("连接数据库失败")
	}
	// 自动迁移模式
	db.AutoMigrate(&StudentInfo{}, &TeacherInfo{}, &ClassInfo{}, &Major{}, &Subject{},
		&Score{}, &Cla2sub{}, &UserInfo{})
	return db
}
