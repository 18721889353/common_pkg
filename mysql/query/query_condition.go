// Package query is a library for mysql query, support for complex conditional paging queries.
package query

import (
	"fmt"
	"strings"
)

const (
	// Eq equal
	Eq = "eq"
	// Neq not equal
	Neq = "neq"
	// Gt greater than
	Gt = "gt"
	// Gte greater than or equal
	Gte = "gte"
	// Lt less than
	Lt = "lt"
	// Lte less than or equal
	Lte = "lte"
	// Like fuzzy lookup
	Like = "like"

	// AND logic and
	AND string = "and"
	// OR logic or
	OR string = "or"
)

var expMap = map[string]string{
	Eq:   " = ",
	Neq:  " <> ",
	Gt:   " > ",
	Gte:  " >= ",
	Lt:   " < ",
	Lte:  " <= ",
	Like: " LIKE ",

	"=":  " = ",
	"!=": " <> ",
	">":  " > ",
	">=": " >= ",
	"<":  " < ",
	"<=": " <= ",
}

var logicMap = map[string]string{
	AND: " AND ",
	OR:  " OR ",

	"&":  " AND ",
	"&&": " AND ",
	"|":  " OR ",
	"||": " OR ",
}

// Params query parameters
type Params struct {
	Page int    `form:"page" binding:"gte=0" json:"page"`
	Size int    `form:"size" binding:"gt=0" json:"size"`
	Sort string `form:"sort" binding:"" json:"sort,omitempty"`

	Columns []Column `json:"columns,omitempty"` // not required
}

// Column search information
type Column struct {
	Name  string      `json:"name"`  // column name
	Exp   string      `json:"exp"`   // expressions, which default to = when the value is null, have =, ! =, >, >=, <, <=, like
	Value interface{} `json:"value"` // column value
	Logic string      `json:"logic"` // logical type, defaults to and when the value is null, with &(and), ||(or)
}

func (c *Column) checkValid() error {
	if c.Name == "" {
		return fmt.Errorf("field 'name' cannot be empty")
	}
	if c.Value == nil {
		return fmt.Errorf("field 'value' cannot be nil")
	}
	return nil
}

// converting ExpType to sql expressions and LogicType to sql using characters
func (c *Column) convert() error {
	if c.Exp == "" {
		c.Exp = Eq
	}
	if v, ok := expMap[strings.ToLower(c.Exp)]; ok { //nolint
		c.Exp = v
		if c.Exp == " LIKE " {
			c.Value = fmt.Sprintf("%%%v%%", c.Value)
		}
	} else {
		return fmt.Errorf("unknown c expression type '%s'", c.Exp)
	}

	if c.Logic == "" {
		c.Logic = AND
	}
	if v, ok := logicMap[strings.ToLower(c.Logic)]; ok { //nolint
		c.Logic = v
	} else {
		return fmt.Errorf("unknown logic type '%s'", c.Logic)
	}

	return nil
}

// ConvertToPage converted to conform to gorm rules based on the page size sort parameter
func (p *Params) ConvertToPage() (order string, limit int, offset int) { //nolint
	page := NewPage(p.Page, p.Size, p.Sort)
	order = page.sort
	limit = page.size
	offset = page.page * page.size
	return //nolint
}

// ConvertToGormConditions conversion to gorm-compliant parameters based on the Columns parameter
// ignore the logical type of the last column, whether it is a one-column or multi-column query
func (p *Params) ConvertToGormConditions() (string, []interface{}, error) {
	str := ""
	args := []interface{}{}
	l := len(p.Columns)
	if l == 0 {
		return "", nil, nil
	}

	isUseIN := true
	if l == 1 {
		isUseIN = false
	}
	field := p.Columns[0].Name

	for i, column := range p.Columns {
		if err := column.checkValid(); err != nil {
			return "", nil, err
		}

		err := column.convert()
		if err != nil {
			return "", nil, err
		}

		if i == l-1 { // ignore the logical type of the last column
			str += column.Name + column.Exp + "?"
		} else {
			str += column.Name + column.Exp + "?" + column.Logic
		}
		args = append(args, column.Value)

		if isUseIN {
			if field != column.Name {
				isUseIN = false
				continue
			}
			if column.Exp != expMap[Eq] {
				isUseIN = false
			}
		}
	}

	if isUseIN {
		str = field + " IN (?)"
		args = []interface{}{args}
	}

	return str, args, nil
}

func getExpsAndLogics(keyLen int, paramSrc string) ([]string, []string) { //nolint
	exps, logics := []string{}, []string{}
	param := strings.Replace(paramSrc, " ", "", -1)
	sps := strings.SplitN(param, "?", 2)
	if len(sps) == 2 {
		param = sps[1]
	}

	num := keyLen
	if num == 0 {
		return exps, logics
	}

	fields := []string{}
	kvs := strings.Split(param, "&")
	for _, kv := range kvs {
		if strings.Contains(kv, "page=") || strings.Contains(kv, "size=") || strings.Contains(kv, "sort=") {
			continue
		}
		fields = append(fields, kv)
	}

	// divide into num groups based on non-repeating keys, and determine in each group whether exp and logic exist
	group := map[string]string{}
	for _, field := range fields {
		split := strings.SplitN(field, "=", 2)
		if len(split) != 2 {
			continue
		}

		if _, ok := group[split[0]]; ok {
			// if exp does not exist, the default value of null is filled, and if logic does not exist, the default value of null is filled.
			exps = append(exps, group["exp"])
			logics = append(logics, group["logic"])

			group = map[string]string{}
			continue
		}
		group[split[0]] = split[1]
	}

	// handling the last group
	exps = append(exps, group["exp"])
	logics = append(logics, group["logic"])

	return exps, logics
}
