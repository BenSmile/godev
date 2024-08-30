package types

import (
	"fmt"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 2
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsAdmin   bool   `json:"isAdmin"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (p UpdateUserParams) ToBson() bson.M {
	m := bson.M{}

	if len(p.FirstName) > 0 {
		m["firstName"] = p.FirstName
	}

	if len(p.LastName) > 0 {
		m["lastName"] = p.LastName
	}

	if len(p.Email) > 0 {
		m["email"] = p.Email
	}

	return m

}

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LasttName string             `json:"lastName" bson:"lastName"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"-" bson:"password"`
	IsAdmin   bool               `json:"isAdmin" bson:"is_admin"`
}

func (params CreateUserParams) Validate() fiber.Map {
	errors := fiber.Map{}
	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
	}
	if !isValidEmail(params.Email) {
		errors["email"] = "email is invalid"
	}
	return errors
}

func isValidEmail(email string) bool {
	// Define the regular expression for validating an email address
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func IsValidPassword(hashedPwd, pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd)) == nil
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encryptedPwd, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName: params.FirstName,
		LasttName: params.LastName,
		Email:     params.Email,
		Password:  string(encryptedPwd),
		IsAdmin:   params.IsAdmin,
	}, nil
}
