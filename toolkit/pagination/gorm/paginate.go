package gorm

import (
	"crypto/md5"
	"fmt"
	"log"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/sliceutils"

	"github.com/morkid/gocache"
	"github.com/wubin1989/gorm"
)

var json = sonic.ConfigDefault

// ResponseContext interface
type ResponseContext interface {
	Cache(string) ResponseContext
	Fields([]string) ResponseContext
	Distinct([]string) ResponseContext
	Response(interface{}) Page
	BuildWhereClause() (string, []interface{})
	Error() error
}

// RequestContext interface
type RequestContext interface {
	Request(IParameter) ResponseContext
}

// Pagination gorm paginate struct
type Pagination struct {
	Config *Config
}

// With func
func (p *Pagination) With(stmt *gorm.DB) RequestContext {
	return reqContext{
		Statement:  stmt,
		Pagination: p,
	}
}

// ClearCache clear cache contains prefix
func (p Pagination) ClearCache(keyPrefixes ...string) {
	if len(keyPrefixes) > 0 && nil != p.Config && nil != p.Config.CacheAdapter {
		adapter := *p.Config.CacheAdapter
		for i := range keyPrefixes {
			if err := adapter.ClearPrefix(keyPrefixes[i]); nil != err {
				log.Println(err)
			}
		}
	}
}

// ClearAllCache clear all existing cache
func (p Pagination) ClearAllCache() {
	if nil != p.Config && nil != p.Config.CacheAdapter {
		adapter := *p.Config.CacheAdapter
		if err := adapter.ClearAll(); nil != err {
			log.Println(err)
		}
	}
}

type reqContext struct {
	Statement  *gorm.DB
	Pagination *Pagination
}

func (r reqContext) Request(parameter IParameter) ResponseContext {
	var response ResponseContext = &resContext{
		Statement:  r.Statement,
		Parameter:  parameter,
		Pagination: r.Pagination,
	}

	return response
}

type resContext struct {
	Pagination   *Pagination
	Statement    *gorm.DB
	Parameter    IParameter
	cachePrefix  string
	fieldList    []string
	customSelect string
	distinct     bool
	error        error
}

func (r *resContext) Error() error {
	return r.error
}

func (r *resContext) Cache(prefix string) ResponseContext {
	r.cachePrefix = prefix
	return r
}

func (r *resContext) Fields(fields []string) ResponseContext {
	r.fieldList = fields
	return r
}

// CustomSelect currently used for distinct on clause
func (r *resContext) Distinct(fields []string) ResponseContext {
	r.fieldList = fields
	r.distinct = true
	return r
}

