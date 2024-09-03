package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/template"
)

// 定义关联项数据结构
type RelatedItem struct {
	FunctionName string // 关联项获取函数名，例如 GetOrderItemsByOrderNoList
	Key          string // 关联项映射的 key
	KeyField     string // 关联项在主对象中的字段，例如 OrderNo
	FieldName    string // 映射到最终响应中的字段名，例如 Items
	ReturnType   string // 关联函数的返回类型，如 []*OrderItemRes
}

// Define a struct for template data
type QueryCondition struct {
	FieldName   string
	DBFieldName string
	ZeroValue   string
}

type CombinedTemplateData struct {
	StructName            string
	LowerStructName       string
	HasIsDeletedField     bool
	HasGetListFunction    bool
	HasGetMapFunction     bool
	HasGetListMapFunction bool
	QueryStructName       string
	KeyField              string
	DBFieldName           string
	RelatedItems          []RelatedItem
	Conditions            []QueryCondition
	HasDateRange          bool
}

// Helper function to make the first letter of a string lowercase
func lower(input string) string {
	return strings.ToLower(input[:1]) + input[1:]
}

// Check if the struct has an "is_deleted" field
func hasIsDeletedField(fields []string) bool {
	for _, field := range fields {
		if strings.ToLower(field) == "isdeleted" || strings.ToLower(field) == "is_deleted" {
			return true
		}
	}
	return false
}

// Helper to generate zero value for different types
func getZeroValue(fieldType string) string {
	switch fieldType {
	case "string":
		return `""`
	case "int", "int32", "int64", "float32", "float64":
		return "0"
	default:
		return "nil"
	}
}

// Load template file
func loadTemplate(filename string) (*template.Template, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return template.New("combined").Funcs(template.FuncMap{
		"lower": lower,
	}).Parse(string(content))
}

// Generate code based on struct and query struct
func GenerateCombinedCode(structName string, structType interface{}, queryStruct interface{}, DBFieldName string, keyField string, relatedItems []RelatedItem, templateFile string) string {
	val := reflect.TypeOf(structType)
	if val.Kind() != reflect.Struct {
		return "Error: Input is not a struct."
	}

	fields := make([]string, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		fields[i] = val.Field(i).Name
	}

	// Check if the query struct has conditions
	queryVal := reflect.TypeOf(queryStruct)
	conditions := []QueryCondition{}
	hasDateRange := false

	if queryVal.Kind() == reflect.Struct {
		for i := 0; i < queryVal.NumField(); i++ {
			field := queryVal.Field(i)
			fieldName := field.Name
			fieldType := field.Type.String()
			tag := field.Tag.Get("json")
			if tag == "start_time" || tag == "end_time" {
				hasDateRange = true
			} else if tag != "" {
				conditions = append(conditions, QueryCondition{
					FieldName:   fieldName,
					DBFieldName: tag,
					ZeroValue:   getZeroValue(fieldType),
				})
			}
		}
	}

	data := CombinedTemplateData{
		StructName:            structName,
		LowerStructName:       lower(structName),
		HasIsDeletedField:     hasIsDeletedField(fields),
		HasGetListFunction:    true, // Set this based on your requirements
		HasGetMapFunction:     true, // Set this based on your requirements
		HasGetListMapFunction: true, // Set this based on your requirements
		QueryStructName:       queryVal.Name(),
		Conditions:            conditions,
		HasDateRange:          hasDateRange,
		RelatedItems:          relatedItems,
		KeyField:              keyField,
		DBFieldName:           DBFieldName,
	}

	tmpl, err := loadTemplate(templateFile)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		panic(err)
	}

	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("Error formatting code:", err)
		return buf.String() // Return unformatted code in case of error
	}

	return string(formattedCode)
}

func main() {
	// Example struct
	type User struct {
		ID        int
		Username  string
		Email     string
		Password  string
		CreatedAt string
		UpdatedAt string
		Uuid      string
	}

	// Define query parameter struct
	type ReqUserQueryParam struct {
		Email    string `json:"email"`    // 邮箱
		Phone    string `json:"phone"`    // 手机号
		Nickname string `json:"nickname"` // 昵称
		Sex      int    `json:"sex"`      // 性别
		Username string `json:"username"` // 用户名
		Status   int    `json:"status"`   // 状态
		Uuid     string `json:"uuid"`     // uuid
	}

	// 定义 Order 结构体和查询参数结构体
	type Order struct {
		ID        int
		UserID    string
		Status    string
		OrderNo   string
		CreatedAt string
		UpdatedAt string
	}

	relatedItems := []RelatedItem{
		{
			FunctionName: "GetOrderItemsByOrderNoList",
			Key:          "items",
			KeyField:     "OrderNo",
			FieldName:    "Items",
			ReturnType:   "[]*OrderItemRes",
		},
		{
			FunctionName: "GetOrderItemsByOrderNoList",
			Key:          "items",
			KeyField:     "OrderNo",
			FieldName:    "Items",
			ReturnType:   "[]*OrderItemRes",
		},
		// 可以继续添加其他关联项
		// {
		//     FunctionName: "GetOrderLogisticsByOrderNoList",
		//     Key:          "logistics",
		//     KeyField:     "OrderNo",
		//     FieldName:    "Logistics",
		//     ReturnType:   "*OrderLogisticsRes",
		// },
	}

	// code := GenerateCombinedCode("User", User{}, ReqUserQueryParam{}, "service_template.go.tpl")
	// fmt.Println(code)
	// 生成代码
	code := GenerateCombinedCode("Order", Order{}, ReqUserQueryParam{}, "order_no", "OrderNo", relatedItems, "service_template.go.tpl")
	fmt.Println(code)

	// Optional: Save to file
	err := os.WriteFile("user_service.go", []byte(code), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
