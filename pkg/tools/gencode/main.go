package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sgin/model"
	"strings"
	"text/template"
)

// 定义关联项数据结构
type RelatedItem struct {
	FunctionName string // 关联项获取函数名，例如 GetOrderItemsByOrderNoList
	Key          string // 关联项映射的 key
	KeyField     string // 关联项在主对象中的字段，例如 OrderNo
	TypeName     string // 关联项的类型名，例如 OrderItem
	FieldName    string // 映射到最终响应中的字段名，例如 Items
	ReturnType   string // 关联函数的返回类型，如 []*OrderItemRes
}

// 定义查询条件数据结构
type QueryCondition struct {
	FieldName   string
	DBFieldName string
	ZeroValue   string
}

// 配置结构体，用于加载配置文件
type Config struct {
	ServiceDir    string `json:"serviceDir"`
	ControllerDir string `json:"controllerDir"`
	ReponseFile   string `json:"responseFile"`
	RouterFile    string `json:"routerFile"`
}

// Define a struct for combined template data
type CombinedTemplateData struct {
	StructName            string
	LowerStructName       string
	QueryStructName       string
	LowerQueryStructName  string
	ReqCreateStructName   string
	ResStructName         string
	IsResStruct           bool
	ModuleName            string // 新增字段，模块名
	HasIsDeletedField     bool
	HasGetListFunction    bool
	HasGetMapFunction     bool
	HasGetListMapFunction bool
	StructType            interface{}
	QueryStructType       interface{}
	ReqCreateStructType   interface{}
	ResStructType         interface{}
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

// Extract struct information and conditions
func extractStructInfo(data *CombinedTemplateData) {
	// 获取 StructName 和 LowerStructName
	structVal := reflect.TypeOf(data.StructType)
	data.StructName = structVal.Name()
	data.LowerStructName = lower(data.StructName)

	// 如果 ModuleName 为空，则使用 StructName
	if data.ModuleName == "" {
		data.ModuleName = data.StructName
	}

	// 获取 QueryStructName 和 LowerQueryStructName
	queryVal := reflect.TypeOf(data.QueryStructType)
	data.QueryStructName = queryVal.Name()
	data.LowerQueryStructName = lower(data.QueryStructName)

	createVal := reflect.TypeOf(data.ReqCreateStructType)
	data.ReqCreateStructName = createVal.Name()

	resVal := reflect.TypeOf(data.ResStructType)
	data.ResStructName = resVal.Name()

	// 判断 ResStructName 是否是 Res 结尾
	if strings.HasSuffix(data.ResStructName, "Res") {
		data.IsResStruct = true
	}

	// 查找第一个带 unique_index 的字段作为 KeyField 和 DBFieldName
	for i := 0; i < structVal.NumField(); i++ {
		field := structVal.Field(i)
		gormTag := field.Tag.Get("gorm")
		jsonTag := field.Tag.Get("json")

		// 查找 unique_index 字段并设置 KeyField 和 DBFieldName
		if strings.Contains(gormTag, "unique_index") {
			data.KeyField = field.Name
			if jsonTag != "" {
				data.DBFieldName = jsonTag
			} else {
				// 如果没有 JSON 标签，则使用字段名的小写形式作为 DBFieldName
				data.DBFieldName = strings.ToLower(field.Name)
			}
			break
		}
	}

	// Extract fields and conditions from query struct type
	if queryVal.Kind() == reflect.Struct {
		for i := 0; i < queryVal.NumField(); i++ {
			field := queryVal.Field(i)
			fieldName := field.Name
			fieldType := field.Type.String()
			tag := field.Tag.Get("json")
			if tag == "start_time" || tag == "end_time" {
				data.HasDateRange = true
			} else if tag != "" {
				data.Conditions = append(data.Conditions, QueryCondition{
					FieldName:   fieldName,
					DBFieldName: tag,
					ZeroValue:   getZeroValue(fieldType),
				})
			}
		}
	}

	// Check if the struct type has an "is_deleted" field
	fields := make([]string, structVal.NumField())
	for i := 0; i < structVal.NumField(); i++ {
		fields[i] = structVal.Field(i).Name
	}
	data.HasIsDeletedField = hasIsDeletedField(fields)

	// 自动生成 RelatedItems
	generateRelatedItems(data)
}

// 自动生成 RelatedItems
func generateRelatedItems(data *CombinedTemplateData) {
	fmt.Println("generateRelatedItems")
	resVal := reflect.TypeOf(data.ResStructType)
	if resVal.Kind() == reflect.Struct {
		for i := 0; i < resVal.NumField(); i++ {
			field := resVal.Field(i)
			//if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Ptr {
			if field.Type.Kind() == reflect.Ptr {

				elemType := field.Type.Elem()

				typeName := elemType.Name()

				fieldName := field.Name
				tag := field.Tag.Get("json")
				// 如果tag以info结尾，则去掉info
				if strings.HasSuffix(tag, "_info") {
					tag = tag[:len(tag)-5]
				}
				keyField := fieldName
				// 如果keyField以Info结尾，则去掉Info
				if strings.HasSuffix(keyField, "Info") {
					keyField = keyField[:len(keyField)-4]
				}
				relatedItem := RelatedItem{
					FunctionName: fmt.Sprintf("New%sService().Get%sByUuidList", typeName, typeName),
					Key:          strings.ToLower(tag),
					KeyField:     keyField,
					FieldName:    fieldName,
					TypeName:     typeName,
					ReturnType:   fmt.Sprintf("[]*%s", elemType.Name()),
				}
				data.RelatedItems = append(data.RelatedItems, relatedItem)
			}
		}
	}
	fmt.Println(data.RelatedItems)
}

// Generate combined code based on CombinedTemplateData
func GenerateCombinedCode(data *CombinedTemplateData, templateFile string) string {
	// Extract struct information and conditions
	// extractStructInfo(data)

	// Load and execute the template
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
		fmt.Println("Error formatting code:", err, buf.String())
		return buf.String() // Return unformatted code in case of error
	}

	return string(formattedCode)
}

