package resp

import (
	"github.com/Cospk/go-mall/pkg/config"
	"github.com/gin-gonic/gin"
	"strconv"
)

type PageInfo struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

func GetPageInfo(c *gin.Context) *PageInfo {
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	if pageNum <= 0 {
		pageNum = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	if pageSize <= config.AppConfig.PageInfo.DefaultSize {
		pageSize = config.AppConfig.PageInfo.DefaultSize
	}

	if pageSize > config.AppConfig.PageInfo.MaxSize {
		pageSize = config.AppConfig.PageInfo.MaxSize
	}
	return &PageInfo{
		PageNum:  pageNum,
		PageSize: pageSize,
	}
}
