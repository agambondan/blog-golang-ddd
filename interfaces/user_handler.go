package interfaces

import (
	"Repository-Pattern/application"
	"Repository-Pattern/domain/model"
	"Repository-Pattern/infrastructure/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
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
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "invalid json",
		})
		return
	}
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
	c.JSON(http.StatusCreated, newUser.PublicUser())
}

func (s *Users) GetUsers(c *gin.Context) {
	users := model.Users{} //customize user
	var err error
	//us, err = application.UserApp.GetUsers()
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
	po
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	user.Posts = posts
	c.JSON(http.StatusOK, user)
}
