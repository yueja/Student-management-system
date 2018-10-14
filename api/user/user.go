package user

import (
	"database/sql"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"study0/StructureType"
	"xiangmu/Student/data_conn"
	"xiangmu/Student/structure_type"
)

type UserAPi struct {
	db *gorm.DB
}

func MakeDb(db *gorm.DB) *UserAPi {
	DB := &UserAPi{db}
	return DB
}

func (user *UserAPi) RegisterStuUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userNu := r.Form["stuNu"][0]
	userPwd := r.Form["stuPwd"][0]
	var stuId, userId int

	if userNu == "" || userPwd == "" {
		s := structure_type.Things{"请将信息输入完整", false}
		render.JSON(w, r, s)
		return
	}
	//判断该学号是否已经注册过
	rows, err := user.db.Table("UserInfos").Where("UserName=?", userNu).Select("UserId").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			return
		}
	}
	if userId != 0 {
		s := structure_type.Things{"该学号已注册", false}
		render.JSON(w, r, s)
		return
	}

	//判断学号是否存在
	rows, err = user.db.Table("StudentInfos").Where("StuNu=?", userNu).Select("Stu_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&stuId)
		if err != nil {
			return
		}
	}
	if stuId == 0 {
		s := StructureType.Things{"该学号不存在"}
		render.JSON(w, r, s)
		return
	}

	err = user.db.Create(&data_conn.UserInfo{UserName: userNu, UserPwd: userPwd, RoleName: "学生"}).Error
	if err != nil {
		return
	}
	s := structure_type.Things{"注册学生用户信息成功", false}
	render.JSON(w, r, s)
}

func (user *UserAPi) RegisterTeaUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userName := r.Form["tecName"][0]
	userPwd := r.Form["tecPaw"][0]
	var userId, tecId int

	if userName == "" || userPwd == "" {
		s := structure_type.Things{"请将信息输入完整", false}
		render.JSON(w, r, s)
		return
	}

	//判断该教师是否已经注册过
	rows, err := user.db.Table("UserInfos").Where("UserName=?", userName).Select("UserId").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			return
		}
	}
	if userId != 0 {
		s := structure_type.Things{"该教师信息已注册", false}
		render.JSON(w, r, s)
		return
	}

	rows, err = user.db.Table("TeacherInfos").Where("TecName=?", userName).Select("TecId").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&tecId)
		if err != nil {
			return
		}
	}
	if tecId == 0 {
		s := structure_type.Things{"该教师名字不存在", false}
		render.JSON(w, r, s)
		return
	}

	err = user.db.Create(&data_conn.UserInfo{UserName: userName, UserPwd: userPwd, RoleName: "教师"}).Error
	if err != nil {
		return
	}
	s := structure_type.Things{"注册教师用户信息成功", true}
	render.JSON(w, r, s)
}

func (user *UserAPi) LoginStuUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userName := r.Form["userName"][0]
	userPwd := r.Form["userPwd"][0]
	var userPwd_1 string

	if userName == "" || userPwd == "" {
		s := structure_type.Things{"请将信息输入完整", false}
		render.JSON(w, r, s)
		return
	}

	rows, err := user.db.Model(&data_conn.UserInfo{}).Where("UserName=? and RoleName=?", userName, "学生").Select("UserPwd").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&userPwd_1)
		if err != nil {
			return
		}
	}
	if userPwd_1 == "" {
		s := structure_type.Things{"该学生账号不存在", false}
		render.JSON(w, r, s)
		return
	}
	if userPwd_1 != userPwd {
		s := structure_type.Things{"密码错误", false}
		render.JSON(w, r, s)
		return
	}
	s := structure_type.Things{"登录成功", true}
	render.JSON(w, r, s)
}

func (user *UserAPi) LoginTeaUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userName := r.Form["userName"][0]
	userPwd := r.Form["userPwd"][0]
	var userPwd_1 string

	if userName == "" || userPwd == "" {
		s := structure_type.Things{"请将信息输入完整", false}
		render.JSON(w, r, s)
		return
	}

	rows, err := user.db.Model(&data_conn.UserInfo{}).Where("UserName=? and RoleName=?", userName, "教师").Select("UserPwd").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&userPwd_1)
		if err != nil {
			return
		}
	}
	if userPwd_1 == "" {
		s := structure_type.Things{"该教师账号不存在", false}
		render.JSON(w, r, s)
		return
	}
	if userPwd_1 != userPwd {
		s := structure_type.Things{"密码错误", false}
		render.JSON(w, r, s)
		return
	}
	s := structure_type.Things{"登录成功", true}
	render.JSON(w, r, s)
}

func (user *UserAPi) UserPwdModify(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userName := r.Form["userName"][0]
	newUserPwd := r.Form["newUserPwd"][0]
	oldUserPwd := r.Form["oldUserPwd"][0]

	var userPwd string
	if newUserPwd == "" || oldUserPwd == "" {
		s := structure_type.Things{"请将信息输入完整", false}
		render.JSON(w, r, s)
		return
	}

	rows, err := user.db.Model(&data_conn.UserInfo{}).Where("UserName=?", userName).Select("UserPwd").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&userPwd)
		if err != nil {
			return
		}
	}
	if userPwd != oldUserPwd {
		s := structure_type.Things{"旧密码输入错误", false}
		render.JSON(w, r, s)
		return
	}
	err = user.db.Model(&data_conn.UserInfo{}).Where("UserName=?", userName).Update(&data_conn.UserInfo{UserPwd: newUserPwd}).Error
	if err != nil {
		return
	}
	s := structure_type.Things{"密码修改成功", true}
	render.JSON(w, r, s)
}

func (user *UserAPi) BrowUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	roleName := r.Form["roleName"][0]
	u := structure_type.UserInfoTotal{}
	tem := structure_type.UserInfo{}

	var rows *sql.Rows
	var err error

	//按角色查询
	if roleName != "" {
		rows, err = user.db.Model(&data_conn.UserInfo{}).Where("RoleName=?", roleName).Select("UserId,UserName,UserPwd,RoleName").Rows()
	}
	//查询全部
	if roleName == "" {
		rows, err = user.db.Model(&data_conn.UserInfo{}).Select("UserId,UserName,UserPwd,RoleName").Rows()
	}
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&tem.UserId, &tem.UserName, &tem.UserPwd, &tem.RoleName)
		if err != nil {
			return
		}
		u.UserInfoList = append(u.UserInfoList, tem)
	}
	render.JSON(w, r, u)
}

func (user *UserAPi) DelUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId := r.Form["userId"][0]

	err := user.db.Model(&data_conn.UserInfo{}).Where("UserId=?", userId).Delete(&data_conn.UserInfo{}).Error
	if err != nil {
		return
	}
	s := structure_type.Things{"删除用户信息成功", true}
	render.JSON(w, r, s)
}
