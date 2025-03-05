package model

type Product struct {
	ID          uint32  `gorm:"primary_key;auto_increment" json:"id"`
	Name        string  `gorm:"column:name;type:varchar(255)" json:"name"`
	Description string  `gorm:"column:description;type:text" json:"description"`
	Picture     string  `gorm:"column:picture;type:varchar(255)" json:"picture"`
	Price       float32 `gorm:"column:price;type:decimal(10,2)" json:"price,omitempty"`
	Category    string  `gorm:"column:category;type:varchar(255)" json:"category,omitempty"`
}

/*
type stringarray []string

// Scan 实现 sql.Scanner 接口， 从数据库中读取 JSON 数据并反序列化为 []string 类型
func (s *stringarray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("Failed to unmarshal JSONB value: %v", value))
	}
	var result []string
	err := json.Unmarshal(bytes, &result)
	*s = stringarray(result)
	return err
}

// Value 用于将 []string 数据序列化为 JSON 格式，来存储到数据库中
func (s stringarray) Value() (driver.Value, error) {

	return json.Marshal(s)
}
*/
