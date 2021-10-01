package esutils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"io"
	"strings"
	"time"

	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/constants"
	"github.com/unionj-cloud/go-doudou/logutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
)

//go:generate go-doudou name --file $GOFILE

type queryType int
type queryLogic int

const (
	// SHOULD represents should query
	SHOULD queryLogic = iota + 1
	// MUST represents must query
	MUST
	// MUSTNOT represents must_not query
	MUSTNOT
)

const (
	// TERMS represents terms query
	TERMS queryType = iota + 1
	// MATCHPHRASE represents match_phrase query
	MATCHPHRASE
	// RANGE represents range query
	RANGE
	// PREFIX represents prefix query
	PREFIX
	// WILDCARD https://www.elastic.co/guide/en/elasticsearch/reference/6.8/query-dsl-wildcard-query.html
	WILDCARD
	// EXISTS https://www.elastic.co/guide/en/elasticsearch/reference/6.8/query-dsl-prefix-query.html
	EXISTS
)

type esFieldType string

const (
	// TEXT represents text field type
	TEXT esFieldType = "text"
	// KEYWORD represents keyword field type
	KEYWORD esFieldType = "keyword"
	// DATE represents date field type
	DATE esFieldType = "date"
	// LONG represents long field type
	LONG esFieldType = "long"
	// INTEGER represents integer field type
	INTEGER esFieldType = "integer"
	// SHORT represents short field type
	SHORT esFieldType = "short"
	// DOUBLE represents double field type
	DOUBLE esFieldType = "double"
	// FLOAT represents float field type
	FLOAT esFieldType = "float"
	// BOOL represents bool field type
	BOOL esFieldType = "boolean"
)

// Es defines properties for connecting to an es instance
type Es struct {
	client   *elastic.Client `json:"client"`
	esIndex  string          `json:"esIndex"`
	esType   string          `json:"esType"`
	username string          `json:"username"`
	password string          `json:"password"`
	urls     []string        `json:"urls"`
	logger   *logrus.Logger  `json:"logger"`
}

// SetIndex sets index
func (e *Es) SetIndex(index string) {
	e.esIndex = index
}

// SetType sets type
func (e *Es) SetType(estype string) {
	e.esType = estype
}

func (e *Es) newDefaultClient() {
	client, err := elastic.NewSimpleClient(
		elastic.SetErrorLog(e.logger),
		elastic.SetURL(e.urls...),
		elastic.SetBasicAuth(e.username, e.password),
		elastic.SetGzip(true),
	)
	if err != nil {
		panic(fmt.Errorf("NewSimpleClient() error: %+v\n", err))
	}
	e.client = client
}

func (es *Es) fetchAll(boolQuery *elastic.BoolQuery, callback func(message json.RawMessage) (interface{}, error)) ([]interface{}, error) {
	var (
		rets []interface{}
		err  error
	)
	hits := make(chan *elastic.SearchHit)
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		defer close(hits)
		scroll := es.client.Scroll().Index(es.esIndex).Type(es.esType).Query(boolQuery).Size(1000).KeepAlive("1m")
		for {
			results, err := scroll.Do(ctx)
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return errors.Wrap(err, "call Scroll() error")
			}
			for _, hit := range results.Hits.Hits {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					hits <- hit
				}
			}
		}
		return nil
	})

	c := make(chan interface{})
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			for hit := range hits {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					var ret interface{}
					if callback == nil {
						var p map[string]interface{}
						json.Unmarshal(*hit.Source, &p)
						p["_id"] = hit.Id
						ret = p
					} else {
						if ret, err = callback(*hit.Source); err != nil {
							return errors.Wrap(err, "call callback() error")
						}
					}
					c <- ret
				}
			}
			return nil
		})
	}

	go func() {
		g.Wait()
		close(c)
	}()

	for s := range c {
		rets = append(rets, s)
	}

	if err := g.Wait(); err != nil {
		return nil, errors.Wrap(err, "call Wait() error")
	}
	return rets, nil
}

