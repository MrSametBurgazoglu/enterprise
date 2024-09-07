package models

import (
	"context"
	"github.com/MrSametBurgazoglu/enterprise/client"

	"github.com/google/uuid"
)

const DenemeTableName = "deneme"

const (
	DenemeIDField         string = "id"
	DenemeTestIDField     string = "test_id"
	DenemeCountField      string = "count"
	DenemeIsActiveField   string = "is_active"
	DenemeDenemeTypeField string = "deneme_type"
)

var databaseDenemeOperationHook = func(operationInfo *client.OperationInfo, model *Deneme, operationFunc func() error) error {
	return operationFunc()
}

var databaseDenemeListOperationHook = func(operationInfo *client.OperationInfo, model *DenemeList, operationFunc func() error) error {
	return operationFunc()
}

func SetDatabaseDenemeOperationHook(f func(operationInfo *client.OperationInfo, model *Deneme, operationFunc func() error) error) {
	databaseDenemeOperationHook = f
}

func SetDatabaseDenemeListOperationHook(f func(operationInfo *client.OperationInfo, model *DenemeList, operationFunc func() error) error) {
	databaseDenemeListOperationHook = f
}

type DenemeType string

const (
	DenemeTypeTestType   DenemeType = "Test"
	DenemeTypeDenemeType DenemeType = "Deneme"
)

func NewDeneme(ctx context.Context, dc client.DatabaseClient) *Deneme {
	v := &Deneme{client: client.NewClient(dc), ctx: ctx}
	v.relations = new(client.RelationList)
	v.relations.RelationMap = make(map[string]*client.Relation)
	v.changedFields = make(map[string]any)
	v.result.Init()
	return v
}

func NewRelationDeneme(ctx context.Context, dc client.DatabaseClient) *Deneme {
	v := &Deneme{client: client.NewClient(dc), ctx: ctx}
	v.relations = new(client.RelationList)
	v.relations.RelationMap = make(map[string]*client.Relation)
	v.changedFields = make(map[string]any)
	v.result.Init()
	return v
}

type Deneme struct {
	id uuid.UUID

	testid *uuid.UUID

	count int

	isactive bool

	denemetype DenemeType

	changedFields     map[string]any
	changedFieldsList []string
	serialFields      []*client.SelectedField

	ctx    context.Context
	client *client.Client
	DenemePredicate
	relations *client.RelationList

	Test        *Test
	AccountList *AccountList

	result DenemeResult
}

func (t *Deneme) GetDBName() string {
	return DenemeTableName
}

func (t *Deneme) GetSelector() *DenemeResult {
	t.result.selectedFields = nil
	return &t.result
}

func (t *Deneme) GetRelationList() *client.RelationList {
	return t.relations
}

func (t *Deneme) IsExist() bool {
	var v uuid.UUID
	return t.id != v
}

func (t *Deneme) GetPrimaryKey() uuid.UUID {
	return t.id
}

func NewDenemeList(ctx context.Context, dc client.DatabaseClient) *DenemeList {
	v := &DenemeList{client: client.NewClient(dc), ctx: ctx}
	v.relations = new(client.RelationList)
	v.relations.RelationMap = make(map[string]*client.Relation)
	v.result.Init()
	return v
}

func NewRelationDenemeList(ctx context.Context, dc client.DatabaseClient) *DenemeList {
	v := &DenemeList{client: client.NewClient(dc), ctx: ctx}
	v.relations = new(client.RelationList)
	v.relations.RelationMap = make(map[string]*client.Relation)
	v.result.Init()
	return v
}

type DenemeList struct {
	Items []*Deneme

	ctx    context.Context
	client *client.Client
	DenemePredicate
	order     []*client.Order
	paging    *client.Paging
	relations *client.RelationList
	result    DenemeResult
}

func (t *DenemeList) GetDBName() string {
	return DenemeTableName
}

func (t *DenemeList) GetRelationList() *client.RelationList {
	return t.relations
}

func (t *DenemeList) IsExist() bool {
	return t.Items[len(t.Items)-1].IsExist()
}

func (t *Deneme) SetID(v uuid.UUID) {
	t.id = v
	t.SetIDField()
}
func (t *Deneme) SetTestID(v *uuid.UUID) {
	t.testid = v
	t.SetTestIDField()
}
func (t *Deneme) SetCount(v int) {
	t.count = v
	t.SetCountField()
}
func (t *Deneme) SetIsActive(v bool) {
	t.isactive = v
	t.SetIsActiveField()
}
func (t *Deneme) SetDenemeType(v DenemeType) {
	t.denemetype = v
	t.SetDenemeTypeField()
}

