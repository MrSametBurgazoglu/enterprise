package client

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

func (a *Aggregate) GroupBy(tableName string) {
	a.groupByList = append(a.groupByList, tableName)
}

func (a *Aggregate) addField(field string, value any) {
	a.aggregateFields = append(a.aggregateFields, field)
	a.aggregateValues = append(a.aggregateValues, value)
}

func (a *Aggregate) Min(field string, value any) {
	a.aggregateFormats = append(a.aggregateFormats, aggregateMin)
	a.addField(field, value)
}

func (a *Aggregate) Max(field string, value any) {
	a.aggregateFormats = append(a.aggregateFormats, aggregateMax)
	a.addField(field, value)
}

func (a *Aggregate) Count(field string, value any) {
	a.aggregateFormats = append(a.aggregateFormats, aggregateCount)
	a.addField(field, value)
}

func (a *Aggregate) Sum(field string, value any) {
	a.aggregateFormats = append(a.aggregateFormats, aggregateSum)
	a.addField(field, value)
}

func (a *Aggregate) Avg(field string, value any) {
	a.aggregateFormats = append(a.aggregateFormats, aggregateAvg)
	a.addField(field, value)
}

func (a *Aggregate) Field(field string, value any) {
	a.aggregateFormats = append(a.aggregateFormats, aggregateCustom)
	a.addField(field, value)
}
