package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
)

//todo add return serial values from create

type DatabaseClient interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	BeginHook()
	EndHook()
}

type DatabaseTransactionClient interface {
	DatabaseClient
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	SavePoint(ctx context.Context) (DatabaseTransactionClient, error)
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
		res := item.Parse(dbName)
		for i, s := range res.Names {
			args[s] = res.Values[i]
		}
		whereStrings = append(whereStrings, res.SqlString)
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
			res := rel.parseWhere()
			relationWhereStrings = append(relationWhereStrings, res.SqlString)
			for i, s := range res.Names {
				args[s] = res.Values[i]
			}
		}
		a, b := createTableRelationWhereSql(rel.RelationModel, args)
		sqlString += " " + a
		relationWhereStrings = append(relationWhereStrings, b...)
	}
	return sqlString, relationWhereStrings
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
	summedRelationWhereStrings := strings.Join(allWhereStrings, " AND ")

	names := strings.Join(selectedNames, ", ")
	sqlString := fmt.Sprintf("SELECT %s FROM %s %s WHERE (%s);",
		names,
		model.GetDBName(),
		relationSqlString,
		summedRelationWhereStrings)

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
	summedRelationWhereStrings := ""
	if len(relationWhereStrings) > 0 {
		summedRelationWhereStrings = fmt.Sprintf(" AND %s", strings.Join(relationWhereStrings, " AND "))
	}

	var orderStrings []string
	for _, order := range orders {
		orderStrings = append(orderStrings, order.String())
	}
	orderString := strings.Join(orderStrings, ", ")

	pagingString := ""
	if paging != nil {
		pagingString = paging.String()
	}

	names := strings.Join(selectedNames, ", ")
	sqlString := fmt.Sprintf("SELECT %s FROM %s %s WHERE (%s %s) %s %s;",
		names,
		model.GetDBName(),
		relationSqlString,
		mainTableWhereString,
		summedRelationWhereStrings,
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
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

	sqlString, selectedAddress, args := CreateSelectQuery(list, model, result)

	rows, err := receiver.Database.Query(ctx, sqlString, args)
	if err != nil {
		return err
	}

	err = ScanFirstRow(rows, model, selectedAddress)
	if err != nil {
		return err
	}

	err = ScanNextRows(rows, model, selectedAddress)
	if err != nil {
		return err
	}

	return nil
}

func (receiver *Client) Refresh(ctx context.Context, model Model, result Result, idName string, idValue any) error {
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

	var selectedNames []string
	var selectedAddresses []any

	selectedNames, selectedAddresses = createTableNameAndAddresses(result, selectedNames, selectedAddresses)
	names := strings.Join(selectedNames, ", ")
	args := pgx.NamedArgs{
		"idvalue": idValue,
	}

	sqlString := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s = @idvalue",
		names,
		model.GetDBName(),
		idName,
	)

	rows, err := receiver.Database.Query(ctx, sqlString, args)
	if err != nil {
		return err
	}

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
		names = append(names, n)
		values = append(values, "@"+n)
		args[n] = v
	}
	nameString := fmt.Sprintf("(%s)", strings.Join(names, ","))
	valueString := fmt.Sprintf("(%s)", strings.Join(values, ","))
	return nameString, valueString, args
}

func (receiver *Client) Create(ctx context.Context, tableName string, fields map[string]any, fieldsList []string, serialFields []*SelectedField) error {
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

	nameString, valueString, args := CreateInsertQuery(fields, fieldsList)

	var serialSql string
	var serialFieldAddresses []any
	if len(serialFields) > 0 {
		serialSql = "RETURNING %s"
		var serialNames []string
		for _, field := range serialFields {
			serialNames = append(serialNames, field.Name)
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
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

	statements, args := CreateUpdateQuery(fields, fieldslist)

	sqlString := fmt.Sprintf("UPDATE \"%s\" SET %s WHERE %s = @idvalue", tableName, statements, idName)

	args["idvalue"] = idValue
	_, err := receiver.Database.Exec(ctx, sqlString, args)
	if err != nil {
		return err
	}
	return nil
}

func (receiver *Client) Delete(ctx context.Context, tableName string, idName string, idValue any) error {
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

	sqlString := fmt.Sprintf("DELETE FROM \"%s\" WHERE %s = @idvalue;", tableName, idName)

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
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

	sqlString, selectedAddress, args := CreateSelectListQuery(list, model, result, orders, paging)

	rows, err := receiver.Database.Query(ctx, sqlString, args)
	if err != nil {
		return err
	}

	err = ScanListFirstRow(rows, model, selectedAddress)
	if err != nil {
		return err
	}

	err = ScanListNextRows(rows, model, selectedAddress)
	if err != nil {
		return err
	}

	return nil
}

func ScanValues(rows pgx.Rows, selectedAddress []any) error {
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
	summedRelationWhereStrings := ""
	if len(relationWhereStrings) > 0 {
		summedRelationWhereStrings = fmt.Sprintf(" AND %s", strings.Join(relationWhereStrings, " AND "))
	}

	selectedNames := make([]string, len(aggregate.aggregateFields))
	for i := 0; i < len(aggregate.aggregateFields); i++ {
		selectedNames[i] = fmt.Sprintf(aggregate.aggregateFormats[i], aggregate.aggregateFields[i])
	}

	groupBys := make([]string, len(aggregate.groupByList))
	for i := 0; i < len(aggregate.groupByList); i++ {
		groupBys[i] = fmt.Sprintf("GROUP BY %s", aggregate.groupByList[i])
	}

	names := strings.Join(selectedNames, ", ")
	groupBy := strings.Join(groupBys, ", ")
	sqlString := fmt.Sprintf("SELECT %s FROM %s %s WHERE (%s %s) %s;",
		names,
		model.GetDBName(),
		relationSqlString,
		mainTableWhereString,
		summedRelationWhereStrings,
		groupBy)
	return sqlString, args
}

func (receiver *Client) Aggregate(ctx context.Context, list []*WhereList, model Model, aggregate *Aggregate) (func() error, error) {
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

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
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

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
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

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
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

	sqlString, args := IsExistRelationQuery(relationshipTable, id, relationshipID, idValue, relationshipIDValue)

	rows, err := receiver.Database.Query(ctx, sqlString, args)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func CreateBulkInsertQuery(args pgx.NamedArgs, fieldsList []map[string]any, fieldsListList [][]string) (string, string) {
	var names []string
	var values [][]string
	for _, n := range fieldsListList[0] {
		names = append(names, n)
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
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

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
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

	statements, args := CreateUpdateQuery(fields, fieldsList)

	sqlString := fmt.Sprintf("UPDATE \"%s\" SET %s WHERE %s IN (@idvalue)", tableName, statements, idName)

	args["idvalue"] = idValue
	_, err := receiver.Database.Exec(ctx, sqlString, args)
	if err != nil {
		return err
	}
	return nil
}

func (receiver *Client) BulkDelete(ctx context.Context, tableName string, idName string, idValue []any) error {
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

	sqlString := fmt.Sprintf("DELETE FROM \"%s\" WHERE %s IN (@idvalue);", tableName, idName)

	args := pgx.NamedArgs{}
	args["idvalue"] = idValue
	_, err := receiver.Database.Exec(ctx, sqlString, args)
	if err != nil {
		return err
	}
	return nil
}
