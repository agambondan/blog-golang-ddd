package interfaces

import (
	"Repository-Pattern/application"
	"Repository-Pattern/domain/model"
	"Repository-Pattern/infrastructure/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Category struct {
	categoryApp application.CategoryAppInterface
	roleApp     application.RoleAppInterface
	userApp     application.UserAppInterface
	tk          auth.TokenInterface
	rd          auth.AuthInterface
}

//Category constructor
func NewCategory(fApp application.CategoryAppInterface, rApp application.RoleAppInterface, uApp application.UserAppInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Category {
	return &Category{
		categoryApp: fApp,
		roleApp:     rApp,
		userApp:     uApp,
		rd:          rd,
		tk:          tk,
	}
}

func (ca *Category) SaveCategory(c *gin.Context) {
	//check is the user is authenticated first
	metadata, err := ca.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := ca.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	var saveCategoryError = make(map[string]string)
	name := c.PostForm("name")
	if fmt.Sprintf("%T", name) != "string" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "Invalid json",
		})
		return
	}
	emptyCategory := model.Category{}
	emptyCategory.Name = name
	saveCategoryError = emptyCategory.Validate("")
	if len(saveCategoryError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveCategoryError)
		return
	}
	//check if the user exist
	user, err := ca.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}
	role, err := ca.roleApp.GetRole(user.RoleId)
	user.Role = role
	if err != nil || user.RoleId != 1 || user.Role.Name != "admin" {
		var Category = model.Category{}
		Category.Name = "admin"
		_, _ = ca.categoryApp.SaveCategory(&Category)
		c.JSON(http.StatusBadRequest, "unauthorized, your not admin")
		return
	}
	var Category = model.Category{}
	Category.Name = name
	saveCategory, err := ca.categoryApp.SaveCategory(&Category)
	c.JSON(http.StatusCreated, saveCategory)
}

func (ca *Category) UpdateCategory(c *gin.Context) {
}

func (ca *Category) GetCategory(c *gin.Context) {
}

func (ca *Category) GetAllCategory(c *gin.Context) {
	//check is the user is authenticated first
	metadata, err := ca.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := ca.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	user, err := ca.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found")
		return
	}
	if user.RoleId != 1 {
		c.JSON(http.StatusInternalServerError, "your not admin, unauthorized")
		return
	}
	categories, err := ca.categoryApp.GetAllCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (ca *Category) DeleteCategory(c *gin.Context) {
}
