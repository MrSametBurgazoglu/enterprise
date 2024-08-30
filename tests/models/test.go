package models

import (
	"context"
	"github.com/MrSametBurgazoglu/enterprise/client"

	"github.com/google/uuid"
	"time"
)

const TestTableName = "test"

const (
	TestIDField        string = "id"
	TestNameField      string = "name"
	TestCreatedAtField string = "created_at"
)

var databaseTestOperationHook = func(operationInfo *client.OperationInfo, model *Test, operationFunc func() error) error {
	return operationFunc()
}

var databaseTestListOperationHook = func(operationInfo *client.OperationInfo, model *TestList, operationFunc func() error) error {
	return operationFunc()
}

func SetDatabaseTestOperationHook(f func(operationInfo *client.OperationInfo, model *Test, operationFunc func() error) error) {
	databaseTestOperationHook = f
}

func SetDatabaseTestListOperationHook(f func(operationInfo *client.OperationInfo, model *TestList, operationFunc func() error) error) {
	databaseTestListOperationHook = f
}

func NewTest(ctx context.Context, dc client.DatabaseClient) *Test {
	v := &Test{client: client.NewClient(dc), ctx: ctx}
	v.relations = new(client.RelationList)
	v.relations.RelationMap = make(map[string]*client.Relation)
	v.changedFields = make(map[string]any)
	v.result.Init()
	v.Default()
	return v
}

func NewRelationTest(ctx context.Context, dc client.DatabaseClient) *Test {
	v := &Test{client: client.NewClient(dc), ctx: ctx}
	v.relations = new(client.RelationList)
	v.relations.RelationMap = make(map[string]*client.Relation)
	v.changedFields = make(map[string]any)
	v.result.Init()
	return v
}

type Test struct {
	id uuid.UUID

	name string

	createdat time.Time

	changedFields     map[string]any
	changedFieldsList []string
	serialFields      []*client.SelectedField

	ctx    context.Context
	client *client.Client
	TestPredicate
	relations *client.RelationList

	DenemeList *DenemeList

	result TestResult
}

func (t *Test) GetDBName() string {
	return TestTableName
}

func (t *Test) GetSelector() *TestResult {
	t.result.selectedFields = nil
	return &t.result
}

func (t *Test) GetRelationList() *client.RelationList {
	return t.relations
}

func (t *Test) IsExist() bool {
	var v uuid.UUID
	return t.id != v
}

func (t *Test) GetPrimaryKey() uuid.UUID {
	return t.id
}

func NewTestList(ctx context.Context, dc client.DatabaseClient) *TestList {
	v := &TestList{client: client.NewClient(dc), ctx: ctx}
	v.relations = new(client.RelationList)
	v.relations.RelationMap = make(map[string]*client.Relation)
	v.result.Init()
	return v
}

func NewRelationTestList(ctx context.Context, dc client.DatabaseClient) *TestList {
	v := &TestList{client: client.NewClient(dc), ctx: ctx}
	v.relations = new(client.RelationList)
	v.relations.RelationMap = make(map[string]*client.Relation)
	v.result.Init()
	return v
}

type TestList struct {
	Items []*Test

	ctx    context.Context
	client *client.Client
	TestPredicate
	order     []*client.Order
	paging    *client.Paging
	relations *client.RelationList
	result    TestResult
}

func (t *TestList) GetDBName() string {
	return TestTableName
}

func (t *TestList) GetRelationList() *client.RelationList {
	return t.relations
}

func (t *TestList) IsExist() bool {
	return t.Items[len(t.Items)-1].IsExist()
}

func (t *Test) SetID(v uuid.UUID) {
	t.id = v
	t.SetIDField()
}
func (t *Test) SetName(v string) {
	t.name = v
	t.SetNameField()
}
func (t *Test) SetCreatedAt(v time.Time) {
	t.createdat = v
	t.SetCreatedAtField()
}

func (t *Test) SetIDNillable(v *uuid.UUID) {
	if v == nil {
		return
	}
	t.SetID(*v)
}
func (t *Test) SetNameNillable(v *string) {
	if v == nil {
		return
	}
	t.SetName(*v)
}
func (t *Test) SetCreatedAtNillable(v *time.Time) {
	if v == nil {
		return
	}
	t.SetCreatedAt(*v)
}

func (t *Test) ParseID(v string) error {
	parsedID, err := uuid.Parse(v)
	if err != nil {
		return err
	}
	t.id = parsedID
	t.SetIDField()
	return nil
}

func (t *Test) FormatCreatedAt(v string) string {

	return t.createdat.Format(v)
}

func (t *Test) ParseCreatedAt(layout, value string) error {
	parsedTime, err := time.Parse(layout, value)
	if err != nil {
		return err
	}
	t.createdat = parsedTime
	t.SetCreatedAtField()
	return nil
}

func (t *Test) IDIN(v ...uuid.UUID) bool {
	for _, x := range v {
		if t.id == x {
			return true
		}
	}
	return false
}

func (t *Test) NameIN(v ...string) bool {
	for _, x := range v {
		if t.name == x {
			return true
		}
	}
	return false
}

