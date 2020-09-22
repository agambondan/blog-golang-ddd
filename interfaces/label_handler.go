package interfaces

import (
	"Repository-Pattern/application"
	"Repository-Pattern/domain/model"
	"Repository-Pattern/infrastructure/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Label struct {
	labelApp application.LabelAppInterface
	roleApp  application.RoleAppInterface
	userApp  application.UserAppInterface
	tk       auth.TokenInterface
	rd       auth.AuthInterface
}

//Label constructor
func NewLabel(fApp application.LabelAppInterface, rApp application.RoleAppInterface, uApp application.UserAppInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Label {
	return &Label{
		labelApp: fApp,
		roleApp: rApp,
		userApp:  uApp,
		rd:       rd,
		tk:       tk,
	}
}

func (ro *Label) SaveLabel(c *gin.Context) {
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
	var saveLabelError = make(map[string]string)
	name := c.PostForm("name")
	if fmt.Sprintf("%T", name) != "string" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "Invalid json",
		})
		return
	}
	emptyLabel := model.Label{}
	emptyLabel.Name = name
	saveLabelError = emptyLabel.Validate("")
	if len(saveLabelError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveLabelError)
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
		var Label = model.Label{}
		Label.Name = "admin"
		_, _ = ro.labelApp.SaveLabel(&Label)
		c.JSON(http.StatusBadRequest, "unauthorized, your not admin")
		return
	}
	var Label = model.Label{}
	Label.Name = name
	saveLabel, err := ro.labelApp.SaveLabel(&Label)
	c.JSON(http.StatusCreated, saveLabel)
}

func (ro *Label) UpdateLabel(c *gin.Context) {
}

func (ro *Label) GetLabel(c *gin.Context) {
}

func (ro *Label) GetAllLabel(c *gin.Context) {
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
	labels, err := ro.labelApp.GetAllLabel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, labels)
}

func (ro *Label) DeleteLabel(c *gin.Context) {
}
