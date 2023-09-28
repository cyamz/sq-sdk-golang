## 基本使用

```go
package main

import (
	"fmt"
	"sqapi"
)

func main() {
	SqApi := sqapi.NewSqApi(clientCode, clientSecret)

	params := map[string]interface{}{
		"consignee": map[string]interface{}{
			"name":     "Raul Ortega",
			"phone":    "7086764422",
			"state":    "TX",
			"province": "TX",
			"city":     "Fort Hood",
			"address1": "51337 Jumano Ct",
			"address2": "unit 1",
			"postcode": "76544-1161",
			"country":  "US",
		},
		"order":  orderNumber,
		"length": 10,
		"width":  20,
		"height": 30,
		"goods": []map[string]interface{}{
			{
				"name":     "鞋子",
				"name_en":  "Shoes",
				"hscode":   "98765432",
				"price":    5,
				"quantity": 1,
			},
			{
				"name":     "衣服",
				"name_en":  "Dresses",
				"hscode":   "12345678",
				"price":    2.5,
				"quantity": 2,
			},
		},
		"declare":      10,
		"weight":       0.5,
		"weight_unit":  "kg",
		"length_unit":  "cm",
		"channel_id":   channelId,
		"warehouse_id": warehouseId,
		"remark":       "测试单",
		"port":         "LAX",
	}

	res, _ := SqApi.Request("CreateBulkOrder", params)

	fmt.Println(res)
}
```