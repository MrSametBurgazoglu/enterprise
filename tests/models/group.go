package models

import (
	"context"
	"github.com/MrSametBurgazoglu/enterprise/client"

	"github.com/google/uuid"
)

const GroupTableName = "group"

const (
	GroupIDField      string = "id"
	GroupNameField    string = "name"
	GroupSurnameField string = "surname"
)

var databaseGroupOperationHook = func(operationInfo *client.OperationInfo, model *Group, operationFunc func() error) error {
	return operationFunc()
}

var databaseGroupListOperationHook = func(operationInfo *client.OperationInfo, model *GroupList, operationFunc func() error) error {
	return operationFunc()
}

func SetDatabaseGroupOperationHook(f func(operationInfo *client.OperationInfo, model *Group, operationFunc func() error) error) {
	databaseGroupOperationHook = f
}

func SetDatabaseGroupListOperationHook(f func(operationInfo *client.OperationInfo, model *GroupList, operationFunc func() error) error) {
	databaseGroupListOperationHook = f
}

func NewGroup(ctx context.Context, dc client.DatabaseClient) *Group {
	v := &Group{client: client.NewClient(dc), ctx: ctx}
	v.relations = new(client.RelationList)
	v.relations.RelationMap = make(map[string]*client.Relation)
	v.changedFields = make(map[string]any)
	v.result.Init()
	return v
}

func NewRelationGroup(ctx context.Context, dc client.DatabaseClient) *Group {
	v := &Group{client: client.NewClient(dc), ctx: ctx}
	v.relations = new(client.RelationList)
	v.relations.RelationMap = make(map[string]*client.Relation)
	v.changedFields = make(map[string]any)
	v.result.Init()
	return v
}

type Group struct {
	id uuid.UUID

	name string

	surname string

	changedFields     map[string]any
	changedFieldsList []string
	serialFields      []*client.SelectedField

	ctx    context.Context
	client *client.Client
	GroupPredicate
	relations *client.RelationList

	AccountList *AccountList

	result GroupResult
}

func (t *Group) GetDBName() string {
	return GroupTableName
}

func (t *Group) GetSelector() *GroupResult {
	t.result.selectedFields = nil
	return &t.result
}

func (t *Group) GetRelationList() *client.RelationList {
	return t.relations
}

func (t *Group) IsExist() bool {
	var v uuid.UUID
	return t.id != v
}

func (t *Group) GetPrimaryKey() uuid.UUID {
	return t.id
}

func NewGroupList(ctx context.Context, dc client.DatabaseClient) *GroupList {
	v := &GroupList{client: client.NewClient(dc), ctx: ctx}
	v.relations = new(client.RelationList)
	v.relations.RelationMap = make(map[string]*client.Relation)
	v.result.Init()
	return v
}

func NewRelationGroupList(ctx context.Context, dc client.DatabaseClient) *GroupList {
	v := &GroupList{client: client.NewClient(dc), ctx: ctx}
	v.relations = new(client.RelationList)
	v.relations.RelationMap = make(map[string]*client.Relation)
	v.result.Init()
	return v
}

type GroupList struct {
	Items []*Group

	ctx    context.Context
	client *client.Client
	GroupPredicate
	order     []*client.Order
	paging    *client.Paging
	relations *client.RelationList
	result    GroupResult
}

func (t *GroupList) GetDBName() string {
	return GroupTableName
}

func (t *GroupList) GetRelationList() *client.RelationList {
	return t.relations
}

func (t *GroupList) IsExist() bool {
	return t.Items[len(t.Items)-1].IsExist()
}

func (t *Group) SetID(v uuid.UUID) {
	t.id = v
	t.SetIDField()
}
func (t *Group) SetName(v string) {
	t.name = v
	t.SetNameField()
}
func (t *Group) SetSurname(v string) {
	t.surname = v
	t.SetSurnameField()
}

func (t *Group) SetIDNillable(v *uuid.UUID) {
	if v == nil {
		return
	}
	t.SetID(*v)
}
func (t *Group) SetNameNillable(v *string) {
	if v == nil {
		return
	}
	t.SetName(*v)
}
func (t *Group) SetSurnameNillable(v *string) {
	if v == nil {
		return
	}
	t.SetSurname(*v)
}

func (t *Group) ParseID(v string) error {
	parsedID, err := uuid.Parse(v)
	if err != nil {
		return err
	}
	t.id = parsedID
	t.SetIDField()
	return nil
}

func (t *Group) IDIN(v ...uuid.UUID) bool {
	for _, x := range v {
		if t.id == x {
			return true
		}
	}
	return false
}

func (t *Group) NameIN(v ...string) bool {
	for _, x := range v {
		if t.name == x {
			return true
		}
	}
	return false
}

func (t *Group) SurnameIN(v ...string) bool {
	for _, x := range v {
		if t.surname == x {
			return true
		}
	}
	return false
}

func (t *Group) IDNotIN(v ...uuid.UUID) bool {
	for _, x := range v {
		if t.id == x {
			return false
		}
	}
	return true
}

