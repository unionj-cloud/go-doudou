package query

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/kit/ddl/arithsymbol"
	"github.com/unionj-cloud/go-doudou/kit/ddl/logicsymbol"
	"github.com/unionj-cloud/go-doudou/kit/ddl/sortenum"
	"strings"
)

type Q interface {
	Sql() string
	And(q Q) Q
	Or(q Q) Q
}

type criteria struct {
	col  string
	val  interface{}
	asym arithsymbol.ArithSymbol
}

func (c criteria) Sql() string {
	return fmt.Sprintf("%s %s '%s'", c.col, c.asym, c.val)
}

func C() criteria {
	return criteria{}
}

func (c criteria) Col(col string) criteria {
	c.col = col
	return c
}

func (c criteria) Eq(val interface{}) criteria {
	c.val = val
	c.asym = arithsymbol.Eq
	return c
}

func (c criteria) Ne(val interface{}) criteria {
	c.val = val
	c.asym = arithsymbol.Ne
	return c
}

func (c criteria) Gt(val interface{}) criteria {
	c.val = val
	c.asym = arithsymbol.Gt
	return c
}

func (c criteria) Lt(val interface{}) criteria {
	c.val = val
	c.asym = arithsymbol.Lt
	return c
}

func (c criteria) Gte(val interface{}) criteria {
	c.val = val
	c.asym = arithsymbol.Gte
	return c
}

func (c criteria) Lte(val interface{}) criteria {
	c.val = val
	c.asym = arithsymbol.Lte
	return c
}

func (c criteria) Is(val interface{}) criteria {
	c.val = val
	c.asym = arithsymbol.Is
	return c
}

func (c criteria) Not(val interface{}) criteria {
	c.val = val
	c.asym = arithsymbol.Not
	return c
}

func (c criteria) And(cri Q) Q {
	w := where{
		children: make([]Q, 0),
	}
	w.children = append(w.children, c, cri)
	w.lsym = logicsymbol.And
	return w
}

func (c criteria) Or(cri Q) Q {
	w := where{
		children: make([]Q, 0),
	}
	w.children = append(w.children, c, cri)
	w.lsym = logicsymbol.Or
	return w
}

type where struct {
	lsym     logicsymbol.LogicSymbol
	children []Q
}

func (w where) Sql() string {
	return fmt.Sprintf("(%s %s %s)", w.children[0].Sql(), w.lsym, w.children[1].Sql())
}

func (w where) And(whe Q) Q {
	parentW := where{
		children: make([]Q, 0),
	}
	parentW.children = append(parentW.children, w, whe)
	parentW.lsym = logicsymbol.And
	return parentW
}

func (w where) Or(whe Q) Q {
	parentW := where{
		children: make([]Q, 0),
	}
	parentW.children = append(parentW.children, w, whe)
	parentW.lsym = logicsymbol.Or
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

// order by age desc limit 2,1
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
