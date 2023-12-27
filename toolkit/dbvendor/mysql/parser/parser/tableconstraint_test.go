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
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/dbvendor/mysql/parser/gen"
)

func TestVisitor_VisitTableConstraint(t *testing.T) {
	p := NewParser(WithDebugMode(true))
	accept := func(p *gen.MySqlParser, visitor *visitor) interface{} {
		ctx := p.TableConstraint()
		return visitor.visitTableConstraint(ctx)
	}

	t.Run("uniqueKeyTableConstraint", func(t *testing.T) {
		v, err := p.testMysqlSyntax("test.sql", accept, "UNIQUE INDEX `data_UNIQUE` (`data` ASC)")
		assert.Nil(t, err)
		tc, ok := v.(*TableConstraint)
		assert.True(t, ok)
		assertEqualStringSlice(t, []string{"data"}, tc.ColumnUniqueKey)

		v, err = p.testMysqlSyntax("test.sql", accept, "UNIQUE INDEX `data_UNIQUE` (`data` ASC) INVISIBLE VISIBLE")
		assert.Nil(t, err)
		tc, ok = v.(*TableConstraint)
		assert.True(t, ok)
		assertEqualStringSlice(t, []string{"data"}, tc.ColumnUniqueKey)

		v, err = p.testMysqlSyntax("test.sql", accept, "UNIQUE INDEX `data_UNIQUE` (`data` ASC) INVISIBLE VISIBLE")
		assert.Nil(t, err)
		tc, ok = v.(*TableConstraint)
		assert.True(t, ok)
		assertEqualStringSlice(t, []string{"data"}, tc.ColumnUniqueKey)

		v, err = p.testMysqlSyntax("test.sql", accept, "UNIQUE INDEX `data_UNIQUE` (`column1` ASC, `column2`) INVISIBLE VISIBLE")
		assert.Nil(t, err)
		tc, ok = v.(*TableConstraint)
		assert.True(t, ok)
		assertEqualStringSlice(t, []string{"column1", "column2"}, tc.ColumnUniqueKey)
	})

	t.Run("primaryKeyTableConstraint", func(t *testing.T) {
		v, err := p.testMysqlSyntax("test.sql", accept, "PRIMARY KEY (`description_id`)")
		assert.Nil(t, err)
		tc, ok := v.(*TableConstraint)
		assert.True(t, ok)
		assertEqualStringSlice(t, []string{"description_id"}, tc.ColumnPrimaryKey)
	})

}

func assertEqualStringSlice(t *testing.T, expected, actual []string) {
	sort.Strings(expected)
	sort.Strings(actual)
	assert.Equal(t, expected, actual)
}
