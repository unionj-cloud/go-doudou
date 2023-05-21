package gorm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/morkid/gocache"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var format = "%s doesn't match. Expected: %v, Result: %v"

func TestGetNetHttp(t *testing.T) {
	size := 20
	page := 1
	sort := "user.name,-id"
	avg := "seventy %"

	queryFilter := fmt.Sprintf(`[["user.average_point","like","%s"]]`, avg)
	query := fmt.Sprintf(`page=%d&size=%d&sort=%s&filters=%s`, page, size, sort, url.QueryEscape(queryFilter))

	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: query,
		},
	}

	parsed := parseRequest(req, Config{})
	if parsed.Size != size {
		t.Errorf(format, "Size", size, parsed.Size)
	}
	if parsed.Page != page {
		t.Errorf(format, "Page", page, parsed.Page)
	}
	if len(parsed.Sorts) != 2 {
		t.Errorf(format, "Sort length", 2, len(parsed.Sorts))
	} else {
		if parsed.Sorts[0].Column != "user.name" {
			t.Errorf(format, "Sort field 0", "user.name", parsed.Sorts[0].Column)
		}
		if parsed.Sorts[0].Direction != "ASC" {
			t.Errorf(format, "Sort direction 0", "ASC", parsed.Sorts[0].Direction)
		}
		if parsed.Sorts[1].Column != "id" {
			t.Errorf(format, "Sort field 1", "id", parsed.Sorts[1].Column)
		}
		if parsed.Sorts[1].Direction != "DESC" {
			t.Errorf(format, "Sort direction 1", "DESC", parsed.Sorts[1].Direction)
		}
	}

	filters, ok := parsed.Filters.Value.([]pageFilters)
	if ok {
		if filters[0].Column != "user.average_point" {
			t.Errorf(format, "Filter field for user.average_point", "user.average_point", filters[0].Column)
		}
		if filters[0].Operator != "LIKE" {
			t.Errorf(format, "Filter operator for user.average_point", "LIKE", filters[0].Operator)
		}
		value, isValid := filters[0].Value.(string)
		expected := "%" + avg + "%"
		if !isValid || value != expected {
			t.Errorf(format, "Filter operator for user.average_point", expected, value)
		}
	} else {
		t.Log(parsed.Filters)
		t.Errorf(format, "pageFilters class", "paginate.pageFilters", "null")
	}
}

func TestPostNetHttp(t *testing.T) {
	size := 20
	page := 1
	sort := "user.name,-id"
	avg := "seventy %"

	data := `
		{
			"page": "%d",
			"size": "%d",
			"sort": "%s",
			"filters": %s
		}
	`

	queryFilter := fmt.Sprintf(`[["user.average_point","like","%s"]]`, avg)
	query := fmt.Sprintf(data, page, size, sort, queryFilter)

	body := ioutil.NopCloser(bytes.NewReader([]byte(query)))

	req := &http.Request{
		Method: "POST",
		Body:   body,
	}

	parsed := parseRequest(req, Config{})
	if parsed.Size != size {
		t.Errorf(format, "Size", size, parsed.Size)
	}
	if parsed.Page != page {
		t.Errorf(format, "Page", page, parsed.Page)
	}
	if len(parsed.Sorts) != 2 {
		t.Errorf(format, "Sort length", 2, len(parsed.Sorts))
	} else {
		if parsed.Sorts[0].Column != "user.name" {
			t.Errorf(format, "Sort field 0", "user.name", parsed.Sorts[0].Column)
		}
		if parsed.Sorts[0].Direction != "ASC" {
			t.Errorf(format, "Sort direction 0", "ASC", parsed.Sorts[0].Direction)
		}
		if parsed.Sorts[1].Column != "id" {
			t.Errorf(format, "Sort field 1", "id", parsed.Sorts[1].Column)
		}
		if parsed.Sorts[1].Direction != "DESC" {
			t.Errorf(format, "Sort direction 1", "DESC", parsed.Sorts[1].Direction)
		}
	}

	filters, ok := parsed.Filters.Value.([]pageFilters)
	if ok {
		if filters[0].Column != "user.average_point" {
			t.Errorf(format, "Filter field for user.average_point", "user.average_point", filters[0].Column)
		}
		if filters[0].Operator != "LIKE" {
			t.Errorf(format, "Filter operator for user.average_point", "LIKE", filters[0].Operator)
		}
		value, isValid := filters[0].Value.(string)
		expected := "%" + avg + "%"
		if !isValid || value != expected {
			t.Errorf(format, "Filter operator for user.average_point", expected, value)
		}
	} else {
		t.Log(parsed.Filters)
		t.Errorf(format, "pageFilters class", "paginate.pageFilters", "null")
	}
}

