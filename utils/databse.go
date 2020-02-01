package utils

import (
	"irisProject/models"
	"strconv"

	"irisProject/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	db.DropTableIfExists(&models.User{}, &models.Teacher{}, &models.Profile{})
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
		// random choose a teacher record
		db.First(&teacher)
		// UPDATE teachers SET students_id = array_append(students_id, 'ae1a635b-b9a5-4bd5-a3cb-5fbd42e15f1c')
		// WHERE teachers.id = '698b4b89-9607-4fed-9283-7584a8e94fc4';
		db.Model(&u).Where("id = ?", u.ID).Update("teachers_id", append(user.TeachersID, teacher.ID.String()))
	}

	// 手动添加外键
	db.Model(&models.User{}).
		AddForeignKey("profile_id", "profiles(id)", "RESTRICT", "CASCADE")
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
