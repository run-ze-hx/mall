package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"log"
	"mall/kitex_gen/product"
	"mall/model"
	"mall/rpc/product/productdao/mysql"
	"mall/rpc/product/productdao/redis"
	"time"
)

// ProductCatalogServiceImpl implements the last service interface defined in the IDL.
type ProductCatalogServiceImpl struct{}

// ListProducts by Page PageSize CategoryName
// ListProducts implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) ListProducts(ctx context.Context, req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	var modelproducts []*model.Product

	modelproducts, err = mysql.ListProducts(req.GetPage(), req.GetPageSize(), req.GetCategoryName())
	if err != nil {
		return nil, kerrors.NewBizStatusError(500601, err.Error())
	}

	//用数据库查询的数据填充 serviceProducts
	var serviceProducts []*product.Product
	for _, modelPdt := range modelproducts {
		pdt := &product.Product{
			Id:          modelPdt.ID,
			Name:        modelPdt.Name,
			Description: modelPdt.Description,
			Picture:     modelPdt.Picture,
			Price:       modelPdt.Price,
			Categories:  modelPdt.Category,
		}
		serviceProducts = append(serviceProducts, pdt)
	}

	return &product.ListProductsResp{Products: serviceProducts}, nil
}

// GetProduct by id
// GetProduct implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) GetProduct(ctx context.Context, req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	var pdt *product.Product //保存返回值
	pt := &model.Product{}
	// 尝试从缓存中获取商品,如果缓存中没有，自动去mysql查询，并把该商品缓存
	pt, err = redis.GetById(ctx, req.Id)
	if err != nil {
		return &product.GetProductResp{}, kerrors.NewBizStatusError(500602, err.Error())
	}

	pdt = &product.Product{
		Id:          pt.ID,
		Name:        pt.Name,
		Description: pt.Description,
		Picture:     pt.Picture,
		Price:       pt.Price,
		Categories:  pt.Category,
	}

	return &product.GetProductResp{Product: pdt}, nil
}

// SearchProducts by query
// SearchProducts implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) SearchProducts(ctx context.Context, req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	var products []*product.Product
	var modelPdt []*model.Product
	modelPdt, err = mysql.SearchProducts(req.GetQuery())
	if err != nil {
		return nil, kerrors.NewBizStatusError(500603, err.Error())
	}

	//用数据库查询的数据填充 service products
	for _, Pdt := range modelPdt {
		pdt := &product.Product{
			Id:          Pdt.ID,
			Name:        Pdt.Name,
			Description: Pdt.Description,
			Picture:     Pdt.Picture,
			Price:       Pdt.Price,
			Categories:  Pdt.Category,
		}
		products = append(products, pdt)
	}

	return &product.SearchProductsResp{Results: products}, nil
}

// CreateProduct implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) CreateProduct(ctx context.Context, req *product.CreateProductReq) (resp *product.CreateProductResp, err error) {
	var modelPdt *model.Product
	modelPdt, err = mysql.CreateProduct(req.Name, req.Description, req.Price, req.Picture, req.Categories)
	if err != nil {
		return nil, kerrors.NewBizStatusError(500604, err.Error())
	}

	//用数据库查询数据填充 service product
	pdt := &product.Product{
		Id:          modelPdt.ID,
		Name:        modelPdt.Name,
		Description: modelPdt.Description,
		Picture:     modelPdt.Picture,
		Price:       modelPdt.Price,
		Categories:  modelPdt.Category,
	}
	return &product.CreateProductResp{Product: pdt}, nil
}

// UpdateProduct implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) UpdateProduct(ctx context.Context, req *product.UpdateProductReq) (resp *product.UpdateProductResp, err error) {

	return
}

// DeleteProduct implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) DeleteProduct(ctx context.Context, req *product.DeleteProductReq) (resp *product.DeleteProductResp, err error) {

	return
}

func startCacheRefresher(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cacheHotProducts(ctx)
		}
	}
}

func cacheHotProducts(ctx context.Context) {
	// 查询数据库中的热门商品（假设按访问量排序）
	var hotProducts []*model.Product
	err := mysql.DB.Order("views DESC").Limit(10).Find(&hotProducts).Error
	if err != nil {
		log.Printf("查询热门商品失败: %v", err)
		return
	}

	// 将热门商品缓存到 Redis
	for _, product := range hotProducts {
		productJSON, err := json.Marshal(product)
		if err != nil {
			log.Printf("序列化商品失败: %v", err)
			continue
		}
		err = redis.RedisClient.Set(ctx, fmt.Sprintf("product_%d", product.ID), productJSON, 5*time.Minute).Err()
		if err != nil {
			log.Printf("缓存商品失败: %v", err)
		}
	}
}
