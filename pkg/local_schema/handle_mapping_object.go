package local_schema

import (
	"fmt"
	"github.com/KoNekoD/gormite/pkg/assets"
	"github.com/KoNekoD/gormite/pkg/types"
	"github.com/fatih/structtag"
	"go/ast"
	"golang.org/x/exp/maps"
	"slices"
	"strings"
)

func (s *store) newTable(name string) *assets.Table {
	name = getName(s, name)

	for _, table := range s.tables {
		if table.GetName() == name {
			return table
		}
	}

	t := assets.NewTable(name, nil, nil, nil, nil, nil)

	s.tables = append(s.tables, t)

	return t
}

type tableBag struct {
	store       *store
	table       *assets.Table
	primaryKeys []string

	uniqColumnsMap    map[string][]string
	uniqConditionsMap map[string]string

	indexColumnsMap    map[string][]string
	indexConditionsMap map[string]string
}

func newTableBag(store *store, table *assets.Table) *tableBag {
	bag := &tableBag{
		store:              store,
		table:              table,
		primaryKeys:        make([]string, 0),
		uniqColumnsMap:     make(map[string][]string),
		uniqConditionsMap:  make(map[string]string),
		indexColumnsMap:    make(map[string][]string),
		indexConditionsMap: make(map[string]string),
	}

	return bag
}

func (t *tableBag) colIdent(fieldType *ast.Ident, tags *structtag.Tags) {
	objectsKeys := maps.Keys(t.store.objectsMap)

	columnTagsData := t.parseColumnTags(tags, fieldType, objectsKeys)

	if columnTagsData.ColumnType != nil {
		t.table.AddColumn(
			columnTagsData.ColumnName,
			columnTagsData.ColumnType,
			columnTagsData.Options...,
		)
	} else {
		if !columnTagsData.IsForeignKey {
			panic(fmt.Sprintf("unknown type %s", columnTagsData.TypeName))
		}

		// Maybe we need rewrite it to allow use non integer ids...
		t.table.AddColumn(
			columnTagsData.ColumnName,
			types.NewIntegerType(),
			columnTagsData.Options...,
		)
	}

	applyMetadataMutatorsForNewColumn(columnTagsData, t)
}

func (t *tableBag) colSel(
	fType *ast.SelectorExpr,
	tags *structtag.Tags,
	mustBeNullable bool,
) {
	objectsKeys := maps.Keys(t.store.objectsMap)

	selPackage := fType.X.(*ast.Ident).Name
	selType := fType.Sel.Name

	columnTagsData := t.parseColumnTags(tags, fType.Sel, objectsKeys)

	if selPackage == "time" && selType == "Time" {
		t.table.AddColumn(
			columnTagsData.ColumnName,
			types.NewDateTimeImmutableType(),
			columnTagsData.Options...,
		)
		applyMetadataMutatorsForNewColumn(columnTagsData, t)
		return
	}

	if mustBeNullable && columnTagsData.IsNotNull {
		panic("column " + columnTagsData.ColumnName + " of table " + t.table.GetName() + " cannot be not null")
	}

	if columnTagsData.ColumnType != nil {
		t.table.AddColumn(
			columnTagsData.ColumnName,
			columnTagsData.ColumnType,
			columnTagsData.Options...,
		)
		applyMetadataMutatorsForNewColumn(columnTagsData, t)
	} else {
		if found, ok := t.store.structNamesIdentsMap[selType]; ok {
			t.colIdent(found, tags)
		} else {
			panic("unknown type for " + columnTagsData.ColumnName)
		}
	}
}

func (t *tableBag) colStar(fType *ast.StarExpr, tags *structtag.Tags) {
	switch fieldTypeX := fType.X.(type) {
	case *ast.Ident:
		t.colIdent(fieldTypeX, tags)
	case *ast.SelectorExpr:
		t.colSel(
			fieldTypeX,
			tags,
			true,
		) // TODO: must be nullable fool protection not work, idk why
	default:
		panic(fmt.Sprintf("Unknown star %T", fieldTypeX))
	}
}

func (t *tableBag) colArray(fieldType *ast.ArrayType, tags *structtag.Tags) {
	objectsKeys := maps.Keys(t.store.objectsMap)

	ident, ok := fieldType.Elt.(*ast.Ident)
	if !ok {
		panic("Only primitive array types are supported")
	}

	if ident.Name != "string" {
		panic("Array type is not supported for type: " + ident.Name)
	}

	columnTagsData := t.parseColumnTags(tags, ident, objectsKeys)

	if "string" == columnTagsData.TypeName {
		panic("please add type tag to string array for property: " + columnTagsData.ColumnName)
	}

	if !slices.Contains([]string{"json", "jsonb"}, columnTagsData.TypeName) {
		panic("Only json/jsonb array types are supported")
	}

	t.table.AddColumn(
		columnTagsData.ColumnName,
		columnTagsData.ColumnType,
		columnTagsData.Options...,
	)
}

// OneToOne - uniq_cd4f5a305067c3d4
// OneToMany - virtual, not owner
// ManyToOne - idx_79bd4a955067c3d4
// ManyToMany - table1_table2_pkey, idx_f9cb7c79afc2b591, idx_f9cb7c79810212b - create separate table

func handleMappingObject(objectName string, store *store) (err error) {
	t := store.newTable(objectName)

	object := store.objectsMap[objectName]

	typeSpec := object.Decl.(*ast.TypeSpec)
	structType := typeSpec.Type.(*ast.StructType)

	bag := newTableBag(store, t)

	for _, field := range structType.Fields.List {
		tag := field.Tag.Value
		tag = strings.Trim(field.Tag.Value, "`")

		tags, err := structtag.Parse(tag)
		if err != nil {
			panic(err)
		}

		switch fType := field.Type.(type) {
		case *ast.Ident:
			bag.colIdent(fType, tags)
		case *ast.StarExpr:
			bag.colStar(fType, tags)
		case *ast.SelectorExpr:
			bag.colSel(fType, tags, false)
		case *ast.ArrayType:
			bag.colArray(fType, tags)
		default:
			panic(fmt.Sprintf("Unknown fieldType %T", fType))
		}
	}

	applyMetadataMutatorsAfterColumnsIntrospection(bag)

	return nil
}
