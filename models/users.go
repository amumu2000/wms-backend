package models

type User struct {
	ID       int64  `gorm:"column:id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	Email    string `gorm:"column:email"`
	Role     int    `gorm:"column:role"`
}

func GetUserByID(id int64) *User {
	var count int64
	user := User{}
	db.Table("users").Where("id = ?", id).Count(&count).First(&user)

	if count == 0 {
		return nil
	} else {
		return &user
	}
}

func GetUsers(username, password, email string, role int) []User {
	query := db.Table("users")

	if username != "" {
		query = query.Where("username = ?", username)
	}

	if password != "" {
		query = query.Where("password = ?", password)
	}

	if email != "" {
		query = query.Where("email = ?", email)
	}

	if role >= 0 {
		query = query.Where("role = ?", role)
	}

	var users []User
	query.Find(&users)

	return users
}

func FindUsers(userID int64, username, email string, role int) []User {
	query := db.Table("users")

	if userID >= 0 {
		query = query.Where("id = ?", userID)
	}

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}

	if role >= 0 {
		query = query.Where("role = ?", role)
	}

	var users []User
	query.Find(&users)

	return users
}

func EditUser(user User) {
	db.Table("users").Save(&user)
}

func DeleteUser(id []int64) {
	db.Table("users").Delete(&User{}, id)
}

func AddUser(username, password, email string, role int) int64 {
	user := User{
		Username: username,
		Password: password,
		Email:    email,
		Role:     role,
	}

	db.Table("users").Create(&user)

	return user.ID
}
