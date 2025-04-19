package framework

import (
	"golang.org/x/exp/maps"
	"sync"

	"github.com/common-nighthawk/go-figure"
	"github.com/unionj-cloud/go-doudou/v2/framework/config"
)

var annotationStoreInstance = make(AnnotationStore)

func RegisterAnnotationStore(store AnnotationStore) {
	maps.Copy(annotationStoreInstance, store)
}

func HasAnnotation(key string, annotationName string) bool {
	for _, item := range annotationStoreInstance[key] {
		if item.Name == annotationName {
			return true
		}
	}
	return false
}

func GetAnnotation(key string, annotationName string) (Annotation, bool) {
	for _, item := range annotationStoreInstance[key] {
		if item.Name == annotationName {
			return item, true
		}
	}
	return Annotation{}, false
}

type Annotation struct {
	Name   string
	Params []string
}

type AnnotationStore map[string][]Annotation

func (receiver AnnotationStore) HasAnnotation(key string, annotationName string) bool {
	for _, item := range receiver[key] {
		if item.Name == annotationName {
			return true
		}
	}
	return false
}

func (receiver AnnotationStore) GetParams(key string, annotationName string) []string {
	for _, item := range receiver[key] {
		if item.Name == annotationName {
			return item.Params
		}
	}
	return nil
}

var PrintLock sync.Mutex

var once sync.Once

func PrintBanner() {
	once.Do(func() {
		if !config.CheckDev() {
			return
		}
		if config.GddConfig.Banner {
			figure.NewColorFigure(config.GddConfig.BannerText, "doom", "green", true).Print()
		}
	})
}
