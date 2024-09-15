package models

import "github.com/MrSametBurgazoglu/enterprise/client"

import "github.com/google/uuid"
import "time"
import "github.com/MrSametBurgazoglu/enterprise/tests/custom_data_type"
import "strings"

type TestPredicate struct {
	where []*client.WhereList
}

func (t *TestPredicate) Where(w ...*client.Where) {
	t.where = nil
	wl := &client.WhereList{}
	wl.Items = append(wl.Items, w...)
	t.where = append(t.where, wl)
}

func (t *TestPredicate) ORWhere(w ...*client.Where) {
	wl := &client.WhereList{}
	wl.Items = append(wl.Items, w...)
	t.where = append(t.where, wl)
}

func (t *TestPredicate) IsIDEqual(v uuid.UUID) *client.Where {
	return &client.Where{
		Type:     client.EQ,
		Name:     TestIDField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) IsNameEqual(v string) *client.Where {
	return &client.Where{
		Type:     client.EQ,
		Name:     TestNameField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) IsCreatedAtEqual(v time.Time) *client.Where {
	return &client.Where{
		Type:     client.EQ,
		Name:     TestCreatedAtField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) IsInfoEqual(v custom_data_type.Custom) *client.Where {
	return &client.Where{
		Type:     client.EQ,
		Name:     TestInfoField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) IsIDNotEqual(v uuid.UUID) *client.Where {
	return &client.Where{
		Type:     client.NEQ,
		Name:     TestIDField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) IsNameNotEqual(v string) *client.Where {
	return &client.Where{
		Type:     client.NEQ,
		Name:     TestNameField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) IsCreatedAtNotEqual(v time.Time) *client.Where {
	return &client.Where{
		Type:     client.NEQ,
		Name:     TestCreatedAtField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) IsInfoNotEqual(v custom_data_type.Custom) *client.Where {
	return &client.Where{
		Type:     client.NEQ,
		Name:     TestInfoField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) IsIDIN(v ...uuid.UUID) *client.Where {
	return &client.Where{
		Type:     client.ANY,
		Name:     TestIDField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) IsNameIN(v ...string) *client.Where {
	return &client.Where{
		Type:     client.ANY,
		Name:     TestNameField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) IsCreatedAtIN(v ...time.Time) *client.Where {
	return &client.Where{
		Type:     client.ANY,
		Name:     TestCreatedAtField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) IsInfoIN(v ...custom_data_type.Custom) *client.Where {
	return &client.Where{
		Type:     client.ANY,
		Name:     TestInfoField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) IsIDNotIN(v ...uuid.UUID) *client.Where {
	return &client.Where{
		Type:     client.NANY,
		Name:     TestIDField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) IsNameNotIN(v ...string) *client.Where {
	return &client.Where{
		Type:     client.NANY,
		Name:     TestNameField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) IsCreatedAtNotIN(v ...time.Time) *client.Where {
	return &client.Where{
		Type:     client.NANY,
		Name:     TestCreatedAtField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) IsInfoNotIN(v ...custom_data_type.Custom) *client.Where {
	return &client.Where{
		Type:     client.NANY,
		Name:     TestInfoField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) IsInfoNil() *client.Where {
	return &client.Where{
		Type: client.NIL,
		Name: TestInfoField,
	}
}

func (t *TestPredicate) IsInfoNotNil() *client.Where {
	return &client.Where{
		Type: client.NNIL,
		Name: TestInfoField,
	}
}

func (t *TestPredicate) CreatedAtGreaterThan(v time.Time) *client.Where {
	return &client.Where{
		Type:     client.GT,
		Name:     TestCreatedAtField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) CreatedAtGreaterEqualThan(v time.Time) *client.Where {
	return &client.Where{
		Type:     client.GTE,
		Name:     TestCreatedAtField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) CreatedAtLowerThan(v time.Time) *client.Where {
	return &client.Where{
		Type:     client.LT,
		Name:     TestCreatedAtField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) CreatedAtLowerEqualThan(v time.Time) *client.Where {
	return &client.Where{
		Type:     client.LTE,
		Name:     TestCreatedAtField,
		HasValue: true,
		Value:    v,
	}
}

func (t *TestPredicate) GetWhereInfoString() string {
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
