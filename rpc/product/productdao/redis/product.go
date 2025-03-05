package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"mall/model"
	"mall/rpc/product/productdao/mysql"
	"time"
)

func GetById(ctx context.Context, productId uint32) (*model.Product, error) {
	var pdt *model.Product
	cachedkey := fmt.Sprintf("%s-%s-%d", "product by id", productId)
	cachedResult, err := RedisClient.Get(ctx, cachedkey).Result()

	if err == nil {
		pdt = &model.Product{}
		cachedResultBytes := []byte(cachedResult)
		err := json.Unmarshal(cachedResultBytes, pdt)
		if err != nil {
			return &model.Product{}, err
		}
		return pdt, nil
	}

	//缓存查不到，用mysql查
	pdt, err = mysql.GetProduct(productId)
	if err != nil {
		return &model.Product{}, err
	}

	// 将查询结果序列化并存入缓存
	cachedValue, err := json.Marshal(pdt)
	if err != nil {
		return &model.Product{}, err
	}
	err = RedisClient.Set(ctx, cachedkey, string(cachedValue), 30*time.Minute).Err()
	if err != nil {
		return &model.Product{}, err
	}

	return pdt, nil
}
