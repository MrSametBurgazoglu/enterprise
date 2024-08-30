package models

import (
"context"
"github.com/MrSametBurgazoglu/enterprise/client"
{{range .RequiredPackages}}
"{{.}}"{{end}}

)

const {{.TableName}}TableName = "{{.DBName}}"

const (
    {{range .Fields}}{{$.TableName}}{{.GetName}}Field string = "{{.GetDBName}}"
    {{end}}
)

var database{{.TableName}}OperationHook = func(operationInfo *client.OperationInfo, model *{{.TableName}}, operationFunc func() error) error {
    return operationFunc()
}

var database{{.TableName}}ListOperationHook = func(operationInfo *client.OperationInfo, model *{{.TableName}}List, operationFunc func() error) error {
    return operationFunc()
}

func SetDatabase{{.TableName}}OperationHook(f func(operationInfo *client.OperationInfo, model *{{.TableName}}, operationFunc func() error) error){
    database{{.TableName}}OperationHook = f
}

func SetDatabase{{.TableName}}ListOperationHook(f func(operationInfo *client.OperationInfo, model *{{.TableName}}List, operationFunc func() error) error){
    database{{.TableName}}ListOperationHook = f
}


{{range .Fields}}{{if .IsCustomType}}
type {{.GetName}} string
{{ $name := .GetName }}
const ({{range .Values}}
    {{$name}}{{.}}Type {{$name}} = "{{.}}"{{end}}
)
{{end}}{{end}}

func New{{.TableName}}(ctx context.Context, dc client.DatabaseClient) *{{.TableName}}{
    v := &{{.TableName}}{client: client.NewClient(dc), ctx: ctx}
    v.relations = new(client.RelationList)
    v.relations.RelationMap = make(map[string]*client.Relation)
    v.changedFields = make(map[string]any)
    v.result.Init()
    v.Default()
    return v
}

func NewRelation{{.TableName}}(ctx context.Context, dc client.DatabaseClient) *{{.TableName}}{
    v := &{{.TableName}}{client: client.NewClient(dc), ctx: ctx}
    v.relations = new(client.RelationList)
    v.relations.RelationMap = make(map[string]*client.Relation)
    v.changedFields = make(map[string]any)
    v.result.Init()
    return v
}

type {{.TableName}} struct {
    {{range .Fields}}
    {{.GetNameLower}} {{.GetType}}
    {{end}}

    changedFields map[string]any
    changedFieldsList []string
    serialFields []*client.SelectedField

    ctx context.Context
    client *client.Client
    {{.TableName}}Predicate
    relations *client.RelationList
    {{range .Relations}}
    {{.GetRelationField}} *{{.GetRelationField}}{{end}}

    result {{.TableName}}Result
}

func (t *{{$.TableName}}) GetDBName() string{
    return {{.TableName}}TableName
}

func (t *{{$.TableName}}) GetSelector() *{{$.TableName}}Result {
    t.result.selectedFields = nil
    return &t.result
}

func (t *{{$.TableName}}) GetRelationList() *client.RelationList{
    return t.relations
}

func (t *{{$.TableName}}) IsExist() bool{
    var v {{.IDFieldType}}
    return t.{{.IDFieldLower}} != v
}

func (t *{{$.TableName}}) GetPrimaryKey() {{.IDFieldType}}{
    return t.{{.IDFieldLower}}
}



func New{{.TableName}}List(ctx context.Context, dc client.DatabaseClient) *{{.TableName}}List{
    v := &{{.TableName}}List{client: client.NewClient(dc), ctx: ctx}
    v.relations = new(client.RelationList)
    v.relations.RelationMap = make(map[string]*client.Relation)
    v.result.Init()
    return v
}

func NewRelation{{.TableName}}List(ctx context.Context, dc client.DatabaseClient) *{{.TableName}}List{
    v := &{{.TableName}}List{client: client.NewClient(dc), ctx: ctx}
    v.relations = new(client.RelationList)
    v.relations.RelationMap = make(map[string]*client.Relation)
    v.result.Init()
    return v
}

type {{.TableName}}List struct {
    Items []*{{.TableName}}

    ctx context.Context
    client *client.Client
    {{.TableName}}Predicate
    order []*client.Order
    paging *client.Paging
    relations *client.RelationList
    result {{.TableName}}Result
}