func (t *Test) CreatedAtIN(v ...time.Time) bool {
	for _, x := range v {
		if t.createdat == x {
			return true
		}
	}
	return false
}

func (t *Test) IDNotIN(v ...uuid.UUID) bool {
	for _, x := range v {
		if t.id == x {
			return false
		}
	}
	return true
}

func (t *Test) NameNotIN(v ...string) bool {
	for _, x := range v {
		if t.name == x {
			return false
		}
	}
	return true
}

func (t *Test) CreatedAtNotIN(v ...time.Time) bool {
	for _, x := range v {
		if t.createdat == x {
			return false
		}
	}
	return true
}

func (t *Test) GetID() uuid.UUID {
	return t.id
}
func (t *Test) GetName() string {
	return t.name
}
func (t *Test) GetCreatedAt() time.Time {
	return t.createdat
}

func (t *Test) SetIDField() {
	if _, exist := t.changedFields[TestIDField]; !exist {
		t.changedFields[TestIDField] = t.id
		t.changedFieldsList = append(t.changedFieldsList, TestIDField)
	}

}
func (t *Test) SetNameField() {
	if _, exist := t.changedFields[TestNameField]; !exist {
		t.changedFields[TestNameField] = t.name
		t.changedFieldsList = append(t.changedFieldsList, TestNameField)
	}

}
func (t *Test) SetCreatedAtField() {
	if _, exist := t.changedFields[TestCreatedAtField]; !exist {
		t.changedFields[TestCreatedAtField] = t.createdat
		t.changedFieldsList = append(t.changedFieldsList, TestCreatedAtField)
	}

}

func (t *Test) WithDenemeList(opts ...func(*DenemeList)) {
	t.DenemeList = NewRelationDenemeList(t.ctx, t.client.Database)
	for _, opt := range opts {
		opt(t.DenemeList)
	}
	t.result.Deneme = new(DenemeResult)
	t.result.Deneme.Init()
	t.result.relations = append(t.result.relations, t.result.Deneme)
	t.result.relationsMap["deneme"] = t.result.Deneme
	for _, Relation := range t.DenemeList.relations.Relations {
		t.result.Deneme.relations = append(t.result.Deneme.relations, Relation.RelationResult)
		t.result.Deneme.relationsMap[Relation.RelationTable] = Relation.RelationResult
	}
	t.relations.Relations = append(t.relations.Relations,
		&client.Relation{
			RelationModel:  t.DenemeList,
			RelationTable:  "deneme",
			RelationResult: t.result.Deneme,
			Where:          t.DenemeList.where,

			RelationWhere: &client.RelationCondition{
				RelationValue: "test_id",
				TableValue:    "id",
			},
		},
	)
	t.relations.RelationMap["deneme"] = t.relations.Relations[len(t.relations.Relations)-1]
}

func (t *TestList) WithDenemeList(opts ...func(*DenemeList)) {
	v := NewRelationDenemeList(t.ctx, t.client.Database)
	for _, opt := range opts {
		opt(v)
	}
	t.result.Deneme = new(DenemeResult)
	t.result.Deneme.Init()
	t.result.relations = append(t.result.relations, t.result.Deneme)
	t.result.relationsMap["deneme"] = t.result.Deneme
	for _, Relation := range v.relations.Relations {
		t.result.Deneme.relations = append(t.result.Deneme.relations, Relation.RelationResult)
		t.result.Deneme.relationsMap[Relation.RelationTable] = Relation.RelationResult
	}
	t.relations.Relations = append(t.relations.Relations,
		&client.Relation{
			RelationModel:  v,
			RelationTable:  "deneme",
			RelationResult: t.result.Deneme,
			Where:          v.where,
			RelationWhere: &client.RelationCondition{
				RelationValue: "test_id",
				TableValue:    "id",
			},
		},
	)
	t.relations.RelationMap["deneme"] = t.relations.Relations[len(t.relations.Relations)-1]
}

func (t *TestList) cleanDenemeList() {
	Relation := t.Items[len(t.Items)-1].relations
	p := 0
	for i, v := range Relation.Relations {
		if v.RelationTable == "deneme" {
			p = i
		}
	}
	Relation.Relations = append(Relation.Relations[:p], Relation.Relations[p+1:]...)
}

func (t *Test) Default() {
	t.id = uuid.New()
	t.changedFields[TestIDField] = t.id
	t.changedFieldsList = append(t.changedFieldsList, TestIDField)

}

func (t *Test) SetResult(result client.Result) {
	if t == nil {
		v := NewRelationTest(t.ctx, t.client.Database)
		*t = *v
	}
	t.result = *result.(*TestResult)
}

func (t *TestList) SetResult(result client.Result) {
	if t == nil {
		v := NewRelationTestList(t.ctx, t.client.Database)
		*t = *v
	}
	t.result = *result.(*TestResult)
}

