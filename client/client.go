package client

import (
	"context"
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
}

type Client struct {
	Database DatabaseClient
}

func NewClient(d DatabaseClient) *Client {
	return &Client{Database: d}
}

func CreateTableNameAndAddresses(result Result, selectedNames []string, selectedAddress []any) ([]string, []any) {
	for _, field := range result.GetSelectedFields() {
		selectedNames = append(selectedNames, fmt.Sprintf("\"%s\".\"%s\"", result.GetDBName(), field.Name))
		selectedAddress = append(selectedAddress, field.Value)
	}
	return selectedNames, selectedAddress
}

func CreateTableRelationNameAndAddresses(result Result, selectedNames []string, selectedAddress []any) ([]string, []any) {
	for _, field := range result.GetSelectedFields() {
		selectedNames = append(selectedNames, fmt.Sprintf("\"%s\".\"%s\"", result.GetDBName(), field.Name))
		selectedAddress = append(selectedAddress, field.Value)
	}
	return selectedNames, selectedAddress
}

func CreateTableRelationsNameAndAddresses(result Result, selectedNames []string, selectedAddress []any) ([]string, []any) {
	for _, relation := range result.GetRelations() {
		selectedNames, selectedAddress = CreateTableRelationNameAndAddresses(relation, selectedNames, selectedAddress)
		selectedNames, selectedAddress = CreateTableRelationsNameAndAddresses(relation, selectedNames, selectedAddress)
	}
	return selectedNames, selectedAddress
}

