CREATE TABLE products (
                          id INT PRIMARY KEY AUTO_INCREMENT,
                          name VARCHAR(255) NOT NULL, -- 商品名称
                          description TEXT,           -- 商品描述
                          picture VARCHAR(255),       -- 商品图片路径
                          price DECIMAL(10, 2),       -- 商品价格
                          category VARCHAR(255)       -- 商品类别
);