func (t *Deneme) SetTestIDValue(v uuid.UUID) {
	t.SetTestID(&v)
	t.SetTestIDField()
}

func (t *Deneme) SetIDNillable(v *uuid.UUID) {
	if v == nil {
		return
	}
	t.SetID(*v)
}

func (t *Deneme) SetCountNillable(v *int) {
	if v == nil {
		return
	}
	t.SetCount(*v)
}
func (t *Deneme) SetIsActiveNillable(v *bool) {
	if v == nil {
		return
	}
	t.SetIsActive(*v)
}
func (t *Deneme) SetDenemeTypeNillable(v *DenemeType) {
	if v == nil {
		return
	}
	t.SetDenemeType(*v)
}

func (t *Deneme) ParseID(v string) error {
	parsedID, err := uuid.Parse(v)
	if err != nil {
		return err
	}
	t.id = parsedID
	t.SetIDField()
	return nil
}

func (t *Deneme) ParseTestID(v string) error {
	parsedID, err := uuid.Parse(v)
	if err != nil {
		return err
	}
	t.testid = &parsedID
	t.SetTestIDField()
	return nil
}

func (t *Deneme) IDIN(v ...uuid.UUID) bool {
	for _, x := range v {
		if t.id == x {
			return true
		}
	}
	return false
}

func (t *Deneme) TestIDIN(v ...uuid.UUID) bool {
	if t.testid == nil {
		return false
	}
	for _, x := range v {
		if *t.testid == x {
			return true
		}
	}
	return false
}

func (t *Deneme) CountIN(v ...int) bool {
	for _, x := range v {
		if t.count == x {
			return true
		}
	}
	return false
}

func (t *Deneme) IsActiveIN(v ...bool) bool {
	for _, x := range v {
		if t.isactive == x {
			return true
		}
	}
	return false
}

func (t *Deneme) DenemeTypeIN(v ...DenemeType) bool {
	for _, x := range v {
		if t.denemetype == x {
			return true
		}
	}
	return false
}

func (t *Deneme) IDNotIN(v ...uuid.UUID) bool {
	for _, x := range v {
		if t.id == x {
			return false
		}
	}
	return true
}

func (t *Deneme) TestIDNotIN(v ...uuid.UUID) bool {
	if t.testid == nil {
		return true
	}
	for _, x := range v {
		if *t.testid == x {
			return false
		}
	}
	return true
}

func (t *Deneme) CountNotIN(v ...int) bool {
	for _, x := range v {
		if t.count == x {
			return false
		}
	}
	return true
}

func (t *Deneme) IsActiveNotIN(v ...bool) bool {
	for _, x := range v {
		if t.isactive == x {
			return false
		}
	}
	return true
}

func (t *Deneme) DenemeTypeNotIN(v ...DenemeType) bool {
	for _, x := range v {
		if t.denemetype == x {
			return false
		}
	}
	return true
}

func (t *Deneme) GetID() uuid.UUID {
	return t.id
}
func (t *Deneme) GetTestID() *uuid.UUID {
	return t.testid
}
func (t *Deneme) GetCount() int {
	return t.count
}
func (t *Deneme) GetIsActive() bool {
	return t.isactive
}
func (t *Deneme) GetDenemeType() DenemeType {
	return t.denemetype
}

func (t *Deneme) SetIDField() {
	if _, exist := t.changedFields[DenemeIDField]; !exist {
		t.changedFields[DenemeIDField] = t.id
		t.changedFieldsList = append(t.changedFieldsList, DenemeIDField)
	}

}
func (t *Deneme) SetTestIDField() {
	if _, exist := t.changedFields[DenemeTestIDField]; !exist {
		t.changedFields[DenemeTestIDField] = t.testid
		t.changedFieldsList = append(t.changedFieldsList, DenemeTestIDField)
	}

}
func (t *Deneme) SetCountField() {
	if _, exist := t.changedFields[DenemeCountField]; !exist {
		t.changedFields[DenemeCountField] = t.count
		t.changedFieldsList = append(t.changedFieldsList, DenemeCountField)
	}

}
func (t *Deneme) SetIsActiveField() {
	if _, exist := t.changedFields[DenemeIsActiveField]; !exist {
		t.changedFields[DenemeIsActiveField] = t.isactive
		t.changedFieldsList = append(t.changedFieldsList, DenemeIsActiveField)
	}

}
func (t *Deneme) SetDenemeTypeField() {
	if _, exist := t.changedFields[DenemeDenemeTypeField]; !exist {
		t.changedFields[DenemeDenemeTypeField] = t.denemetype
		t.changedFieldsList = append(t.changedFieldsList, DenemeDenemeTypeField)
	}

}

