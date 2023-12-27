/*
 * MIT License
 *
 * Copyright (c) 2021 zeromicro
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 */

package parser

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"io/ioutil"
	"path/filepath"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/dbvendor/mysql/parser/console"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/dbvendor/mysql/parser/gen"
)

var (
	empty []*Table
)

// Parser is the syntax entry to parse sql as AST, you can use NewParser to create
// an instance with options, WithDebugMode option can parse sql with debug, WithLogger
// option can print logs while parsing.
type Parser struct {
	antlr.DefaultErrorListener
	debug  bool
	logger console.Console
	prefix string
}

// Option is the alias of function.
type Option func(p *Parser)

// Acceptor is the alias of function
type Acceptor func(p *gen.MySqlParser, visitor *visitor) interface{}

// NewParser creates an instance of Parser.
func NewParser(options ...Option) *Parser {
	p := &Parser{}
	for _, opt := range options {
		opt(p)
	}

	if p.logger == nil {
		p.logger = console.NewColorConsole()
	}

	return p
}

// SyntaxError overrides SyntaxError from antlr.DefaultErrorListener, which could catch error from this function, and panic,
// the parser would catch the panic by testMysqlSyntax and returns.
func (p *Parser) SyntaxError(_ antlr.Recognizer, _ interface{}, line, column int, msg string, _ antlr.RecognitionException) {
	str := fmt.Sprintf(`%s line %d:%d  %s`, p.prefix, line, column, msg)
	if p.debug {
		p.logger.Error(str)
	}

	panic(str)
}

// WithDebugMode is a Parser option to set debug mode.
func WithDebugMode(debug bool) Option {
	return func(p *Parser) {
		p.debug = debug
	}
}

// WithConsole is a Parser option to set console.
func WithConsole(logger console.Console) Option {
	return func(p *Parser) {
		p.logger = logger
	}
}

func (p *Parser) From(filename string) (ret []*Table, err error) {
	if !filepath.IsAbs(filename) {
		return nil, fmt.Errorf("%s is not a valid path", filename)
	}

	defer func() {
		p := recover()
		if p != nil {
			switch e := p.(type) {
			case error:
				err = e
			default:
				err = fmt.Errorf("%+v", p)
			}
		}
	}()

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	prefix := filepath.Base(filename)
	p.prefix = prefix
	inputStream := antlr.NewInputStream(string(bytes))
	caseChangingStream := newCaseChangingStream(inputStream, true)
	lexer := gen.NewMySqlLexer(caseChangingStream)
	lexer.RemoveErrorListeners()
	tokens := antlr.NewCommonTokenStream(lexer, antlr.LexerDefaultTokenChannel)
	mysqlParser := gen.NewMySqlParser(tokens)
	mysqlParser.RemoveErrorListeners()
	mysqlParser.AddErrorListener(p)

	visitor := &visitor{
		prefix: prefix,
		debug:  p.debug,
		logger: p.logger,
	}
	v := mysqlParser.Root().Accept(visitor)
	if v == nil {
		return empty, nil
	}

	createTables, ok := v.([]*CreateTable)
	if !ok {
		return empty, nil
	}

	for _, e := range createTables {
		ret = append(ret, e.Convert())
	}

	return
}

// testMysqlSyntax tests the mysql syntax with unit test.
func (p *Parser) testMysqlSyntax(prefix string, acceptor Acceptor, sql string) (v interface{}, err error) {
	defer func() {
		p := recover()
		if p != nil {
			switch e := p.(type) {
			case error:
				err = e
			default:
				err = fmt.Errorf("%+v", p)
			}
		}
	}()

	p.prefix = prefix
	inputStream := antlr.NewInputStream(sql)
	caseChangingStream := newCaseChangingStream(inputStream, true)
	lexer := gen.NewMySqlLexer(caseChangingStream)
	lexer.RemoveErrorListeners()
	tokens := antlr.NewCommonTokenStream(lexer, antlr.LexerDefaultTokenChannel)
	mysqlParser := gen.NewMySqlParser(tokens)
	mysqlParser.RemoveErrorListeners()
	mysqlParser.AddErrorListener(p)

	visitor := &visitor{
		prefix: prefix,
		debug:  p.debug,
		logger: p.logger,
	}
	v = acceptor(mysqlParser, visitor)
	return
}

func (p *Parser) ParseDDL(sql string) (ret interface{}, err error) {
	if stringutils.IsEmpty(sql) {
		return nil, errors.New("Parameter sql should not be empty")
	}

	defer func() {
		p := recover()
		if p != nil {
			switch e := p.(type) {
			case error:
				err = e
			default:
				err = fmt.Errorf("%+v", p)
			}
		}
	}()

	inputStream := antlr.NewInputStream(sql)
	caseChangingStream := newCaseChangingStream(inputStream, true)
	lexer := gen.NewMySqlLexer(caseChangingStream)
	lexer.RemoveErrorListeners()
	tokens := antlr.NewCommonTokenStream(lexer, antlr.LexerDefaultTokenChannel)
	mysqlParser := gen.NewMySqlParser(tokens)
	mysqlParser.RemoveErrorListeners()
	mysqlParser.AddErrorListener(p)

	visitor := &visitor{
		debug:  p.debug,
		logger: p.logger,
	}
	ret = mysqlParser.Root().Accept(visitor)
	return ret, nil
}
