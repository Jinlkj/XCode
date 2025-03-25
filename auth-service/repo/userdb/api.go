package userdb

import (
	"code-search/auth-service/entity/config"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string `gorm:"password"`
}

func NewClient() Client {
	cfg := config.DefaultConfig
	// 连接MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(dsn)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// 自动迁移数据库表
	if err = db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	return &impl{
		db: db,
	}
}

type Client interface {
	CreateUser(ctx context.Context, username, email, password string) error
	CheckUserPassword(ctx context.Context, username, password string) error
}

type impl struct {
	db *gorm.DB
}

func (i *impl) CreateUser(_ context.Context, username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("hash password error: %v", err)
		return fmt.Errorf("hash error")
	}
	user := User{
		Password: string(hashedPassword),
		Email:    email,
		Username: username,
	}
	if err = i.db.Create(&user).Error; err != nil {
		log.Printf("create user error: %v", err)
		return fmt.Errorf("craete user error")
	}
	return nil
}

func (i *impl) CheckUserPassword(_ context.Context, username, password string) error {
	var user User
	if err := i.db.Where("username = ?", username).First(&user).Error; err != nil {
		log.Printf("find user error: %v", err)
		return fmt.Errorf("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("password not match error: %v", err)
		return fmt.Errorf("password %s not match user %s", password, user)
	}
	return nil
}
