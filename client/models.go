package client

type Res struct {
	SqlString string
	Names     []string
	Values    []any
}

type Model interface {
	GetDBName() string
	GetRelationList() *RelationList
	ScanResult()
	SetResult(result Result)
}

type Result interface {
	GetDBName() string
	GetSelectedFields() []*SelectedField
	GetRelations() []Result
	IsExist() bool
}
