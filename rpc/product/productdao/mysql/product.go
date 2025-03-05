package mysql

import (
	"fmt"

	"mall/model"
)

func ListProducts(page int32, pageSize int64, category string) ([]*model.Product, error) {
	var products []*model.Product
	DB.Where("category = ?", category).Offset(int(page-1) * int(pageSize)).Limit(int(pageSize)).Find(&products)
	return products, nil
}

func GetProduct(id uint32) (*model.Product, error) {
	var pdt model.Product
	DB.First(&pdt, id)
	return &pdt, nil
}

func SearchProducts(query string) ([]*model.Product, error) {
	var products []*model.Product
	DB.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").Find(&products)
	return products, nil
}

func CreateProduct(name string, description string, price float32, picture string, category string) (*model.Product, error) {
	product := &model.Product{Name: name, Description: description, Price: price, Picture: picture, Category: category}
	// 执行数据库插入
	if err := DB.Create(product).Error; err != nil {
		return nil, fmt.Errorf("创建商品失败: %w", err)
	}
	return product, nil
}

func UpdateProduct() {

}
