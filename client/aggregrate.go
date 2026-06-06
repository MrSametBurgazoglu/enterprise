package client

import "fmt"

const (
	aggregateMin    = "MIN(%s)"
	aggregateMax    = "MAX(%s)"
	aggregateCount  = "COUNT(%s)"
	aggregateSum    = "SUM(%s)"
	aggregateAvg    = "AVG(%s)"
	aggregateCustom = "%s"
)

type Aggregate struct {
	aggregateFormats []string
	aggregateFields  []string
	aggregateValues  []any
	groupByList      []string
}

func (a *Aggregate) GroupBy(tableName string) *Aggregate {
	a.groupByList = append(a.groupByList, tableName)
	return a
}

func (a *Aggregate) addField(field string, value any) {
	a.aggregateFields = append(a.aggregateFields, field)
	a.aggregateValues = append(a.aggregateValues, value)
}

func (a *Aggregate) Min(field string, value any) *Aggregate {
	a.aggregateFormats = append(a.aggregateFormats, aggregateMin)
	a.addField(field, value)
	return a
}

func (a *Aggregate) Max(field string, value any) *Aggregate {
	a.aggregateFormats = append(a.aggregateFormats, aggregateMax)
	a.addField(field, value)
	return a
}

func (a *Aggregate) Count(field string, value any) *Aggregate {
	a.aggregateFormats = append(a.aggregateFormats, aggregateCount)
	a.addField(field, value)
	return a
}

func (a *Aggregate) Sum(field string, value any) *Aggregate {
	a.aggregateFormats = append(a.aggregateFormats, aggregateSum)
	a.addField(field, value)
	return a
}

func (a *Aggregate) SumCast(field string, cast string, value any) *Aggregate {
	a.aggregateFormats = append(a.aggregateFormats, fmt.Sprintf("SUM(CAST(%%s AS %s))", cast))
	a.addField(field, value)
	return a
}

func (a *Aggregate) Avg(field string, value any) *Aggregate {
	a.aggregateFormats = append(a.aggregateFormats, aggregateAvg)
	a.addField(field, value)
	return a
}

func (a *Aggregate) AvgCast(field string, cast string, value any) *Aggregate {
	a.aggregateFormats = append(a.aggregateFormats, fmt.Sprintf("AVG(CAST(%%s AS %s))", cast))
	a.addField(field, value)
	return a
}

func (a *Aggregate) Field(field string, value any) *Aggregate {
	a.aggregateFormats = append(a.aggregateFormats, aggregateCustom)
	a.addField(field, value)
	return a
}
