package enums

type QueryType string

const (
	QueryTypeSelect QueryType = "SELECT"
	QueryTypeDelete QueryType = "DELETE"
	QueryTypeUpdate QueryType = "UPDATE"
	QueryTypeInsert QueryType = "INSERT"
	QueryTypeUnion  QueryType = "UNION"
)