func (t *Deneme) WithTest(opts ...func(*Test)) {
	t.Test = NewRelationTest(t.ctx, t.client.Database)
	for _, opt := range opts {
		opt(t.Test)
	}
	t.result.Test = new(TestResult)
	t.result.Test.Init()
	t.result.relations = append(t.result.relations, t.result.Test)
	t.result.relationsMap["test"] = t.result.Test
	for _, Relation := range t.Test.relations.Relations {
		t.result.Test.relations = append(t.result.Test.relations, Relation.RelationResult)
		t.result.Test.relationsMap[Relation.RelationTable] = Relation.RelationResult
	}
	t.relations.Relations = append(t.relations.Relations,
		&client.Relation{
			RelationModel:  t.Test,
			RelationTable:  "test",
			RelationResult: t.result.Test,
			Where:          t.Test.where,

			RelationWhere: &client.RelationCondition{
				RelationValue: "id",
				TableValue:    "test_id",
			},
		},
	)
	t.relations.RelationMap["test"] = t.relations.Relations[len(t.relations.Relations)-1]
}
func (t *Deneme) WithAccountList(opts ...func(*AccountList)) {
	t.AccountList = NewRelationAccountList(t.ctx, t.client.Database)
	for _, opt := range opts {
		opt(t.AccountList)
	}
	t.result.Account = new(AccountResult)
	t.result.Account.Init()
	t.result.relations = append(t.result.relations, t.result.Account)
	t.result.relationsMap["account"] = t.result.Account
	for _, Relation := range t.AccountList.relations.Relations {
		t.result.Account.relations = append(t.result.Account.relations, Relation.RelationResult)
		t.result.Account.relationsMap[Relation.RelationTable] = Relation.RelationResult
	}
	t.relations.Relations = append(t.relations.Relations,
		&client.Relation{
			RelationModel:  t.AccountList,
			RelationTable:  "account",
			RelationResult: t.result.Account,
			Where:          t.AccountList.where,

			RelationWhere: &client.RelationCondition{
				RelationValue: "deneme_id",
				TableValue:    "id",
			},
		},
	)
	t.relations.RelationMap["account"] = t.relations.Relations[len(t.relations.Relations)-1]
}

func (t *DenemeList) WithTest(opts ...func(*Test)) {
	v := NewRelationTest(t.ctx, t.client.Database)
	for _, opt := range opts {
		opt(v)
	}
	t.result.Test = new(TestResult)
	t.result.Test.Init()
	t.result.relations = append(t.result.relations, t.result.Test)
	t.result.relationsMap["test"] = t.result.Test
	for _, Relation := range v.relations.Relations {
		t.result.Test.relations = append(t.result.Test.relations, Relation.RelationResult)
		t.result.Test.relationsMap[Relation.RelationTable] = Relation.RelationResult
	}
	t.relations.Relations = append(t.relations.Relations,
		&client.Relation{
			RelationModel:  v,
			RelationTable:  "test",
			RelationResult: t.result.Test,
			Where:          v.where,
			RelationWhere: &client.RelationCondition{
				RelationValue: "id",
				TableValue:    "test_id",
			},
		},
	)
	t.relations.RelationMap["test"] = t.relations.Relations[len(t.relations.Relations)-1]
}

func (t *DenemeList) cleanTest() {
	Relation := t.Items[len(t.Items)-1].relations
	p := 0
	for i, v := range Relation.Relations {
		if v.RelationTable == "test" {
			p = i
		}
	}
	Relation.Relations = append(Relation.Relations[:p], Relation.Relations[p+1:]...)
}
func (t *DenemeList) WithAccountList(opts ...func(*AccountList)) {
	v := NewRelationAccountList(t.ctx, t.client.Database)
	for _, opt := range opts {
		opt(v)
	}
	t.result.Account = new(AccountResult)
	t.result.Account.Init()
	t.result.relations = append(t.result.relations, t.result.Account)
	t.result.relationsMap["account"] = t.result.Account
	for _, Relation := range v.relations.Relations {
		t.result.Account.relations = append(t.result.Account.relations, Relation.RelationResult)
		t.result.Account.relationsMap[Relation.RelationTable] = Relation.RelationResult
	}
	t.relations.Relations = append(t.relations.Relations,
		&client.Relation{
			RelationModel:  v,
			RelationTable:  "account",
			RelationResult: t.result.Account,
			Where:          v.where,
			RelationWhere: &client.RelationCondition{
				RelationValue: "deneme_id",
				TableValue:    "id",
			},
		},
	)
	t.relations.RelationMap["account"] = t.relations.Relations[len(t.relations.Relations)-1]
}

