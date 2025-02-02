package assets

// View - Representation of a Database View.
type View struct {
	*AbstractAsset
	sql string
}

func NewView(name string, sql string) *View {
	v := &View{AbstractAsset: NewAbstractAsset(), sql: sql}

	v.SetName(name)

	return v
}

func (v *View) GetSQL() string {
	return v.sql
}
