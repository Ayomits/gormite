package local_schema

import (
	"github.com/KoNekoD/gormite/pkg/assets"
	"github.com/KoNekoD/gormite/pkg/dtos"
	"go/ast"
)

type store struct {
	// Root
	path         string
	config       *dtos.ConfigData
	tables       []*assets.Table
	sequences    []*assets.Sequence
	schemaConfig *dtos.SchemaConfig
	namespaces   []string

	// Mapping key level
	objectsMap           map[string]*ast.Object
	namesMap             map[string]string
	importsMap           map[string][]*ast.ImportSpec
	structNamesIdentsMap map[string]*ast.Ident
}

func newStore(path string) *store {
	return &store{
		path:                 path,
		config:               nil,
		tables:               make([]*assets.Table, 0),
		sequences:            make([]*assets.Sequence, 0),
		schemaConfig:         dtos.NewSchemaConfig(),
		namespaces:           make([]string, 0),
		objectsMap:           make(map[string]*ast.Object),
		namesMap:             make(map[string]string),
		importsMap:           make(map[string][]*ast.ImportSpec),
		structNamesIdentsMap: make(map[string]*ast.Ident),
	}
}