func (t *DenemeList) cleanAccountList() {
	Relation := t.Items[len(t.Items)-1].relations
	p := 0
	for i, v := range Relation.Relations {
		if v.RelationTable == "account" {
			p = i
		}
	}
	Relation.Relations = append(Relation.Relations[:p], Relation.Relations[p+1:]...)
}

func (t *Deneme) SetDefaults() {
	t.id = uuid.New()
	t.changedFields[DenemeIDField] = t.id
	t.changedFieldsList = append(t.changedFieldsList, DenemeIDField)

	t.isactive = true
	t.changedFields[DenemeIsActiveField] = t.isactive
	t.changedFieldsList = append(t.changedFieldsList, DenemeIsActiveField)

}

func (t *Deneme) SetResult(result client.Result) {
	if t == nil {
		v := NewRelationDeneme(t.ctx, t.client.Database)
		*t = *v
	}
	t.result = *result.(*DenemeResult)
}

func (t *DenemeList) SetResult(result client.Result) {
	if t == nil {
		v := NewRelationDenemeList(t.ctx, t.client.Database)
		*t = *v
	}
	t.result = *result.(*DenemeResult)
}

func (t *Deneme) ScanResult() {
	t.id = t.result.id
	t.testid = t.result.testid
	t.count = t.result.count
	t.isactive = t.result.isactive
	t.denemetype = t.result.denemetype

	if _, ok := t.relations.RelationMap["test"]; ok {
		if t.Test == nil {
			t.Test = NewRelationTest(t.ctx, t.client.Database)
		}
		t.Test.relations = t.relations.RelationMap["test"].RelationModel.GetRelationList()
		t.Test.SetResult(t.result.relationsMap["test"])
		t.Test.ScanResult()
	}
	if _, ok := t.relations.RelationMap["account"]; ok {
		if t.AccountList == nil {
			t.AccountList = NewRelationAccountList(t.ctx, t.client.Database)
		}
		t.AccountList.relations = t.relations.RelationMap["account"].RelationModel.GetRelationList()
		t.AccountList.SetResult(t.result.relationsMap["account"])
		t.AccountList.ScanResult()
	}
}

func (t *Deneme) CheckPrimaryKey(v uuid.UUID) bool {
	return t.id == v
}

func (t *DenemeList) ScanResult() {
	var v *Deneme
	if len(t.Items) == 0 {
		v = NewRelationDeneme(t.ctx, t.client.Database)
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
		v = NewRelationDeneme(t.ctx, t.client.Database)
		t.Items = append(t.Items, v)
	}
	v.result = t.result
	v.relations = t.relations
	v.ScanResult()
}

func (t *Deneme) GetContext() context.Context {
	return t.ctx
}

func (t *Deneme) Get() error {
	return databaseDenemeOperationHook(
		client.NewOperationInfo(
			DenemeTableName,
			client.OperationTypeGet,
		),
		t,
		func() error {
			return t.client.Get(t.ctx, t.where, t, &t.result)
		},
	)
}

func (t *Deneme) Refresh() error {
	return databaseDenemeOperationHook(
		client.NewOperationInfo(
			DenemeTableName,
			client.OperationTypeRefresh,
		),
		t,
		func() error {
			return t.client.Refresh(t.ctx, t, &t.result, DenemeIDField, t.id)
		},
	)
}

func (t *Deneme) Create() error {
	return databaseDenemeOperationHook(
		client.NewOperationInfo(
			DenemeTableName,
			client.OperationTypeCreate,
		),
		t,
		func() error {
			return t.client.Create(t.ctx, DenemeTableName, t.changedFields, t.changedFieldsList, t.serialFields)
		},
	)
}

func (t *Deneme) Update() error {
	return databaseDenemeOperationHook(
		client.NewOperationInfo(
			DenemeTableName,
			client.OperationTypeUpdate,
		),
		t,
		func() error {
			return t.client.Update(t.ctx, DenemeTableName, t.changedFields, t.changedFieldsList, DenemeIDField, t.id)
		},
	)
}

func (t *Deneme) Delete() error {
	return databaseDenemeOperationHook(
		client.NewOperationInfo(
			DenemeTableName,
			client.OperationTypeDelete,
		),
		t,
		func() error {
			return t.client.Delete(t.ctx, DenemeTableName, DenemeIDField, t.id)
		},
	)
}

func (t *DenemeList) List() error {
	return databaseDenemeListOperationHook(
		client.NewOperationInfo(
			DenemeTableName,
			client.OperationTypeList,
		),
		t,
		func() error {
			return t.client.List(t.ctx, t.where, t, &t.result, t.order, t.paging)
		},
	)
}

