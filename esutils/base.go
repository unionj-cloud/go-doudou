package esutils

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/constants"
	"github.com/unionj-cloud/go-doudou/logutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/test"
	"strings"
	"sync"
	"time"
)

var lock sync.Mutex

type queryType int
type queryLogic int

const (
	SHOULD queryLogic = iota + 1
	MUST
	MUSTNOT
)

const (
	TERMS queryType = iota + 1
	MATCHPHRASE
	RANGE
	PREFIX
	WILDCARD
	EXISTS
)

type esFieldType string

const (
	TEXT    esFieldType = "text"
	KEYWORD esFieldType = "keyword"
	DATE    esFieldType = "date"
	LONG    esFieldType = "long"
	INTEGER esFieldType = "integer"
	SHORT   esFieldType = "short"
	DOUBLE  esFieldType = "double"
	FLOAT   esFieldType = "float"
	BOOL    esFieldType = "boolean"
)

type Es struct {
	client *elastic.Client
	esIndex  string
	esType   string
	username string
	password string
	urls     []string
	logger   *logrus.Logger
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

type EsOption func(*Es)

func WithClient(client *elastic.Client) EsOption {
	return func(es *Es) {
		es.client = client
	}
}

func WithUsername(username string) EsOption {
	return func(es *Es) {
		es.username = username
	}
}

func WithPassword(password string) EsOption {
	return func(es *Es) {
		es.password = password
	}
}

func WithLogger(logger *logrus.Logger) EsOption {
	return func(es *Es) {
		es.logger = logger
	}
}

func WithUrls(urls []string) EsOption {
	return func(es *Es) {
		es.urls = urls
	}
}

func NewEs(esIndex, esType string, opts ...EsOption) *Es {
	es := &Es{
		esIndex: esIndex,
		esType:  esType,
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

type IBase interface {
	GetIndex() string
	GetType() string
	SetType(s string)
}

type Base struct {
	Index string `json:"index"`
	Type  string `json:"type"`
}

func (b *Base) GetIndex() string {
	return b.Index
}

func (b *Base) GetType() string {
	return b.Type
}

func (b *Base) SetType(s string) {
	b.Type = s
}

type Field struct {
	Name   string      `json:"name"`
	Type   esFieldType `json:"type"`
	Format string      `json:"format"`
}

type QueryCond struct {
	Pair       map[string][]interface{} `json:"pair"`
	QueryLogic queryLogic               `json:"query_logic"`
	QueryType  queryType                `json:"query_type"`
	Children   []QueryCond              `json:"children"`
}

type Sort struct {
	Field     string
	Ascending bool
}

type Paging struct {
	StartDate  string      `json:"start_date"`
	EndDate    string      `json:"end_date"`
	DateField  string      `json:"date_field"`
	QueryConds []QueryCond `json:"query_conds"`
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
			termsQuery := elastic.NewTermsQuery(field, value...)
			if qc.QueryLogic == SHOULD {
				boolQuery.Should(termsQuery)
			} else if qc.QueryLogic == MUST {
				boolQuery.Must(termsQuery)
			} else if qc.QueryLogic == MUSTNOT {
				boolQuery.MustNot(termsQuery)
			}
		} else if qc.QueryType == RANGE {
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
		} else if qc.QueryType == MATCHPHRASE {
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
		} else if qc.QueryType == PREFIX {
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
		} else if qc.QueryType == WILDCARD {
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
		} else if qc.QueryType == EXISTS {
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

func setupSubTest() (*Es, func()) {
	es, terminator := prepareTestEnvironment()
	prepareTestIndex(es)
	prepareTestData(es)
	return es, terminator
}

func prepareTestEnvironment() (*Es, func()) {
	logger := logutils.NewLogger()
	var terminateContainer func() // variable to store function to terminate container
	var host string
	var port int
	var err error
	terminateContainer, host, port, err = test.SetupEs6Container(logger)
	if err != nil {
		logger.Panicln("failed to setup Elasticsearch container")
	}
	return NewEs("test", "test", WithLogger(logger), WithUrls([]string{fmt.Sprintf("http://%s:%d", host, port)})), terminateContainer
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

	err := es.BulkSaveOrUpdate(context.Background(), []map[string]interface{}{
		{
			"createAt": createAt1.UTC().Format(constants.FORMATES),
			"text":     "2020年7月8日11时25分，高考文科综合/理科综合科目考试将要结束时，平顶山市一中考点一考生突然情绪失控，先后抓其右边、后边考生答题卡，造成两位考生答题卡损毁。",
		},
		{
			"createAt": createAt2.UTC().Format(constants.FORMATES),
			"text":     "考场两位监考教师及时制止，并稳定了考场秩序，市一中考点按程序启用备用答题卡，按规定补足答题卡被损毁的两位考生耽误的考试时间，两位考生将损毁卡的内容誊写在新答题卡上。",
		},
		{
			"createAt": createAt3.UTC().Format(constants.FORMATES),
			"text":     "目前，我办已将损毁其他考生答题卡的考生违规情况上报河南省招生办公室，将依规对该考生进行处理。平顶山市招生考试委员会办公室",
		},
	})
	if err != nil {
		panic(err)
	}
}
