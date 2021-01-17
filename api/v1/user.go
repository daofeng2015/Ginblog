package v1

import (
	"ginblog/api/server"
	"ginblog/model"
	"ginblog/utils/errmsg"
	"ginblog/utils/validator"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var code int

// 添加用户
func AddUser(c *gin.Context) {
	var data model.User
	var msg string
	_ = c.ShouldBindJSON(&data)

	msg, code = validator.Validate(&data)
	if code != errmsg.SUCCSE {
		c.JSON(
			http.StatusOK, &server.Message1{
				Code:  code,
				Message: msg,
			},
		)
		c.Abort()
	}else {
		code = model.CheckUser(data.Username)
		if code == errmsg.SUCCSE {
			model.CreateUser(&data)
		}
		if code == errmsg.ERROR_USERNAME_USED {
			code = errmsg.ERROR_USERNAME_USED
		}
		c.JSON(
			http.StatusOK, &server.Message1{
				Code:  code,
				Message: errmsg.GetErrMsg(code),
			},
		)
	}
}

// 查询单个用户
func GetUserInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var maps = make(map[string]interface{})
	data, code := model.GetUser(id)
	maps["username"] = data.Username
	maps["role"] = data.Role
	c.JSON(
		http.StatusOK, &server.Message1{
			Code:  code,
			Data:    map[string]interface{}{
				"list":data,
				"total":1,
			},
			Message: errmsg.GetErrMsg(code),
		},
	)

}

// 查询用户列表
func GetUsers(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	username := c.Query("username")

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	data, total := model.GetUsers(username, pageSize, pageNum)

	code = errmsg.SUCCSE
	c.JSON(
		http.StatusOK, &server.Message1{
			Code:  code,
			Data:    map[string]interface{}{
				"list":data,
				"total":total,
			},
			Message: errmsg.GetErrMsg(code),
		},
	)
}

// 编辑用户
func EditUser(c *gin.Context) {
	var data model.User
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)

	code = model.CheckUpUser(id, data.Username)
	if code == errmsg.SUCCSE {
		model.EditUser(id, &data)
	}
	if code == errmsg.ERROR_USERNAME_USED {
		c.Abort()
	}

	c.JSON(
		http.StatusOK, &server.Message1{
			Code:  code,
			Message: errmsg.GetErrMsg(code),
		},
	)
}
func ChangeUserPassword(c *gin.Context) {
	var data model.User
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)

	code = model.ChangePassword(id, &data)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// 删除用户
func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	code = model.DeleteUser(id)

	c.JSON(
		http.StatusOK, &server.Message1{
			Code:  code,
			Message: errmsg.GetErrMsg(code),
		},
	)
}
