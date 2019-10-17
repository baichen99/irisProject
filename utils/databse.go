package utils

import (
	"irisProject/models"
	"strconv"

	"irisProject/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
)

// ConnectDB connect to a psql database
func ConnectDB() (db *gorm.DB) {
	DatabaseURI := "host=" + config.Conf.Postgres.Host + " port=" +
		config.Conf.Postgres.Port + " user=" + config.Conf.Postgres.User + " dbname=" +
		config.Conf.Postgres.Database + " password=" + config.Conf.Postgres.Password + " sslmode=disable"
	db, err := gorm.Open("postgres", DatabaseURI)
	if err != nil {
		panic(err)
	}
	return
}

// InitDB initialize database
func InitDB(db *gorm.DB) *gorm.DB {
	// 创建迁移表
	db.DropTableIfExists("student_teacher")
	db.DropTableIfExists(&models.User{}, &models.Teacher{}, &models.Profile{})
	db.Table("student_teacher").CreateTable(&StudentTeacher{})
	db.AutoMigrate(&models.User{}, &models.Teacher{}, &models.Profile{})

	var user models.User
	var users []models.User
	var teacher models.Teacher
	var profile models.Profile
	user.Password, _ = HashPassword("password")
	// 生成数据	 建立关系
	for i := 1; i <= 10; i++ {
		user.Username = "user_" + strconv.Itoa(i)
		teacher.Name = "teacher_" + strconv.Itoa(i)
		profile.Content = "profile_" + strconv.Itoa(i)
		db.Create(&user)
		db.Create(&teacher)
		db.Create(&profile)
		db.Model(&models.User{}).Where("id = ?", user.ID).Update("profile_id", profile.ID)
	}

	db.Find(&users)
	for _, u := range users {
		db.First(&teacher)
		db.Table("student_teacher").Create(&StudentTeacher{StudentID: u.ID, TeacherID: teacher.ID})
	}

	// 手动添加外键
	db.Model(&models.User{}).AddForeignKey("profile_id", "profiles(id)", "RESTRICT", "CASCADE")
	db.Table("student_teacher").AddForeignKey("student_id", "users(id)", "RESTRICT", "CASCADE").AddForeignKey("teacher_id", "teachers(id)", "RESTRICT", "CASCADE")

	return db
}

// InitAdmin add admin users
func InitAdmin(db *gorm.DB) {
	var profile models.Profile
	password, _ := HashPassword("password")
	db.Create(&profile)
	user1 := models.User{
		Username:  "admin",
		Password:  password,
		Role:      "admin",
		ProfileID: profile.ID,
	}

	db.Create(&user1)
}

type StudentTeacher struct {
	StudentID uuid.UUID `gorm:"uuid;PRIMARYKEY"`
	TeacherID uuid.UUID `gorm:"uuid;PRIMARYKEY"`
}
