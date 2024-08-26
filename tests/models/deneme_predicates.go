package models

import "github.com/MrSametBurgazoglu/enterprise/client"

import "github.com/google/uuid"

type DenemePredicate struct {
	where []*client.WhereList
}

func (t *DenemePredicate) Where(w ...*client.Where) {
	t.where = nil
	wl := &client.WhereList{}
	wl.Items = append(wl.Items, w...)
	t.where = append(t.where, wl)
}

func (t *DenemePredicate) ORWhere(w ...*client.Where) {
	wl := &client.WhereList{}
	wl.Items = append(wl.Items, w...)
	t.where = append(t.where, wl)
}

func (t *DenemePredicate) IsIDEqual(v uuid.UUID) *client.Where {
	return &client.Where{
		Type:     client.EQ,
		Name:     DenemeIDField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsTestIDEqual(v uuid.UUID) *client.Where {
	return &client.Where{
		Type:     client.EQ,
		Name:     DenemeTestIDField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsCountEqual(v int) *client.Where {
	return &client.Where{
		Type:     client.EQ,
		Name:     DenemeCountField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsIsActiveEqual(v bool) *client.Where {
	return &client.Where{
		Type:     client.EQ,
		Name:     DenemeIsActiveField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsDenemeTypeEqual(v DenemeType) *client.Where {
	return &client.Where{
		Type:     client.EQ,
		Name:     DenemeDenemeTypeField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsIDNotEqual(v uuid.UUID) *client.Where {
	return &client.Where{
		Type:     client.NEQ,
		Name:     DenemeIDField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsTestIDNotEqual(v uuid.UUID) *client.Where {
	return &client.Where{
		Type:     client.NEQ,
		Name:     DenemeTestIDField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsCountNotEqual(v int) *client.Where {
	return &client.Where{
		Type:     client.NEQ,
		Name:     DenemeCountField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsIsActiveNotEqual(v bool) *client.Where {
	return &client.Where{
		Type:     client.NEQ,
		Name:     DenemeIsActiveField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsDenemeTypeNotEqual(v DenemeType) *client.Where {
	return &client.Where{
		Type:     client.NEQ,
		Name:     DenemeDenemeTypeField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsIDIN(v ...uuid.UUID) *client.Where {
	return &client.Where{
		Type:     client.ANY,
		Name:     DenemeIDField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsTestIDIN(v ...uuid.UUID) *client.Where {
	return &client.Where{
		Type:     client.ANY,
		Name:     DenemeTestIDField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsCountIN(v ...int) *client.Where {
	return &client.Where{
		Type:     client.ANY,
		Name:     DenemeCountField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsIsActiveIN(v ...bool) *client.Where {
	return &client.Where{
		Type:     client.ANY,
		Name:     DenemeIsActiveField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsDenemeTypeIN(v ...DenemeType) *client.Where {
	return &client.Where{
		Type:     client.ANY,
		Name:     DenemeDenemeTypeField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsIDNotIN(v ...uuid.UUID) *client.Where {
	return &client.Where{
		Type:     client.NANY,
		Name:     DenemeIDField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsTestIDNotIN(v ...uuid.UUID) *client.Where {
	return &client.Where{
		Type:     client.NANY,
		Name:     DenemeTestIDField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsCountNotIN(v ...int) *client.Where {
	return &client.Where{
		Type:     client.NANY,
		Name:     DenemeCountField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsIsActiveNotIN(v ...bool) *client.Where {
	return &client.Where{
		Type:     client.NANY,
		Name:     DenemeIsActiveField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsDenemeTypeNotIN(v ...DenemeType) *client.Where {
	return &client.Where{
		Type:     client.NANY,
		Name:     DenemeDenemeTypeField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) IsTestIDNil() *client.Where {
	return &client.Where{
		Type: client.NIL,
		Name: DenemeTestIDField,
	}
}

func (t *DenemePredicate) IsTestIDNotNil() *client.Where {
	return &client.Where{
		Type: client.NNIL,
		Name: DenemeTestIDField,
	}
}

func (t *DenemePredicate) CountGreaterThan(v int) *client.Where {
	return &client.Where{
		Type:     client.GT,
		Name:     DenemeCountField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) CountGreaterEqualThan(v int) *client.Where {
	return &client.Where{
		Type:     client.GTE,
		Name:     DenemeCountField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) CountLowerThan(v int) *client.Where {
	return &client.Where{
		Type:     client.LT,
		Name:     DenemeCountField,
		HasValue: true,
		Value:    v,
	}
}

func (t *DenemePredicate) CountLowerEqualThan(v int) *client.Where {
	return &client.Where{
		Type:     client.LTE,
		Name:     DenemeCountField,
		HasValue: true,
		Value:    v,
	}
}
