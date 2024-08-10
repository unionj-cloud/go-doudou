## Benchmarks

All Benchmarks Last Run Nov 17, 2019

Run on MacBook Pro (15-inch, 2017) using go version go1.13.4 darwin/amd64
go test -run=NONE -bench=. -benchmem=true

### go-playground/form
```go
BenchmarkSimpleUserDecodeStruct-8                                4447569               255 ns/op              64 B/op          1 allocs/op
BenchmarkSimpleUserDecodeStructParallel-8                       14087551                77.2 ns/op            64 B/op          1 allocs/op
BenchmarkSimpleUserEncodeStruct-8                                1863354               645 ns/op             485 B/op         10 allocs/op
BenchmarkSimpleUserEncodeStructParallel-8                        5554753               208 ns/op             485 B/op         10 allocs/op
BenchmarkPrimitivesDecodeStructAllPrimitivesTypes-8              1345276               881 ns/op              96 B/op          1 allocs/op
BenchmarkPrimitivesDecodeStructAllPrimitivesTypesParallel-8      4729965               259 ns/op              96 B/op          1 allocs/op
BenchmarkPrimitivesEncodeStructAllPrimitivesTypes-8               303967              3331 ns/op            2977 B/op         35 allocs/op
BenchmarkPrimitivesEncodeStructAllPrimitivesTypesParallel-8      1094600              1077 ns/op            2978 B/op         35 allocs/op
BenchmarkComplexArrayDecodeStructAllTypes-8                        76928             14567 ns/op            2248 B/op        121 allocs/op
BenchmarkComplexArrayDecodeStructAllTypesParallel-8               292060              4355 ns/op            2249 B/op        121 allocs/op
BenchmarkComplexArrayEncodeStructAllTypes-8                        94536             11334 ns/op            7113 B/op        104 allocs/op
BenchmarkComplexArrayEncodeStructAllTypesParallel-8               298318              3633 ns/op            7112 B/op        104 allocs/op
BenchmarkComplexMapDecodeStructAllTypes-8                          58084             18635 ns/op            5306 B/op        130 allocs/op
BenchmarkComplexMapDecodeStructAllTypesParallel-8                 187159              5454 ns/op            5308 B/op        130 allocs/op
BenchmarkComplexMapEncodeStructAllTypes-8                         101962             11763 ns/op            6971 B/op        129 allocs/op
BenchmarkComplexMapEncodeStructAllTypesParallel-8                 312925              4185 ns/op            6970 B/op        129 allocs/op
BenchmarkDecodeNestedStruct-8                                     469940              2547 ns/op             384 B/op         14 allocs/op
BenchmarkDecodeNestedStructParallel-8                            1486963               810 ns/op             384 B/op         14 allocs/op
BenchmarkEncodeNestedStruct-8                                     796798              1501 ns/op             693 B/op         16 allocs/op
BenchmarkEncodeNestedStructParallel-8                            2290203               520 ns/op             693 B/op         16 allocs/op

```