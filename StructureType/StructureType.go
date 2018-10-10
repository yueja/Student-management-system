package StructureType

import "time"

type Things struct {
	Thing string `json:"err_1"`
}

type Major struct {
	Maj_id    int
	Maj_name  string
	Maj_prin  string
	Maj_link  string
	Maj_phone string
}

type Majortotal struct {
	MajorList []Major
}

type Classinfo struct {
	Class_id   int
	Class_name string
	Class_tec  string
	Maj_id     int
}

type Class struct {
	Class_name string
	Class_tec  string
	Maj_name   string
}
type Classtotal struct {
	ClassList []Class
}

type Studentinfo struct {
	Stu_id    int
	Ope_id    int
	Stu_nu    string
	Stu_name  string
	Stu_sex   string
	Stu_birth time.Time
	Stu_pic   string
	Maj_name   string
	Cla_name    string
}
type Studenttotal struct {
	StudentList []Studentinfo
}

type Subject struct {
	Sub_id    int
	Sub_name  string
	Sub_type  string
	Sub_times int
}
type Subjecttotal struct {
	SubjectList []Subject
}

type Teacherinfo struct {
	Tec_id    int
	Ope_id    int
	Tec_name  string
	Tec_sex   string
	Tec_birth time.Time
	Tec_major string
	Tec_phone string
}

type Teachertotal struct {
	TeacherList []Teacherinfo
}

type Cla2sub struct {
	Cla2sub_id int
	Cla_name   string
	Sub_name   string
	Tec_name   string
}
type Cla2subtotal struct {
	Cla2subList []Cla2sub
}

type Userinfo struct {
	User_id   int
	User_name string
	User_pwd  string
	Role_name string
}
type Userinfototal struct {
	UserinfoList []Userinfo
}


type Score struct {
	Sco_id     int
	Sco_daily  string
	Sub_exam   string
	Wco_count  string
	Stu_name    string
	Sub_name     string
	Cla_name    string
}

type Scoretotal struct {
	ScoreList [] Score
}