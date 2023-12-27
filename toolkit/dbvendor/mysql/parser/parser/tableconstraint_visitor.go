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
	"strings"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/dbvendor/mysql/parser/gen"
)

type TableConstraint struct {
	// ColumnPrimaryKey describes the name of columns
	ColumnPrimaryKey []string
	// ColumnUniqueKey describes the name of columns
	ColumnUniqueKey []string
}

// visitTableConstraint visits a parse tree produced by MySqlParser#tableConstraint.
func (v *visitor) visitTableConstraint(ctx gen.ITableConstraintContext) *TableConstraint {
	v.trace("VisitTableConstraint")
	var ret TableConstraint
	switch tx := ctx.(type) {
	case *gen.PrimaryKeyTableConstraintContext:
		if tx.IndexColumnNames() != nil {
			indexColumnNamesCtx, ok := tx.IndexColumnNames().(*gen.IndexColumnNamesContext)
			if ok {
				ret.ColumnPrimaryKey = v.visitIndexColumnNames(indexColumnNamesCtx)
			}
		}
	case *gen.UniqueKeyTableConstraintContext:
		if tx.IndexColumnNames() != nil {
			indexColumnNamesCtx, ok := tx.IndexColumnNames().(*gen.IndexColumnNamesContext)
			if ok {
				ret.ColumnUniqueKey = v.visitIndexColumnNames(indexColumnNamesCtx)
			}
		}
	case *gen.ForeignKeyTableConstraintContext:
		v.panicWithExpr(tx.GetStart(), "Unsupported foreign key constraint")
	}

	return &ret
}

// visitIndexColumnNames visits a parse tree produced by MySqlParser#indexColumnNames.
func (v *visitor) visitIndexColumnNames(ctx *gen.IndexColumnNamesContext) []string {
	v.trace("VisitIndexColumnNames")
	var columns []string
	for _, e := range ctx.AllIndexColumnName() {
		indexCtx, ok := e.(*gen.IndexColumnNameContext)
		if !ok {
			continue
		}

		columns = append(columns, v.visitIndexColumnName(indexCtx))
	}

	return columns
}

// visitIndexColumnName visits a parse tree produced by MySqlParser#indexColumnName.
func (v *visitor) visitIndexColumnName(ctx *gen.IndexColumnNameContext) string {
	v.trace("VisitIndexColumnName")
	var column string
	if ctx.Uid() != nil {
		column = v.visitUid(ctx.Uid())
	} else {
		column = parseTerminalNode(
			ctx.STRING_LITERAL(),
			withTrim("`"),
			withTrim("'"),
			withReplacer("\r", "", "\n", ""),
		)
	}

	return column
}

func (v *visitor) visitUid(ctx gen.IUidContext) string {
	str := ctx.GetText()
	str = strings.Trim(str, "`")
	str = strings.Trim(str, "'")
	str = strings.NewReplacer("\r", "", "\n", "").Replace(str)
	return str
}
