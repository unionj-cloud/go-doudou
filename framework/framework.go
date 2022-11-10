package framework

import (
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"os"
)

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

func CheckDev() bool {
	return stringutils.IsEmpty(os.Getenv("GDD_ENV")) || os.Getenv("GDD_ENV") == "dev"
}
