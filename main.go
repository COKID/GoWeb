package main

import (
	"fmt"
	"github.com/COKID/GoWeb/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)



func main() {
	fmt.Println("hello world")

	db:=InitDB()//初始化数据库

	r := gin.Default()//初始化Gin
	r.GET("/api/auth/register", func(ctx *gin.Context) {
		//获取参数
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")

		//数据验证
		if len(telephone) != 11 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "手机号码必须为11位",
			})
			return
		}
		if len(password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "密码不能少于6位",
			})
			return
		}
		//没有名称的随机给10位的字符串
		if len(name) == 0 {
			name = RandomString(10)
		}

		log.Println(name, telephone, password)
		//判断手机号是否存在
		if isTelephoneExist(db,telephone){
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "用户已经存在",
			})
			return
		}

		//创建用户
		newUser:=models.User{
			Name: name,
			Telephone: telephone,
			Password: password,
		}
		db.Create(&newUser)

		//返回结果
		ctx.JSON(200, gin.H{
			"msg": "注册成功",
		})
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user models.User
	db.Where("telephone=?",telephone).First(&user)
	if user.ID!=0{
		return true
	}
	return false
}

func RandomString(n int) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().Unix())

	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
func InitDB() *gorm.DB {
	//driverName:="mysql"
	username := "root"
	password := "000000"
	host := "localhost"
	port := "3306"
	database := "cokid_db"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset,
	)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database:" + err.Error())
	}

	//这个是gorm自动创建数据表的函数。它会自动在数据库中创建一个名为users的数据表
	_ = db.AutoMigrate(&models.User{})
	return db
}
