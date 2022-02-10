package query

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/toolkit/sqlext/arithsymbol"
	"github.com/unionj-cloud/go-doudou/toolkit/sqlext/logicsymbol"
	"github.com/unionj-cloud/go-doudou/toolkit/sqlext/sortenum"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"reflect"
	"strings"
)

// Base sql expression
type Base interface {
	Sql() (string, []interface{})
	//NamedSql() (string, []interface{})
}

// Q used for building sql expression
type Q interface {
	Base
	And(q Base) Where
	Or(q Base) Where
	Append(q Base) Where
}

// Criteria wrap a group of column, value and operator such as name = 20
type Criteria struct {
	// table alias
	talias string
	col    string
	val    interface{}
	asym   arithsymbol.ArithSymbol
}

// Sql implement Base interface, return sql expression
func (c Criteria) Sql() (string, []interface{}) {
	if c.asym == arithsymbol.In {
		var args []interface{}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("`%s` %s (", c.col, c.asym))

		var vals []string
		switch reflect.TypeOf(c.val).Kind() {
		case reflect.Slice:
			data := reflect.ValueOf(c.val)
			for i := 0; i < data.Len(); i++ {
				vals = append(vals, "?")
				args = append(args, data.Index(i).Interface())
			}
		default:
			vals = append(vals, "?")
			args = append(args, c.val)
		}

		sb.WriteString(strings.Join(vals, ","))
		sb.WriteString(")")

		return sb.String(), args
	}
	if stringutils.IsNotEmpty(c.talias) {
		if c.asym == arithsymbol.Is || c.asym == arithsymbol.Not {
			return fmt.Sprintf("%s.`%s` %s null", c.talias, c.col, c.asym), nil
		}
		return fmt.Sprintf("%s.`%s` %s ?", c.talias, c.col, c.asym), []interface{}{c.val}
	}
	if c.asym == arithsymbol.Is || c.asym == arithsymbol.Not {
		return fmt.Sprintf("`%s` %s null", c.col, c.asym), nil
	}
	return fmt.Sprintf("`%s` %s ?", c.col, c.asym), []interface{}{c.val}
}

// C new a Criteria
func C() Criteria {
	return Criteria{}
}

// Col set column name
func (c Criteria) Col(col string) Criteria {
	if strings.Contains(col, ".") {
		i := strings.Index(col, ".")
		c.talias = col[:i]
		c.col = col[i+1:]
	} else {
		c.col = col
	}
	return c
}

// Eq set = operator and column value
func (c Criteria) Eq(val interface{}) Criteria {
	c.val = val
	c.asym = arithsymbol.Eq
	return c
}

// Ne set != operator and column value
func (c Criteria) Ne(val interface{}) Criteria {
	c.val = val
	c.asym = arithsymbol.Ne
	return c
}

// Gt set > operator and column value
func (c Criteria) Gt(val interface{}) Criteria {
	c.val = val
	c.asym = arithsymbol.Gt
	return c
}

// Lt set < operator and column value
func (c Criteria) Lt(val interface{}) Criteria {
	c.val = val
	c.asym = arithsymbol.Lt
	return c
}

// Gte set >= operator and column value
func (c Criteria) Gte(val interface{}) Criteria {
	c.val = val
	c.asym = arithsymbol.Gte
	return c
}

// Lte set <= operator and column value
func (c Criteria) Lte(val interface{}) Criteria {
	c.val = val
	c.asym = arithsymbol.Lte
	return c
}

// IsNull set is null
func (c Criteria) IsNull() Criteria {
	c.asym = arithsymbol.Is
	return c
}

// IsNotNull set is not null
func (c Criteria) IsNotNull() Criteria {
	c.asym = arithsymbol.Not
	return c
}

// In set in operator and column value, val should be a slice type value
func (c Criteria) In(val interface{}) Criteria {
	c.val = val
	c.asym = arithsymbol.In
	return c
}

// Like set like operator and column value, val should be a slice type value
func (c Criteria) Like(val interface{}) Criteria {
	c.val = val
	c.asym = arithsymbol.Like
	return c
}

// And concat another sql expression builder with And
func (c Criteria) And(cri Base) Where {
	w := Where{
		children: make([]Base, 0),
	}
	w.children = append(w.children, c, cri)
	w.lsym = logicsymbol.And
	return w
}

