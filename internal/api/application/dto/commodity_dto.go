package dto

import "time"

// http 请求结构

// CommoditySearch 商品搜索请求, --用Gin的BindQuery把URL参数绑定到结构体
type CommoditySearch struct {
	Keyword string `form:"keyword" binding:"required"`
	// 下面两个参数由Pagination组件使用
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"max=100"`
}

// http 响应

type HierarchicCommodityCategory struct {
	ID            int64                          `json:"id"`
	Level         int                            `json:"level"`
	ParentId      int64                          `json:"parent_id"`
	Name          string                         `json:"name"`
	IconImg       string                         `json:"icon_img"`
	Rank          int                            `json:"rank"`
	SubCategories []*HierarchicCommodityCategory `json:"sub_categories"` // 分类的子分类
}

type CommodityCategory struct {
	ID       int64  `json:"id"`
	Level    int    `json:"level"`
	ParentId int64  `json:"parent_id"`
	Name     string `json:"name"`
	IconImg  string `json:"icon_img"`
	Rank     int    `json:"rank"`
}

type Commodity struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Intro         string    `json:"intro"`
	CategoryId    int64     `json:"category_id"`
	CoverImg      string    `json:"cover_img"`
	Images        string    `json:"images"`
	DetailContent string    `json:"detail_content"`
	OriginalPrice int       `json:"original_price"`
	SellingPrice  int       `json:"selling_price"`
	StockNum      int       `json:"stock_num"`
	Tag           string    `json:"tag"`
	SellStatus    int       `json:"sell_status"`
	CreatedAt     time.Time `json:"created_at"`
}

type CommodityListElem struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Intro         string `json:"intro"`
	CategoryId    int64  `json:"category_id"`
	CoverImg      string `json:"cover_img"`
	OriginalPrice int    `json:"original_price"`
	SellingPrice  int    `json:"selling_price"`
	Tag           string `json:"tag"`
	SellStatus    int    `json:"sell_status"`
	CreatedAt     string `json:"created_at"`
}