func (es *Es) doPaging(ctx context.Context, paging *Paging, boolQuery *elastic.BoolQuery, callback func(message json.RawMessage) (interface{}, error)) ([]interface{}, error) {
	var (
		rets         []interface{}
		searchResult *elastic.SearchResult
		err          error
	)
	ss := es.client.Search().Index(es.esIndex).Type(es.esType).Query(boolQuery)
	if paging.Sortby != nil && len(paging.Sortby) > 0 {
		for _, v := range paging.Sortby {
			ss = ss.Sort(v.Field, v.Ascending)
		}
	}
	if searchResult, err = ss.From(paging.Skip).Size(paging.Limit).Do(ctx); err != nil {
		return nil, errors.Wrap(err, "call Search() error")
	}
	for _, hit := range searchResult.Hits.Hits {
		var ret interface{}
		if callback == nil {
			var p map[string]interface{}
			json.Unmarshal(*hit.Source, &p)
			p["_id"] = hit.Id
			ret = p
		} else {
			if ret, err = callback(*hit.Source); err != nil {
				return nil, errors.Wrap(err, "call callback() error")
			}
		}

		rets = append(rets, ret)
	}
	return rets, nil
}

// EsOption represents functions for changing Es properties
type EsOption func(*Es)

// WithClient sets client
func WithClient(client *elastic.Client) EsOption {
	return func(es *Es) {
		es.client = client
	}
}

// WithUsername sets username
func WithUsername(username string) EsOption {
	return func(es *Es) {
		es.username = username
	}
}

// WithPassword sets password
func WithPassword(password string) EsOption {
	return func(es *Es) {
		es.password = password
	}
}

// WithLogger sets logger
func WithLogger(logger *logrus.Logger) EsOption {
	return func(es *Es) {
		es.logger = logger
	}
}

// WithUrls set urls
func WithUrls(urls []string) EsOption {
	return func(es *Es) {
		es.urls = urls
	}
}

// NewEs creates an Es instance
func NewEs(esIndex, esType string, opts ...EsOption) *Es {
	_esType := esType
	if stringutils.IsEmpty(_esType) {
		_esType = esIndex
	}
	es := &Es{
		esIndex: esIndex,
		esType:  _esType,
	}
	for _, opt := range opts {
		opt(es)
	}
	if es.logger == nil {
		es.logger = logutils.NewLogger()
	}
	if len(es.urls) == 0 && es.client == nil {
		panic("NewEs() error: you must provide urls or elastic client")
	}
	if es.client == nil {
		es.newDefaultClient()
	}
	return es
}

// IBase wraps functions for getting es index and type
type IBase interface {
	GetIndex() string
	GetType() string
	SetType(s string)
}

// Base defines es index and type
type Base struct {
	Index string `json:"index"`
	Type  string `json:"type"`
}

// GetIndex return index name
func (b *Base) GetIndex() string {
	return b.Index
}

// GetType return es type name
func (b *Base) GetType() string {
	return b.Type
}

// SetType sets es type
func (b *Base) SetType(s string) {
	b.Type = s
}

// Field defines a es field
type Field struct {
	Name   string      `json:"name"`
	Type   esFieldType `json:"type"`
	Format string      `json:"format"`
}

// QueryCond defines query conditions
type QueryCond struct {
	Pair       map[string][]interface{} `json:"pair"`
	QueryLogic queryLogic               `json:"queryLogic"`
	QueryType  queryType                `json:"queryType"`
	Children   []QueryCond              `json:"children"`
}

// Sort defines sort condition
type Sort struct {
	Field     string `json:"field"`
	Ascending bool   `json:"ascending"`
}

// Paging defines pagination query conditions
type Paging struct {
	StartDate  string      `json:"startDate"`
	EndDate    string      `json:"endDate"`
	DateField  string      `json:"dateField"`
	QueryConds []QueryCond `json:"queryConds"`
	Skip       int         `json:"skip"`
	Limit      int         `json:"limit"`
	Sortby     []Sort      `json:"sortby"`
}

func querynode(boolQuery *elastic.BoolQuery, qc QueryCond) {
	for field, value := range qc.Pair {
		if len(value) == 0 && qc.QueryType != EXISTS {
			continue
		}
		if qc.QueryType == TERMS {
			terms(boolQuery, qc, field, value)
		} else if qc.QueryType == RANGE {
			rangeQ(boolQuery, qc, field, value)
		} else if qc.QueryType == MATCHPHRASE {
			matchPhrase(boolQuery, qc, field, value)
		} else if qc.QueryType == PREFIX {
			prefixQ(boolQuery, qc, field, value)
		} else if qc.QueryType == WILDCARD {
			wildcard(boolQuery, qc, field, value)
		} else if qc.QueryType == EXISTS {
			exists(boolQuery, qc, field, value)
		}
	}
}

