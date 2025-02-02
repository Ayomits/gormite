package types

import (
	"github.com/KoNekoD/gormite/pkg/gormite/enums"
	"github.com/KoNekoD/gormite/pkg/gormite/utils"
)

type TypeRegistry struct {
	instances             map[enums.TypesType]AbstractTypeInterface
	instancesReverseIndex map[string]enums.TypesType
}

func (t *TypeRegistry) NewTypeRegistry(instances map[enums.TypesType]AbstractTypeInterface) *TypeRegistry {
	t.instances = make(map[enums.TypesType]AbstractTypeInterface)
	t.instancesReverseIndex = make(map[string]enums.TypesType)
	for name, typeVar := range instances {
		t.Register(name, typeVar)
	}

	return t
}
func (t *TypeRegistry) Get(name enums.TypesType) AbstractTypeInterface {
	typeVar, ok := t.instances[name]
	if !ok {
		panic("Unknown column type " + name)
	}
	return typeVar
}
func (t *TypeRegistry) LookupName(typeVar AbstractTypeInterface) enums.TypesType {
	name := t.findTypeName(typeVar)

	if name == nil {
		panic("TypeNotRegistered")
	}

	return *name
}
func (t *TypeRegistry) Has(name enums.TypesType) bool {
	_, ok := t.instances[name]

	return ok
}
func (t *TypeRegistry) Register(
	name enums.TypesType,
	typeVar AbstractTypeInterface,
) {
	_, ok := t.instances[name]
	if ok {
		panic("TypeAlreadyRegistered " + name)
	}

	if t.findTypeName(typeVar) != nil {
		panic("TypeAlreadyRegistered " + name)
	}

	t.instances[name] = typeVar
	t.instancesReverseIndex[utils.SplObjectID(typeVar)] = name
}
func (t *TypeRegistry) Override(
	name enums.TypesType,
	typeVar AbstractTypeInterface,
) {
	origType, ok := t.instances[name]
	if !ok {
		panic("TypeNotFound " + name)
	}

	typeName := t.findTypeName(typeVar)
	if typeName != nil && *typeName != name {
		panic("TypeAlreadyRegistered " + name)
	}

	delete(t.instancesReverseIndex, utils.SplObjectID(origType))
	t.instances[name] = typeVar
	t.instancesReverseIndex[utils.SplObjectID(typeVar)] = name
}
func (t *TypeRegistry) GetMap() map[enums.TypesType]AbstractTypeInterface {
	return t.instances
}
func (t *TypeRegistry) findTypeName(typeVar AbstractTypeInterface) *enums.TypesType {
	v, ok := t.instancesReverseIndex[utils.SplObjectID(typeVar)]

	if !ok {
		return nil
	}

	return &v
}
