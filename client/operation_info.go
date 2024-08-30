package client

type OperationType string

const (
	OperationTypeGet     OperationType = "GET"
	OperationTypeRefresh OperationType = "REFRESH"
	OperationTypeList    OperationType = "LIST"
	OperationTypeCreate  OperationType = "CREATE"
	OperationTypeUpdate  OperationType = "UPDATE"
	OperationTypeDelete  OperationType = "DELETE"
)

type OperationInfo struct {
	TableName     string
	OperationType OperationType
}

func NewOperationInfo(tableName string, operationType OperationType) *OperationInfo {
	return &OperationInfo{TableName: tableName, OperationType: operationType}
}
