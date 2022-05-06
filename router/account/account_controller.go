package account

import (
	"bytes"
	"crypto/md5"
	"e-gov/global"
	"e-gov/middleware"
	"e-gov/model"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/tealeg/xlsx"
	"golang.org/x/crypto/bcrypt"
)

// @Summary 账号登录
// @Tags Account
// @version 1.0
// @Accept application/json
// @Param user body model.UserLoginTemp true "账号信息"
// @Router /login [post]
func login(r *gin.Engine) {
	r.POST("/login", func(c *gin.Context) {
		var inp_login model.User
		if err := c.ShouldBindJSON(&inp_login); err != nil {
			middleware.Err(c, model.CODE_401, "请求出现异常，请稍后重试", nil)
			return
		}

		var query_users model.User
		// select * from users where id = inp_login.Username
		global.DB.Where("username = ?", inp_login.Username).First(&query_users)
		err := bcrypt.CompareHashAndPassword([]byte(query_users.Password), []byte(inp_login.Password))
		if err != nil {
			middleware.Err(c, model.CODE_600, "用户名或密码错误，请重新登录", nil)
			return
		}

		var query_roles model.Role
		// select * from role where flag = query_roles.Flag
		global.DB.Where("flag = ?", query_users.Role).First(&query_roles)

		var query_role_menus []model.RoleMenu
		// select * from role_menu where role = query_roles.Flag
		global.DB.Where("role_id = ?", query_roles.ID).Find(&query_role_menus)

		// struct slice
		var user_query_menus []model.Menu
		// var user_query_menus = make([]model.Menu, 30)

		for count := 0; count < len(query_role_menus); count++ {
			var query_menus model.Menu

			global.DB.Where("id = ?", query_role_menus[count].MenuId).First(&query_menus)
			// user_query_menus = append(user_query_menus, query_menus)

			// if exist parent
			// cur_pid := user_query_menus[i].Pid
			cur_pid := query_menus.Pid
			if cur_pid != "" {
				cur_pid_num, _ := strconv.Atoi(cur_pid)
				cur_pid_num = cur_pid_num - 1

				last := len(user_query_menus) - 1
				user_query_menus[last].Children = append(user_query_menus[last].Children, query_menus)

			} else {
				user_query_menus = append(user_query_menus, query_menus)
			}
		}

		var user_auth model.UserAuth
		user_auth.User = query_users
		user_auth.Token = generateToken(c, query_users)
		user_auth.Menus = user_query_menus

		middleware.Success(c, model.CODE_200, "登录成功", user_auth)
	})
}

func generateToken(c *gin.Context, user model.User) string {
	j := middleware.NewJWT()
	claims := middleware.CustomClaims{
		ID:       uint(user.ID),
		Username: user.Username,
		Role:     user.Role,
		// AuthorityId: uint(user.Identity),
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),        // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 60*60*24*30), // 签名过期时间
			Issuer:    "grezkl-noir",                          // 签名发布者
		},
	}

	token, err := j.CreateToken(claims)
	if err != nil {
		middleware.Err(c, model.CODE_401, "token生成失败，请重试", nil)
		return ""
	}
	return token
}

// @Summary 账号注册
// @Tags Account
// @version 1.0
// @Accept json
// @Param user body model.UserRegisterTemp true "账号信息"
// @Router /register [post]
func register(r *gin.Engine) {
	r.POST("/register", func(c *gin.Context) {
		var cmt_reg model.User
		if err := c.ShouldBindJSON(&cmt_reg); err != nil {
			middleware.Err(c, model.CODE_401, "请求出现异常，请稍后重试", nil)
			return
		}

		// 账号名称已存在
		if isUsernameExist(cmt_reg.Username) {
			middleware.Err(c, model.CODE_600, "该用户已存在", nil)
			return
		}

		// default role
		cmt_reg.Role = "ROLE_USER"

		// bcrypt encode
		hash, err := bcrypt.GenerateFromPassword(
			[]byte(cmt_reg.Password),
			bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
		}
		encodePWD := string(hash)
		cmt_reg.Password = encodePWD

		global.DB.Create(&cmt_reg)
		middleware.Success(c, model.CODE_200, "注册成功", nil)

	})
}

