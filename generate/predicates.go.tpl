package models

import "github.com/MrSametBurgazoglu/enterprise/client"
{{range .RequiredPackages}}
import "{{.}}"{{end}}
import "strings"

type {{$.TableName}}Predicate struct{
    where []*client.WhereList
}

func (t *{{.TableName}}Predicate) Where(w ...*client.Where){
    t.where = nil
    wl := &client.WhereList{}
    wl.Items = append(wl.Items, w...)
    t.where = append(t.where, wl)
}

func (t *{{.TableName}}Predicate) ORWhere(w ...*client.Where){
    wl := &client.WhereList{}
    wl.Items = append(wl.Items, w...)
    t.where = append(t.where, wl)
}

{{range .Fields}}
func (t *{{$.TableName}}Predicate) Is{{.GetName}}Equal(v {{.GetBaseType}}) *client.Where{
    return &client.Where{
        Type: client.EQ,
        Name: {{$.TableName}}{{.GetName}}Field,
        HasValue: true,
        Value: v,
    }
}
{{end}}

{{range .Fields}}
func (t *{{$.TableName}}Predicate) Is{{.GetName}}NotEqual(v {{.GetBaseType}}) *client.Where{
    return &client.Where{
        Type: client.NEQ,
        Name: {{$.TableName}}{{.GetName}}Field,
        HasValue: true,
        Value: v,
    }
}
{{end}}

{{range .Fields}}
func (t *{{$.TableName}}Predicate) Is{{.GetName}}IN(v ...{{.GetBaseType}}) *client.Where{
    return &client.Where{
        Type: client.ANY,
        Name: {{$.TableName}}{{.GetName}}Field,
        HasValue: true,
        Value: v,
    }
}
{{end}}

{{range .Fields}}
func (t *{{$.TableName}}Predicate) Is{{.GetName}}NotIN(v ...{{.GetBaseType}}) *client.Where{
    return &client.Where{
        Type: client.NANY,
        Name: {{$.TableName}}{{.GetName}}Field,
        HasValue: true,
        Value: v,
    }
}
{{end}}

{{range .Fields}}
{{if .IsNillable}}
func (t *{{$.TableName}}Predicate) Is{{.GetName}}Nil() *client.Where{
   return &client.Where{
           Type: client.NIL,
           Name: {{$.TableName}}{{.GetName}}Field,
       }
}{{end}}{{end}}

{{range .Fields}}
{{if .IsNillable}}
func (t *{{$.TableName}}Predicate) Is{{.GetName}}NotNil() *client.Where{
   return &client.Where{
           Type: client.NNIL,
           Name: {{$.TableName}}{{.GetName}}Field,
       }
}{{end}}{{end}}

{{range .Fields}}
{{if .CanBeGreater}}
func (t *{{$.TableName}}Predicate) {{.GetName}}GreaterThan(v {{.GetBaseType}}) *client.Where{
   return &client.Where{
           Type: client.GT,
           Name: {{$.TableName}}{{.GetName}}Field,
           HasValue: true,
           Value: v,
       }
}

func (t *{{$.TableName}}Predicate) {{.GetName}}GreaterEqualThan(v {{.GetBaseType}}) *client.Where{
   return &client.Where{
           Type: client.GTE,
           Name: {{$.TableName}}{{.GetName}}Field,
           HasValue: true,
           Value: v,
       }
}

func (t *{{$.TableName}}Predicate) {{.GetName}}LowerThan(v {{.GetBaseType}}) *client.Where{
   return &client.Where{
           Type: client.LT,
           Name: {{$.TableName}}{{.GetName}}Field,
           HasValue: true,
           Value: v,
       }
}

func (t *{{$.TableName}}Predicate) {{.GetName}}LowerEqualThan(v {{.GetBaseType}}) *client.Where{
   return &client.Where{
           Type: client.LTE,
           Name: {{$.TableName}}{{.GetName}}Field,
           HasValue: true,
           Value: v,
       }
}{{end}}{{end}}

func (t *{{$.TableName}}Predicate) GetWhereInfoString() string {
	var whereString []string
	for _, list := range t.where {
        whereAnd := ""
		var whereAndString []string
		for _, item := range list.Items {
			whereAndString = append(whereAndString, item.Name)
		}
        whereAnd = strings.Join(whereAndString, "_AND_")
        whereString = append(whereString, whereAnd)
	}
    return strings.Join(whereString, "__OR__")
}