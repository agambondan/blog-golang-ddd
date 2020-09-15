package interfaces

import (
	"Repository-Pattern/application"
	"Repository-Pattern/domain/model"
	"Repository-Pattern/infrastructure/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Role struct {
	roleApp application.RoleAppInterface
	userApp application.UserAppInterface
	tk      auth.TokenInterface
	rd      auth.AuthInterface
}

//Role constructor
func NewRole(fApp application.RoleAppInterface, uApp application.UserAppInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Role {
	return &Role{
		roleApp: fApp,
		userApp: uApp,
		rd:      rd,
		tk:      tk,
	}
}

func (ro *Role) SaveRole(c *gin.Context) {
	//check is the user is authenticated first
	metadata, err := ro.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := ro.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	var saveRoleError = make(map[string]string)
	name := c.PostForm("name")
	if fmt.Sprintf("%T", name) != "string" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "Invalid json",
		})
		return
	}
	emptyRole := model.Role{}
	emptyRole.Name = name
	saveRoleError = emptyRole.Validate("")
	if len(saveRoleError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveRoleError)
		return
	}
	//check if the user exist
	user, err := ro.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}
	role, err := ro.roleApp.GetRole(user.RoleId)
	user.Role = role
	if err != nil || user.RoleId != 1 || user.Role.Name != "admin" {
		var Role = model.Role{}
		Role.Name = "admin"
		_, _ = ro.roleApp.SaveRole(&Role)
		c.JSON(http.StatusBadRequest, "unauthorized, your not admin")
		return
	}
	var Role = model.Role{}
	Role.Name = name
	saveRole, saveErr := ro.roleApp.SaveRole(&Role)
	for saveError := range saveErr {
		if saveError == "" {
			c.JSON(http.StatusInternalServerError, saveErr)
			return
		} else {
			c.JSON(http.StatusInternalServerError, saveErr)
			return
		}
	}
	c.JSON(http.StatusCreated, saveRole)
}

func (ro *Role) UpdateRole(c *gin.Context) {
}

func (ro *Role) GetRole(c *gin.Context) {
}

func (ro *Role) GetAllRole(c *gin.Context) {
	//check is the user is authenticated first
	metadata, err := ro.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := ro.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	user, err := ro.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found")
		return
	}
	if user.RoleId != 1 {
		c.JSON(http.StatusInternalServerError, "your not admin, unauthorized")
		return
	}
	roles, err := ro.roleApp.GetAllRole()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (ro *Role) DeleteRole(c *gin.Context) {
}
