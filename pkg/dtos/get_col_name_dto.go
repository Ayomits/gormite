package dtos

type GetColNameDto struct {
	Attnum  string `db:"attnum"`
	Attname string `db:"attname"`
}
