package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"log"
	"mall/kitex_gen/cart"
	"mall/kitex_gen/cart/cartservice"
	"mall/utils"
	"net/http"
	"strconv"
)

// CartRequest 用于获取参数
type CartRequest struct {
	UserId    uint32 `form:"userid" binding:"required"`
	ProductId uint32 `form:"productid" binding:"required"`
	Quantity  int32  `form:"quantity" binding:"required,min=1"`
}

// CartItemResponse 用于返回购物车里的商品信息
type CartItemResponse struct {
	ProductID uint32 `json:"productid"` // 商品ID
	Quantity  int32  `json:"quantity"`  // 商品数量
}

var CartCli cartservice.Client

// CartCliInit 创建 cart 服务的 Client，并连接到 etcd
func CartCliInit(r discovery.Resolver) {
	cc, err := cartservice.NewClient("cart",
		client.WithResolver(r), // 使用 etcd 注册中心
	)
	if err != nil {
		log.Fatal(err)
	}
	CartCli = cc
}

func GetCart(c context.Context, ctx *app.RequestContext) {
	uidStr := ctx.PostForm("userid")
	var userId uint32
	//string 转uint32
	if u, err := strconv.ParseUint(uidStr, 10, 32); err != nil {
		// 转换失败
		ctx.JSON(200, utils.Response{
			Code: 4001001,
			Msg:  err.Error(),
		})
		return
	} else {
		userId = uint32(u)
	}

	//调用cart微服务
	Req := &cart.GetCartReq{UserId: userId}
	resp, err := CartCli.GetCart(c, Req)
	if err != nil {
		ctx.JSON(200, utils.Response{
			Code: 5001002,
			Msg:  err.Error()})
		return
	}
	//返回商品信息
	var cartItemsResponse []CartItemResponse
	for _, item := range resp.Cart.GetItems() {
		cartItemsResponse = append(cartItemsResponse, CartItemResponse{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}
	ctx.JSON(200, utils.Response{Data: cartItemsResponse})
}

func Empty(c context.Context, ctx *app.RequestContext) {
	var creq CartRequest
	if err := ctx.Bind(&creq); err != nil {
		ctx.JSON(200, utils.Response{
			Code: 4001003,
			Msg:  err.Error(),
		})
		return
	}
	req := &cart.EmptyCartReq{UserId: creq.UserId}
	_, err := CartCli.EmptyCart(c, req)
	if err != nil {
		ctx.JSON(200, utils.Response{
			Code: 5001004,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(200, utils.Response{Msg: "清空购物车成功！"})
}

func AddItem(c context.Context, ctx *app.RequestContext) {
	var creq CartRequest
	if err := ctx.Bind(&creq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Code: 4001005,
			Msg:  err.Error(),
		})
		return
	}
	if creq.UserId <= 0 || creq.ProductId <= 0 || creq.Quantity <= 0 {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Code: 4001006,
			Msg:  "参数不合法",
		})
		return
	}

	var ca = &cart.CartItem{
		ProductId: creq.ProductId,
		Quantity:  creq.Quantity,
	}

	//调用cart 微服务
	req := &cart.AddItemReq{UserId: creq.UserId, Item: ca}
	_, err := CartCli.AddItem(c, req)
	if err != nil {
		ctx.JSON(200, utils.Response{
			Code: 5001007,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(200, utils.Response{Msg: "为购物车增加商品成功"})
}
