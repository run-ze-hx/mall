package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"log"
	"mall/kitex_gen/product"
	"mall/kitex_gen/product/productcatalogservice"
	"mall/model"
	"mall/utils"
	"net/http"
	"strconv"
)

var ProductCli productcatalogservice.Client

// ProductCliInit 创建 product 服务的 Client，并连接到 etcd
func ProductCliInit(r discovery.Resolver) {
	pc, err := productcatalogservice.NewClient("product",
		client.WithResolver(r), // 使用 etcd 注册中心
	)
	if err != nil {
		log.Fatal(err)
	}
	ProductCli = pc
}

// 用来接受http请求传递的参数
type list struct {
	Page         int32  `json:"page,omitempty"`
	PageSize     int64  `json:"pagesize,omitempty"`
	CategoryName string `json:"categoryname,omitempty"`
}

func ListProducts(c context.Context, ctx *app.RequestContext) {

	body, err := ctx.Body()
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Response{Code: 4003001, Msg: "get body failed:" + err.Error()})
		return
	}

	var List list
	err = json.Unmarshal(body, &List)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Response{Code: 5003002, Msg: "json unmarshal failed:" + err.Error()})
		return
	}
	//用从请求拿到的参数填充rpc函数的参数
	req := &product.ListProductsReq{
		Page:         List.Page,
		PageSize:     List.PageSize,
		CategoryName: List.CategoryName,
	}
	resp, err1 := ProductCli.ListProducts(c, req)
	if err1 != nil {
		ctx.JSON(5003003, err1.Error())
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{Data: resp})

}

func GetProduct(c context.Context, ctx *app.RequestContext) {

	idStr := ctx.Param("id")
	//string转换成int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Response{Code: 4003003, Msg: "transform failed，invalid id"})

		return
	}
	//rpc_cli
	req := &product.GetProductReq{Id: uint32(id)}
	resp, err1 := ProductCli.GetProduct(c, req)
	if err1 != nil {
		ctx.JSON(http.StatusOK, utils.Response{Code: 5003004, Msg: err1.Error()})
		return
	}
	ctx.JSON(http.StatusOK, utils.Response{Data: resp})
}

func SearchProduct(c context.Context, ctx *app.RequestContext) {
	query := ctx.Query("query")
	//rpc_cli
	req := &product.SearchProductsReq{Query: query}
	resp, err := ProductCli.SearchProducts(c, req)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Response{Code: 5003005, Msg: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, utils.Response{Data: resp.Results})
}

func CreateProduct(c context.Context, ctx *app.RequestContext) {
	body, err := ctx.Body()
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Response{Code: 5003006, Msg: "get body failed: " + err.Error()})
		return
	}

	var modelpdt *model.Product
	err = json.Unmarshal(body, &modelpdt)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Response{Code: 5003007, Msg: "json Unmarshal error:" + err.Error()})
		return
	}
	//用拿到的请求参数给rpc函数req赋值
	req := &product.CreateProductReq{
		Name:        modelpdt.Name,
		Description: modelpdt.Description,
		Picture:     modelpdt.Picture,
		Price:       modelpdt.Price,
		Categories:  modelpdt.Category,
	}
	resp, err1 := ProductCli.CreateProduct(c, req)
	if err1 != nil {
		ctx.JSON(http.StatusOK, utils.Response{Code: 5003008, Msg: err1.Error()})
		return
	}
	ctx.JSON(http.StatusOK, utils.Response{Data: resp.Product, Msg: "create product success"})
}

func UpdateProduct(c context.Context, ctx *app.RequestContext) {
	body, err := ctx.Body()
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Response{Code: 5003009, Msg: "get body failed: " + err.Error()})
		return
	}
	var modelpdt *model.Product
	err = json.Unmarshal(body, &modelpdt)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Response{Code: 50030010, Msg: "Unmarshal error:" + err.Error()})
		return
	}
	req := &product.UpdateProductReq{
		Name:        &modelpdt.Name,
		Description: &modelpdt.Description,
		Picture:     &modelpdt.Picture,
		Price:       &modelpdt.Price,
		Categories:  modelpdt.Category,
	}
	resp, err1 := ProductCli.UpdateProduct(c, req)
	if err1 != nil {
		ctx.JSON(http.StatusOK, utils.Response{Code: 50030011, Msg: err1.Error()})
		return
	}
	ctx.JSON(http.StatusOK, utils.Response{Code: 200, Data: resp.Product})

}
func DeleteProduct(c context.Context, ctx *app.RequestContext) {

}