func (t *{{$.TableName}}List) GetDBName() string{
    return {{.TableName}}TableName
}

func (t *{{$.TableName}}List) GetRelationList() *client.RelationList{
    return t.relations
}

func (t *{{$.TableName}}List) IsExist() bool{
    return t.Items[len(t.Items)-1].IsExist()
}


{{range .Fields}}
func (t *{{$.TableName}}) Set{{.GetName}}(v {{.GetType}}){
    t.{{.GetNameLower}} = v
    t.Set{{.GetName}}Field()
}{{end}}
{{range .Fields}}
{{if .IsNillable}}
func (t *{{$.TableName}}) Set{{.GetName}}Value(v {{.GetBaseType}}){
   t.Set{{.GetName}}(&v)
   t.Set{{.GetName}}Field()
}{{end}}
{{end}}

{{range .Fields}}
{{- if eq .IsNillable false -}}
func (t *{{$.TableName}}) Set{{.GetName}}Nillable(v *{{.GetBaseType}}){
    if v == nil{
        return
    }
    t.Set{{.GetName}}(*v)
}{{end}}
{{end}}

{{range .Fields}}
{{if .CanTime}}
func (t *{{$.TableName}}) Format{{.GetName}}(v string) string{
    {{if .IsNillable}}
    if t.{{.GetNameLower}} == nil{
        return ""
    }
    {{end}}
   return t.{{.GetNameLower}}.Format(v)
}

func (t *{{$.TableName}}) Parse{{.GetName}}(layout, value string) error{
    parsedTime, err := time.Parse(layout, value)
    if err != nil {
        return err
    }
    {{if .IsNillable}}t.{{.GetNameLower}} = &parsedTime{{else}}t.{{.GetNameLower}} = parsedTime{{end}}
    t.Set{{.GetName}}Field()
    return nil
}
{{end}}

{{if .CanUUID}}
func (t *{{$.TableName}}) Parse{{.GetName}}(v string) error{
    parsedID, err := uuid.Parse(v)
    if err != nil {
        return err
    }
    {{if .IsNillable}}t.{{.GetNameLower}} = &parsedID{{else}}t.{{.GetNameLower}} = parsedID{{end}}
    t.Set{{.GetName}}Field()
    return nil
}
{{end}}
{{end}}

{{range .Fields}}
func (t *{{$.TableName}}) {{.GetName}}IN(v ...{{.GetBaseType}}) bool{ {{if .IsNillable}}
    if t.{{.GetNameLower}} == nil{
        return false
    }{{end}}
    for _, x := range v{
        {{if .IsNillable}}if *t.{{.GetNameLower}} == x{{else}}if t.{{.GetNameLower}} == x{{end}}{
            return true
        }
    }
    return false
}
{{end}}

{{range .Fields}}
func (t *{{$.TableName}}) {{.GetName}}NotIN(v ...{{.GetBaseType}}) bool{ {{if .IsNillable}}
    if t.{{.GetNameLower}} == nil{
        return true
    }{{end}}
    for _, x := range v{
        {{if .IsNillable}}if *t.{{.GetNameLower}} == x{{else}}if t.{{.GetNameLower}} == x{{end}}{
            return false
        }
    }
    return true
}
{{end}}

{{range .Fields}}
func (t *{{$.TableName}}) Get{{.GetName}}() {{.GetType}}{
    return t.{{.GetNameLower}}
}{{end}}

{{range .Fields}}
func (t *{{$.TableName}}) Set{{.GetName}}Field(){
    if _, exist := t.changedFields[{{$.TableName}}{{.GetName}}Field]; !exist{
        t.changedFields[{{$.TableName}}{{.GetName}}Field] = t.{{.GetNameLower}}
        t.changedFieldsList = append(t.changedFieldsList, {{$.TableName}}{{.GetName}}Field)
    }

}{{end}}

