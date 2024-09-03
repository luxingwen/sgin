package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type {{.StructName}}Service struct {
}

func New{{.StructName}}Service() *{{.StructName}}Service {
	return &{{.StructName}}Service{}
}

func (s *{{.StructName}}Service) Create{{.StructName}}(ctx *app.Context, {{.StructName | lower}} *model.{{.StructName}}) error {
	{{.StructName | lower}}.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	{{.StructName | lower}}.UpdatedAt = {{.StructName | lower}}.CreatedAt
	{{.StructName | lower}}.Uuid = uuid.New().String()

	err := ctx.DB.Create({{.StructName | lower}}).Error
	if err != nil {
		ctx.Logger.Error("Failed to create {{.StructName | lower}}", err)
		return errors.New("failed to create {{.StructName | lower}}")
	}
	return nil
}

func (s *{{.StructName}}Service) Get{{.StructName}}ByUUID(ctx *app.Context, uuid string) (*model.{{.StructName}}, error) {
	{{.StructName | lower}} := &model.{{.StructName}}{}
	err := ctx.DB.Where("uuid = ?", uuid).First({{.StructName | lower}}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("{{.StructName | lower}} not found")
		}
		ctx.Logger.Error("Failed to get {{.StructName | lower}} by UUID", err)
		return nil, errors.New("failed to get {{.StructName | lower}} by UUID")
	}
	return {{.StructName | lower}}, nil
}

func (s *{{.StructName}}Service) Update{{.StructName}}(ctx *app.Context, {{.StructName | lower}} *model.{{.StructName}}) error {
	{{.StructName | lower}}.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	err := ctx.DB.Where("uuid = ?", {{.StructName | lower}}.Uuid).Updates({{.StructName | lower}}).Error
	if err != nil {
		ctx.Logger.Error("Failed to update {{.StructName | lower}}:", err)
		return errors.New("failed to update {{.StructName | lower}}")
	}

	return nil
}

{{if .HasIsDeletedField}}
func (s *{{.StructName}}Service) Delete{{.StructName}}(ctx *app.Context, uuid string) error {
	err := ctx.DB.Model(&model.{{.StructName}}{}).Where("uuid = ?", uuid).Update("is_deleted", 1).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete {{.StructName | lower}}", err)
		return errors.New("failed to delete {{.StructName | lower}}")
	}

	return nil
}
{{else}}
func (s *{{.StructName}}Service) Delete{{.StructName}}(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.{{.StructName}}{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete {{.StructName | lower}}", err)
		return errors.New("failed to delete {{.StructName | lower}}")
	}

	return nil
}
{{end}}

{{if .HasGetListFunction}}
// 获取{{.StructName}}列表
func (s *{{.StructName}}Service) Get{{.StructName}}List(ctx *app.Context, params *model.{{.QueryStructName}}) (r *model.PagedResponse, err error) {
	var (
		{{.LowerStructName}}s []*model.{{.StructName}}
		total                 int64
	)

	db := ctx.DB.Model(&model.{{.StructName}}{})

	// 动态生成查询条件
	{{range .Conditions}}
	if params.{{.FieldName}} != {{.ZeroValue}} {
		db = db.Where("{{.DBFieldName}} = ?", params.{{.FieldName}})
	}
	{{end}}

	// 获取记录总数
	err = db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get {{.LowerStructName}} count", err)
		return nil, errors.New("failed to get {{.LowerStructName}} count")
	}

	// 获取数据列表
	err = db.Order("id DESC").Offset(params.GetOffset()).Limit(params.PageSize).Find(&{{.LowerStructName}}s).Error
	if err != nil {
		ctx.Logger.Error("Failed to get {{.LowerStructName}} list", err)
		return nil, errors.New("failed to get {{.LowerStructName}} list")
	}

	// 提取关键字段列表
	keyFieldList := make([]string, 0)
	for _, item := range {{.LowerStructName}}s {
		keyFieldList = append(keyFieldList, item.{{.KeyField}})
	}

	// 动态获取多个关联项
	relatedData := make(map[string]interface{})
	{{range .RelatedItems}}
	related{{.FunctionName}}, err := s.{{.FunctionName}}(ctx, keyFieldList)
	if err != nil {
		ctx.Logger.Error("Failed to get related {{.FunctionName | lower}}", err)
		return nil, errors.New("failed to get related {{.FunctionName | lower}}")
	}
	relatedData["{{.Key}}"] = related{{.FunctionName}}
	{{end}}

	// 构造最终响应列表
	res := make([]*model.{{.StructName}}Res, 0)
	for _, item := range {{.LowerStructName}}s {
		itemRes := &model.{{.StructName}}Res{
			{{.StructName}}: *item,
		}
		{{range .RelatedItems}}
		if relatedItems, ok := relatedData["{{.Key}}"].(map[string]{{.ReturnType}})[item.{{.KeyField}}]; ok {
			itemRes.{{.FieldName}} = relatedItems
		}
		{{end}}
		res = append(res, itemRes)
	}

	return &model.PagedResponse{
		Total: total,
		Data:  res,
		Current: params.Current,
		PageSize: params.PageSize,
	}, nil
}
{{end}}



