一、项目介绍
本项目是一个基于微服务架构的电商系统，包含用户管理、商品管理、购物车、订单处理和支付等功能模块。
项目服务地址：localhost:8889（API 网关服务地址）
https://github.com/run-ze-hx/mall







二、项目实现
2.1 技术选型与相关开发文档
技术选型：
微服务框架：CloudWego Kitex
服务注册与发现：Etcd
数据库：MySQL（使用 GORM 作为 ORM 框架）,Redis
HTTP 框架：CloudWego Hertz
身份验证：JWT（JSON Web Token）


2.2 架构设计
本项目采用微服务架构，包含以下主要服务：
API 网关服务：负责处理外部 HTTP 请求，路由到相应的微服务。
用户服务：负责用户注册、登录和管理。
商品服务：负责商品的创建、查询和管理。
购物车服务：负责用户的购物车操作，包括添加商品、查询购物车和清空购物车。
订单服务：负责订单的创建、查询和管理。
支付服务：负责处理支付请求，生成支付记录。
架构图
[图片]



api网关
api
├─ handler
│    ├─ auth.go
│    ├─ cart_handler.go
│    ├─ checkout_handler.go
│    ├─ order_handler.go
│    ├─ payment_handler.go
│    ├─ product_handler.go
│    └─ user_handle.go
├─ main.go
├─ middleware
│    └─ jwtauth.go
└─ router
       └─ router.go
微服务模块
checkout
├─ build.sh
├─ handler.go
├─ kitex_info.yaml
├─ main.go
├─ rpc_cli
│    ├─ cart_cli.go
│    ├─ order_cli.go
│    ├─ payment_cli.go
│    └─ product_cli.go
└─ script
       └─ bootstrap.sh

数据库设计
[图片]

2.3 项目代码介绍
主要模块：
用户服务：实现用户注册、登录和管理功能。
认证服务: 生成jwt，校验jwt 。
商品服务：实现商品的创建、查询和管理功能。
购物车服务：实现购物车的添加商品、查询和清空功能。
订单服务：实现订单的创建、查询和管理功能。
支付服务：实现支付请求的处理和支付记录的生成。
结算服务：调用购物车,订单,支付,商品服务完成订单上所有商品的结算。
代码结构：
API 网关服务：mall/api
认证服务：mall/rpc/auth
用户服务：mall/rpc/user
商品服务：mall/rpc/product
购物车服务：mall/rpc/cart
订单服务：mall/rpc/order
支付服务：mall/rpc/payment
结算服务：mall/rpc/checkout
mall
├─ .idea
├─ api  //api网关
├─ go.mod
├─ go.sum
├─ idl  //接口定义
├─ kitex_gen
├─ model //数据表定义
├─ rpc
│    ├─ auth  //鉴权微服务
│    ├─ cart   //购物车微服务
│    ├─ checkout //结算微服务
│    ├─ order//订单微服务
│    ├─ payment//支付微服务
│    ├─ product//产品微服务
│    └─ user//用户微服务
└─ utils
       └─ response.go //网关统一返回结构






三、测试结果
功能测试
user
用户注册
Endpoint: POST /register
[图片]
 用户登录
Endpoint: POST /login
[图片]

后续，均携带token。
[图片]

product
创建商品
Endpoint: POST /product/create
[图片]
查询商品详情
Endpoint: GET /products/:id
[图片]
查询商品列表
Endpoint: GET /products/list
[图片]
搜索商品
Endpoint:/products/search
[图片]
cart
为购物车增加商品
Endpoint:/cart/additem
[图片]
查看购物车
Endpoint:/cart/get
[图片]
清空购物车
Endpoint:/cart/empty
[图片]
order
创建订单
Endpoint: POST /order/place
[图片]

查询订单列表
Endpoint: POST /order/listorder
[图片]
标记订单为已支付
[图片]
payment
处理支付请求
Endpoint: POST /payment/charge
[图片]
checkout
结算服务
Endpoint：/checkout
用例：
{
  "user_id": 5,
  "firstname": "John",
  "lastname": "Doe",
  "email": "test@example.com",
  "address": {
    "street_address": "123 Main St",
    "city": "City",
    "state": "State",
    "country": "Country",
    "zip_code": "12345"
  },
  "credit_card": {
    "credit_card_number": "4111111111111111",
    "credit_card_cvv": 123,
    "credit_card_expiration_year": 2025,
    "credit_card_expiration_month": 3
  }
}
[图片]


性能分析报告
api网关
Cpu
[图片]
url.ParseRequestURI 和 url.escape 是主要的 CPU 使用热点，各占总 CPU 时间的 50%。

内存
[图片]
goroutine
[图片]

Product  
cpu
[图片]
[图片]



四、项目总结与反思
1. 目前仍存在的问题：部分服务的错误处理不够完善，可能导致用户收到不友好的错误信息。
2. 已识别出的优化项：服务拆分：进一步拆分服务，减少服务之间的依赖关系。
3. 架构演进的可能性：.
引入消息队列：使用 Kafka或 RabbitMQ 处理异步任务，削减流量高峰。
4. 项目过程中的反思与总结：完善的错误处理很重要，有利于快速发现问题。






