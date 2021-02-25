package es

import (
	"github.com/unionj-cloud/papilio/kit/stringutils"
	"gopkg.in/olivere/elastic.v5"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
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

var G_EsClient *elastic.Client

func InitG_EsClient(hosts []string, username, password string) {
	errorlog := log.New(os.Stdout, "nlp", log.LstdFlags)
	var err error
	G_EsClient, err = elastic.NewSimpleClient(
		elastic.SetErrorLog(errorlog),
		elastic.SetURL(hosts...),
		elastic.SetBasicAuth(username, password))
	if err != nil {
		log.Fatal(err)
	}
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
		for field, value := range qc.Pair {
			if len(value) == 0 {
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
			}
		}
	}
	if hasShould {
		boolQuery.MinimumNumberShouldMatch(1)
	}
	return boolQuery
}

func SetupSubTest(index string, t *testing.T) func(t *testing.T) {
	t.Log("setup sub test")
	mapping := NewMapping(MappingPayload{
		Base{
			Index: index,
			Type:  index,
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
	_, err := NewIndex(index, mapping)
	if err != nil {
		panic(err)
	}
	return func(t *testing.T) {
		t.Log("teardown sub test")
		err := DeleteIndex(index)
		if err != nil {
			panic(err)
		}
	}
}
