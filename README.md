# 一、项目介绍
* 本项目是一个基于微服务架构的电商系统，包含用户管理、商品管理、购物车、订单处理和支付等功能模块。
* 项目服务地址：localhost:8889（API 网关服务地址）
  
https://github.com/run-ze-hx/mall







# 二、项目实现
## 2.1 技术选型与相关开发文档
技术选型：  

微服务框架：CloudWego Kitex  

服务注册与发现：Etcd  

数据库：MySQL（使用 GORM 作为 ORM 框架）,Redis  

HTTP 框架：CloudWego Hertz  

身份验证：JWT（JSON Web Token）  



## 2.2 架构设计
本项目采用微服务架构，包含以下主要服务：  

API 网关服务：负责处理外部 HTTP 请求，路由到相应的微服务。  

用户服务：负责用户注册、登录和管理。  

商品服务：负责商品的创建、查询和管理。  

购物车服务：负责用户的购物车操作，包括添加商品、查询购物车和清空购物车。  

订单服务：负责订单的创建、查询和管理。  

支付服务：负责处理支付请求，生成支付记录。  

架构图
[图片]



api网关<br>
api<br>
├─ handler<br>
│    ├─ auth.go<br>
│    ├─ cart_handler.go<br>
│    ├─ checkout_handler.go<br>
│    ├─ order_handler.go<br>
│    ├─ payment_handler.go<br>
│    ├─ product_handler.go<br>
│    └─ user_handle.go<br>
├─ main.go<br>
├─ middleware<br>
│    └─ jwtauth.go<br>
└─ router<br>
       └─ router.go<br>
<br>
微服务模块<br>

checkout<br>
├─ build.sh<br>
├─ handler.go<br>
├─ kitex_info.yaml<br>
├─ main.go<br>
├─ rpc_cli<br>
│    ├─ cart_cli.go<br>
│    ├─ order_cli.go<br>
│    ├─ payment_cli.go<br>
│    └─ product_cli.go<br>
└─ script<br>
       └─ bootstrap.sh<br>
<br>
数据库设计<br>

[图片]

## 2.3 项目代码介绍
主要模块：<br>
用户服务：实现用户注册、登录和管理功能。<br>
认证服务: 生成jwt，校验jwt 。<br>
商品服务：实现商品的创建、查询和管理功能。<br>
购物车服务：实现购物车的添加商品、查询和清空功能。<br>
订单服务：实现订单的创建、查询和管理功能。<br>
支付服务：实现支付请求的处理和支付记录的生成。<br>
结算服务：调用购物车,订单,支付,商品服务完成订单上所有商品的结算。<br>
代码结构：<br>
API 网关服务：mall/api<br>
认证服务：mall/rpc/auth<br>
用户服务：mall/rpc/user<br>
商品服务：mall/rpc/product<br>
购物车服务：mall/rpc/cart<br>
订单服务：mall/rpc/order<br>
支付服务：mall/rpc/payment<br>
结算服务：mall/rpc/checkout<br>
mall<br>
├─ .idea<br>
├─ api  //api网关<br>
├─ go.mod<br>
├─ go.sum<br>
├─ idl  //接口定义<br>
├─ kitex_gen<br>
├─ model //数据表定义<br>
├─ rpc<br>
│    ├─ auth  //鉴权微服务<br>
│    ├─ cart   //购物车微服务<br>
│    ├─ checkout //结算微服务<br>
│    ├─ order//订单微服务<br>
│    ├─ payment//支付微服务<br>
│    ├─ product//产品微服务<br>
│    └─ user//用户微服务<br>
└─ utils<br>
       └─ response.go //网关统一返回结构<br>






# 三、测试结果
功能测试<br>
user<br>
用户注册<br>
Endpoint: POST /register<br>
[图片]
 用户登录<br>
Endpoint: POST /login<br>
[图片]

后续，均携带token。<br>
[图片]

product<br>
创建商品<br>
Endpoint: POST /product/create<br>
[图片]
查询商品详情<br>
Endpoint: GET /products/:id<br>
[图片]
查询商品列表<br>
Endpoint: GET /products/list<br>
[图片]
搜索商品<br>
Endpoint:/products/search<br>
[图片]
cart<br>
为购物车增加商品<br>
Endpoint:/cart/additem<br>
[图片]
查看购物车<br>
Endpoint:/cart/get<br>
[图片]
清空购物车<br>
Endpoint:/cart/empty<br>
[图片]
order<br>
创建订单<br>
Endpoint: POST /order/place<br>
[图片]

查询订单列表<br>
Endpoint: POST /order/listorder<br>
[图片]
标记订单为已支付<br>
[图片]
payment<br>
处理支付请求<br>
Endpoint: POST /payment/charge<br>
[图片]
checkout<br>
结算服务<br>
Endpoint：/checkout<br>
用例：<br>
{<br>
  "user_id": 5,<br>
  "firstname": "John",<br>
  "lastname": "Doe",<br>
  "email": "test@example.com",<br>
  "address": {<br>
    "street_address": "123 Main St",<br>
    "city": "City",<br>
    "state": "State",<br>
    "country": "Country",<br>
    "zip_code": "12345"<br>
  },<br>
  "credit_card": {<br>
    "credit_card_number": "4111111111111111",<br>
    "credit_card_cvv": 123,<br>
    "credit_card_expiration_year": 2025,<br>
    "credit_card_expiration_month": 3<br>
  }<br>
}<br>
[图片]


性能分析报告<br>
api网关<br>
Cpu<br>
[图片]
url.ParseRequestURI 和 url.escape 是主要的 CPU 使用热点，各占总 CPU 时间的 50%。

内存<br>
[图片]
goroutine<br>
[图片]

Product  <br>
cpu<br>
[图片]
[图片]



# 四、项目总结与反思
1. 目前仍存在的问题：部分服务的错误处理不够完善，可能导致用户收到不友好的错误信息。<br>
2. 已识别出的优化项：服务拆分：进一步拆分服务，减少服务之间的依赖关系。

 架构演进的可能性：.<br>

引入消息队列：使用 Kafka或 RabbitMQ 处理异步任务，削减流量高峰。<br>

 项目过程中的反思与总结：完善的错误处理很重要，有利于快速发现问题。<br>
   






