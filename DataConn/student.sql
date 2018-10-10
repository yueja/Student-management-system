create database student;
use student;

create table  Roles(
	role_id   int primary key auto_increment not null,
	role_name varchar(255) not null
);

create table Operators(
	ope_id  int  primary key auto_increment not null ,
	ope_name varchar(255)  not null,
	ope_pwd  varchar(255)  not null,            
	rol_id   int,   
    CONSTRAINT fk_Roles_Operator FOREIGN KEY (rol_id) REFERENCES roles (role_id)
);

create table Privileges(
	pri_id    int  primary key auto_increment not null ,
	pri_name  varchar(255)  not null,
	pri_url   varchar(255)  not null,
	menu_name varchar(255)  not null,
	rol_id    int  not null,
	CONSTRAINT fk_Roles_Privilege FOREIGN KEY (rol_id) REFERENCES roles (role_id)
);

create table Subjects (
	sub_id    int  primary key auto_increment not null ,
	sub_name  varchar(255)  not null,           
	sub_type  varchar(255)  not null,                
	sub_times int  not null    
);

create table Majors (
	maj_id    int    primary key auto_increment not null ,
	maj_name  varchar(255)  not null,         
	maj_prin  varchar(255)  not null,         
	maj_link  varchar(255)  not null,         
	maj_phone varchar(255)  not null        
);

create table  Classinfos (
	class_id   int    primary key auto_increment not null ,
	class_name varchar(255)  not null,        
	class_tec  varchar(255)  not null,        
	maj_id     int ,
    CONSTRAINT fk_Classinfo_Major FOREIGN KEY (maj_id) REFERENCES Majors (maj_id)
);

create table  Teacherinfos (
	tec_id    int   primary key auto_increment not null ,
	ope_id    int , 
	tec_name  varchar(255)  not null,
	tec_sex   varchar(255)  not null,
	tec_birth varchar(255)  not null,
	tec_major varchar(255)  not null,
	tec_phone varchar(255)  not null,
    CONSTRAINT fk_Teacherinfo_Operator FOREIGN KEY (ope_id) REFERENCES Operators (ope_id)
);

create table  Studentinfos (
	stu_id    int   primary key auto_increment not null ,
	ope_id    int,
	stu_nu    varchar(255)  not null,
	stu_name  varchar(255)  not null,
	stu_sex   varchar(255)  not null,
	stu_birth varchar(255)  not null,
	stu_pic   varchar(255)  not null,
	cla_id    int,
    CONSTRAINT fk_Studentinfo_Operator FOREIGN KEY (ope_id) REFERENCES Operators (ope_id),
    CONSTRAINT fk_Studentinfo_Classinfo FOREIGN KEY (cla_id) REFERENCES Classinfos  (class_id)
);

create table  Scores(
	sco_id     int primary key auto_increment not null ,
	sco_daily  float  not null,
	sub_exam   float  not null,
	wco_count   float  not null,
	stu_id     int ,
	sub_id     int ,
	cla2sub_id int ,
	cla_id     int ,
    CONSTRAINT fk_Score_Studentinfo FOREIGN KEY (stu_id) REFERENCES Studentinfos (stu_id),
    CONSTRAINT fk_Score_Subject FOREIGN KEY (sub_id) REFERENCES Subjects  (sub_id),
    CONSTRAINT fk_Score_Cla2sub FOREIGN KEY (cla2sub_id) REFERENCES Cla2subs (cla2sub_id),
    CONSTRAINT fk_Score_Classinfo FOREIGN KEY (cla_id) REFERENCES Classinfos  (class_id)
);

create table  Cla2subs (
	cla2sub_id int primary key auto_increment not null ,
	cla_id     int,
	sub_id     int ,
	tec_id     int ,
    CONSTRAINT fk_Cla2sub_Classinfo FOREIGN KEY (Cla_id) REFERENCES Classinfos  (class_id),
    CONSTRAINT fk_Cla2sub_Subject FOREIGN KEY (sub_id) REFERENCES Subjects  (sub_id),
    CONSTRAINT fk_Score_Teacherinfo FOREIGN KEY (tec_id) REFERENCES Teacherinfos (tec_id)
)