func (t *DenemeList) Aggregate(f func(aggregate *client.Aggregate)) (func() error, error) {
	a := new(client.Aggregate)
	f(a)
	return t.client.Aggregate(t.ctx, t.where, t, a)
}

func (t *DenemeList) Create(list ...*Deneme) error {
	return databaseDenemeListOperationHook(
		client.NewOperationInfo(
			DenemeTableName,
			client.OperationTypeBulkCreate,
		),
		t,
		func() error {
			var changedFieldsList []map[string]any
			var changedFieldsListList [][]string
			for _, item := range list {
				changedFieldsList = append(changedFieldsList, item.changedFields)
				changedFieldsListList = append(changedFieldsListList, item.changedFieldsList)
			}
			return t.client.BulkCreate(t.ctx, DenemeTableName, changedFieldsList, changedFieldsListList)
		},
	)
}

func (t *DenemeList) Update(list ...*Deneme) error {
	return databaseDenemeListOperationHook(
		client.NewOperationInfo(
			DenemeTableName,
			client.OperationTypeBulkUpdate,
		),
		t,
		func() error {
			var valueList []any
			for _, item := range list {
				valueList = append(valueList, item.id)
			}
			return t.client.BulkUpdate(t.ctx, DenemeTableName, list[0].changedFields, list[0].changedFieldsList, DenemeIDField, valueList)
		},
	)
}

func (t *DenemeList) Delete(list ...*Deneme) error {
	return databaseDenemeListOperationHook(
		client.NewOperationInfo(
			DenemeTableName,
			client.OperationTypeBulkDelete,
		),
		t,
		func() error {
			var valueList []any
			for _, item := range list {
				valueList = append(valueList, item.id)
			}
			return t.client.BulkDelete(t.ctx, DenemeTableName, DenemeIDField, valueList)
		},
	)
}

func (t *DenemeList) Order(field string) *DenemeList {
	t.order = append(t.order, &client.Order{Field: field})
	return t
}

func (t *DenemeList) OrderDesc(field string) *DenemeList {
	t.order = append(t.order, &client.Order{Field: field, Desc: true})
	return t
}

func (t *DenemeList) Paging(skip, limit int) *DenemeList {
	t.paging = &client.Paging{Skip: skip, Limit: limit}
	return t
}

type DenemeResult struct {
	id         uuid.UUID
	testid     *uuid.UUID
	count      int
	isactive   bool
	denemetype DenemeType

	selectedFields []*client.SelectedField

	Test    *TestResult
	Account *AccountResult

	relations    []client.Result
	relationsMap map[string]client.Result
}

func (t *DenemeResult) Init() {
	t.selectedFields = []*client.SelectedField{}
	t.relationsMap = make(map[string]client.Result)
	t.prepare()
	t.SelectAll()
}

func (t *DenemeResult) GetSelectedFields() []*client.SelectedField {
	return t.selectedFields
}

func (t *DenemeResult) GetRelations() []client.Result {
	return t.relations
}

func (t *DenemeResult) prepare() {
	t.testid = &uuid.UUID{}

}

func (t *DenemeResult) SelectID() {
	v := &client.SelectedField{Name: DenemeIDField, Value: &t.id}
	t.selectedFields = append(t.selectedFields, v)
}

func (t *DenemeResult) SelectTestID() {

	v := &client.SelectedField{Name: DenemeTestIDField, Value: t.testid}
	t.selectedFields = append(t.selectedFields, v)
}

func (t *DenemeResult) SelectCount() {
	v := &client.SelectedField{Name: DenemeCountField, Value: &t.count}
	t.selectedFields = append(t.selectedFields, v)
}

func (t *DenemeResult) SelectIsActive() {
	v := &client.SelectedField{Name: DenemeIsActiveField, Value: &t.isactive}
	t.selectedFields = append(t.selectedFields, v)
}

func (t *DenemeResult) SelectDenemeType() {
	v := &client.SelectedField{Name: DenemeDenemeTypeField, Value: &t.denemetype}
	t.selectedFields = append(t.selectedFields, v)
}

func (t *DenemeResult) GetDBName() string {
	return DenemeTableName
}

func (t *DenemeResult) SelectAll() {
	t.SelectID()
	t.SelectTestID()
	t.SelectCount()
	t.SelectIsActive()
	t.SelectDenemeType()

}

func (t *DenemeResult) IsExist() bool {
	if t == nil {
		return false
	}
	var v uuid.UUID
	return t.id != v
}
