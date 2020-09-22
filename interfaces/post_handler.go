package interfaces

import (
	"Repository-Pattern/application"
	"Repository-Pattern/domain/model"
	"Repository-Pattern/helper"
	"Repository-Pattern/infrastructure/auth"
	"Repository-Pattern/interfaces/fileupload"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Post struct {
	postApp         application.PostAppInterface
	userApp         application.UserAppInterface
	postLabelApp    application.PostLabelAppInterface
	postCategoryApp application.PostCategoryAppInterface
	fileUpload      fileupload.UploadFileInterface
	tk              auth.TokenInterface
	rd              auth.AuthInterface
}

//Post constructor
func NewPost(fApp application.PostAppInterface, uApp application.UserAppInterface, lApp application.PostLabelAppInterface, cApp application.PostCategoryAppInterface, fd fileupload.UploadFileInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Post {
	return &Post{
		postApp:         fApp,
		userApp:         uApp,
		postLabelApp:    lApp,
		postCategoryApp: cApp,
		fileUpload:      fd,
		rd:              rd,
		tk:              tk,
	}
}

func (po *Post) SavePost(c *gin.Context) {
	//check is the user is authenticated first
	metadata, err := po.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := po.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//We we are using a frontend(vuejs), our errors need to have keys for easy checking, so we use a map to hold our errors
	var savePostError = make(map[string]string)
	title := c.PostForm("title")
	description := c.PostForm("description")
	labelArray := c.PostFormArray("labels")
	categoryArray := c.PostFormArray("categories")
	label := model.Label{}
	category := model.Category{}
	//var categories []model.Category
	//var labels []model.Label
	//for key, value := range c.Request.PostForm {
	//	fmt.Println(key, value)
	//}
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
	//check if the user exist
	user, err := po.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}
	formPost, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	postImages := formPost.File["post_images"]
	for _, file := range postImages {
		basename := filepath.Base(file.Filename)
		regex := helper.After(basename, ".")
		if regex == "png" || regex == "jpg" {
			dir := filepath.Join("./assets/images/", userId.String(), "/post")
			if dir != "" {
				err := os.Mkdir("./assets/images/"+userId.String()+"/post", os.ModePerm)
				if err != nil {
					fmt.Println(err.Error())
					_ = os.Mkdir("./assets/images/"+userId.String(), os.ModePerm)
					_ = os.Mkdir("./assets/images/"+userId.String()+"/post", os.ModePerm)
				}
			}
		}
		filename := filepath.Join("./assets/images/", userId.String(), "/post", basename)
		err := c.SaveUploadedFile(file, filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
	}
	var filenames []string
	for _, file := range postImages {
		filenames = append(filenames, file.Filename)
	}
	var Post = model.Post{}
	Post.UserUUID = userId
	Post.Title = title
	Post.Description = description
	Post.Author = user
	Post.PostImage = strings.Join(filenames, "")
	for i := 0; i < len(labelArray); i++ {
		label.Name = labelArray[i]
		//labels = append(labels, label)
	}
	for j := 0; j < len(categoryArray); j++ {
		category.Name = categoryArray[j]
		//categories = append(categories, category)
	}
	savedPost, saveErr := po.postApp.SavePost(&Post)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": savedPost, "filenames": filenames})
}

func (po *Post) UpdatePost(c *gin.Context) {
	//Check if the user is authenticated first
	metadata, err := po.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := po.rd.FetchAuth(metadata.TokenUuid)
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
	//Since it is a multipart formPost data we sent, we will do a manual check on each item
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
	user, err := po.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	//check if the Post exist:
	Post, err := po.postApp.GetPost(PostId)
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
		Post.PostImage, err = po.fileUpload.UploadFile(file)
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
	updatedPost, dbUpdateErr := po.postApp.UpdatePost(Post)
	if dbUpdateErr != nil {
		c.JSON(http.StatusInternalServerError, dbUpdateErr)
		return
	}
	c.JSON(http.StatusOK, updatedPost)
}

func (po *Post) GetAllPost(c *gin.Context) {
	posts := model.Posts{}
	var err error
	posts, err = po.postApp.GetAllPost()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, posts.PublicPosts())
}

func (po *Post) GetPostAndCreator(c *gin.Context) {
	postId, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	post, err := po.postApp.GetPost(postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	user, err := po.userApp.GetUser(post.UserUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	postAndUser := map[string]interface{}{
		"post":    post,
		"creator": user.PublicUser(),
	}
	c.JSON(http.StatusOK, postAndUser)
}

func (po *Post) DeletePost(c *gin.Context) {
	metadata, err := po.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	postId, err := strconv.ParseUint(c.Param("Post_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	_, err = po.userApp.GetUser(metadata.UserUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	err = po.postApp.DeletePost(postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "Post deleted")
}
