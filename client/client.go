package client

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DatabaseClient interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type DatabaseTransactionClient interface {
	DatabaseClient
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

var NotFoundError = errors.New("not found")

type Client struct {
	Database DatabaseClient
}

func NewClient(d DatabaseClient) *Client {
	return &Client{Database: d}
}

func createTableNameAndAddresses(result Result, selectedNames []string, selectedAddress []any) ([]string, []any) {
	for _, field := range result.GetSelectedFields() {
		selectedNames = append(selectedNames, fmt.Sprintf("\"%s\".\"%s\"", result.GetDBName(), field.Name))
		selectedAddress = append(selectedAddress, field.Value)
	}
	return selectedNames, selectedAddress
}

func createTableRelationNameAndAddresses(result Result, selectedNames []string, selectedAddress []any) ([]string, []any) {
	for _, field := range result.GetSelectedFields() {
		selectedNames = append(selectedNames, fmt.Sprintf("\"%s\".\"%s\"", result.GetDBName(), field.Name))
		selectedAddress = append(selectedAddress, field.Value)
	}
	return selectedNames, selectedAddress
}

func createTableRelationsNameAndAddresses(result Result, selectedNames []string, selectedAddress []any) ([]string, []any) {
	for _, relation := range result.GetRelations() {
		selectedNames, selectedAddress = createTableRelationNameAndAddresses(relation, selectedNames, selectedAddress)
		selectedNames, selectedAddress = createTableRelationsNameAndAddresses(relation, selectedNames, selectedAddress)
	}
	return selectedNames, selectedAddress
}

func createTableWhereSql(list []*WhereList, args pgx.NamedArgs, dbName string) []string {
	var whereStrings []string
	for _, item := range list {
		res := item.Parse(dbName, args)
		if res.SqlString != "" {
			whereStrings = append(whereStrings, res.SqlString)
		}
	}
	return whereStrings
}

func createTableRelationWhereSql(model Model, args pgx.NamedArgs) (string, []string) {
	sqlString := ""
	var relationWhereStrings []string
	for _, rel := range model.GetRelationList().Relations {
		sql := rel.getJoinString(model.GetDBName())
		sqlString += sql
		if rel.isRelationHaveWhereClause() {
			res := rel.parseWhere(args)
			relationWhereStrings = append(relationWhereStrings, res.SqlString)
		}
		a, b := createTableRelationWhereSql(rel.RelationModel, args)
		sqlString += " " + a
		relationWhereStrings = append(relationWhereStrings, b...)
	}
	return sqlString, relationWhereStrings
}

func isHasManyRelation(rel *Relation) bool {
	if rel.ManyToManyTable != "" {
		return true
	}
	if rel.RelationModel == nil {
		return false
	}
	_, ok := rel.RelationModel.(ListModel)
	return ok
}

func CreateSelectQuery(list []*WhereList, model Model, result Result) (string, []any, pgx.NamedArgs) {
	var selectedNames []string
	var selectedAddress []any
	selectedNames, selectedAddress = createTableNameAndAddresses(result, selectedNames, selectedAddress)
	selectedNames, selectedAddress = createTableRelationsNameAndAddresses(result, selectedNames, selectedAddress)

	args := pgx.NamedArgs{}
	whereStrings := createTableWhereSql(list, args, result.GetDBName())
	relationSqlString, relationWhereStrings := createTableRelationWhereSql(model, args)
	mainTableWhereString := strings.Join(whereStrings, " OR ")
	var allWhereStrings []string
	if mainTableWhereString != "" {
		allWhereStrings = append(allWhereStrings, mainTableWhereString)
	}
	if len(relationWhereStrings) > 0 {
		allWhereStrings = append(allWhereStrings, relationWhereStrings...)
	}

	summedWhereString := ""
	if len(allWhereStrings) > 0 {
		summedWhereString = fmt.Sprintf(" WHERE (%s)", strings.Join(allWhereStrings, " AND "))
	}

	hasHasMany := false
	if model.GetRelationList() != nil {
		for _, rel := range model.GetRelationList().Relations {
			if isHasManyRelation(rel) {
				hasHasMany = true
				break
			}
		}
	}

	limitStr := " LIMIT 1"
	if hasHasMany {
		limitStr = ""
	}

	names := strings.Join(selectedNames, ", ")
	sqlString := fmt.Sprintf("SELECT %s FROM \"%s\" %s%s%s;",
		names,
		model.GetDBName(),
		relationSqlString,
		summedWhereString,
		limitStr)

	return sqlString, selectedAddress, args
}

func CreateSelectListQuery(list []*WhereList, model Model, result Result, orders []*Order, paging *Paging) (string, []any, pgx.NamedArgs) {
	var selectedNames []string
	var selectedAddress []any
	selectedNames, selectedAddress = createTableNameAndAddresses(result, selectedNames, selectedAddress)
	selectedNames, selectedAddress = createTableRelationsNameAndAddresses(result, selectedNames, selectedAddress)

	args := pgx.NamedArgs{}
	whereStrings := createTableWhereSql(list, args, result.GetDBName())
	relationSqlString, relationWhereStrings := createTableRelationWhereSql(model, args)
	mainTableWhereString := strings.Join(whereStrings, " OR ")
	var allWhereStrings []string
	if mainTableWhereString != "" {
		allWhereStrings = append(allWhereStrings, mainTableWhereString)
	}
	if len(relationWhereStrings) > 0 {
		allWhereStrings = append(allWhereStrings, relationWhereStrings...)
	}

	summedWhereString := ""
	if len(allWhereStrings) > 0 {
		summedWhereString = fmt.Sprintf(" WHERE (%s)", strings.Join(allWhereStrings, " AND "))
	}

	var orderStrings []string
	for _, order := range orders {
		orderStrings = append(orderStrings, order.String())
	}
	orderString := ""
	if len(orderStrings) > 0 {
		orderString = "ORDER BY " + strings.Join(orderStrings, ", ")
	}

	pagingString := ""
	if paging != nil {
		pagingString = paging.String()
	}

	names := strings.Join(selectedNames, ", ")
	sqlString := fmt.Sprintf("SELECT %s FROM \"%s\" %s%s %s %s;",
		names,
		model.GetDBName(),
		relationSqlString,
		summedWhereString,
		orderString,
		pagingString)

	return sqlString, selectedAddress, args
}

func ScanFirstRow(rows pgx.Rows, model Model, selectedAddress []any) error {
	if err := rows.Err(); err != nil {
		return err
	}
	if rows.Next() {
		err := rows.Scan(selectedAddress...)
		if err != nil {
			return err
		}
		model.ScanResult()
		return nil
	} else {
		return NotFoundError
	}
}

func ScanNextRows(rows pgx.Rows, model Model, selectedAddress []any) error {
	for rows.Next() {
		err := rows.Scan(selectedAddress...)
		if err != nil {
			return err
		}
		model.ScanResult()
	}
	return nil
}

func (receiver *Client) Get(ctx context.Context, list []*WhereList, model Model, result Result) error {
	sqlString, selectedAddress, args := CreateSelectQuery(list, model, result)

	rows, err := receiver.Database.Query(ctx, sqlString, args)
	if err != nil {
		return err
	}
	defer rows.Close()

	err = ScanFirstRow(rows, model, selectedAddress)
	if err != nil {
		return err
	}

	hasHasMany := false
	if model.GetRelationList() != nil {
		for _, rel := range model.GetRelationList().Relations {
			if isHasManyRelation(rel) {
				hasHasMany = true
				break
			}
		}
	}

	if hasHasMany {
		err = ScanNextRows(rows, model, selectedAddress)
		if err != nil {
			return err
		}
	}

	return nil
}

func (receiver *Client) Refresh(ctx context.Context, model Model, result Result, idName string, idValue any) error {
	var selectedNames []string
	var selectedAddresses []any

	selectedNames, selectedAddresses = createTableNameAndAddresses(result, selectedNames, selectedAddresses)
	names := strings.Join(selectedNames, ", ")
	args := pgx.NamedArgs{
		"idvalue": idValue,
	}

	sqlString := fmt.Sprintf(
		"SELECT %s FROM \"%s\" WHERE \"%s\".\"%s\" = @idvalue",
		names,
		model.GetDBName(),
		model.GetDBName(),
		idName,
	)

	rows, err := receiver.Database.Query(ctx, sqlString, args)
	if err != nil {
		return err
	}
	defer rows.Close()

	err = ScanFirstRow(rows, model, selectedAddresses)
	if err != nil {
		return err
	}

	return nil
}

func CreateInsertQuery(fields map[string]any, fieldsList []string) (string, string, pgx.NamedArgs) {
	args := pgx.NamedArgs{}
	var names []string
	var values []string
	for _, n := range fieldsList {
		v := fields[n]
		names = append(names, fmt.Sprintf("\"%s\"", n))
		values = append(values, "@"+n)
		args[n] = v
	}
	nameString := fmt.Sprintf("(%s)", strings.Join(names, ","))
	valueString := fmt.Sprintf("(%s)", strings.Join(values, ","))
	return nameString, valueString, args
}

func (receiver *Client) Create(ctx context.Context, tableName string, fields map[string]any, fieldsList []string, serialFields []*SelectedField) error {
	nameString, valueString, args := CreateInsertQuery(fields, fieldsList)

	var serialSql string
	var serialFieldAddresses []any
	if len(serialFields) > 0 {
		serialSql = "RETURNING %s"
		var serialNames []string
		for _, field := range serialFields {
			serialNames = append(serialNames, fmt.Sprintf("\"%s\"", field.Name))
			serialFieldAddresses = append(serialFieldAddresses, field.Value)
		}
		serialSql = fmt.Sprintf(serialSql, strings.Join(serialNames, ", "))
	}

	if serialSql == "" {
		sqlString := fmt.Sprintf("INSERT INTO \"%s\" %s VALUES %s;", tableName, nameString, valueString)

		_, err := receiver.Database.Exec(ctx, sqlString, args)
		if err != nil {
			return err
		}
	} else {
		sqlString := fmt.Sprintf("INSERT INTO \"%s\" %s VALUES %s %s;", tableName, nameString, valueString, serialSql)
		row := receiver.Database.QueryRow(ctx, sqlString, args)
		err := row.Scan(serialFieldAddresses...)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateUpdateQuery(fields map[string]any, fieldsList []string) (string, pgx.NamedArgs) {
	args := pgx.NamedArgs{}
	var statements []string
	for _, n := range fieldsList {
		v := fields[n]
		statements = append(statements, fmt.Sprintf("\"%s\" = @%s", n, n))
		args[n] = v
	}
	return strings.Join(statements, ", "), args
}

func (receiver *Client) Update(ctx context.Context, tableName string, fields map[string]any, fieldslist []string, idName string, idValue any) error {
	statements, args := CreateUpdateQuery(fields, fieldslist)

	sqlString := fmt.Sprintf("UPDATE \"%s\" SET %s WHERE \"%s\" = @idvalue", tableName, statements, idName)

	args["idvalue"] = idValue
	_, err := receiver.Database.Exec(ctx, sqlString, args)
	if err != nil {
		return err
	}
	return nil
}

func (receiver *Client) Delete(ctx context.Context, tableName string, idName string, idValue any) error {
	sqlString := fmt.Sprintf("DELETE FROM \"%s\" WHERE \"%s\" = @idvalue;", tableName, idName)

	args := pgx.NamedArgs{}
	args["idvalue"] = idValue
	_, err := receiver.Database.Exec(ctx, sqlString, args)
	if err != nil {
		return err
	}
	return nil
}

func ScanListFirstRow(rows pgx.Rows, model Model, selectedAddress []any) error {
	if err := rows.Err(); err != nil {
		return err
	}
	if rows.Next() {
		err := rows.Scan(selectedAddress...)
		if err != nil {
			return err
		}
		model.ScanResult()
		return nil
	} else {
		return NotFoundError
	}
}

func ScanListNextRows(rows pgx.Rows, model Model, selectedAddress []any) error {
	if err := rows.Err(); err != nil {
		return err
	}
	a := 0
	for rows.Next() {
		err := rows.Scan(selectedAddress...)
		a++
		if err != nil {
			return err
		}
		model.ScanResult()
	}
	return nil
}

func (receiver *Client) List(ctx context.Context, list []*WhereList, model Model, result Result, orders []*Order, paging *Paging) error {
	sqlString, selectedAddress, args := CreateSelectListQuery(list, model, result, orders, paging)

	rows, err := receiver.Database.Query(ctx, sqlString, args)
	if err != nil {
		return err
	}
	defer rows.Close()

	err = ScanListNextRows(rows, model, selectedAddress)
	if err != nil {
		return err
	}

	return nil
}

func ScanValues(rows pgx.Rows, selectedAddress []any) error {
	defer rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	if rows.Next() {
		err := rows.Scan(selectedAddress...)
		return err
	} else {
		return ErrFinalRow
	}
}

func CreateAggregateQuery(list []*WhereList, model Model, aggregate *Aggregate) (string, pgx.NamedArgs) {
	args := pgx.NamedArgs{}
	whereStrings := createTableWhereSql(list, args, model.GetDBName())
	relationSqlString, relationWhereStrings := createTableRelationWhereSql(model, args)
	mainTableWhereString := strings.Join(whereStrings, " OR ")

	var allWhereStrings []string
	if mainTableWhereString != "" {
		allWhereStrings = append(allWhereStrings, mainTableWhereString)
	}
	if len(relationWhereStrings) > 0 {
		allWhereStrings = append(allWhereStrings, relationWhereStrings...)
	}

	summedWhereString := ""
	if len(allWhereStrings) > 0 {
		summedWhereString = fmt.Sprintf(" WHERE (%s)", strings.Join(allWhereStrings, " AND "))
	}

	selectedNames := make([]string, len(aggregate.aggregateFields))
	for i := 0; i < len(aggregate.aggregateFields); i++ {
		field := ValidateIdentifier(aggregate.aggregateFields[i])
		if !strings.Contains(field, "\"") && field != "*" && !strings.ContainsAny(field, "() ,") {
			field = fmt.Sprintf("\"%s\"", field)
		}
		selectedNames[i] = fmt.Sprintf(aggregate.aggregateFormats[i], field)
	}

	groupBys := make([]string, len(aggregate.groupByList))
	for i := 0; i < len(aggregate.groupByList); i++ {
		gb := ValidateIdentifier(aggregate.groupByList[i])
		if !strings.Contains(gb, "\"") && !strings.ContainsAny(gb, "() ,") {
			gb = fmt.Sprintf("\"%s\"", gb)
		}
		groupBys[i] = fmt.Sprintf("GROUP BY %s", gb)
	}

	names := strings.Join(selectedNames, ", ")
	groupBy := strings.Join(groupBys, ", ")
	sqlString := fmt.Sprintf("SELECT %s FROM \"%s\" %s%s %s;",
		names,
		model.GetDBName(),
		relationSqlString,
		summedWhereString,
		groupBy)
	return sqlString, args
}

func (receiver *Client) Aggregate(ctx context.Context, list []*WhereList, model Model, aggregate *Aggregate) (func() error, error) {
	sqlString, args := CreateAggregateQuery(list, model, aggregate)

	rows, err := receiver.Database.Query(ctx, sqlString, args)
	if err != nil {
		return nil, err
	}

	return func() error {
		return ScanValues(rows, aggregate.aggregateValues)
	}, nil
}

func CreateAddRelationQuery(relationshipTable, id, relationshipID string, idValue, relationshipIDValue any) (string, pgx.NamedArgs) {
	sqlString := fmt.Sprintf(
		"INSERT INTO \"%s\" (%s, %s) VALUES (@%s, @%s) ;",
		relationshipTable,
		id,
		relationshipID,
		id,
		relationshipID,
	)
	args := pgx.NamedArgs{id: idValue, relationshipID: relationshipIDValue}
	return sqlString, args
}

func (receiver *Client) AddManyToManyRelation(ctx context.Context, relationshipTable, id, relationshipID string, idValue, relationshipIDValue any) error {
	sqlString, args := CreateAddRelationQuery(relationshipTable, id, relationshipID, idValue, relationshipIDValue)

	_, err := receiver.Database.Exec(ctx, sqlString, args)
	return err
}

func DeleteRelationQuery(relationshipTable, id, relationshipID string, idValue, relationshipIDValue any) (string, pgx.NamedArgs) {
	sqlString := fmt.Sprintf(
		"DELETE FROM \"%s\" WHERE %s = @%s AND %s = @%s;",
		relationshipTable,
		id,
		id,
		relationshipID,
		relationshipID,
	)
	args := pgx.NamedArgs{id: idValue, relationshipID: relationshipIDValue}
	return sqlString, args
}

func (receiver *Client) DeleteManyToManyRelation(ctx context.Context, relationshipTable, id, relationshipID string, idValue, relationshipIDValue any) error {
	sqlString, args := DeleteRelationQuery(relationshipTable, id, relationshipID, idValue, relationshipIDValue)

	_, err := receiver.Database.Exec(ctx, sqlString, args)
	return err
}

func IsExistRelationQuery(relationshipTable, id, relationshipID string, idValue, relationshipIDValue any) (string, pgx.NamedArgs) {
	sqlString := fmt.Sprintf(
		"SELECT 1 FROM \"%s\" WHERE %s = @%s AND %s = @%s;",
		relationshipTable,
		id,
		id,
		relationshipID,
		relationshipID,
	)
	args := pgx.NamedArgs{id: idValue, relationshipID: relationshipIDValue}
	return sqlString, args
}

func (receiver *Client) ExistManyToManyRelation(ctx context.Context, relationshipTable, id, relationshipID string, idValue, relationshipIDValue any) (bool, error) {
	sqlString, args := IsExistRelationQuery(relationshipTable, id, relationshipID, idValue, relationshipIDValue)

	rows, err := receiver.Database.Query(ctx, sqlString, args)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func CreateBulkInsertQuery(args pgx.NamedArgs, fieldsList []map[string]any, fieldsListList [][]string) (string, string) {
	if len(fieldsListList) == 0 || len(fieldsList) == 0 {
		return "", ""
	}
	var names []string
	var values [][]string
	for _, n := range fieldsListList[0] {
		names = append(names, fmt.Sprintf("\"%s\"", n))
	}

	for i, fieldListItem := range fieldsListList {
		var currentValues []string
		for _, n := range fieldListItem {
			v := fieldsList[i][n]
			valueName := fmt.Sprintf("%d%s", i, n)
			currentValues = append(currentValues, fmt.Sprintf("@%s", valueName))
			args[valueName] = v
		}
		values = append(values, currentValues)
	}

	nameString := fmt.Sprintf("(%s)", strings.Join(names, ","))
	var valueString []string
	for _, value := range values {
		valueString = append(valueString, fmt.Sprintf("(%s)", strings.Join(value, ",")))
	}
	valuesString := strings.Join(valueString, ", ")
	return nameString, valuesString
}

func (receiver *Client) BulkCreate(ctx context.Context, tableName string, fieldsList []map[string]any, fieldsListList [][]string) error {
	if len(fieldsList) == 0 || len(fieldsListList) == 0 {
		return nil
	}
	args := pgx.NamedArgs{}

	nameString, valueString := CreateBulkInsertQuery(args, fieldsList, fieldsListList)

	sqlString := fmt.Sprintf("INSERT INTO \"%s\" %s VALUES %s;", tableName, nameString, valueString)

	_, err := receiver.Database.Exec(ctx, sqlString, args)
	if err != nil {
		return err
	}
	return nil
}

func (receiver *Client) BulkUpdate(ctx context.Context, tableName string, fields map[string]any, fieldsList []string, idName string, idValue []any) error {
	statements, args := CreateUpdateQuery(fields, fieldsList)

	sqlString := fmt.Sprintf("UPDATE \"%s\" SET %s WHERE \"%s\" = ANY(@idvalue)", tableName, statements, idName)

	args["idvalue"] = idValue
	_, err := receiver.Database.Exec(ctx, sqlString, args)
	if err != nil {
		return err
	}
	return nil
}

func (receiver *Client) BulkDelete(ctx context.Context, tableName string, idName string, idValue []any) error {
	sqlString := fmt.Sprintf("DELETE FROM \"%s\" WHERE \"%s\" = ANY(@idvalue);", tableName, idName)

	args := pgx.NamedArgs{}
	args["idvalue"] = idValue
	_, err := receiver.Database.Exec(ctx, sqlString, args)
	if err != nil {
		return err
	}
	return nil
}

func (receiver *Client) DistinctString(ctx context.Context, list []*WhereList, model Model, field string, orders []*Order, paging *Paging) ([]string, error) {
	args := pgx.NamedArgs{}
	whereStrings := createTableWhereSql(list, args, model.GetDBName())
	relationSqlString, relationWhereStrings := createTableRelationWhereSql(model, args)
	mainTableWhereString := strings.Join(whereStrings, " OR ")

	var allWhereStrings []string
	if mainTableWhereString != "" {
		allWhereStrings = append(allWhereStrings, mainTableWhereString)
	}
	if len(relationWhereStrings) > 0 {
		allWhereStrings = append(allWhereStrings, relationWhereStrings...)
	}

	summedWhereString := ""
	if len(allWhereStrings) > 0 {
		summedWhereString = fmt.Sprintf(" WHERE (%s)", strings.Join(allWhereStrings, " AND "))
	}

	var orderStrings []string
	for _, order := range orders {
		orderStrings = append(orderStrings, order.String())
	}
	orderString := ""
	if len(orderStrings) > 0 {
		orderString = "ORDER BY " + strings.Join(orderStrings, ", ")
	}

	pagingString := ""
	if paging != nil {
		pagingString = paging.String()
	}

	sqlString := fmt.Sprintf("SELECT DISTINCT (\"%s\".\"%s\")::text FROM \"%s\" %s%s %s %s;",
		model.GetDBName(),
		field,
		model.GetDBName(),
		relationSqlString,
		summedWhereString,
		orderString,
		pagingString)

	rows, err := receiver.Database.Query(ctx, sqlString, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var val *string
		if err := rows.Scan(&val); err != nil {
			return nil, err
		}
		if val != nil {
			result = append(result, *val)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (receiver *Client) Count(ctx context.Context, list []*WhereList, model Model) (int, error) {
	args := pgx.NamedArgs{}
	whereStrings := createTableWhereSql(list, args, model.GetDBName())
	relationSqlString, relationWhereStrings := createTableRelationWhereSql(model, args)
	mainTableWhereString := strings.Join(whereStrings, " OR ")

	var allWhereStrings []string
	if mainTableWhereString != "" {
		allWhereStrings = append(allWhereStrings, mainTableWhereString)
	}
	if len(relationWhereStrings) > 0 {
		allWhereStrings = append(allWhereStrings, relationWhereStrings...)
	}

	summedWhereString := ""
	if len(allWhereStrings) > 0 {
		summedWhereString = fmt.Sprintf(" WHERE (%s)", strings.Join(allWhereStrings, " AND "))
	}

	sqlString := fmt.Sprintf("SELECT COUNT(*) FROM \"%s\" %s%s;",
		model.GetDBName(),
		relationSqlString,
		summedWhereString)

	var count int
	// QueryRow handles closing the underlying connection automatically
	err := receiver.Database.QueryRow(ctx, sqlString, args).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (receiver *Client) DeleteWhere(ctx context.Context, tableName string, list []*WhereList, model Model, idName string) (int64, error) {
	args := pgx.NamedArgs{}
	whereStrings := createTableWhereSql(list, args, tableName)
	relationSqlString, relationWhereStrings := createTableRelationWhereSql(model, args)
	mainTableWhereString := strings.Join(whereStrings, " OR ")

	var allWhereStrings []string
	if mainTableWhereString != "" {
		allWhereStrings = append(allWhereStrings, mainTableWhereString)
	}
	if len(relationWhereStrings) > 0 {
		allWhereStrings = append(allWhereStrings, relationWhereStrings...)
	}

	summedWhereString := ""
	if len(allWhereStrings) > 0 {
		summedWhereString = fmt.Sprintf(" WHERE (%s)", strings.Join(allWhereStrings, " AND "))
	}

	var sqlString string
	if relationSqlString != "" {
		sqlString = fmt.Sprintf("DELETE FROM \"%s\" WHERE \"%s\" IN (SELECT \"%s\".\"%s\" FROM \"%s\" %s%s);",
			tableName, idName, tableName, idName, tableName, relationSqlString, summedWhereString)
	} else {
		sqlString = fmt.Sprintf("DELETE FROM \"%s\"%s;", tableName, summedWhereString)
	}

	res, err := receiver.Database.Exec(ctx, sqlString, args)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func (receiver *Client) UpdateWhere(ctx context.Context, tableName string, set map[string]any, list []*WhereList, model Model, idName string) (int64, error) {
	if len(set) == 0 {
		return 0, nil
	}

	var fieldsList []string
	for k := range set {
		fieldsList = append(fieldsList, k)
	}

	statements, args := CreateUpdateQuery(set, fieldsList)

	whereStrings := createTableWhereSql(list, args, tableName)
	relationSqlString, relationWhereStrings := createTableRelationWhereSql(model, args)
	mainTableWhereString := strings.Join(whereStrings, " OR ")

	var allWhereStrings []string
	if mainTableWhereString != "" {
		allWhereStrings = append(allWhereStrings, mainTableWhereString)
	}
	if len(relationWhereStrings) > 0 {
		allWhereStrings = append(allWhereStrings, relationWhereStrings...)
	}

	summedWhereString := ""
	if len(allWhereStrings) > 0 {
		summedWhereString = fmt.Sprintf(" WHERE (%s)", strings.Join(allWhereStrings, " AND "))
	}

	var sqlString string
	if relationSqlString != "" {
		sqlString = fmt.Sprintf("UPDATE \"%s\" SET %s WHERE \"%s\" IN (SELECT \"%s\".\"%s\" FROM \"%s\" %s%s);",
			tableName, statements, idName, tableName, idName, tableName, relationSqlString, summedWhereString)
	} else {
		sqlString = fmt.Sprintf("UPDATE \"%s\" SET %s%s;", tableName, statements, summedWhereString)
	}

	res, err := receiver.Database.Exec(ctx, sqlString, args)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func (receiver *Client) AggregateRows(ctx context.Context, list []*WhereList, model Model, aggregate *Aggregate) (func() error, func(), error) {
	sqlString, args := CreateAggregateQuery(list, model, aggregate)

	rows, err := receiver.Database.Query(ctx, sqlString, args)
	if err != nil {
		return nil, nil, err
	}

	scanNext := func() error {
		if err := rows.Err(); err != nil {
			rows.Close()
			return err
		}
		if rows.Next() {
			err := rows.Scan(aggregate.aggregateValues...)
			if err != nil {
				rows.Close()
				return err
			}
			return nil
		}
		rows.Close()
		return ErrFinalRow
	}

	closeRows := func() {
		rows.Close()
	}

	return scanNext, closeRows, nil
}

func (receiver *Client) AggregateRowsSeq(ctx context.Context, list []*WhereList, model Model, aggregate *Aggregate) iter.Seq2[int, error] {
	return func(yield func(int, error) bool) {
		sqlString, args := CreateAggregateQuery(list, model, aggregate)
		rows, err := receiver.Database.Query(ctx, sqlString, args)
		if err != nil {
			yield(0, err)
			return
		}
		defer rows.Close()

		idx := 0
		for rows.Next() {
			err := rows.Scan(aggregate.aggregateValues...)
			if !yield(idx, err) {
				return
			}
			if err != nil {
				return
			}
			idx++
		}
		if err := rows.Err(); err != nil {
			yield(idx, err)
		}
	}
}
