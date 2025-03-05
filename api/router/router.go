package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"mall/api/handler"
	"mall/api/middleware"
	"mall/kitex_gen/auth/authservice"
)

func Register(hz *server.Hertz, authClient authservice.Client) {

	// 公开路由
	publicGroup := hz.Group("/")
	{
		publicGroup.POST("/register", handler.UserRegister)
		publicGroup.POST("/login", handler.UserLogin)
	}
	// 注册 JWT 中间件
	hz.Use(middleware.JwtAuthMiddleware(authClient))
	//  JWT
	authenticatedGroup := hz.Group("/")
	{
		authenticatedGroup.GET("/products/list", handler.ListProducts)
		authenticatedGroup.GET("/products/:id", handler.GetProduct)
		authenticatedGroup.GET("/products/search", handler.SearchProduct)
		authenticatedGroup.POST("/product/create", handler.CreateProduct)

		authenticatedGroup.POST("/cart/additem", handler.AddItem)
		authenticatedGroup.POST("/cart/get", handler.GetCart)
		authenticatedGroup.POST("/cart/empty", handler.Empty)

		authenticatedGroup.POST("/checkout", handler.Checkout)

		authenticatedGroup.POST("/payment/charge", handler.Charge)

		authenticatedGroup.POST("/order/place", handler.PlaceOrder)
		authenticatedGroup.POST("/order/mark-paid", handler.MarkOrderPaid)
		authenticatedGroup.POST("/order/listorder", handler.ListOrder)

	}
}
