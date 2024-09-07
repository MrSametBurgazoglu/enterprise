package client

type OperationType string

const (
	OperationTypeGet        OperationType = "GET"
	OperationTypeRefresh    OperationType = "REFRESH"
	OperationTypeCreate     OperationType = "CREATE"
	OperationTypeUpdate     OperationType = "UPDATE"
	OperationTypeDelete     OperationType = "DELETE"
	OperationTypeList       OperationType = "LIST"
	OperationTypeBulkCreate OperationType = "BULK_CREATE"
	OperationTypeBulkUpdate OperationType = "BULK_UPDATE"
	OperationTypeBulkDelete OperationType = "BULK_DELETE"
)

type OperationInfo struct {
	TableName     string
	OperationType OperationType
}

func NewOperationInfo(tableName string, operationType OperationType) *OperationInfo {
	return &OperationInfo{TableName: tableName, OperationType: operationType}
}