{{range .Relations}}
func (t *{{$.TableName}}) With{{.GetRelationField}}(opts ...func(*{{.GetRelationField}})){
    t.{{.GetRelationField}} = NewRelation{{.GetRelationField}}(t.ctx, t.client.Database)
    for _, opt := range opts {
        opt(t.{{.GetRelationField}})
    }
    t.result.{{.RelationTable}} = new({{.RelationTable}}Result)
    t.result.{{.RelationTable}}.Init()
    t.result.relations = append(t.result.relations, t.result.{{.RelationTable}})
    t.result.relationsMap["{{.GetRelationTableLower}}"] = t.result.{{.RelationTable}}
    for _, Relation := range t.{{.GetRelationField}}.relations.Relations {
        t.result.{{.RelationTable}}.relations = append(t.result.{{.RelationTable}}.relations, Relation.RelationResult)
        t.result.{{.RelationTable}}.relationsMap[Relation.RelationTable] = Relation.RelationResult
    }
    t.relations.Relations = append(t.relations.Relations,
    	&client.Relation{
    	    RelationModel: t.{{.GetRelationField}},
    	    RelationTable: "{{.RelationTableDBName}}",
    	    RelationResult: t.result.{{.RelationTable}},
    	    Where: t.{{.GetRelationField}}.where,
    	    {{if .IsManyToMany}}ManyToManyTable: "{{.ManyTableDBName}}",{{end}}
    	    RelationWhere: &client.RelationCondition{
    	        RelationValue: "{{.RelationField}}",
    	        TableValue: "{{.OnField}}",
    	        {{if .IsManyToMany}}RelationTableValue: "{{.RelationTableField}}",{{end}}
    	    },
    	},
    )
    t.relations.RelationMap["{{.GetRelationTableLower}}"] = t.relations.Relations[len(t.relations.Relations)-1]
}{{end}}

{{range .Relations}}
func (t *{{$.TableName}}List) With{{.GetRelationField}}(opts ...func(*{{.GetRelationField}})){
    v := NewRelation{{.GetRelationField}}(t.ctx, t.client.Database)
    for _, opt := range opts {
        opt(v)
    }
    t.result.{{.RelationTable}} = new({{.RelationTable}}Result)
    t.result.{{.RelationTable}}.Init()
    t.result.relations = append(t.result.relations, t.result.{{.RelationTable}})
    t.result.relationsMap["{{.GetRelationTableLower}}"] = t.result.{{.RelationTable}}
    for _, Relation := range v.relations.Relations {
    	t.result.{{.RelationTable}}.relations = append(t.result.{{.RelationTable}}.relations, Relation.RelationResult)
    	t.result.{{.RelationTable}}.relationsMap[Relation.RelationTable] = Relation.RelationResult
    }
    t.relations.Relations = append(t.relations.Relations,
    	&client.Relation{
    	    RelationModel: v,
    	    RelationTable: "{{.RelationTableDBName}}",
    	    RelationResult: t.result.{{.RelationTable}},
    	    Where: v.where,
    	    RelationWhere: &client.RelationCondition{
    	        RelationValue: "{{.RelationField}}",
    	        TableValue: "{{.OnField}}",
    	    },
    	},
    )
    t.relations.RelationMap["{{.GetRelationTableLower}}"] = t.relations.Relations[len(t.relations.Relations)-1]
}

func (t *{{$.TableName}}List) clean{{.GetRelationField}}(){
    Relation := t.Items[len(t.Items)-1].relations
    p := 0
    for i, v := range Relation.Relations {
        if v.RelationTable == "{{.RelationTableDBName}}" {
            p = i
        }
    }
    Relation.Relations = append(Relation.Relations[:p], Relation.Relations[p+1:]...)
}{{end}}

func (t *{{$.TableName}}) Default(){
    {{range .Fields}}{{if .IsDefault}}t.{{.GetNameLower}} = {{.GetDefault}}
    t.changedFields[{{$.TableName}}{{.GetName}}Field] = t.{{.GetNameLower}}
    t.changedFieldsList = append(t.changedFieldsList, {{$.TableName}}{{.GetName}}Field){{end}}
    {{end}}

    {{range .Fields}}{{if .IsSerial}}{{if .IsNillable}}
    v := &client.SelectedField{Name:{{$.TableName}}{{.GetName}}Field, Value:t.{{.GetNameLower}}}{{else}}v := &client.SelectedField{Name:{{$.TableName}}{{.GetName}}Field, Value:&t.{{.GetNameLower}}}{{end}}
    t.serialFields = append(t.serialFields, v)
    {{end}}
    {{end}}
}

