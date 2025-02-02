package local_schema

import (
	"fmt"
	"github.com/KoNekoD/gormite/pkg/gormite/assets"
	"github.com/KoNekoD/gormite/pkg/gormite/types"
	"github.com/fatih/structtag"
	"go/ast"
	"slices"
	"strconv"
)

func (t *tableBag) parseColumnTags(
	tags *structtag.Tags,
	fieldType *ast.Ident,
	objectsKeys []string,
) *columnData {
	typeName := fieldType.Name

	colNameTag, _ := tags.Get("db")
	columnName := colNameTag.Value()

	isForeignKey := slices.Contains(objectsKeys, typeName)

	pk, _ := tags.Get("pk")
	isPrimaryKey := pk != nil

	nullableTag, _ := tags.Get("nullable")
	isNullable := nullableTag != nil
	isNotNull := !isNullable

	lengthTag, _ := tags.Get("length")
	length := 255
	if lengthTag != nil {
		length, _ = strconv.Atoi(lengthTag.Value())
	}

	uniqTag, _ := tags.Get("uniq")
	isUnique := uniqTag != nil
	var uniqueName *string
	if isUnique {
		uniqTagName := uniqTag.Value()
		uniqueName = &uniqTagName
	}
	uniqCondTag, _ := tags.Get("uniq_cond")
	isUniqueCondition := uniqCondTag != nil
	var uniqueCondition *string
	if isUniqueCondition {
		uniqCondTagName := uniqCondTag.Value()
		uniqueCondition = &uniqCondTagName
	}

	indexTag, _ := tags.Get("index")
	isIndex := indexTag != nil
	var indexName *string
	if isIndex && indexTag != nil {
		indexTagName := indexTag.Value()
		indexName = &indexTagName
	}
	indexCondTag, _ := tags.Get("index_cond")
	isIndexCondition := indexCondTag != nil
	var indexCondition *string
	if isIndexCondition && indexCondTag != nil {
		indexCondTagName := indexCondTag.Value()
		indexCondition = &indexCondTagName
	}

	defaultTag, _ := tags.Get("default")
	var defaultValue *string
	if defaultTag != nil {
		defaultValueTmp := defaultTag.Value()
		defaultValue = &defaultValueTmp
	}

	var columnType types.AbstractTypeInterface
	options := make([]assets.ColumnOption, 0)

	typeTag, _ := tags.Get("type")
	if typeTag != nil {
		typeTagValue := typeTag.Value()
		switch typeTagValue {
		case "text":
			columnType = types.NewTextType()
		case "varchar":
			columnType = types.NewStringType()
		case "json":
			columnType = types.NewJsonType()
			typeName = "json"
		case "jsonb":
			typeName = "jsonb"
			columnType = types.NewJsonType()
			options = append(
				options,
				func(c *assets.Column) { c.SetPlatformOption("jsonb", true) },
			)
		case "integer":
			columnType = types.NewIntegerType()
		default:
			panic(
				fmt.Sprintf(
					"unknown tag type %s for type %s on table %s",
					typeTagValue,
					typeName,
					t.table.GetName(),
				),
			)
		}
	}
	if columnType == nil {
		switch typeName {
		case "int":
			columnType = types.NewIntegerType()
		case "string":
			columnType = types.NewStringType()
		case "bool":
			columnType = types.NewBooleanType()
		}
	}

	if isNotNull {
		options = append(options, assets.WithColumnNotNull())
	}
	if defaultValue != nil {
		options = append(options, assets.WithColumnDefault(*defaultValue))
	}
	isColumnTypeVarchar := false
	if columnType != nil {
		_, isColumnTypeVarchar = columnType.(*types.StringType)
	}
	if typeName == "string" || isColumnTypeVarchar {
		options = append(options, assets.WithColumnLength(&length))
	}

	return &columnData{
		ColumnName:        columnName,
		IsPrimaryKey:      isPrimaryKey,
		IsForeignKey:      isForeignKey,
		IsNotNull:         isNotNull,
		TypeName:          typeName,
		IsUnique:          isUnique,
		UniqueName:        uniqueName,
		IsUniqueCondition: isUniqueCondition,
		UniqueCondition:   uniqueCondition,
		IsIndex:           isIndex,
		IndexName:         indexName,
		IsIndexCondition:  isIndexCondition,
		IndexCondition:    indexCondition,
		Length:            length,
		DefaultValue:      defaultValue,
		ColumnType:        columnType,
		Options:           options,
	}
}