func exists(boolQuery *elastic.BoolQuery, qc QueryCond, field string, value []interface{}) {
	if stringutils.IsNotEmpty(field) {
		prefixQuery := elastic.NewExistsQuery(field)
		if qc.QueryLogic == SHOULD {
			boolQuery.Should(prefixQuery)
		} else if qc.QueryLogic == MUST {
			boolQuery.Must(prefixQuery)
		} else if qc.QueryLogic == MUSTNOT {
			boolQuery.MustNot(prefixQuery)
		}
	}
}

func wildcard(boolQuery *elastic.BoolQuery, qc QueryCond, field string, value []interface{}) {
	var wild string
	if len(value) > 0 && value[0] != nil {
		wild = value[0].(string)
	}
	if stringutils.IsNotEmpty(wild) {
		prefixQuery := elastic.NewWildcardQuery(field, wild)
		if qc.QueryLogic == SHOULD {
			boolQuery.Should(prefixQuery)
		} else if qc.QueryLogic == MUST {
			boolQuery.Must(prefixQuery)
		} else if qc.QueryLogic == MUSTNOT {
			boolQuery.MustNot(prefixQuery)
		}
	}
}

func prefixQ(boolQuery *elastic.BoolQuery, qc QueryCond, field string, value []interface{}) {
	var prefix string
	if len(value) > 0 && value[0] != nil {
		prefix = value[0].(string)
	}
	if stringutils.IsNotEmpty(prefix) {
		prefixQuery := elastic.NewPrefixQuery(field, prefix)
		if qc.QueryLogic == SHOULD {
			boolQuery.Should(prefixQuery)
		} else if qc.QueryLogic == MUST {
			boolQuery.Must(prefixQuery)
		} else if qc.QueryLogic == MUSTNOT {
			boolQuery.MustNot(prefixQuery)
		}
	}
}

func matchPhrase(boolQuery *elastic.BoolQuery, qc QueryCond, field string, value []interface{}) {
	bQuery := elastic.NewBoolQuery()
	for _, item := range value {
		keyword := item.(string)
		words := strings.Split(keyword, "+")
		if len(words) > 1 {
			nestedBoolQuery := elastic.NewBoolQuery()
			for _, word := range words {
				word = strings.TrimSpace(word)
				if word != "" {
					if word[0] != '-' {
						nestedBoolQuery.Must(elastic.NewMatchPhraseQuery(field, word))
					} else {
						word = word[1:]
						nestedBoolQuery.MustNot(elastic.NewMatchPhraseQuery(field, word))
					}
				}
			}
			bQuery.Should(nestedBoolQuery)
		} else {
			word := words[0]
			if word != "" {
				if word[0] != '-' {
					bQuery.Should(elastic.NewMatchPhraseQuery(field, word))
				} else {
					nestedBoolQuery := elastic.NewBoolQuery()
					word = word[1:]
					nestedBoolQuery.MustNot(elastic.NewMatchPhraseQuery(field, word))
					bQuery.Should(nestedBoolQuery)
				}
			}
		}
	}
	if qc.QueryLogic == SHOULD {
		boolQuery.Should(bQuery)
	} else if qc.QueryLogic == MUST {
		boolQuery.Must(bQuery)
	} else if qc.QueryLogic == MUSTNOT {
		boolQuery.MustNot(bQuery)
	}
}

func rangeQ(boolQuery *elastic.BoolQuery, qc QueryCond, field string, value []interface{}) {
	if paramsMap, ok := value[0].(map[string]interface{}); ok {
		rangeQuery := elastic.NewRangeQuery(field)
		if paramsMap["from"] != nil || paramsMap["to"] != nil {
			if paramsMap["from"] != nil {
				rangeQuery.From(paramsMap["from"])
			}
			if paramsMap["to"] != nil {
				rangeQuery.To(paramsMap["to"])
			}
			if paramsMap["include_lower"] != nil {
				rangeQuery.IncludeLower(paramsMap["include_lower"].(bool))
			}
			if paramsMap["include_upper"] != nil {
				rangeQuery.IncludeUpper(paramsMap["include_upper"].(bool))
			}
			if qc.QueryLogic == SHOULD {
				boolQuery.Should(rangeQuery)
			} else if qc.QueryLogic == MUST {
				boolQuery.Must(rangeQuery)
			} else if qc.QueryLogic == MUSTNOT {
				boolQuery.MustNot(rangeQuery)
			}
		}
	}
}