func (t *{{$.TableName}}) SetResult(result client.Result){
    if t == nil {
        v := NewRelation{{.TableName}}(t.ctx, t.client.Database)
        *t = *v
    }
    t.result = *result.(*{{$.TableName}}Result)
}

func (t *{{$.TableName}}List) SetResult(result client.Result){
    if t == nil {
        v := NewRelation{{.TableName}}List(t.ctx, t.client.Database)
        *t = *v
    }
    t.result = *result.(*{{$.TableName}}Result)
}

func (t *{{$.TableName}}) ScanResult(){
    {{range .Fields}}t.{{.GetNameLower}} = t.result.{{.GetNameLower}}
    {{end}}
    {{range .Relations}}
    if _, ok := t.relations.RelationMap["{{.GetRelationTableLower}}"]; ok {
        if t.{{.GetRelationField}} == nil{
            t.{{.GetRelationField}} = NewRelation{{.GetRelationField}}(t.ctx, t.client.Database)
        }
        t.{{.GetRelationField}}.relations = t.relations.RelationMap["{{.GetRelationTableLower}}"].RelationModel.GetRelationList()
        t.{{.GetRelationField}}.SetResult(t.result.relationsMap["{{.GetRelationTableLower}}"])
        t.{{.GetRelationField}}.ScanResult()
    }{{end}}
}

func (t *{{.TableName}}) CheckPrimaryKey(v {{.IDFieldType}}) bool{
    return t.{{.IDFieldLower}} == v
}

func (t *{{$.TableName}}List) ScanResult(){
    var v *{{$.TableName}}
    if len(t.Items) == 0 {
        v = NewRelation{{.TableName}}(t.ctx, t.client.Database)
        t.Items = append(t.Items, v)
    } else {
        for _, item := range t.Items {
            if item.CheckPrimaryKey(t.result.{{.IDFieldLower}}) {
                v = item
                break
            }
        }
    }
    if v == nil{
        v = NewRelation{{.TableName}}(t.ctx, t.client.Database)
        t.Items = append(t.Items, v)
    }
    v.result = t.result
    v.relations = t.relations
    v.ScanResult()
}

func (t *{{.TableName}}) GetContext() context.Context{
    return t.ctx
}

func (t *{{.TableName}}) Get() error{
    return database{{.TableName}}OperationHook(
        client.NewOperationInfo(
            {{.TableName}}TableName,
            client.OperationTypeGet,
        ),
        t,
        func() error {
            return t.client.Get(t.ctx, t.where, t, &t.result)
        },
    )
}

func (t *{{.TableName}}) Refresh() error{
    return database{{.TableName}}OperationHook(
        client.NewOperationInfo(
            {{.TableName}}TableName,
            client.OperationTypeRefresh,
        ),
        t,
        func() error {
            return t.client.Refresh(t.ctx, t, &t.result, {{.TableName}}{{.IDField}}Field, t.{{.IDFieldLower}})
        },
    )
}

func (t *{{.TableName}}) Create() error{
    return database{{.TableName}}OperationHook(
        client.NewOperationInfo(
            {{.TableName}}TableName,
            client.OperationTypeCreate,
        ),
        t,
        func() error {
            return t.client.Create(t.ctx, {{.TableName}}TableName, t.changedFields, t.changedFieldsList, t.serialFields)
        },
    )
}

func (t *{{.TableName}}) Update() error{
    return database{{.TableName}}OperationHook(
        client.NewOperationInfo(
            {{.TableName}}TableName,
            client.OperationTypeUpdate,
        ),
        t,
        func() error {
            return t.client.Update(t.ctx, {{.TableName}}TableName, t.changedFields, t.changedFieldsList, {{.TableName}}{{.IDField}}Field, t.{{.IDFieldLower}})
        },
    )
}

func (t *{{.TableName}}) Delete() error{
    return database{{.TableName}}OperationHook(
        client.NewOperationInfo(
            {{.TableName}}TableName,
            client.OperationTypeDelete,
        ),
        t,
        func() error {
            return t.client.Delete(t.ctx, {{.TableName}}TableName, {{.TableName}}{{.IDField}}Field, t.{{.IDFieldLower}})
        },
    )
}

