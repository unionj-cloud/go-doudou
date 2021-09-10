package query

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/ddl/arithsymbol"
	"github.com/unionj-cloud/go-doudou/ddl/logicsymbol"
	"github.com/unionj-cloud/go-doudou/ddl/sortenum"
	"github.com/unionj-cloud/go-doudou/ddl/valtypeenum"
	"github.com/unionj-cloud/go-doudou/reflectutils"
	"reflect"
	"strings"
)

type Base interface {
	Sql() string
}

type Q interface {
	Base
	And(q Base) Where
	Or(q Base) Where
	Append(q Base) Where
}

type Val struct {
	Data interface{}
	Type valtypeenum.ValType
}

func Literal(data interface{}) Val {
	return Val{
		Data: data,
		Type: valtypeenum.Literal,
	}
}

func Func(data string) Val {
	return Val{
		Data: data,
		Type: valtypeenum.Func,
	}
}

func null() Val {
	return Val{
		Data: "null",
		Type: valtypeenum.Null,
	}
}

type Criteria struct {
	col  string
	val  Val
	asym arithsymbol.ArithSymbol
}

func (c Criteria) Sql() string {
	if c.asym == arithsymbol.In {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("`%s` %s (", c.col, c.asym))

		var vals []string
		switch reflect.TypeOf(c.val.Data).Kind() {
		case reflect.Slice:
			data := reflect.ValueOf(c.val.Data)
			for i := 0; i < data.Len(); i++ {
				if c.val.Type != valtypeenum.Literal {
					vals = append(vals, fmt.Sprintf("%v", reflectutils.ValueOfValue(data.Index(i))))
				} else {
					vals = append(vals, fmt.Sprintf("'%v'", reflectutils.ValueOfValue(data.Index(i))))
				}
			}
		default:
			if c.val.Type != valtypeenum.Literal {
				vals = append(vals, fmt.Sprintf("%v", reflectutils.ValueOf(c.val.Data)))
			} else {
				vals = append(vals, fmt.Sprintf("'%v'", reflectutils.ValueOf(c.val.Data)))
			}
		}

		sb.WriteString(strings.Join(vals, ","))
		sb.WriteString(")")

		return sb.String()
	} else {
		if c.val.Type != valtypeenum.Literal {
			return fmt.Sprintf("`%s` %s %v", c.col, c.asym, reflectutils.ValueOf(c.val.Data))
		}
		return fmt.Sprintf("`%s` %s '%v'", c.col, c.asym, reflectutils.ValueOf(c.val.Data))
	}
}

func C() Criteria {
	return Criteria{}
}

func (c Criteria) Col(col string) Criteria {
	c.col = col
	return c
}

func (c Criteria) Eq(val Val) Criteria {
	c.val = val
	c.asym = arithsymbol.Eq
	return c
}

func (c Criteria) Ne(val Val) Criteria {
	c.val = val
	c.asym = arithsymbol.Ne
	return c
}

func (c Criteria) Gt(val Val) Criteria {
	c.val = val
	c.asym = arithsymbol.Gt
	return c
}

func (c Criteria) Lt(val Val) Criteria {
	c.val = val
	c.asym = arithsymbol.Lt
	return c
}

func (c Criteria) Gte(val Val) Criteria {
	c.val = val
	c.asym = arithsymbol.Gte
	return c
}

func (c Criteria) Lte(val Val) Criteria {
	c.val = val
	c.asym = arithsymbol.Lte
	return c
}

func (c Criteria) IsNull() Criteria {
	c.val = null()
	c.asym = arithsymbol.Is
	return c
}

func (c Criteria) IsNotNull() Criteria {
	c.val = null()
	c.asym = arithsymbol.Not
	return c
}

func (c Criteria) In(val Val) Criteria {
	c.val = val
	c.asym = arithsymbol.In
	return c
}

func (c Criteria) And(cri Base) Where {
	w := Where{
		children: make([]Base, 0),
	}
	w.children = append(w.children, c, cri)
	w.lsym = logicsymbol.And
	return w
}

func (c Criteria) Or(cri Base) Where {
	w := Where{
		children: make([]Base, 0),
	}
	w.children = append(w.children, c, cri)
	w.lsym = logicsymbol.Or
	return w
}

func (c Criteria) Append(cri Base) Where {
	w := Where{
		children: make([]Base, 0),
	}
	w.children = append(w.children, c, cri)
	w.lsym = logicsymbol.Append
	return w
}

type Where struct {
	lsym     logicsymbol.LogicSymbol
	children []Base
}

func (w Where) Sql() string {
	if w.lsym != logicsymbol.Append {
		return fmt.Sprintf("(%s %s %s)", w.children[0].Sql(), w.lsym, w.children[1].Sql())
	}
	return fmt.Sprintf("%s%s%s", w.children[0].Sql(), w.lsym, w.children[1].Sql())
}

func (w Where) And(whe Base) Where {
	parentW := Where{
		children: make([]Base, 0),
	}
	parentW.children = append(parentW.children, w, whe)
	parentW.lsym = logicsymbol.And
	return parentW
}

func (w Where) Or(whe Base) Where {
	parentW := Where{
		children: make([]Base, 0),
	}
	parentW.children = append(parentW.children, w, whe)
	parentW.lsym = logicsymbol.Or
	return parentW
}

func (w Where) Append(whe Base) Where {
	parentW := Where{
		children: make([]Base, 0),
	}
	parentW.children = append(parentW.children, w, whe)
	parentW.lsym = logicsymbol.Append
	return parentW
}

type Order struct {
	Col  string
	Sort sortenum.Sort
}

type Page struct {
	Orders []Order
	Offset int
	Size   int
}

func P() Page {
	return Page{
		Orders: make([]Order, 0),
	}
}

func (p Page) Order(o Order) Page {
	p.Orders = append(p.Orders, o)
	return p
}

func (p Page) Limit(offset, size int) Page {
	p.Offset = offset
	p.Size = size
	return p
}

// Sql order by age desc limit 2,1
func (p Page) Sql() string {
	var sb strings.Builder

	if len(p.Orders) > 0 {
		sb.WriteString("order by ")

		for i, order := range p.Orders {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(fmt.Sprintf("%s %s", order.Col, order.Sort))
		}
	}

	sb.WriteString(" ")

	if p.Size > 0 {
		sb.WriteString(fmt.Sprintf("limit %d,%d", p.Offset, p.Size))
	}

	return strings.TrimSpace(sb.String())
}

type PageRet struct {
	Items    interface{}
	PageNo   int
	PageSize int
	Total    int
	HasNext  bool
}

func NewPageRet(page Page) PageRet {
	pageNo := page.Offset/page.Size + 1
	return PageRet{
		PageNo:   pageNo,
		PageSize: page.Size,
	}
}
