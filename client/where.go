package client

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

const (
	EQ    = "\"%s\".\"%s\" = %s"
	NEQ   = "\"%s\".\"%s\" != %s"
	GT    = "\"%s\".\"%s\" > %s"
	GTE   = "\"%s\".\"%s\" >= %s"
	LT    = "\"%s\".\"%s\" < %s"
	LTE   = "\"%s\".\"%s\" <= %s"
	NIL   = "\"%s\".\"%s\" IS NULL"
	NNIL  = "\"%s\".\"%s\" IS NOT NULL"
	ANY   = "\"%s\".\"%s\" = ANY(%s)"
	NANY  = "\"%s\".\"%s\" != ANY(%s)"
	LIKE  = "\"%s\".\"%s\" LIKE %s"
	ILIKE = "\"%s\".\"%s\" ILIKE %s"
)

type PredicateI interface {
	Parse(tableName string, args pgx.NamedArgs) string
	GetName() string
}

type LogicalOperator string

const (
	OperatorAnd LogicalOperator = "AND"
	OperatorOr  LogicalOperator = "OR"
)

type LogicalPredicate struct {
	Operator   LogicalOperator
	Predicates []PredicateI
}

func (l *LogicalPredicate) Parse(tableName string, args pgx.NamedArgs) string {
	var parts []string
	for _, p := range l.Predicates {
		s := p.Parse(tableName, args)
		if s != "" {
			parts = append(parts, s)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return "(" + strings.Join(parts, fmt.Sprintf(" %s ", l.Operator)) + ")"
}

func (l *LogicalPredicate) GetName() string {
	var names []string
	for _, p := range l.Predicates {
		names = append(names, p.GetName())
	}
	return "(" + strings.Join(names, fmt.Sprintf("_%s_", l.Operator)) + ")"
}

func Or(preds ...PredicateI) PredicateI {
	return &LogicalPredicate{Operator: OperatorOr, Predicates: preds}
}

func And(preds ...PredicateI) PredicateI {
	return &LogicalPredicate{Operator: OperatorAnd, Predicates: preds}
}

type Where struct {
	Type     string
	Name     string
	HasValue bool
	Value    any
}

func (w *Where) GetName() string {
	return w.Name
}

func (w *Where) Parse(tableName string, args pgx.NamedArgs) string {
	if !w.HasValue {
		return fmt.Sprintf(w.Type, tableName, w.Name)
	}

	baseName := fmt.Sprintf("%s__%s", tableName, w.Name)
	paramName := ""
	for idx := 1; ; idx++ {
		candidate := fmt.Sprintf("%s_%d", baseName, idx)
		if _, exists := args[candidate]; !exists {
			paramName = candidate
			break
		}
	}

	args[paramName] = w.Value
	return fmt.Sprintf(w.Type, tableName, w.Name, fmt.Sprintf("@%s", paramName))
}

func (w *Where) GetSqlString(tableName string) string {
	// Deprecated: kept for compatibility
	if w.HasValue {
		return fmt.Sprintf(w.Type, tableName, w.Name, fmt.Sprintf("@%s__%s", tableName, w.Name))
	} else {
		return fmt.Sprintf(w.Type, tableName, w.Name)
	}
}