// Or concat another sql expression builder with Or
func (c Criteria) Or(cri Base) Where {
	w := Where{
		children: make([]Base, 0),
	}
	w.children = append(w.children, c, cri)
	w.lsym = logicsymbol.Or
	return w
}

// Append concat another sql expression builder with Append
func (c Criteria) Append(cri Base) Where {
	w := Where{
		children: make([]Base, 0),
	}
	w.children = append(w.children, c, cri)
	w.lsym = logicsymbol.Append
	return w
}

// Where concat children clauses with one of logic operators And, Or, Append
type Where struct {
	lsym     logicsymbol.LogicSymbol
	children []Base
}

// Sql implement Base interface, return string sql expression
func (w Where) Sql() (string, []interface{}) {
	var args []interface{}
	w0, args0 := w.children[0].Sql()
	args = append(args, args0...)
	w1, args1 := w.children[1].Sql()
	args = append(args, args1...)
	if w.lsym != logicsymbol.Append {
		return fmt.Sprintf("(%s %s %s)", w0, w.lsym, w1), args
	}
	return fmt.Sprintf("%s%s%s", w0, w.lsym, w1), args
}

// And concat another sql expression builder with And
func (w Where) And(whe Base) Where {
	parentW := Where{
		children: make([]Base, 0),
	}
	parentW.children = append(parentW.children, w, whe)
	parentW.lsym = logicsymbol.And
	return parentW
}

// Or concat another sql expression builder with Or
func (w Where) Or(whe Base) Where {
	parentW := Where{
		children: make([]Base, 0),
	}
	parentW.children = append(parentW.children, w, whe)
	parentW.lsym = logicsymbol.Or
	return parentW
}

// Append concat another sql expression builder with Append
func (w Where) Append(whe Base) Where {
	parentW := Where{
		children: make([]Base, 0),
	}
	parentW.children = append(parentW.children, w, whe)
	parentW.lsym = logicsymbol.Append
	return parentW
}

// Order by Col Sort
type Order struct {
	Col  string
	Sort sortenum.Sort
}

// Page a sql expression builder for order by clause
type Page struct {
	Orders []Order
	Offset int
	Size   int
}

// P new a Page
func P() Page {
	return Page{
		Orders: make([]Order, 0),
	}
}

// Order append an Order
func (p Page) Order(o Order) Page {
	p.Orders = append(p.Orders, o)
	return p
}

// Limit set Offset and Size
func (p Page) Limit(offset, size int) Page {
	p.Offset = offset
	p.Size = size
	return p
}

// Sql implement Base interface, order by age desc limit 2,1
func (p Page) Sql() (string, []interface{}) {
	var sb strings.Builder
	var args []interface{}
	if len(p.Orders) > 0 {
		sb.WriteString("order by ")

		for i, order := range p.Orders {
			if i > 0 {
				sb.WriteString(",")
			}
			var (
				alias string
				col   string
			)
			if strings.Contains(order.Col, ".") {
				idx := strings.Index(order.Col, ".")
				alias = order.Col[:idx]
				col = order.Col[idx+1:]
			} else {
				col = order.Col
			}
			if stringutils.IsNotEmpty(alias) {
				sb.WriteString(fmt.Sprintf("%s.`%s` %s", alias, col, order.Sort))
			} else {
				sb.WriteString(fmt.Sprintf("`%s` %s", col, order.Sort))
			}
		}
	}

	sb.WriteString(" ")

	if p.Size > 0 {
		sb.WriteString("limit ?,?")
		args = append(args, p.Offset, p.Size)
	}

	return strings.TrimSpace(sb.String()), args
}

// PageRet wrap page query result
type PageRet struct {
	Items    interface{}
	PageNo   int
	PageSize int
	Total    int
	HasNext  bool
}

// NewPageRet new a PageRet
func NewPageRet(page Page) PageRet {
	pageNo := 1
	if page.Size > 0 {
		pageNo = page.Offset/page.Size + 1
	}
	return PageRet{
		PageNo:   pageNo,
		PageSize: page.Size,
	}
}

// String is an alias of string
type String string

// Sql implements Base
func (s String) Sql() (string, []interface{}) {
	return string(s), nil
}
