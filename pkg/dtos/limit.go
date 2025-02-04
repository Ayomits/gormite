package dtos

type Limit struct {
	MaxResults  *int
	FirstResult int
}

func NewLimit(maxResults *int, firstResult int) *Limit {
	return &Limit{MaxResults: maxResults, FirstResult: firstResult}
}

func (l *Limit) IsDefined() bool {
	return l.MaxResults != nil || l.FirstResult != 0
}

func (l *Limit) GetMaxResults() *int {
	return l.MaxResults
}

func (l *Limit) GetFirstResult() int {
	return l.FirstResult
}
