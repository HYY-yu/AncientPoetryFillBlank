package setdb

import (
	"github.com/jinzhu/gorm"
	"os"
	"log"
	"fmt"
)

//local 数据库
const (
	DBUserName = "root"
	DBPassword = "yufeng"
	DBHost     = "127.0.0.1:3306"
	DBName     = "dreamEbagPaperTest"

	DBMaxIdle = 20
	DBMaxConn = 20
)

type QuestionD struct {
	QuestionId          int64   `gorm:"primary_key;column:F_question_id;type:BIGINT(20)" json:"id"`
	CourseId            uint    `gorm:"column:F_course_id;type:TINYINT(2) UNSIGNED" json:"courseId"`        //问题对应的课程
	Content             string  `gorm:"column:F_content;type:LONGTEXT" json:"content"`                      //问题的内容
	CorrectAnswer       string  `gorm:"column:F_correct_answer;type:TEXT" json:"correctAnswer"` //正确答案 （不一定有）
	Type                int     `gorm:"column:F_type" json:"type"`                              //问题的类型分类
}

func InitGorm() {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "t_" + defaultTableName
	}
}

func SetToDB() {
	//db
	db, _ := gorm.Open("mysql", DBUserName+":"+DBPassword+"@tcp("+DBHost+")/"+DBName+"?charset=utf8&parseTime=True&loc=Asia%2FShanghai")
	db.DB().SetMaxIdleConns(DBMaxIdle)
	db.DB().SetMaxOpenConns(DBMaxConn)
	defer db.Close()
	InitGorm()

	// 启用Logger，显示详细日志
	//db.LogMode(MyConfig.DBlog) // 只显示错误日志
	logFile, _ := os.OpenFile("log/db.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	db.SetLogger(log.New(logFile, "\r\n", log.LstdFlags))


	var allQuestion []QuestionD
	db.Table("t_questions").Where("F_type = ? and F_correct_answer = '' and F_course_id IN (12,42,53,70)",61).Scan(&allQuestion)


	fmt.Println(len(allQuestion))

}