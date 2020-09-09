package interfaces

import (
	"Repository-Pattern/application"
	"Repository-Pattern/domain/model"
	"Repository-Pattern/infrastructure/auth"
	"Repository-Pattern/interfaces/fileupload"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Post struct {
	PostApp    application.PostAppInterface
	userApp    application.UserAppInterface
	fileUpload fileupload.UploadFileInterface
	tk         auth.TokenInterface
	rd         auth.AuthInterface
}

//Post constructor
func NewPost(fApp application.PostAppInterface, uApp application.UserAppInterface, fd fileupload.UploadFileInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Post {
	return &Post{
		PostApp:    fApp,
		userApp:    uApp,
		fileUpload: fd,
		rd:         rd,
		tk:         tk,
	}
}

func (fo *Post) SavePost(c *gin.Context) {
	//check is the user is authenticated first
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := fo.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//We we are using a frontend(vuejs), our errors need to have keys for easy checking, so we use a map to hold our errors
	var savePostError = make(map[string]string)

	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "Invalid json",
		})
		return
	}
	//We initialize a new Post for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyPost := model.Post{}
	emptyPost.Title = title
	emptyPost.Description = description
	savePostError = emptyPost.Validate("")
	if len(savePostError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, savePostError)
		return
	}
	file, err := c.FormFile("Post_image")
	if err != nil {
		savePostError["invalid_file"] = "a valid file is required"
		c.JSON(http.StatusUnprocessableEntity, savePostError)
		return
	}
	//check if the user exist
	_, err = fo.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}
	uploadedFile, err := fo.fileUpload.UploadFile(file)
	if err != nil {
		savePostError["upload_err"] = err.Error() //this error can be any we defined in the UploadFile method
		c.JSON(http.StatusUnprocessableEntity, savePostError)
		return
	}
	var Post = model.Post{}
	Post.UserUUID = userId
	Post.Title = title
	Post.Description = description
	Post.PostImage = uploadedFile
	savedPost, saveErr := fo.PostApp.SavePost(&Post)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}
	c.JSON(http.StatusCreated, savedPost)
}

func (fo *Post) UpdatePost(c *gin.Context) {
	//Check if the user is authenticated first
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := fo.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//We we are using a frontend(vuejs), our errors need to have keys for easy checking, so we use a map to hold our errors
	var updatePostError = make(map[string]string)

	PostId, err := strconv.ParseUint(c.Param("Post_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	//Since it is a multipart form data we sent, we will do a manual check on each item
	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json")
	}
	//We initialize a new Post for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyPost := model.Post{}
	emptyPost.Title = title
	emptyPost.Description = description
	updatePostError = emptyPost.Validate("update")
	if len(updatePostError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, updatePostError)
		return
	}
	user, err := fo.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	//check if the Post exist:
	Post, err := fo.PostApp.GetPost(PostId)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	//if the user id doesnt match with the one we have, dont update. This is the case where an authenticated user tries to update someone else post using postman, curl, etc
	if user.UUID != Post.UserUUID {
		c.JSON(http.StatusUnauthorized, "you are not the owner of this Post")
		return
	}
	//Since this is an update request,  a new image may or may not be given.
	// If not image is given, an error occurs. We know this that is why we ignored the error and instead check if the file is nil.
	// if not nil, we process the file by calling the "UploadFile" method.
	// if nil, we used the old one whose path is saved in the database
	file, _ := c.FormFile("Post_image")
	if file != nil {
		Post.PostImage, err = fo.fileUpload.UploadFile(file)
		//since i am using Digital Ocean(DO) Spaces to save image, i am appending my DO url here. You can comment this line since you may be using Digital Ocean Spaces.
		Post.PostImage = os.Getenv("DO_SPACES_URL") + Post.PostImage
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"upload_err": err.Error(),
			})
			return
		}
	}
	//we dont need to update user's id
	Post.Title = title
	Post.Description = description
	Post.UpdatedAt = time.Now()
	updatedPost, dbUpdateErr := fo.PostApp.UpdatePost(Post)
	if dbUpdateErr != nil {
		c.JSON(http.StatusInternalServerError, dbUpdateErr)
		return
	}
	c.JSON(http.StatusOK, updatedPost)
}

func (fo *Post) GetAllPost(c *gin.Context) {
	allPost, err := fo.PostApp.GetAllPost()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allPost)
}

func (fo *Post) GetPostAndCreator(c *gin.Context) {
	PostId, err := strconv.ParseUint(c.Param("Post_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	Post, err := fo.PostApp.GetPost(PostId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	user, err := fo.userApp.GetUser(Post.UserUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	PostAndUser := map[string]interface{}{
		"Post":    Post,
		"creator": user.PublicUser(),
	}
	c.JSON(http.StatusOK, PostAndUser)
}

func (fo *Post) DeletePost(c *gin.Context) {
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	PostId, err := strconv.ParseUint(c.Param("Post_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	_, err = fo.userApp.GetUser(metadata.UserUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	err = fo.PostApp.DeletePost(PostId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "Post deleted")
}