func isUsernameExist(username string) bool {
	var user model.User
	global.DB.Where("username = ?", username).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

// @Summary 修改密码
// @Tags Account
// @version 1.0
// @Accept json
// @Param username body string true "账号"
// @Param password body string true "密码"
// @Router /updatePWD [post]
func updatePWD(r *gin.Engine) {
	r.POST("/updatePWD", func(c *gin.Context) {
		var form model.User
		// if err := c.Bind(&form); err != nil {
		if err := c.ShouldBind(&form); err != nil {
			// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			middleware.Err(c, model.CODE_401, "请求出现异常，请稍后重试", nil)
			return
		}

		fmt.Sprintln("New PWD: ")

		// bcrypt 加密处理
		hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}
		encodePWD := string(hash)

		global.DB.Model(&form).Where("username = ?", form.Username).Update("password", encodePWD)

		middleware.Success(c, model.CODE_200, "密码修改成功", nil)
	})
}

// @Summary 获取用户信息
// @Tags Account
// @version 1.0
// @Accept json
// @Param username path string true "账号"
// @Router /username/:name [get]
func username(r *gin.Engine) {
	r.GET("/username/:name", func(c *gin.Context) {
		name := c.Param("name")
		var user_query model.User
		global.DB.Where("username = ?", name).First(&user_query)

		middleware.Success(c, model.CODE_200, "", user_query)
	})
}

type PageForm struct {
	PageNum  int    `form:"pageNum" json:"pageNum" binding:"required"`
	PageSize int    `form:"pageSize" json:"pageSize" binding:"required"`
	Username string `form:"username" json:"username"`
	Email    string `form:"email" json:"email"`
	Address  string `form:"address" json:"address"`
}

type UserPageResp struct {
	Records []model.User `form:"records" json:"records"`
	Total   int64        `form:"total" json:"total"`
}

// @Summary 用户管理分页与搜索
// @Tags Account
// @version 1.0
// @Accept json
// @Param pageForm query PageForm true "分页规格与关键词信息"
// @Router /affairPage [get]
func userPage(r *gin.Engine) {
	r.GET("/userPage", func(c *gin.Context) {

		var page_form PageForm
		if c.ShouldBindQuery(&page_form) != nil {
			middleware.Err(c, model.CODE_500, "参数错误", nil)
			return
		}
		if page_form.PageNum <= 0 {
			page_form.PageNum = 1
		}

		var user_list_query []model.User
		username := page_form.Username
		email := page_form.Email
		address := page_form.Address

		var total int64
		// global.DB.Model(&user_list_query).Count(&total)

		if username != "" && email != "" && address != "" {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("username like ? AND email like ? AND address like ?", username, email, address).Find(&user_list_query).Count(&total)
		} else if username == "" && email != "" && address != "" {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("email like ? AND address like ?", email, address).Find(&user_list_query).Count(&total)
		} else if username != "" && email == "" && address != "" {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("username like ? AND address like ?", username, address).Find(&user_list_query).Count(&total)
		} else if username != "" && email != "" && address == "" {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("username like ? AND email like ?", username, email).Find(&user_list_query).Count(&total)
		} else if username == "" && email == "" && address != "" {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("address like ?", address).Find(&user_list_query).Count(&total)
		} else if username == "" && email != "" && address == "" {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("email like ?", email).Find(&user_list_query).Count(&total)
		} else if username != "" && email == "" && address == "" {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("username like ?", username).Find(&user_list_query).Count(&total)
		} else {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum - 1) * page_form.PageSize).Find(&user_list_query).Count(&total)
		}

		var page_resp UserPageResp
		page_resp.Records = user_list_query
		page_resp.Total = total

		middleware.Success(c, model.CODE_200, "", page_resp)
	})
}

// @Summary 获取用户列表
// @Tags Account
// @version 1.0
// @Accept json
// @Router /pullUserList [get]
func pullUserList(r *gin.Engine) {
	r.GET("/pullUserList", func(c *gin.Context) {
		var user_list_query []model.User
		global.DB.Find(&user_list_query)

		middleware.Success(c, model.CODE_200, "", user_list_query)
	})
}

// @Summary 根据编号查找用户
// @Tags Account
// @version 1.0
// @Accept json
// @Param id path int true "用户编号"
// @Router /pullUserById/{id} [get]
func pullUserById(r *gin.Engine) {
	r.GET("/pullUserById/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user_query model.User
		global.DB.Where("id = ?", id).First(&user_query)

		if user_query.ID != 0 {
			middleware.Success(c, model.CODE_200, "", user_query)
		} else {
			// middleware.Err(c, model.CODE_500, "该 id 查询无结果", nil)
		}
	})
}

// @Summary 注销用户（单个）
// @Tags Account
// @version 1.0
// @Accept json
// @Param id path int true "用户编号"
// @Router /deleteUserById/{id} [delete]
func deleteUserById(r *gin.Engine) {
	r.DELETE("/deleteUserById/:id", middleware.JWTAuth(), func(c *gin.Context) {
		id := c.Param("id")
		global.DB.Where("id = ?", id).Delete(&model.User{})

		middleware.Success(c, model.CODE_200, "", nil)
	})
}

