package structure_type

import "time"

type Things struct {
	Thing     string `json:"thing"`
	IsSuccess bool   `json:"isSuccess"`
}

type Major struct {
	MajId    int
	MajName  string
	MajPrin  string
	MajLink  string
	MajPhone string
}

type MajorTotal struct {
	MajorList []Major
}

type ClassInfo struct {
	ClassId   int
	ClassName string
	ClassTec  string
	MajId     int
}

type Class struct {
	ClassName string
	ClassTec  string
	MajName   string
}
type ClassTotal struct {
	ClassList []Class
}

type StudentInfo struct {
	StuId    int
	StuNu    string
	StuName  string
	StuSex   string
	StuBirth time.Time
	StuPic   string
	MajName  string
	ClaName  string
}
type StudentTotal struct {
	StudentList []StudentInfo
}

type Subject struct {
	SubId    int
	SubName  string
	SubType  string
	SubTimes int
}
type SubjectTotal struct {
	SubjectList []Subject
}

type TeacherInfo struct {
	TecId    int
	TecName  string
	TecSex   string
	TecBirth time.Time
	TecMajor string
	TecPhone string
}

type TeacherTotal struct {
	TeacherList []TeacherInfo
}

type Cla2sub struct {
	Cla2subId int
	ClaName   string
	SubName   string
	TecName   string
}
type Cla2subTotal struct {
	Cla2subList []Cla2sub
}

type UserInfo struct {
	UserId   int
	UserName string
	UserPwd  string
	RoleName string
}
type UserInfoTotal struct {
	UserInfoList []UserInfo
}

type Score struct {
	ScoId    int
	ScoDaily string
	SubExam  string
	WcoCount string
	StuName  string
	SubName  string
	ClaName  string
}

type ScoreTotal struct {
	ScoreList []Score
}