func (t *{{.TableName}}List) List() error{
    return database{{.TableName}}ListOperationHook(
        client.NewOperationInfo(
            {{.TableName}}TableName,
            client.OperationTypeList,
        ),
        t,
        func() error {
            return t.client.List(t.ctx, t.where, t, &t.result, t.order, t.paging)
        },
    )
}

func (t *{{.TableName}}List) Aggregate(f func (aggregate *client.Aggregate)) (func() error,error){
    a := new(client.Aggregate)
    f(a)
    return t.client.Aggregate(t.ctx, t.where, t, a)
}

func (t *{{.TableName}}List) Order(field string) *{{.TableName}}List{
    t.order = append(t.order, &client.Order{Field:field})
    return t
}

func (t *{{.TableName}}List) OrderDesc(field string) *{{.TableName}}List{
    t.order = append(t.order, &client.Order{Field:field, Desc:true})
    return t
}

func (t *{{.TableName}}List) Paging(skip, limit int) *{{.TableName}}List{
    t.paging = &client.Paging{Skip:skip, Limit:limit}
    return t
}


type {{.TableName}}Result struct {
    {{range .Fields}}{{.GetNameLower}} {{.GetType}}
    {{end}}
    selectedFields []*client.SelectedField

    {{range .Relations}}{{.RelationTable}} *{{.RelationTable}}Result
    {{end}}
    relations []client.Result
    relationsMap map[string]client.Result
}

func (t *{{$.TableName}}Result) Init(){
    t.selectedFields = []*client.SelectedField{}
    t.relationsMap = make(map[string]client.Result)
    t.prepare()
    t.SelectAll()
}

func (t *{{$.TableName}}Result) GetSelectedFields() []*client.SelectedField{
    return t.selectedFields
}

func (t *{{$.TableName}}Result) GetRelations() []client.Result{
    return t.relations
}

func (t *{{$.TableName}}Result) prepare(){
    {{range .Fields}}{{if .NeedPrepare}}t.{{.GetNameLower}} = {{.PrepareFunc}}
    {{end}}{{end}}
}

{{range .Fields}}
func (t *{{$.TableName}}Result) Select{{.GetName}}(){
    {{if .IsNillable}}
    v := &client.SelectedField{Name:{{$.TableName}}{{.GetName}}Field, Value:t.{{.GetNameLower}}}{{else}}v := &client.SelectedField{Name:{{$.TableName}}{{.GetName}}Field, Value:&t.{{.GetNameLower}}}{{end}}
    t.selectedFields = append(t.selectedFields, v)
}
{{end}}

func (t *{{$.TableName}}Result) GetDBName() string{
    return {{.TableName}}TableName
}

func (t *{{$.TableName}}Result) SelectAll(){
    {{range .Fields}}t.Select{{.GetName}}()
    {{end}}
}

func (t *{{$.TableName}}Result) IsExist() bool{
    if t == nil{
        return false
    }
    var v {{.IDFieldType}}
    return t.{{.IDFieldLower}} != v
}

{{range .Relations}}{{if .IsManyToMany}}
func (t *{{$.TableName}}) AddInto{{.RelationTable}}(relationship *{{.RelationTable}}) error {
    return t.client.AddManyToManyRelation(t.ctx, "{{.ManyTableDBName}}", "{{.RelationField}}", "{{.OnField}}", t.{{$.IDFieldLower}}, relationship.GetPrimaryKey())
}

func (t *{{$.TableName}}) RemoveFrom{{.RelationTable}}(relationship *{{.RelationTable}}) error {
    return t.client.DeleteManyToManyRelation(t.ctx, "{{.ManyTableDBName}}", "{{.RelationField}}", "{{.OnField}}", t.{{$.IDFieldLower}}, relationship.GetPrimaryKey())
}

func (t *{{$.TableName}}) IsIn{{.RelationTable}}(relationship *{{.RelationTable}}) (bool, error){
    return t.client.ExistManyToManyRelation(t.ctx, "{{.ManyTableDBName}}", "{{.RelationField}}", "{{.OnField}}", t.{{$.IDFieldLower}}, relationship.GetPrimaryKey())
}{{end}}
{{end}}