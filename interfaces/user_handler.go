package interfaces

import (
	"Repository-Pattern/application"
	"Repository-Pattern/domain/model"
	"Repository-Pattern/helper"
	"Repository-Pattern/infrastructure/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

//Users struct defines the dependencies that will be used
type Users struct {
	us application.UserAppInterface
	po application.PostAppInterface
	ro application.RoleAppInterface
	rd auth.AuthInterface
	tk auth.TokenInterface
}

//Users constructor
func NewUsers(us application.UserAppInterface, po application.PostAppInterface, ro application.RoleAppInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Users {
	return &Users{
		us: us,
		po: po,
		ro: ro,
		rd: rd,
		tk: tk,
	}
}

func (s *Users) SaveUser(c *gin.Context) {
	var user model.User
	//if err := c.ShouldBindJSON(&user); err != nil {
	//	c.JSON(http.StatusUnprocessableEntity, gin.H{
	//		"invalid_json": "invalid json",
	//	})
	//	return
	//}
	firstName := c.PostForm("first_name")
	lastName := c.PostForm("last_name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	username := c.PostForm("username")
	phoneNumber := c.PostForm("phone_number")
	roleID := c.PostForm("role_id")

	user.FirstName = firstName
	user.LastName = lastName
	user.Email = email
	user.Password = password
	user.Username = username
	user.PhoneNumber = phoneNumber
	user.RoleId, _ = strconv.ParseUint(roleID, 10, 64)
	//validate the request:
	validateErr := user.Validate("")
	if len(validateErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, validateErr)
		return
	}
	newUser, err := s.us.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	formUser, err2 := c.MultipartForm()
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err2.Error()})
		return
	}
	files := formUser.File["user_images"]
	for _, file := range files {
		basename := filepath.Base(file.Filename)
		regex := helper.After(basename, ".")
		if regex == "png" || regex == "jpg" {
			dir := filepath.Join("./assets/images/", newUser.UUID.String(), "/user")
			if dir != "" {
				err := os.Mkdir("./assets/images/"+newUser.UUID.String()+"/user", os.ModePerm)
				if err != nil {
					fmt.Println(err.Error())
					_ = os.Mkdir("./assets/images/"+newUser.UUID.String(), os.ModePerm)
					_ = os.Mkdir("./assets/images/"+newUser.UUID.String()+"/user", os.ModePerm)
				}
			}
		}
		filename := filepath.Join("./assets/images/", newUser.UUID.String(), "/user", basename)
		err := c.SaveUploadedFile(file, filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
	}
	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Filename)
	}
	c.JSON(http.StatusCreated, newUser.PublicUser())
}

func (s *Users) GetPrivateUsers(c *gin.Context) {
	//check is the user is authenticated first
	metadata, err := s.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := s.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	user, err := s.us.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found")
		return
	}
	if user.RoleId != 1 {
		c.JSON(http.StatusInternalServerError, "your not admin, unauthorized")
		return
	}
	users, err := s.us.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users)
}

func (s *Users) GetUsers(c *gin.Context) {
	users := model.Users{} //customize user
	var err error
	users, err = s.us.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users.PublicUsers())
}

func (s *Users) GetUser(c *gin.Context) {
	uuidParam := c.Param("user_id")
	userUUID, err := uuid.Parse(uuidParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := s.us.GetUser(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	user.Role = nil
	posts, err := s.po.GetPostByIdUser(userUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	user.Posts = posts
	c.JSON(http.StatusOK, user)
}