func TestPaginate(t *testing.T) {
	type User struct {
		gorm.Model
		Name         string `json:"name"`
		AveragePoint string `json:"average_point"`
	}

	type Article struct {
		gorm.Model
		Title   string `json:"title"`
		Content string `json:"content"`
		UserID  uint   `json:"-"`
		User    User   `json:"user"`
	}

	// dsn := "host=127.0.0.1 port=5433 user=postgres password=postgres dbname=postgres sslmode=disable TimeZone=Asia/Jakarta"
	// dsn := "gorm.db"
	dsn := "file::memory:?cache=shared"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	db.AutoMigrate(&User{}, &Article{})

	users := []User{{Name: "John doe", AveragePoint: "Seventy %"}, {Name: "Jane doe", AveragePoint: "one hundred %"}}
	articles := []Article{}
	articles = append(articles, Article{Title: "Written by john", Content: "Example by john", User: users[0]})
	articles = append(articles, Article{Title: "Written by jane", Content: "Example by jane", User: users[1]})

	if nil != err {
		t.Error(err.Error())
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&users).Error; nil != err {
		tx.Rollback()
		t.Error(err.Error())
		return
	} else if err := tx.Create(&articles).Error; nil != err {
		tx.Rollback()
		t.Error(err.Error())
		return
	} else if err := tx.Commit().Error; nil != err {
		tx.Rollback()
		t.Error(err.Error())
		return
	}

	size := 1
	page := 0
	sort := "user.name,-id"
	avg := "y %"
	data := "page=%d&size=%d&sort=%s&filters=%s"

	queryFilter := fmt.Sprintf(`[["user.average_point","like","%s"],["AND"],["user.name","IS NOT",null],["id","like","1"]]`, avg)
	query := fmt.Sprintf(data, page, size, sort, url.QueryEscape(queryFilter))

	request := &http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: query,
		},
	}
	response := []Article{}

	model := db.Joins("User").Model(&Article{})
	result := New().Response(model, request, &response)

	str, err := json.MarshalIndent(result, "", "  ")
	if nil == err {
		t.Log(string(str))
	}
}

type noOpAdapter struct {
	keyValues          map[string]string
	T                  *testing.T
	clearCounter       int
	clearPrefixCounter int
}

func (n *noOpAdapter) Get(key string) (string, error) {
	n.T.Log(key)
	if v, ok := n.keyValues[key]; ok {
		n.T.Log("OK, Cache found! serving data from cache")
		return v, nil
	}

	n.T.Log("Cache not found")

	return "", errors.New("Cache not found")
}
func (n *noOpAdapter) Set(key string, value string) error {
	if _, ok := n.keyValues[key]; !ok {
		n.keyValues = map[string]string{}
	}
	n.keyValues[key] = value
	n.T.Log("Writing cache")
	return nil
}
func (n *noOpAdapter) IsValid(key string) bool {
	if _, ok := n.keyValues[key]; ok {
		n.T.Log("Cache exists and not expired")
		return false
	}
	n.T.Log("Cache doesn't exists or expired")
	return true
}
func (n *noOpAdapter) Clear(key string) error {
	return nil
}
func (n *noOpAdapter) ClearPrefix(keyPrefix string) error {
	if n.clearPrefixCounter > 2 {
		return errors.New("maximum clear")
	}
	n.clearPrefixCounter = n.clearPrefixCounter + 1
	return nil
}
func (n *noOpAdapter) ClearAll() error {
	if n.clearCounter > 0 {
		return errors.New("maximum clear")
	}
	n.clearCounter = n.clearCounter + 1
	return nil
}

func TestCache(t *testing.T) {
	type User struct {
		gorm.Model
		Name         string `json:"name"`
		AveragePoint string `json:"average_point"`
	}

	type Category struct {
		gorm.Model
		Name string `json:"name"`
	}

	type Article struct {
		gorm.Model
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		UserID     uint     `json:"-"`
		CategoryID uint     `json:"-"`
		User       User     `json:"user"`
		Category   Category `json:"category"`
	}
	dsn := "file::memory:"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if nil != err {
		t.Error(err.Error())
		return
	}
	db.AutoMigrate(&User{}, &Article{})
	categories := []Category{{Name: "Blog"}}
	users := []User{{Name: "John doe", AveragePoint: "Seventy %"}, {Name: "Jane doe", AveragePoint: "one hundred %"}}
	articles := []Article{}
	articles = append(articles, Article{Title: "Written by john", Content: "Example by john", UserID: 1, CategoryID: 1})
	articles = append(articles, Article{Title: "Written by jane", Content: "Example by jane", UserID: 2, CategoryID: 1})
	db.Create(&categories)
	db.Create(&users)
	db.Create(&articles)
	request := &http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: "page=0&size=10&fields=id",
		},
	}

	var adapter gocache.AdapterInterface = &noOpAdapter{T: t}
	config := &Config{
		CacheAdapter:         &adapter,
		FieldSelectorEnabled: true,
	}
	pg := New(config)
	// set cache
	model1 := db.Joins("User").Model(&Article{}).Preload(`Category`)
	page1 := pg.With(model1).
		Request(request).
		Fields([]string{"id"}).
		Cache("cache_prefix").
		Response(&[]Article{})

	// get cache
	var cached []Article
	model2 := db.Joins("User").Model(&Article{})
	page2 := pg.With(model2).Request(request).Cache("cache_prefix").Response(&cached)

	if len(cached) < 1 {
		t.Error("Cache pointer not working perfectly")
	}

	if page1.Total != page2.Total {
		t.Error("Total doesn't match")
	}

	pg.ClearCache("cache", "cache_")
	pg.ClearCache("cache", "cache_")
	pg.ClearAllCache()
	pg.ClearAllCache()
}