// GenerateFile generates the code file based on the template and saves it to the specified directory
func GenerateFile(data *CombinedTemplateData, templateFile, outputDir, fileName string) {
	code := GenerateCombinedCode(data, templateFile)
	fullPath := filepath.Join(outputDir, fileName)
	err := os.WriteFile(fullPath, []byte(code), 0644)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", fullPath, err)
	}
}

// AppendToFile appends generated code to a specified file
func AppendToFile(data *CombinedTemplateData, templateFile, filePath string) {
	code := GenerateCombinedCode(data, templateFile)

	// 先以只读模式打开文件，检查内容
	readFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return
	}
	defer readFile.Close()

	// 读取文件内容
	content, err := ioutil.ReadAll(readFile)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return
	}
	key := fmt.Sprintf("type %sInfoResponse", data.StructName)

	// 检查文件内容是否包含关键词
	if strings.Contains(string(content), key) {
		fmt.Printf("Code already exists in file %s\n", filePath)
		return
	}

	key = fmt.Sprintf("Init%sRouter", data.StructName)
	// 检查文件内容是否包含关键词
	if strings.Contains(string(content), key) {
		fmt.Printf("Code already exists in file %s\n", filePath)
		return
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(code); err != nil {
		fmt.Printf("Error writing to file %s: %v\n", filePath, err)
	}
}

// LoadConfig loads the configuration from the specified file
func LoadConfig(configFile string) (*Config, error) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		// 使用默认配置
		return &Config{
			ServiceDir:    "",
			ControllerDir: "",
		}, nil
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func GenCode(data *CombinedTemplateData, conf *Config) {

	// 先提取结构体信息以确保 StructName 正确设置
	extractStructInfo(data)

	//fmt.Println(code)
	// 生成 service 代码
	GenerateFile(data, "service_template.go.tpl", conf.ServiceDir, fmt.Sprintf("gen_%s_service.go", strings.ToLower(data.StructName)))

	// 生成 controller 代码
	GenerateFile(data, "controller_template.go.tpl", conf.ControllerDir, fmt.Sprintf("gen_%s_controller.go", strings.ToLower(data.StructName)))

	AppendToFile(data, "router_template.go.tpl", conf.RouterFile)

	AppendToFile(data, "response_struct_template.go.tpl", conf.ReponseFile)
}

func main() {

	// Load configuration
	config, err := LoadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	data := &CombinedTemplateData{
		ModuleName:          "文章",
		HasGetListFunction:  true,
		StructType:          model.Article{},
		QueryStructType:     model.ReqArticleQueryParam{},
		ReqCreateStructType: model.ReqCreateArticle{},
		ResStructType:       model.ArticleRes{},
	}

	GenCode(data, config)

}