func (r *resContext) Response(res interface{}) Page {
	p := r.Pagination
	query := r.Statement
	p.Config = defaultConfig(p.Config)
	p.Config.Statement = query.Statement
	if p.Config.DefaultSize == 0 {
		p.Config.DefaultSize = 10
	}

	if p.Config.FieldWrapper == "" && p.Config.ValueWrapper == "" {
		defaultWrapper := "LOWER(%s)"
		wrappers := map[string]string{
			"sqlite":   defaultWrapper,
			"mysql":    defaultWrapper,
			"postgres": "LOWER((%s)::text)",
		}
		p.Config.FieldWrapper = defaultWrapper
		if wrapper, ok := wrappers[query.Dialector.Name()]; ok {
			p.Config.FieldWrapper = wrapper
		}
	}

	page := Page{}
	pr, err := parseRequest(r.Parameter, *p.Config)
	if err != nil {
		r.error = err
		return page
	}
	causes := createCauses(pr)
	cKey := ""
	var adapter gocache.AdapterInterface
	var hasAdapter bool = false

	if nil != p.Config.CacheAdapter {
		cKey = createCacheKey(r.cachePrefix, pr)
		adapter = *p.Config.CacheAdapter
		hasAdapter = true
		if cKey != "" && adapter.IsValid(cKey) {
			if cache, err := adapter.Get(cKey); nil == err {
				if err := p.Config.JSONUnmarshal([]byte(cache), &page); nil == err {
					return page
				}
			}
		}
	}

	dbs := query.Statement.DB.Session(&gorm.Session{NewDB: true})
	var selects []string
	if len(r.fieldList) > 0 {
		if len(pr.Fields) > 0 && p.Config.FieldSelectorEnabled {
			for i := range pr.Fields {
				for j := range r.fieldList {
					if r.fieldList[j] == pr.Fields[i] {
						fname := query.Statement.Quote("s." + fieldName(pr.Fields[i]))
						if !contains(selects, fname) {
							selects = append(selects, fname)
						}
						break
					}
				}
			}
		} else {
			for i := range r.fieldList {
				fname := query.Statement.Quote("s." + fieldName(r.fieldList[i]))
				if !contains(selects, fname) {
					selects = append(selects, fname)
				}
			}
		}
	} else if len(pr.Fields) > 0 && p.Config.FieldSelectorEnabled {
		for i := range pr.Fields {
			fname := query.Statement.Quote("s." + fieldName(pr.Fields[i]))
			if !contains(selects, fname) {
				selects = append(selects, fname)
			}
		}
	}

	result := dbs.
		Unscoped().
		Table("(?) AS s", query)

	if len(selects) > 0 {
		if r.distinct {
			result = result.Distinct(selects)
		} else {
			result = result.Select(selects)
		}
	}

	if len(causes.Params) > 0 || len(causes.WhereString) > 0 {
		result = result.Where(causes.WhereString, causes.Params...)
	}

	dbs = query.Statement.DB.Session(&gorm.Session{NewDB: true})
	result = dbs.
		Unscoped().
		Table("(?) AS s1", result)

	if pr.Page >= 0 {
		result = result.Count(&page.Total).
			Limit(int(causes.Limit)).
			Offset(int(causes.Offset))
	}
	if result.Error != nil {
		r.error = result.Error
		return page
	}

	if nil != query.Statement.Preloads {
		for table, args := range query.Statement.Preloads {
			result = result.Preload(table, args...)
		}
	}
	if len(causes.Sorts) > 0 {
		for _, sort := range causes.Sorts {
			result = result.Order(sort.Column + " " + sort.Direction)
		}
	}

	rs := result.Find(res)
	if result.Error != nil {
		r.error = result.Error
		return page
	}

	page.Items, _ = sliceutils.ConvertAny2Interface(res)
	f := float64(page.Total) / float64(causes.Limit)
	if math.Mod(f, 1.0) > 0 {
		f = f + 1
	}
	page.TotalPages = int64(f)
	page.Page = int64(pr.Page)
	page.Size = int64(pr.Size)
	page.MaxPage = 0
	page.Visible = rs.RowsAffected
	if page.TotalPages > 0 {
		page.MaxPage = page.TotalPages - 1
	}
	if page.TotalPages < 1 {
		page.TotalPages = 1
	}
	if page.Total < 1 {
		page.MaxPage = 0
		page.TotalPages = 0
	}
	page.First = causes.Offset < 1
	page.Last = page.MaxPage == page.Page

	if hasAdapter && cKey != "" {
		if cache, err := p.Config.JSONMarshal(page); nil == err {
			if err := adapter.Set(cKey, string(cache)); err != nil {
				log.Println(err)
			}
		}
	}

	return page
}

func (r *resContext) BuildWhereClause() (statement string, args []interface{}) {
	p := r.Pagination
	query := r.Statement
	p.Config = defaultConfig(p.Config)
	p.Config.Statement = query.Statement
	if p.Config.DefaultSize == 0 {
		p.Config.DefaultSize = 10
	}

	if p.Config.FieldWrapper == "" && p.Config.ValueWrapper == "" {
		defaultWrapper := "LOWER(%s)"
		wrappers := map[string]string{
			"sqlite":   defaultWrapper,
			"mysql":    defaultWrapper,
			"postgres": "LOWER((%s)::text)",
		}
		p.Config.FieldWrapper = defaultWrapper
		if wrapper, ok := wrappers[query.Dialector.Name()]; ok {
			p.Config.FieldWrapper = wrapper
		}
	}

	pr, err := parseRequest(r.Parameter, *p.Config)
	if err != nil {
		r.error = err
		return "", nil
	}
	causes := createCauses(pr)
	return causes.WhereString, causes.Params
}

// New Pagination instance
func New(params ...interface{}) *Pagination {
	if len(params) >= 1 {
		var config *Config
		for _, param := range params {
			c, isConfig := param.(*Config)
			if isConfig {
				config = c
				continue
			}
		}

		return &Pagination{Config: defaultConfig(config)}
	}

	return &Pagination{Config: defaultConfig(nil)}
}

// parseRequest func
func parseRequest(param IParameter, config Config) (pageRequest, error) {
	pr := pageRequest{
		Config: *defaultConfig(&config),
	}
	err := parsingQueryString(param, &pr)
	if err != nil {
		return pageRequest{}, err
	}
	return pr, nil
}

