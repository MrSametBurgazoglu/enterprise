package models

import "github.com/MrSametBurgazoglu/enterprise/client"

import "github.com/google/uuid"

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
