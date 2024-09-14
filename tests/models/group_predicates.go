package models

import "github.com/MrSametBurgazoglu/enterprise/client"

import "github.com/google/uuid"
import "strings"

type GroupPredicate struct {
	where []*client.WhereList
}

func (t *GroupPredicate) Where(w ...*client.Where) {
	t.where = nil
	wl := &client.WhereList{}
	wl.Items = append(wl.Items, w...)
	t.where = append(t.where, wl)
}

func (t *GroupPredicate) ORWhere(w ...*client.Where) {
	wl := &client.WhereList{}
	wl.Items = append(wl.Items, w...)
	t.where = append(t.where, wl)
}

func (t *GroupPredicate) IsIDEqual(v uuid.UUID) *client.Where {
	return &client.Where{
		Type:     client.EQ,
		Name:     GroupIDField,
		HasValue: true,
		Value:    v,
	}
}

func (t *GroupPredicate) IsNameEqual(v string) *client.Where {
	return &client.Where{
		Type:     client.EQ,
		Name:     GroupNameField,
		HasValue: true,
		Value:    v,
	}
}

func (t *GroupPredicate) IsSurnameEqual(v string) *client.Where {
	return &client.Where{
		Type:     client.EQ,
		Name:     GroupSurnameField,
		HasValue: true,
		Value:    v,
	}
}

func (t *GroupPredicate) IsDataEqual(v map[string]any) *client.Where {
	return &client.Where{
		Type:     client.EQ,
		Name:     GroupDataField,
		HasValue: true,
		Value:    v,
	}
}

func (t *GroupPredicate) IsIDNotEqual(v uuid.UUID) *client.Where {
	return &client.Where{
		Type:     client.NEQ,
		Name:     GroupIDField,
		HasValue: true,
		Value:    v,
	}
}

func (t *GroupPredicate) IsNameNotEqual(v string) *client.Where {
	return &client.Where{
		Type:     client.NEQ,
		Name:     GroupNameField,
		HasValue: true,
		Value:    v,
	}
}

func (t *GroupPredicate) IsSurnameNotEqual(v string) *client.Where {
	return &client.Where{
		Type:     client.NEQ,
		Name:     GroupSurnameField,
		HasValue: true,
		Value:    v,
	}
}

func (t *GroupPredicate) IsDataNotEqual(v map[string]any) *client.Where {
	return &client.Where{
		Type:     client.NEQ,
		Name:     GroupDataField,
		HasValue: true,
		Value:    v,
	}
}

func (t *GroupPredicate) IsIDIN(v ...uuid.UUID) *client.Where {
	return &client.Where{
		Type:     client.ANY,
		Name:     GroupIDField,
		HasValue: true,
		Value:    v,
	}
}

func (t *GroupPredicate) IsNameIN(v ...string) *client.Where {
	return &client.Where{
		Type:     client.ANY,
		Name:     GroupNameField,
		HasValue: true,
		Value:    v,
	}
}

func (t *GroupPredicate) IsSurnameIN(v ...string) *client.Where {
	return &client.Where{
		Type:     client.ANY,
		Name:     GroupSurnameField,
		HasValue: true,
		Value:    v,
	}
}

func (t *GroupPredicate) IsDataIN(v ...map[string]any) *client.Where {
	return &client.Where{
		Type:     client.ANY,
		Name:     GroupDataField,
		HasValue: true,
		Value:    v,
	}
}

func (t *GroupPredicate) IsIDNotIN(v ...uuid.UUID) *client.Where {
	return &client.Where{
		Type:     client.NANY,
		Name:     GroupIDField,
		HasValue: true,
		Value:    v,
	}
}

func (t *GroupPredicate) IsNameNotIN(v ...string) *client.Where {
	return &client.Where{
		Type:     client.NANY,
		Name:     GroupNameField,
		HasValue: true,
		Value:    v,
	}
}

func (t *GroupPredicate) IsSurnameNotIN(v ...string) *client.Where {
	return &client.Where{
		Type:     client.NANY,
		Name:     GroupSurnameField,
		HasValue: true,
		Value:    v,
	}
}

func (t *GroupPredicate) IsDataNotIN(v ...map[string]any) *client.Where {
	return &client.Where{
		Type:     client.NANY,
		Name:     GroupDataField,
		HasValue: true,
		Value:    v,
	}
}

func (t *GroupPredicate) GetWhereInfoString() string {
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