func (t *Group) NameNotIN(v ...string) bool {
	for _, x := range v {
		if t.name == x {
			return false
		}
	}
	return true
}

func (t *Group) SurnameNotIN(v ...string) bool {
	for _, x := range v {
		if t.surname == x {
			return false
		}
	}
	return true
}

func (t *Group) GetID() uuid.UUID {
	return t.id
}
func (t *Group) GetName() string {
	return t.name
}
func (t *Group) GetSurname() string {
	return t.surname
}

func (t *Group) SetIDField() {
	if _, exist := t.changedFields[GroupIDField]; !exist {
		t.changedFields[GroupIDField] = t.id
		t.changedFieldsList = append(t.changedFieldsList, GroupIDField)
	}

}
func (t *Group) SetNameField() {
	if _, exist := t.changedFields[GroupNameField]; !exist {
		t.changedFields[GroupNameField] = t.name
		t.changedFieldsList = append(t.changedFieldsList, GroupNameField)
	}

}
func (t *Group) SetSurnameField() {
	if _, exist := t.changedFields[GroupSurnameField]; !exist {
		t.changedFields[GroupSurnameField] = t.surname
		t.changedFieldsList = append(t.changedFieldsList, GroupSurnameField)
	}

}

func (t *Group) WithAccountList(opts ...func(*AccountList)) *client.Relation {
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
	r := &client.Relation{
		RelationModel:   t.AccountList,
		RelationTable:   "account",
		RelationResult:  t.result.Account,
		Where:           t.AccountList.where,
		ManyToManyTable: "account_group",
		RelationWhere: &client.RelationCondition{
			RelationValue:      "group_id",
			TableValue:         "account_id",
			RelationTableValue: "id",
		},
	}
	t.relations.Relations = append(t.relations.Relations, r)
	t.relations.RelationMap["account"] = r
	return r
}

func (t *GroupList) WithAccountList(opts ...func(*AccountList)) *client.Relation {
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
	r := &client.Relation{
		RelationModel:  v,
		RelationTable:  "account",
		RelationResult: t.result.Account,
		Where:          v.where,
		RelationWhere: &client.RelationCondition{
			RelationValue: "group_id",
			TableValue:    "account_id",
		},
	}
	t.relations.Relations = append(t.relations.Relations, r)
	t.relations.RelationMap["account"] = r
	return r
}

func (t *GroupList) cleanAccountList() {
	Relation := t.Items[len(t.Items)-1].relations
	p := 0
	for i, v := range Relation.Relations {
		if v.RelationTable == "account" {
			p = i
		}
	}
	Relation.Relations = append(Relation.Relations[:p], Relation.Relations[p+1:]...)
}

func (t *Group) SetDefaults() {
	t.id = uuid.New()
	t.changedFields[GroupIDField] = t.id
	t.changedFieldsList = append(t.changedFieldsList, GroupIDField)

}

func (t *Group) SetResult(result client.Result) {
	if t == nil {
		v := NewRelationGroup(t.ctx, t.client.Database)
		*t = *v
	}
	t.result = *result.(*GroupResult)
}

func (t *GroupList) SetResult(result client.Result) {
	if t == nil {
		v := NewRelationGroupList(t.ctx, t.client.Database)
		*t = *v
	}
	t.result = *result.(*GroupResult)
}

func (t *Group) ScanResult() {
	t.id = t.result.id
	t.name = t.result.name
	t.surname = t.result.surname

	if _, ok := t.relations.RelationMap["account"]; ok {
		if t.AccountList == nil {
			t.AccountList = NewRelationAccountList(t.ctx, t.client.Database)
		}
		t.AccountList.relations = t.relations.RelationMap["account"].RelationModel.GetRelationList()
		t.AccountList.SetResult(t.result.relationsMap["account"])
		t.AccountList.ScanResult()
	}
}

func (t *Group) CheckPrimaryKey(v uuid.UUID) bool {
	return t.id == v
}

func (t *GroupList) ScanResult() {
	var v *Group
	if len(t.Items) == 0 {
		v = NewRelationGroup(t.ctx, t.client.Database)
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
		v = NewRelationGroup(t.ctx, t.client.Database)
		t.Items = append(t.Items, v)
	}
	v.result = t.result
	v.relations = t.relations
	v.ScanResult()
}

func (t *Group) GetContext() context.Context {
	return t.ctx
}

func (t *Group) Get() error {
	return databaseGroupOperationHook(
		client.NewOperationInfo(
			GroupTableName,
			client.OperationTypeGet,
		),
		t,
		func() error {
			return t.client.Get(t.ctx, t.where, t, &t.result)
		},
	)
}

func (t *Group) Refresh() error {
	return databaseGroupOperationHook(
		client.NewOperationInfo(
			GroupTableName,
			client.OperationTypeRefresh,
		),
		t,
		func() error {
			return t.client.Refresh(t.ctx, t, &t.result, GroupIDField, t.id)
		},
	)
}

func (t *Group) Create() error {
	return databaseGroupOperationHook(
		client.NewOperationInfo(
			GroupTableName,
			client.OperationTypeCreate,
		),
		t,
		func() error {
			return t.client.Create(t.ctx, GroupTableName, t.changedFields, t.changedFieldsList, t.serialFields)
		},
	)
}