func terms(boolQuery *elastic.BoolQuery, qc QueryCond, field string, value []interface{}) {
	termsQuery := elastic.NewTermsQuery(field, value...)
	if qc.QueryLogic == SHOULD {
		boolQuery.Should(termsQuery)
	} else if qc.QueryLogic == MUST {
		boolQuery.Must(termsQuery)
	} else if qc.QueryLogic == MUSTNOT {
		boolQuery.MustNot(termsQuery)
	}
}

func querytree(boolQuery *elastic.BoolQuery, cond QueryCond) {
	if len(cond.Children) > 0 {
		bq := elastic.NewBoolQuery()
		for _, qc := range cond.Children {
			querytree(bq, qc)
		}
		if cond.QueryLogic == SHOULD {
			boolQuery.Should(bq)
		} else if cond.QueryLogic == MUST {
			boolQuery.Must(bq)
		} else if cond.QueryLogic == MUSTNOT {
			boolQuery.MustNot(bq)
		}
		return
	}
	querynode(boolQuery, cond)
}

func query(startDate string, endDate string, dateField string, queryConds []QueryCond) *elastic.BoolQuery {
	boolQuery := elastic.NewBoolQuery()
	if dateField != "" && startDate != "" && endDate != "" {
		boolQuery.Must(
			elastic.NewRangeQuery(dateField).
				Gte(startDate).
				Lte(endDate).
				Format("yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis").
				TimeZone("Asia/Shanghai"),
		)
	}
	var hasShould bool
	for _, qc := range queryConds {
		if !hasShould && qc.QueryLogic == SHOULD {
			hasShould = true
		}
		querytree(boolQuery, qc)
	}
	if hasShould {
		boolQuery.MinimumNumberShouldMatch(1)
	}
	return boolQuery
}

var esHost string
var esPort int

func setupSubTest(esindex string) *Es {
	es := NewEs(esindex, esindex, WithLogger(logutils.NewLogger()), WithUrls([]string{fmt.Sprintf("http://%s:%d", esHost, esPort)}))
	prepareTestIndex(es)
	prepareTestData(es)
	return es
}

func prepareTestIndex(es *Es) {
	mapping := NewMapping(MappingPayload{
		Base{
			Index: es.esIndex,
			Type:  es.esType,
		},
		[]Field{
			{
				Name: "createAt",
				Type: DATE,
			},
			{
				Name: "text",
				Type: TEXT,
			},
		},
	})
	_, err := es.NewIndex(context.Background(), mapping)
	if err != nil {
		panic(err)
	}
}

func prepareTestData(es *Es) {
	data1 := "2020-06-01"
	data2 := "2020-06-20"
	data3 := "2020-07-10"

	createAt1, _ := time.Parse(constants.FORMAT2, data1)
	createAt2, _ := time.Parse(constants.FORMAT2, data2)
	createAt3, _ := time.Parse(constants.FORMAT2, data3)

	err := es.BulkSaveOrUpdate(context.Background(), []interface{}{
		map[string]interface{}{
			"id":       "9seTXHoBNx091WJ2QCh5",
			"createAt": createAt1.UTC().Format(constants.FORMATES),
			"type":     "education",
			"text":     "2020年7月8日11时25分，高考文科综合/理科综合科目考试将要结束时，平顶山市一中考点一考生突然情绪失控，先后抓其右边、后边考生答题卡，造成两位考生答题卡损毁。",
		},
		map[string]interface{}{
			"id":       "9seTXHoBNx091WJ2QCh6",
			"createAt": createAt2.UTC().Format(constants.FORMATES),
			"type":     "sport",
			"text":     "考场两位监考教师及时制止，并稳定了考场秩序，市一中考点按程序启用备用答题卡，按规定补足答题卡被损毁的两位考生耽误的考试时间，两位考生将损毁卡的内容誊写在新答题卡上。",
		},
		map[string]interface{}{
			"id":       "9seTXHoBNx091WJ2QCh7",
			"createAt": createAt3.UTC().Format(constants.FORMATES),
			"type":     "culture",
			"text":     "目前，我办已将损毁其他考生答题卡的考生违规情况上报河南省招生办公室，将依规对该考生进行处理。平顶山市招生考试委员会办公室",
		},
	})
	if err != nil {
		panic(err)
	}
}
