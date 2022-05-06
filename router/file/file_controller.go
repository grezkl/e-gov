package file

import (
	"crypto/md5"
	"e-gov/global"
	"e-gov/middleware"
	"e-gov/model"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// @Summary 文件管理分页
// @Tags File
// @version 1.0
// @Accept json
// @Param pageNum query int true "页数"
// @Param pageSize query int true "页大小"
// @Param name query string false "文件名称"
// @Router /filePage [get]
func filePage(r *gin.Engine) {
	r.GET("/filePage", func(c *gin.Context) {

		var page_form model.PageForm
		if c.ShouldBindQuery(&page_form) != nil {
			middleware.Err(c, model.CODE_500, "参数错误", nil)
			return
		}
		if page_form.PageNum <= 0 {
			page_form.PageNum = 1
		}

		var file_list_query []model.File
		name := page_form.Name

		var total int64
		// global.DB.Model(&file_list_query).Count(&total)

		if name != "" {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("name like ?", name).Find(&file_list_query).Count(&total)
		} else {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum - 1) * page_form.PageSize).Find(&file_list_query).Count(&total)
		}

		var page_resp model.PageResp
		page_resp.Records = file_list_query
		page_resp.Total = total

		middleware.Success(c, model.CODE_200, "", page_resp)
	})
}

// @Summary 删除文件（单个）
// @Tags File
// @version 1.0
// @Accept json
// @Param id path int true "文件编号"
// @Router /deleteFileById/{id} [delete]
func deleteFileById(r *gin.Engine) {
	r.DELETE("/deleteFileById/:id", func(c *gin.Context) {
		id := c.Param("id")
		global.DB.Where("id = ?", id).Delete(&model.File{})

		middleware.Success(c, model.CODE_200, "", nil)
	})
}

// @Summary 删除文件（批量）
// @Tags File
// @version 1.0
// @Accept json
// @Param ids path []int true "文件编号数组"
// @Router /delFileBatch [delete]
func delFileBatch(r *gin.Engine) {
	r.DELETE("/delFileBatch", func(c *gin.Context) {
		var ids []int
		c.ShouldBindJSON(&ids)
		global.DB.Where("id in ?", ids).Delete(&model.File{})

		middleware.Success(c, model.CODE_200, "", ids)
	})
}

// @Summary 更新文件信息
// @Tags File
// @version 1.0
// @Accept json
// @Param file_form query model.File true "文件信息"
// @Router /filePage [get]
func updateFile(r *gin.Engine) {
	r.POST("/updateFile", func(c *gin.Context) {
		var file_form model.File
		if err := c.ShouldBindJSON(&file_form); err != nil {
			middleware.Err(c, model.CODE_401, "请求出现异常，请稍后重试", nil)
			return
		}
		global.DB.Save(&file_form)

		middleware.Success(c, model.CODE_200, "", nil)
	})
}

// @Summary 文件上传
// @Tags File
// @version 1.0
// @Accept multipart/form-data
// @Param file formData file true "文件"
// @Router /uploadFile [post]
func uploadFile(r *gin.Engine) {
	r.POST("/uploadFile", middleware.JWTAuth(), func(c *gin.Context) {
		r.MaxMultipartMemory = 8 << 20
		upload_file, file_header, err := c.Request.FormFile("file")
		if err != nil {
			middleware.Err(c, model.CODE_500, "文件上传失败", nil)
		}

		var file model.File
		file.Name = file_header.Filename
		file.Size = file_header.Size
		file.Enable = true

		userId, _ := c.Get("userId")
		file.UserId = uint(userId.(uint))

		// UUID
		ext := path.Ext(file.Name)
		u1 := uuid.NewV4()
		// file.Name = u1.String() + file.Name
		file.Name = u1.String() + ext
		// file.Name = u1.String()

		// MD5
		md5 := md5.New()
		io.Copy(md5, upload_file)
		file.MD5 = hex.EncodeToString(md5.Sum(nil))

		// Content-Type
		// file_content, err := ioutil.ReadAll(upload_file)
		// file.Type = http.DetectContentType(file_content)
		file.Type = ext[1:]

		var stored_file model.File
		global.DB.Debug().Where("md5 = ?", file.MD5).First(&stored_file)

		if stored_file.ID != 0 {
			file.Url = stored_file.Url

			c.String(http.StatusOK, file.Url)
			// middleware.Success(c, model.CODE_200, "存在重复文件，保存成功", stored_file.Url)
		} else {
			c.SaveUploadedFile(file_header, "./storage/file/"+file.Name)
			file.Url = "http://" + global.Settings.Host + ":" + strconv.Itoa(global.Settings.Port) + "/downloadFile/" + file.Name

			// middleware.Success(c, model.CODE_200, "文件保存成功", file.Url)
			c.String(http.StatusOK, file.Url)

		}
		global.DB.Create(&file)
	})
}

// @Summary 文件下载
// @Tags File
// @version 1.0
// @Accept json
// @Param name path int true "文件名称"
// @Router /downloadFile/:name [get]
func downloadFile(r *gin.Engine) {
	r.GET("/downloadFile/:name", func(c *gin.Context) {
		name := c.Param("name")

		file_content, err := ioutil.ReadFile("./storage/file/" + name)
		if err != nil {
			middleware.Err(c, model.CODE_404, "no such file", nil)
		}

		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s", name))
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
		c.Data(http.StatusOK, "", file_content)
	})
}