func (t *Group) Update() error {
	return databaseGroupOperationHook(
		client.NewOperationInfo(
			GroupTableName,
			client.OperationTypeUpdate,
		),
		t,
		func() error {
			return t.client.Update(t.ctx, GroupTableName, t.changedFields, t.changedFieldsList, GroupIDField, t.id)
		},
	)
}

func (t *Group) Delete() error {
	return databaseGroupOperationHook(
		client.NewOperationInfo(
			GroupTableName,
			client.OperationTypeDelete,
		),
		t,
		func() error {
			return t.client.Delete(t.ctx, GroupTableName, GroupIDField, t.id)
		},
	)
}

func (t *GroupList) List() error {
	return databaseGroupListOperationHook(
		client.NewOperationInfo(
			GroupTableName,
			client.OperationTypeList,
		),
		t,
		func() error {
			return t.client.List(t.ctx, t.where, t, &t.result, t.order, t.paging)
		},
	)
}

func (t *GroupList) Aggregate(f func(aggregate *client.Aggregate)) (func() error, error) {
	a := new(client.Aggregate)
	f(a)
	return t.client.Aggregate(t.ctx, t.where, t, a)
}

func (t *GroupList) Create(list ...*Group) error {
	return databaseGroupListOperationHook(
		client.NewOperationInfo(
			GroupTableName,
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
			return t.client.BulkCreate(t.ctx, GroupTableName, changedFieldsList, changedFieldsListList)
		},
	)
}

func (t *GroupList) Update(list ...*Group) error {
	return databaseGroupListOperationHook(
		client.NewOperationInfo(
			GroupTableName,
			client.OperationTypeBulkUpdate,
		),
		t,
		func() error {
			var valueList []any
			for _, item := range list {
				valueList = append(valueList, item.id)
			}
			return t.client.BulkUpdate(t.ctx, GroupTableName, list[0].changedFields, list[0].changedFieldsList, GroupIDField, valueList)
		},
	)
}

func (t *GroupList) Delete(list ...*Group) error {
	return databaseGroupListOperationHook(
		client.NewOperationInfo(
			GroupTableName,
			client.OperationTypeBulkDelete,
		),
		t,
		func() error {
			var valueList []any
			for _, item := range list {
				valueList = append(valueList, item.id)
			}
			return t.client.BulkDelete(t.ctx, GroupTableName, GroupIDField, valueList)
		},
	)
}

func (t *GroupList) Order(field string) *GroupList {
	t.order = append(t.order, &client.Order{Field: field})
	return t
}

func (t *GroupList) OrderDesc(field string) *GroupList {
	t.order = append(t.order, &client.Order{Field: field, Desc: true})
	return t
}

func (t *GroupList) Paging(skip, limit int) *GroupList {
	t.paging = &client.Paging{Skip: skip, Limit: limit}
	return t
}

type GroupResult struct {
	id      uuid.UUID
	name    string
	surname string

	selectedFields []*client.SelectedField

	Account *AccountResult

	relations    []client.Result
	relationsMap map[string]client.Result
}

func (t *GroupResult) Init() {
	t.selectedFields = []*client.SelectedField{}
	t.relationsMap = make(map[string]client.Result)
	t.prepare()
	t.SelectAll()
}

func (t *GroupResult) GetSelectedFields() []*client.SelectedField {
	return t.selectedFields
}

func (t *GroupResult) GetRelations() []client.Result {
	return t.relations
}

func (t *GroupResult) prepare() {

}

func (t *GroupResult) SelectID() {
	v := &client.SelectedField{Name: GroupIDField, Value: &t.id}
	t.selectedFields = append(t.selectedFields, v)
}

func (t *GroupResult) SelectName() {
	v := &client.SelectedField{Name: GroupNameField, Value: &t.name}
	t.selectedFields = append(t.selectedFields, v)
}

func (t *GroupResult) SelectSurname() {
	v := &client.SelectedField{Name: GroupSurnameField, Value: &t.surname}
	t.selectedFields = append(t.selectedFields, v)
}

func (t *GroupResult) GetDBName() string {
	return GroupTableName
}

func (t *GroupResult) SelectAll() {
	t.SelectID()
	t.SelectName()
	t.SelectSurname()

}

func (t *GroupResult) IsExist() bool {
	if t == nil {
		return false
	}
	var v uuid.UUID
	return t.id != v
}

func (t *Group) AddIntoAccount(relationship *Account) error {
	return t.client.AddManyToManyRelation(t.ctx, "account_group", "group_id", "account_id", t.id, relationship.GetPrimaryKey())
}

func (t *Group) RemoveFromAccount(relationship *Account) error {
	return t.client.DeleteManyToManyRelation(t.ctx, "account_group", "group_id", "account_id", t.id, relationship.GetPrimaryKey())
}

func (t *Group) IsInAccount(relationship *Account) (bool, error) {
	return t.client.ExistManyToManyRelation(t.ctx, "account_group", "group_id", "account_id", t.id, relationship.GetPrimaryKey())
}