// @Summary 注销用户（批量）
// @Tags Account
// @version 1.0
// @Accept json
// @Param ids path []int true "用户编号数组"
// @Router /delUserBatch [delete]
func delUserBatch(r *gin.Engine) {
	r.DELETE("/delUserBatch", middleware.JWTAuth(), func(c *gin.Context) {
		var ids []int
		c.ShouldBindJSON(&ids)
		global.DB.Where("id in ?", ids).Delete(&model.User{})

		middleware.Success(c, model.CODE_200, "", ids)
	})
}

// @Summary 获取访问量统计图
// @Tags Account
// @version 1.0
// @Accept json
// @Router /echartsMember [get]
func echartsMember(r *gin.Engine) {
	r.GET("/echartsMember", func(c *gin.Context) {
		var member []int
		member = append(member, 4)
		member = append(member, 5)
		member = append(member, 4)
		member = append(member, 3)
		middleware.Success(c, model.CODE_200, "", member)
	})
}

// @Summary 保存账号信息
// @Tags Account
// @version 1.0
// @Accept json
// @Param user_form body model.User true "保存账号信息"
// @Router /saveUser [post]
func saveUser(r *gin.Engine) {
	r.POST("/saveUser", middleware.JWTAuth(), func(c *gin.Context) {
		var user_form model.User
		if err := c.ShouldBindJSON(&user_form); err != nil {
			middleware.Err(c, model.CODE_401, "请求出现异常，请稍后重试", nil)
			return
		}
		global.DB.Save(&user_form)

		middleware.Success(c, model.CODE_200, "", nil)
	})
}

// @Summary 根据角色查询用户
// @Tags Account
// @version 1.0
// @Accept json
// @Param id path int true "角色"
// @Router /delUserBatch [get]
func pullUserByRole(r *gin.Engine) {
	r.GET("/pullUserByRole/:role", func(c *gin.Context) {
		role := c.Param("role")
		var user_form []model.User
		global.DB.Where("role = ?", role).Find(&user_form)

		middleware.Success(c, model.CODE_200, "", user_form)
	})
}

// @Summary 导入用户数据
// @Tags Account
// @version 1.0
// @Accept multipart/form-data
// @Param file formData file true "文件"
// @Router /userImport [POST]
func userImport(r *gin.Engine) {
	r.POST("/userImport", func(c *gin.Context) {
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

// @Summary 导出用户数据
// @Tags Account
// @version 1.0
// @Accept json
// @Router /userExport [get]
func userExport(r *gin.Engine) {
	r.GET("/userExport", func(c *gin.Context) {
		var user_data []model.User
		global.DB.Find(&user_data)

		var res []interface{}
		for _, user := range user_data {
			res = append(res, &model.User{
				// ID:        user.ID,
				Username:  user.Username,
				Password:  user.Password,
				Nickname:  user.Nickname,
				AvatarUrl: user.AvatarUrl,
				Email:     user.Email,
				Phone:     user.Phone,
				Address:   user.Address,
				Role:      user.Role,
				// PersonId:  user.PersonId,
			})
		}
		// content := ToExcel([]string{`id`, `username`, `password`, `nickname`, `avatar_url`, `email`, `phone`, `address`, `role`, `person_id`}, res)
		content := toExcel([]string{`username`, `password`, `nickname`, `avatar_url`, `email`, `phone`, `address`, `role`}, res)
		responseXls(c, content, "userData")
	})
}

func toExcel(titleList []string, dataList []interface{}) (content io.ReadSeeker) {
	// 生成一个新的文件
	file := xlsx.NewFile()
	// 添加sheet页
	sheet, _ := file.AddSheet("Sheet1")
	// 插入表头
	titleRow := sheet.AddRow()
	for _, v := range titleList {
		cell := titleRow.AddCell()
		cell.Value = v
	}
	// 插入内容
	for _, v := range dataList {
		row := sheet.AddRow()
		row.WriteStruct(v, -1)
	}

	var buffer bytes.Buffer
	_ = file.Write(&buffer)
	content = bytes.NewReader(buffer.Bytes())
	return
}

func responseXls(c *gin.Context, content io.ReadSeeker, fileTag string) {
	fileName := fmt.Sprintf("%s%s%s.xlsx", time.Now().Format("20060102_150405"), `-`, fileTag)
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeContent(c.Writer, c.Request, fileName, time.Now(), content)
}