// createFilters func
func createFilters(filterParams interface{}, p *pageRequest) error {
	s, ok2 := filterParams.(string)
	if reflect.ValueOf(filterParams).Kind() == reflect.Slice {
		f, err := sliceutils.ConvertAny2Interface(filterParams)
		if err != nil {
			return errors.WithStack(err)
		}
		p.Filters = arrayToFilter(f, p.Config)
		p.Filters.Fields = p.Fields
	} else if ok2 {
		iface := []interface{}{}
		if e := p.Config.JSONUnmarshal([]byte(s), &iface); nil == e && len(iface) > 0 {
			p.Filters = arrayToFilter(iface, p.Config)
		}
		p.Filters.Fields = p.Fields
	}
	return nil
}

// createCauses func
func createCauses(p pageRequest) requestQuery {
	query := requestQuery{}
	wheres, params := generateWhereCauses(p.Filters, p.Config)
	sorts := []sortOrder{}

	for _, so := range p.Sorts {
		so.Column = fieldName(so.Column)
		if nil != p.Config.Statement {
			so.Column = p.Config.Statement.Quote(so.Column)
		}
		sorts = append(sorts, so)
	}

	query.Limit = p.Size
	query.Offset = p.Page * p.Size
	query.Wheres = wheres
	query.WhereString = strings.Join(wheres, " ")
	query.Sorts = sorts
	query.Params = params

	return query
}

func parsingQueryString(param IParameter, p *pageRequest) error {
	p.Size = param.GetSize()

	if p.Size == 0 {
		if p.Config.DefaultSize > 0 {
			p.Size = p.Config.DefaultSize
		} else {
			p.Size = 10
		}
	}

	p.Page = param.GetPage()

	if param.GetSort() != "" {
		sorts := strings.Split(param.GetSort(), ",")
		for _, col := range sorts {
			if col == "" {
				continue
			}

			so := sortOrder{
				Column:    col,
				Direction: "ASC",
			}
			if strings.ToUpper(param.GetOrder()) == "DESC" {
				so.Direction = "DESC"
			}

			if string(col[0]) == "-" {
				so.Column = string(col[1:])
				so.Direction = "DESC"
			}

			p.Sorts = append(p.Sorts, so)
		}
	}

	if param.GetFields() != "" {
		re := regexp.MustCompile(`[^A-z0-9_\.,]+`)
		if fields := strings.Split(param.GetFields(), ","); len(fields) > 0 {
			for i := range fields {
				fieldByte := re.ReplaceAll([]byte(fields[i]), []byte(""))
				if field := string(fieldByte); field != "" {
					p.Fields = append(p.Fields, field)
				}
			}
		}
	}

	return createFilters(param.GetFilters(), p)
}

