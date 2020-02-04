package controllers

// import "C"
import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/Neostore/dbconnection"
	"github.com/Neostore/models"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Neostore/form"

	"github.com/Neostore/middlewares"
	"github.com/Neostore/services"

	"net/http"
	"strconv"
)

func RegisterProductRoutes(router *gin.RouterGroup) {
	router.GET("/", ProductList)
	router.GET("/:slug", GetProductDetailsBySlug)

	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.POST("/", CreateProduct)
		router.DELETE("/:slug", ProductDelete)
	}
}

func ProductList(c *gin.Context) {

	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 5
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	productModels, modelCount, commentsCount, err := services.FetchProductsPage(page, pageSize)
	if err != nil {
		c.JSON(http.StatusNotFound, form.CreateDetailedErrorDto("products", errors.New("Invalid param")))
		return
	}

	c.JSON(http.StatusOK, form.CreatedProductPagedResponse(c.Request, productModels, page, pageSize, modelCount, commentsCount))
}

func GetProductDetailsBySlug(c *gin.Context) {
	productSlug := c.Param("slug")

	product := services.FetchProductDetails(&models.Product{Slug: productSlug}, true)
	if product.ID == 0 {
		c.JSON(http.StatusNotFound, form.CreateDetailedErrorDto("products", errors.New("Invalid slug")))
		return
	}
	c.JSON(http.StatusOK, form.CreateProductDetailsDto(product))
}

func CreateProduct(c *gin.Context) {
	// Only admin users can create products
	user := c.Keys["currentUser"].(models.User)
	if user.IsNotAdmin() {
		c.JSON(http.StatusForbidden, form.CreateErrorDtoWithMessage("Permission denied, you must be admin"))
		return
	}

	var formDto form.CreateProduct
	if err := c.ShouldBind(&formDto); err != nil {
		c.JSON(http.StatusBadRequest, form.CreateBadRequestErrorDto(err))
		return
	}

	name := formDto.Name
	description := formDto.Description

	price := formDto.Price
	stock, err := strconv.ParseInt(c.PostForm("stock"), 10, 32)
	form, err := c.MultipartForm()

	tagCount := 0
	catCount := 0
	for key := range form.Value {
		if strings.HasPrefix(key, "tags[") {
			tagCount++
		}
		if strings.HasPrefix(key, "category[") {
			catCount++
		}
	}

	var tags = make([]models.Tag, tagCount)
	var categories = make([]models.Category, catCount)

	var rgx = regexp.MustCompile(`\[(.*?)\]`)
	database := dbconnection.GetDb()
	tagPtr := 0
	catPtr := 0

	for k, v := range form.Value {
		if strings.HasPrefix(k, "tags[") {
			result := rgx.FindStringSubmatch(k)
			var tag models.Tag
			name := result[1]
			description := v[0]
			database.Where(&models.Tag{Slug: slug.Make(name)}).
				Attrs(models.Tag{Name: name, Description: description}).
				FirstOrCreate(&tag)
			tags[tagPtr] = tag
			tagPtr++
		}

		if strings.HasPrefix(k, "category[") {
			result := rgx.FindStringSubmatch(k)
			var category models.Category
			name := result[1]
			description := v[0]
			database.Where(&models.Category{Slug: slug.Make(name)}).
				Attrs(models.Category{Name: name, Description: description}).
				FirstOrCreate(&category)
			categories[catPtr] = category
			catPtr++
		}
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, form.CreateDetailedErrorDto("form_error", err))
		return
	}

	files := form.File["images[]"]
	var productImages = make([]models.FileUpload, len(files))

	for index, file := range files {
		fileName := randomString(16) + ".png"

		dirPath := filepath.Join(".", "static", "images", "products")
		filePath := filepath.Join(dirPath, fileName)
		if _, err = os.Stat(dirPath); os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, os.ModeDir)
			if err != nil {
				c.JSON(http.StatusInternalServerError, form.CreateDetailedErrorDto("io_error", err))
				return
			}
		}
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusBadRequest, form.CreateDetailedErrorDto("upload_error", err))
			return
		}
		fileSize := (uint)(file.Size)
		productImages[index] = models.FileUpload{Filename: fileName, OriginalName: file.Filename, FilePath: string(filepath.Separator) + filePath, FileSize: fileSize}
	}

	product := models.Product{
		Name:        name,
		Description: description,
		Tags:        tags,
		Categories:  categories,
		Price:       (int)(price),
		Stock:       (int)(stock),
		Images:      productImages,
	}

	if err := services.CreateOne(&product); err != nil {
		c.JSON(http.StatusUnprocessableEntity, form.CreateDetailedErrorDto("database", err))
		return
	}

	c.JSON(http.StatusOK, from.CreateProductCreatedDto(product))

}

func ProductDelete(c *gin.Context) {
	slug := c.Param("slug")
	err := services.DeleteProduct(&models.Product{Slug: slug})
	if err != nil {
		c.JSON(http.StatusNotFound, form.CreateDetailedErrorDto("products", errors.New("Invalid slug")))
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": "Delete success"})
}
