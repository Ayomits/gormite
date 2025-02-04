package dtos

type FetchTableOptionsByTableDto struct {
	Relname  string  `db:"relname"`
	Unlogged bool    `db:"unlogged"`
	Comment  *string `db:"comment"`
}

func (f *FetchTableOptionsByTableDto) ToArray() map[string]interface{} {
	return map[string]interface{}{
		"relname":  f.Relname,
		"unlogged": f.Unlogged,
		"comment":  f.Comment,
	}
}