//gocyclo:ignore
func arrayToFilter(arr []interface{}, config Config) pageFilters {
	filters := pageFilters{
		Single: false,
	}

	operatorEscape := regexp.MustCompile(`[^A-z=\<\>\-\+\^/\*%&! ]+`)
	arrayLen := len(arr)

	if len(arr) > 0 {
		subFilters := []pageFilters{}
		for k, i := range arr {
			iface, ok := i.([]interface{})
			if ok && !filters.Single {
				subFilters = append(subFilters, arrayToFilter(iface, config))
			} else if arrayLen == 1 {
				operator, ok := i.(string)
				if ok {
					operator = operatorEscape.ReplaceAllString(operator, "")
					filters.Operator = strings.ToUpper(operator)
					filters.IsOperator = true
					filters.Single = true
				}
			} else if arrayLen == 2 {
				if k == 0 {
					if column, ok := i.(string); ok {
						filters.Column = column
						filters.Operator = "="
						filters.Single = true
					}
				} else if k == 1 {
					filters.Value = i
					if nil == i {
						filters.Operator = "IS"
					}
				}
			} else if arrayLen == 3 {
				if k == 0 {
					if column, ok := i.(string); ok {
						filters.Column = column
						filters.Single = true
					}
				} else if k == 1 {
					if operator, ok := i.(string); ok {
						operator = operatorEscape.ReplaceAllString(operator, "")
						filters.Operator = strings.ToUpper(operator)
						filters.Single = true
					}
				} else if k == 2 {
					switch filters.Operator {
					case "LIKE", "ILIKE", "NOT LIKE", "NOT ILIKE":
						escapeString := ""
						escapePattern := `(%|\\)`
						if nil != config.Statement {
							driverName := config.Statement.Dialector.Name()
							switch driverName {
							case "sqlite", "sqlserver", "postgres":
								escapeString = `\`
								filters.ValueSuffix = "ESCAPE '\\'"
							case "mysql":
								escapeString = `\`
								filters.ValueSuffix = `ESCAPE '\\'`
							}
						}
						value := fmt.Sprintf("%v", i)
						re := regexp.MustCompile(escapePattern)
						value = string(re.ReplaceAll([]byte(value), []byte(escapeString+`$1`)))
						if config.SmartSearch {
							re := regexp.MustCompile(`[\s]+`)
							byt := re.ReplaceAll([]byte(value), []byte("%"))
							value = string(byt)
						}
						filters.Value = fmt.Sprintf("%s%s%s", "%", value, "%")
					default:
						filters.Value = i
					}
				}
			}
		}
		if len(subFilters) > 0 {
			separatedSubFilters := []pageFilters{}
			hasOperator := false
			defaultOperator := config.Operator
			if "" == defaultOperator {
				defaultOperator = "OR"
			}
			for k, s := range subFilters {
				if s.IsOperator && len(subFilters) == (k+1) {
					break
				}
				if !hasOperator && !s.IsOperator && k > 0 {
					separatedSubFilters = append(separatedSubFilters, pageFilters{
						Operator:   defaultOperator,
						IsOperator: true,
						Single:     true,
					})
				}
				hasOperator = s.IsOperator
				separatedSubFilters = append(separatedSubFilters, s)
			}
			filters.Value = separatedSubFilters
			filters.Single = false
		}
	}

	return filters
}

//gocyclo:ignore
func generateWhereCauses(f pageFilters, config Config) ([]string, []interface{}) {
	wheres := []string{}
	params := []interface{}{}

	if !f.Single && !f.IsOperator {
		ifaces, ok := f.Value.([]pageFilters)
		if ok && len(ifaces) > 0 {
			wheres = append(wheres, "(")
			hasOpen := false
			for _, i := range ifaces {
				subs, isSub := i.Value.([]pageFilters)
				regular, isNotSub := i.Value.(pageFilters)
				if isSub && len(subs) > 0 {
					wheres = append(wheres, "(")
					for _, s := range subs {
						subWheres, subParams := generateWhereCauses(s, config)
						wheres = append(wheres, subWheres...)
						params = append(params, subParams...)
					}
					wheres = append(wheres, ")")
				} else if isNotSub {
					subWheres, subParams := generateWhereCauses(regular, config)
					wheres = append(wheres, subWheres...)
					params = append(params, subParams...)
				} else {
					if !hasOpen && !i.IsOperator {
						wheres = append(wheres, "(")
						hasOpen = true
					}
					subWheres, subParams := generateWhereCauses(i, config)
					wheres = append(wheres, subWheres...)
					params = append(params, subParams...)
				}
			}
			if hasOpen {
				wheres = append(wheres, ")")
			}
			wheres = append(wheres, ")")
		}
	} else if f.Single {
		if f.IsOperator {
			wheres = append(wheres, f.Operator)
		} else {
			fname := fieldName(f.Column)
			if nil != config.Statement {
				fname = config.Statement.Quote(fname)
			}
			switch f.Operator {
			case "IS", "IS NOT":
				if nil == f.Value {
					wheres = append(wheres, fname, f.Operator, "NULL")
				} else {
					if strValue, isStr := f.Value.(string); isStr && strings.ToLower(strValue) == "null" {
						wheres = append(wheres, fname, f.Operator, "NULL")
					} else {
						wheres = append(wheres, fname, f.Operator, "?")
						params = append(params, f.Value)
					}
				}
			case "BETWEEN":
				if values, ok := f.Value.([]interface{}); ok && len(values) >= 2 {
					wheres = append(wheres, "(", fname, f.Operator, "? AND ?", ")")
					params = append(params, valueFixer(values[0]), valueFixer(values[1]))
				}
			case "IN", "NOT IN":
				if values, ok := f.Value.([]interface{}); ok {
					wheres = append(wheres, fname, f.Operator, "?")
					params = append(params, valueFixer(values))
				}
			case "LIKE", "NOT LIKE", "ILIKE", "NOT ILIKE":
				if config.FieldWrapper != "" {
					fname = fmt.Sprintf(config.FieldWrapper, fname)
				}
				wheres = append(wheres, fname, f.Operator, "?")
				if f.ValueSuffix != "" {
					wheres = append(wheres, f.ValueSuffix)
				}
				value, isStrValue := f.Value.(string)
				if isStrValue {
					if config.ValueWrapper != "" {
						value = fmt.Sprintf(config.ValueWrapper, value)
					} else {
						value = strings.ToLower(value)
					}
					params = append(params, value)
				} else {
					params = append(params, f.Value)
				}
			default:
				wheres = append(wheres, fname, f.Operator, "?")
				params = append(params, valueFixer(f.Value))
			}
		}
	}

	return wheres, params
}

func valueFixer(n interface{}) interface{} {
	var values []interface{}
	if rawValues, ok := n.([]interface{}); ok {
		for i := range rawValues {
			values = append(values, valueFixer(rawValues[i]))
		}

		return values
	}
	if nil != n && reflect.TypeOf(n).Name() == "float64" {
		strValue := fmt.Sprintf("%v", n)
		if match, e := regexp.Match(`^[0-9]+$`, []byte(strValue)); nil == e && match {
			v, err := strconv.ParseInt(strValue, 10, 64)
			if nil == err {
				return v
			}
		}
	}

	return n
}

func contains(source []string, value string) bool {
	found := false
	for i := range source {
		if source[i] == value {
			found = true
			break
		}
	}

	return found
}

func FieldAs(tableName, colName string) string {
	return fmt.Sprintf("%s_%s", strings.ToLower(tableName), strings.ToLower(colName))
}

func GetLowerColNameFromAlias(alias, tableName string) string {
	return strings.TrimPrefix(alias, strings.ToLower(tableName)+"_")
}

func fieldName(field string) string {
	if strings.HasPrefix(field, `"`) && strings.HasSuffix(field, `"`) {
		return field
	}
	return `"` + strings.ToLower(field) + `"`
}

// Config for customize pagination result
type Config struct {
	Operator             string
	FieldWrapper         string
	ValueWrapper         string
	DefaultSize          int64
	SmartSearch          bool
	Statement            *gorm.Statement `json:"-"`
	CustomParamEnabled   bool
	SortParams           []string
	PageParams           []string
	OrderParams          []string
	SizeParams           []string
	FilterParams         []string
	FieldsParams         []string
	FieldSelectorEnabled bool
	CacheAdapter         *gocache.AdapterInterface              `json:"-"`
	JSONMarshal          func(v interface{}) ([]byte, error)    `json:"-"`
	JSONUnmarshal        func(data []byte, v interface{}) error `json:"-"`
}

// pageFilters struct
type pageFilters struct {
	Column      string
	Operator    string
	Value       interface{}
	ValuePrefix string
	ValueSuffix string
	Single      bool
	IsOperator  bool
	Fields      []string
}

// Page result wrapper
type Page struct {
	Items      []interface{} `json:"items"`
	Page       int64         `json:"page"`
	Size       int64         `json:"size"`
	MaxPage    int64         `json:"max_page"`
	TotalPages int64         `json:"total_pages"`
	Total      int64         `json:"total"`
	Last       bool          `json:"last"`
	First      bool          `json:"first"`
	Visible    int64         `json:"visible"`
}

type IParameter interface {
	GetPage() int64
	GetSize() int64
	GetSort() string
	GetOrder() string
	GetFields() string
	GetFilters() interface{}
	IParameterInstance()
}

// query struct
type requestQuery struct {
	WhereString string
	Wheres      []string
	Params      []interface{}
	Sorts       []sortOrder
	Limit       int64
	Offset      int64
}

// pageRequest struct
type pageRequest struct {
	Size    int64
	Page    int64
	Sorts   []sortOrder
	Filters pageFilters
	Config  Config `json:"-"`
	Fields  []string
}

// sortOrder struct
type sortOrder struct {
	Column    string
	Direction string
}

func createCacheKey(cachePrefix string, pr pageRequest) string {
	key := ""
	if bte, err := pr.Config.JSONMarshal(pr); nil == err && cachePrefix != "" {
		key = fmt.Sprintf("%s%x", cachePrefix, md5.Sum(bte))
	}

	return key
}

func defaultConfig(c *Config) *Config {
	if nil == c {
		return &Config{
			JSONMarshal:   json.Marshal,
			JSONUnmarshal: json.Unmarshal,
		}
	}

	if nil == c.JSONMarshal {
		c.JSONMarshal = json.Marshal
	}

	if nil == c.JSONUnmarshal {
		c.JSONUnmarshal = json.Unmarshal
	}

	return c
}