{{/* 
   生成获取单个对象映射的方法，根据 UUID 列表获取资源映射
*/}}
{{if .HasGetMapFunction}}
// 根据 UUID 列表获取 {{.StructName}} 映射
func (s *{{.StructName}}Service) Get{{.StructName}}ByUUIDList(ctx *app.Context, uuidList []string) (map[string]*model.{{.StructName}}, error) {
	{{.LowerStructName}}Map := make(map[string]*model.{{.StructName}}, 0)
	var {{.LowerStructName}}s []*model.{{.StructName}}
	err := ctx.DB.Where("uuid IN (?)", uuidList).Find(&{{.LowerStructName}}s).Error
	if err != nil {
		ctx.Logger.Error("Failed to get {{.LowerStructName}} list by UUID list", err)
		return nil, errors.New("failed to get {{.LowerStructName}} list by UUID list")
	}

	for _, {{.LowerStructName}} := range {{.LowerStructName}}s {
		{{.LowerStructName}}Map[{{.LowerStructName}}.Uuid] = {{.LowerStructName}}
	}

	return {{.LowerStructName}}Map, nil
}
{{end}}

{{/* 
   生成获取对象列表映射的方法，根据 Variant UUID 列表获取产品变体选项映射
*/}}
{{if .HasGetListMapFunction}}
// 根据 {{.KeyField}} 列表获取 {{.StructName}} 映射
func (s *{{.StructName}}Service) Get{{.StructName}}By{{.KeyField}}List(ctx *app.Context, {{.KeyField | lower}}List []string) (map[string][]*model.{{.StructName}}, error) {
	{{.LowerStructName}}Map := make(map[string][]*model.{{.StructName}})
	if len({{.KeyField | lower}}List) == 0 {
		return {{.LowerStructName}}Map, nil
	}

	var {{.LowerStructName}}s []*model.{{.StructName}}
	err := ctx.DB.Where("{{.DBFieldName}} IN (?)", {{.KeyField | lower}}List).Find(&{{.LowerStructName}}s).Error
	if err != nil {
		ctx.Logger.Error("Failed to get {{.LowerStructName}} list by {{.KeyField | lower}} list", err)
		return nil, errors.New("failed to get {{.LowerStructName}} list by {{.KeyField | lower}} list")
	}

	for _, {{.LowerStructName}} := range {{.LowerStructName}}s {
		if _, ok := {{.LowerStructName}}Map[{{.LowerStructName}}.{{.KeyField}}]; !ok {
			{{.LowerStructName}}Map[{{.LowerStructName}}.{{.KeyField}}] = make([]*model.{{.StructName}}, 0)
		}
		{{.LowerStructName}}Map[{{.LowerStructName}}.{{.KeyField}}] = append({{.LowerStructName}}Map[{{.LowerStructName}}.{{.KeyField}}], {{.LowerStructName}})
	}

	return {{.LowerStructName}}Map, nil
}
{{end}}