func (t *Test) ScanResult() {
	t.id = t.result.id
	t.name = t.result.name
	t.createdat = t.result.createdat

	if _, ok := t.relations.RelationMap["deneme"]; ok {
		if t.DenemeList == nil {
			t.DenemeList = NewRelationDenemeList(t.ctx, t.client.Database)
		}
		t.DenemeList.relations = t.relations.RelationMap["deneme"].RelationModel.GetRelationList()
		t.DenemeList.SetResult(t.result.relationsMap["deneme"])
		t.DenemeList.ScanResult()
	}
}

func (t *Test) CheckPrimaryKey(v uuid.UUID) bool {
	return t.id == v
}

func (t *TestList) ScanResult() {
	var v *Test
	if len(t.Items) == 0 {
		v = NewRelationTest(t.ctx, t.client.Database)
		t.Items = append(t.Items, v)
	} else {
		for _, item := range t.Items {
			if item.CheckPrimaryKey(t.result.id) {
				v = item
				break
			}
		}
	}
	if v == nil {
		v = NewRelationTest(t.ctx, t.client.Database)
		t.Items = append(t.Items, v)
	}
	v.result = t.result
	v.relations = t.relations
	v.ScanResult()
}

func (t *Test) GetContext() context.Context {
	return t.ctx
}

func (t *Test) Get() error {
	return databaseTestOperationHook(
		client.NewOperationInfo(
			TestTableName,
			client.OperationTypeGet,
		),
		t,
		func() error {
			return t.client.Get(t.ctx, t.where, t, &t.result)
		},
	)
}

func (t *Test) Refresh() error {
	return databaseTestOperationHook(
		client.NewOperationInfo(
			TestTableName,
			client.OperationTypeRefresh,
		),
		t,
		func() error {
			return t.client.Refresh(t.ctx, t, &t.result, TestIDField, t.id)
		},
	)
}

func (t *Test) Create() error {
	return databaseTestOperationHook(
		client.NewOperationInfo(
			TestTableName,
			client.OperationTypeCreate,
		),
		t,
		func() error {
			return t.client.Create(t.ctx, TestTableName, t.changedFields, t.changedFieldsList, t.serialFields)
		},
	)
}

func (t *Test) Update() error {
	return databaseTestOperationHook(
		client.NewOperationInfo(
			TestTableName,
			client.OperationTypeUpdate,
		),
		t,
		func() error {
			return t.client.Update(t.ctx, TestTableName, t.changedFields, t.changedFieldsList, TestIDField, t.id)
		},
	)
}

func (t *Test) Delete() error {
	return databaseTestOperationHook(
		client.NewOperationInfo(
			TestTableName,
			client.OperationTypeDelete,
		),
		t,
		func() error {
			return t.client.Delete(t.ctx, TestTableName, TestIDField, t.id)
		},
	)
}

func (t *TestList) List() error {
	return databaseTestListOperationHook(
		client.NewOperationInfo(
			TestTableName,
			client.OperationTypeList,
		),
		t,
		func() error {
			return t.client.List(t.ctx, t.where, t, &t.result, t.order, t.paging)
		},
	)
}

func (t *TestList) Aggregate(f func(aggregate *client.Aggregate)) (func() error, error) {
	a := new(client.Aggregate)
	f(a)
	return t.client.Aggregate(t.ctx, t.where, t, a)
}

func (t *TestList) Order(field string) *TestList {
	t.order = append(t.order, &client.Order{Field: field})
	return t
}

func (t *TestList) OrderDesc(field string) *TestList {
	t.order = append(t.order, &client.Order{Field: field, Desc: true})
	return t
}

func (t *TestList) Paging(skip, limit int) *TestList {
	t.paging = &client.Paging{Skip: skip, Limit: limit}
	return t
}

type TestResult struct {
	id        uuid.UUID
	name      string
	createdat time.Time

	selectedFields []*client.SelectedField

	Deneme *DenemeResult

	relations    []client.Result
	relationsMap map[string]client.Result
}

func (t *TestResult) Init() {
	t.selectedFields = []*client.SelectedField{}
	t.relationsMap = make(map[string]client.Result)
	t.prepare()
	t.SelectAll()
}

func (t *TestResult) GetSelectedFields() []*client.SelectedField {
	return t.selectedFields
}

func (t *TestResult) GetRelations() []client.Result {
	return t.relations
}

func (t *TestResult) prepare() {

}

func (t *TestResult) SelectID() {
	v := &client.SelectedField{Name: TestIDField, Value: &t.id}
	t.selectedFields = append(t.selectedFields, v)
}

func (t *TestResult) SelectName() {
	v := &client.SelectedField{Name: TestNameField, Value: &t.name}
	t.selectedFields = append(t.selectedFields, v)
}

func (t *TestResult) SelectCreatedAt() {
	v := &client.SelectedField{Name: TestCreatedAtField, Value: &t.createdat}
	t.selectedFields = append(t.selectedFields, v)
}

func (t *TestResult) GetDBName() string {
	return TestTableName
}

func (t *TestResult) SelectAll() {
	t.SelectID()
	t.SelectName()
	t.SelectCreatedAt()

}

func (t *TestResult) IsExist() bool {
	if t == nil {
		return false
	}
	var v uuid.UUID
	return t.id != v
}
