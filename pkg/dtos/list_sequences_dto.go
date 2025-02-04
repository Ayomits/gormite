package dtos

import "strconv"

type ListSequencesDto struct {
	Relname     string `db:"relname"`
	Schemaname  string `db:"schemaname"`
	MinValue    string `db:"min_value"`
	IncrementBy string `db:"increment_by"`
}

func (dto *ListSequencesDto) GetMinValue() int {
	n, _ := strconv.Atoi(dto.MinValue)
	return n
}
func (dto *ListSequencesDto) GetIncrementBy() int {
	n, _ := strconv.Atoi(dto.IncrementBy)
	return n
}