func CreateTableWhereSql(list []*WhereList, args pgx.NamedArgs, dbName string) []string {
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

func CreateTableRelationWhereSql(model Model, args pgx.NamedArgs) (string, []string) {
	sqlString := ""
	var relationWhereStrings []string
	for _, rel := range model.GetRelationList().Relations {
		sql := rel.GetJoinString(model.GetDBName())
		sqlString += sql
		if rel.IsRelationHaveWhereClause() {
			res := rel.ParseWhere()
			relationWhereStrings = append(relationWhereStrings, res.SqlString)
			for i, s := range res.Names {
				args[s] = res.Values[i]
			}
		}
		a, b := CreateTableRelationWhereSql(rel.RelationModel, args)
		sqlString += " " + a
		relationWhereStrings = append(relationWhereStrings, b...)
	}
	return sqlString, relationWhereStrings
}

func CreateSelectQuery(list []*WhereList, model Model, result Result) (string, []any, pgx.NamedArgs) {
	var selectedNames []string
	var selectedAddress []any
	selectedNames, selectedAddress = CreateTableNameAndAddresses(result, selectedNames, selectedAddress)
	selectedNames, selectedAddress = CreateTableRelationsNameAndAddresses(result, selectedNames, selectedAddress)

	args := pgx.NamedArgs{}
	whereStrings := CreateTableWhereSql(list, args, result.GetDBName())
	relationSqlString, relationWhereStrings := CreateTableRelationWhereSql(model, args)
	mainTableWhereString := strings.Join(whereStrings, " OR ")
	summedRelationWhereStrings := ""
	if len(relationWhereStrings) > 0 {
		summedRelationWhereStrings = fmt.Sprintf(" AND %s", strings.Join(relationWhereStrings, " AND "))
	}

	names := strings.Join(selectedNames, ", ")
	sqlString := fmt.Sprintf("SELECT %s FROM %s %s WHERE (%s %s);",
		names,
		model.GetDBName(),
		relationSqlString,
		mainTableWhereString,
		summedRelationWhereStrings)

	return sqlString, selectedAddress, args
}

func CreateSelectListQuery(list []*WhereList, model Model, result Result, orders []*Order, paging *Paging) (string, []any, pgx.NamedArgs) {
	var selectedNames []string
	var selectedAddress []any
	selectedNames, selectedAddress = CreateTableNameAndAddresses(result, selectedNames, selectedAddress)
	selectedNames, selectedAddress = CreateTableRelationsNameAndAddresses(result, selectedNames, selectedAddress)

	args := pgx.NamedArgs{}
	whereStrings := CreateTableWhereSql(list, args, result.GetDBName())
	relationSqlString, relationWhereStrings := CreateTableRelationWhereSql(model, args)
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

func ScanFirstRow(rows pgx.Rows, model Model, selectedAddress []any) (error, bool) {
	if err := rows.Err(); err != nil {
		return err, false
	}
	if rows.Next() {
		err := rows.Scan(selectedAddress...)
		if err != nil {
			return err, false
		}
		model.ScanResult()
		return nil, false
	} else {
		return fmt.Errorf("not found"), true
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

func (receiver *Client) Get(ctx context.Context, list []*WhereList, model Model, result Result) (error, bool) {
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

	sqlString, selectedAddress, args := CreateSelectQuery(list, model, result)

	rows, err := receiver.Database.Query(ctx, sqlString, args)
	if err != nil {
		return err, false
	}

	err, notFound := ScanFirstRow(rows, model, selectedAddress)
	if err != nil {
		return err, notFound
	}

	err = ScanNextRows(rows, model, selectedAddress)
	if err != nil {
		return err, false
	}

	return nil, true
}

func (receiver *Client) Refresh(ctx context.Context, model Model, result Result, idName string, idValue any) error {
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

	var selectedNames []string
	var selectedAddresses []any

	selectedNames, selectedAddresses = CreateTableNameAndAddresses(result, selectedNames, selectedAddresses)
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

	err, _ = ScanFirstRow(rows, model, selectedAddresses)
	if err != nil {
		return err
	}

	return nil
}

func CreateInsertQuery(fields map[string]any) (string, string, pgx.NamedArgs) {
	args := pgx.NamedArgs{}
	var names []string
	var values []string
	for n, v := range fields {
		names = append(names, n)
		values = append(values, "@"+n)
		args[n] = v
	}
	nameString := fmt.Sprintf("(%s)", strings.Join(names, ","))
	valueString := fmt.Sprintf("(%s)", strings.Join(values, ","))
	return nameString, valueString, args
}

func (receiver *Client) Create(ctx context.Context, tableName string, fields map[string]any, serialFields []*SelectedField) error {
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

	nameString, valueString, args := CreateInsertQuery(fields)

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

func CreateUpdateQuery(fields map[string]any) (string, pgx.NamedArgs) {
	args := pgx.NamedArgs{}
	var statements []string
	for n, v := range fields {
		args[n] = v
		statements = append(statements, fmt.Sprintf("\"%s\" = @%s", n, n))
	}
	return strings.Join(statements, ","), args
}

func (receiver *Client) Update(ctx context.Context, tableName string, fields map[string]any, idName string, idValue any) error {
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

	statements, args := CreateUpdateQuery(fields)

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

func ScanListFirstRow(rows pgx.Rows, model Model, selectedAddress []any) (error, bool) {
	if err := rows.Err(); err != nil {
		return err, false
	}
	if rows.Next() {
		err := rows.Scan(selectedAddress...)
		if err != nil {
			return err, false
		}
		model.ScanResult()
		return nil, false
	} else {
		return fmt.Errorf("not found"), true
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

func (receiver *Client) List(ctx context.Context, list []*WhereList, model Model, result Result, orders []*Order, paging *Paging) (error, bool) {
	receiver.Database.BeginHook()
	defer receiver.Database.EndHook()

	sqlString, selectedAddress, args := CreateSelectListQuery(list, model, result, orders, paging)

	rows, err := receiver.Database.Query(ctx, sqlString, args)
	if err != nil {
		return err, false
	}

	err, notFound := ScanListFirstRow(rows, model, selectedAddress)
	if err != nil {
		return err, notFound
	}

	err = ScanListNextRows(rows, model, selectedAddress)
	if err != nil {
		return err, false
	}

	return nil, true
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
	whereStrings := CreateTableWhereSql(list, args, model.GetDBName())
	relationSqlString, relationWhereStrings := CreateTableRelationWhereSql(model, args)
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
