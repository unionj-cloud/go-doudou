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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/dbvendor/mysql/parser/gen"
)

func TestVisitor_VisitColumnDefinition(t *testing.T) {
	p := NewParser(WithDebugMode(true))
	accept := func(p *gen.MySqlParser, visitor *visitor) interface{} {
		definition := p.ColumnDefinition()
		ctx := definition.(*gen.ColumnDefinitionContext)
		return visitor.VisitColumnDefinition(ctx)
	}

	v, err := p.testMysqlSyntax("test.sql", accept, `bigint(20) NOT NULL DEFAULT 'test default' PRIMARY KEY COMMENT 'test comment'`)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	columnDefinition := v.(*ColumnDefinition)

	assert.Equal(t, ColumnConstraint{
		NotNull:         true,
		HasDefaultValue: true,
		Primary:         true,
		Comment:         "test comment",
	}, *columnDefinition.ColumnConstraint)

	v, err = p.testMysqlSyntax("test.sql", accept, `bigint(20) NULL KEY`)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	columnDefinition = v.(*ColumnDefinition)

	assert.Equal(t, ColumnConstraint{
		Key: true,
	}, *columnDefinition.ColumnConstraint)

	v, err = p.testMysqlSyntax("test.sql", accept, `bigint(20) NULL AUTO_INCREMENT UNIQUE KEY`)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	columnDefinition = v.(*ColumnDefinition)

	assert.Equal(t, ColumnConstraint{
		AutoIncrement:   true,
		Unique:          true,
		HasDefaultValue: false,
	}, *columnDefinition.ColumnConstraint)

	v, err = p.testMysqlSyntax("test.sql", accept, `bigint(20) NULL DEFAULT NULL AUTO_INCREMENT UNIQUE KEY`)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	columnDefinition = v.(*ColumnDefinition)

	assert.Equal(t, ColumnConstraint{
		AutoIncrement: true,
		Unique:        true,
	}, *columnDefinition.ColumnConstraint)

	v, err = p.testMysqlSyntax("test.sql", accept, `varchar(20) DEFAULT '' AUTO_INCREMENT UNIQUE KEY`)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	columnDefinition = v.(*ColumnDefinition)

	assert.Equal(t, ColumnConstraint{
		HasDefaultValue: true,
		AutoIncrement:   true,
		Unique:          true,
	}, *columnDefinition.ColumnConstraint)
}
