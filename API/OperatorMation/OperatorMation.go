package UserMation

import (
	"database/sql"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"xiangmu/Student/DataConn"
	"xiangmu/Student/StructureType"
)

type UserAPi struct {
	db *gorm.DB
}

func Make_db(db *gorm.DB) *UserAPi {
	DB := &UserAPi{db}
	return DB
}

func (user *UserAPi) RegisterStuUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user_nu := r.Form["Stu_nu"][0]
	user_pwd := r.Form["Stu_pwd"][0]
	var stu_id, user_id int

	if user_nu == "" || user_pwd == "" {
		s := StructureType.Things{"请将信息输入完整"}
		render.JSON(w, r, s)
		return
	}
	//判断该学号是否已经注册过
	rows, err := user.db.Table("Userinfos").Where("User_name=?", user_nu).Select("User_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&user_id)
		if err != nil {
			return
		}
	}
	if user_id != 0 {
		s := StructureType.Things{"该学号已注册"}
		render.JSON(w, r, s)
		return
	}
	//判断学号是否存在
	rows, err = user.db.Table("Studentinfos").Where("Stu_nu=?", user_nu).Select("Stu_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&stu_id)
		if err != nil {
			return
		}
	}
	if stu_id == 0 {
		s := StructureType.Things{"该学号不存在"}
		render.JSON(w, r, s)
		return
	}

	err = user.db.Create(&DataConn.Userinfo{User_name: user_nu, User_pwd: user_pwd, Role_name: "学生"}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"注册学生用户信息成功"}
	render.JSON(w, r, s)
}

func (user *UserAPi) RegisterTeaUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user_name := r.Form["Tec_name"][0]
	user_pwd := r.Form["Tec_paw"][0]
	var user_id, tec_id int

	if user_name == "" || user_pwd == "" {
		s := StructureType.Things{"请将信息输入完整"}
		render.JSON(w, r, s)
		return
	}

	//判断该教师是否已经注册过
	rows, err := user.db.Table("Userinfos").Where("User_name=?", user_name).Select("User_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&user_id)
		if err != nil {
			return
		}
	}
	if user_id != 0 {
		s := StructureType.Things{"该教师信息已注册"}
		render.JSON(w, r, s)
		return
	}

	rows, err = user.db.Table("Teacherinfos").Where("Tec_name=?", user_name).Select("Tec_id").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&tec_id)
		if err != nil {
			return
		}
	}
	if tec_id == 0 {
		s := StructureType.Things{"该教师名字不存在"}
		render.JSON(w, r, s)
		return
	}

	err = user.db.Create(&DataConn.Userinfo{User_name: user_name, User_pwd: user_pwd, Role_name: "教师"}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"注册教师用户信息成功"}
	render.JSON(w, r, s)
}

func (user *UserAPi) LoginStuUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user_name := r.Form["User_name"][0]
	user_pwd := r.Form["User_pwd"][0]
	var user_pwd_1 string

	if user_name == "" || user_pwd == "" {
		s := StructureType.Things{"请将信息输入完整"}
		render.JSON(w, r, s)
		return
	}

	rows, err := user.db.Model(&DataConn.Userinfo{}).Where("User_name=? and Role_name=?", user_name, "学生").Select("User_pwd").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&user_pwd_1)
		if err != nil {
			return
		}
	}
	if user_pwd_1 == "" {
		s := StructureType.Things{"该学生账号不存在"}
		render.JSON(w, r, s)
		return
	}
	if user_pwd_1 != user_pwd {
		s := StructureType.Things{"密码错误"}
		render.JSON(w, r, s)
		return
	}
	s := StructureType.Things{"登录成功"}
	render.JSON(w, r, s)
}

func (user *UserAPi) LoginTeaUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user_name := r.Form["User_name"][0]
	user_pwd := r.Form["User_pwd"][0]
	var user_pwd_1 string

	if user_name == "" || user_pwd == "" {
		s := StructureType.Things{"请将信息输入完整"}
		render.JSON(w, r, s)
		return
	}

	rows, err := user.db.Model(&DataConn.Userinfo{}).Where("User_name=? and Role_name=?", user_name, "教师").Select("User_pwd").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&user_pwd_1)
		if err != nil {
			return
		}
	}
	if user_pwd_1 == "" {
		s := StructureType.Things{"该教师账号不存在"}
		render.JSON(w, r, s)
		return
	}
	if user_pwd_1 != user_pwd {
		s := StructureType.Things{"密码错误"}
		render.JSON(w, r, s)
		return
	}
	s := StructureType.Things{"登录成功"}
	render.JSON(w, r, s)
}

func (user *UserAPi) UserPwdModify(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user_name := r.Form["User_name"][0]
	newuser_pwd := r.Form["NewUser_pwd"][0]
	olduser_pwd := r.Form["OldUser_pwd"][0]

	var user_pwd string
	if newuser_pwd == "" || olduser_pwd == "" {
		s := StructureType.Things{"请将信息输入完整"}
		render.JSON(w, r, s)
		return
	}

	rows, err := user.db.Model(&DataConn.Userinfo{}).Where("User_name=?", user_name).Select("User_pwd").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&user_pwd)
		if err != nil {
			return
		}
	}
	if user_pwd != olduser_pwd {
		s := StructureType.Things{"旧密码输入错误"}
		render.JSON(w, r, s)
		return
	}
	err = user.db.Model(&DataConn.Userinfo{}).Where("User_name=?", user_name).Update(&DataConn.Userinfo{User_pwd: newuser_pwd}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"密码修改成功"}
	render.JSON(w, r, s)
}

func (user *UserAPi) BrowUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	role_name := r.Form["Role_name"][0]
	u := StructureType.Userinfototal{}
	tem := StructureType.Userinfo{}

	var rows *sql.Rows
	var err error


	//按角色查询
	if role_name != "" {
		rows, err = user.db.Model(&DataConn.Userinfo{}).Where("Role_name=?", role_name).Select("User_id,User_name,User_pwd,Role_name").Rows()
	}
	//查询全部
	if role_name == "" {
		rows, err = user.db.Model(&DataConn.Userinfo{}).Select("User_id,User_name,User_pwd,Role_name").Rows()
	}
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&tem.User_id, &tem.User_name, &tem.User_pwd, &tem.Role_name)
		if err != nil {
			return
		}
		u.UserinfoList = append(u.UserinfoList, tem)
	}
	render.JSON(w, r, u)
}

func (user *UserAPi) SubUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user_id := r.Form["User_id"][0]

	err := user.db.Model(&DataConn.Userinfo{}).Where("User_id=?", user_id).Delete(&DataConn.Userinfo{}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"删除用户信息成功"}
	render.JSON(w, r, s)
}
