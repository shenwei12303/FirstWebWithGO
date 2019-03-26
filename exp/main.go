package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"sw.com/FirstWebWithGO/models"
)

const (
	host     = "192.168.1.8"
	port     = 5432
	user     = "postgres"
	password = "310104shenwei"
	dbname   = "website"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.DestructiveReset()

	user := models.User{
		Name:     "Michael Scott",
		Email:    "michael@dundermifflin.com",
		Password: "bestboss",
	}

	if err := us.Create(&user); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", user)
	if user.Remember == "" {
		panic("Invaild remember token")
	}

	user2, err := us.ByRemember(user.Remember)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", user2)
	// user.Name = "Updated Name"
	// if err := us.Update(&user); err != nil {
	// 	panic(err)
	// }

	// foundUser, err := us.ByEmail("michael@dundermifflin.com")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(foundUser)

	// if err := us.Delete(foundUser.ID); err != nil {
	// 	panic(err)
	// }

	// _, err = us.ByID(foundUser.ID)
	// if err != models.ErrNotFound {
	// 	panic("user was not deleted!")
	// }

}

type User struct {
	gorm.Model
	Name   string
	Email  string `gorm:"not null;unique_index"`
	Orders []Order
}

type Order struct {
	gorm.Model
	UserID      uint
	Amout       int
	Description string
}

func createOrder(db *gorm.DB, user User, amount int, des string) {
	db.Create(&Order{
		UserID:      user.ID,
		Amout:       amount,
		Description: des,
	})
	if db.Error != nil {
		panic(db.Error)
	}
}
