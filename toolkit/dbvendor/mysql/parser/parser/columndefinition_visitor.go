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

type ColumnOpsType int

const (
	_ ColumnOpsType = iota
	AddColumn
	ModifyColumn
	DropColumn
)

type ColumnDefinition struct {
	DataType         DataType
	ColumnConstraint *ColumnConstraint
	Type             ColumnOpsType
}

type ColumnConstraint struct {
	NotNull         bool
	HasDefaultValue bool
	AutoIncrement   bool
	Primary         bool
	Key             bool
	Unique          bool
	Comment         string
	// if HasDefaultValue is false, DefaultValue will be ignored
	DefaultValue string
}

type key bool
type primary bool

// VisitColumnDefinition visits a parse tree produced by MySqlParser#columnDefinition.
func (v *visitor) VisitColumnDefinition(ctx *gen.ColumnDefinitionContext) interface{} {
	v.trace("VisitColumnDefinition")

	var (
		constraint ColumnConstraint
		out        ColumnDefinition
	)
	out.DataType = v.visitDataType(ctx.DataType())
	for _, e := range ctx.AllColumnConstraint() {
		switch tx := e.(type) {
		case *gen.NullColumnConstraintContext:
			constraint.NotNull = v.visitNullColumnConstraint(tx)
		case *gen.DefaultColumnConstraintContext:
			constraint.DefaultValue, constraint.HasDefaultValue = v.visitDefaultColumnConstraint(tx)
		case *gen.AutoIncrementColumnConstraintContext:
			constraint.AutoIncrement = v.visitAutoIncrementColumnConstraint(tx)
		case *gen.PrimaryKeyColumnConstraintContext:
			ret := v.VisitPrimaryKeyColumnConstraint(tx)
			if c, ok := ret.(*primary); ok {
				constraint.Primary = bool(*c)
			} else {
				c := ret.(*key)
				constraint.Key = bool(*c)
			}
		case *gen.UniqueKeyColumnConstraintContext:
			constraint.Unique = v.visitUniqueKeyColumnConstraint(tx)
		case *gen.CommentColumnConstraintContext:
			constraint.Comment = v.visitCommentColumnConstraint(tx)
		case *gen.ReferenceColumnConstraintContext:
			v.panicWithExpr(tx.GetStart(), "Unsupported reference definition")
		}
	}

	out.ColumnConstraint = &constraint
	return &out
}

// visitNullColumnConstraint visits a parse tree produced by MySqlParser#nullColumnConstraint.
func (v *visitor) visitNullColumnConstraint(ctx *gen.NullColumnConstraintContext) bool {
	v.trace("VisitNullColumnConstraint")
	if ret, ok := ctx.NullNotnull().(*gen.NullNotnullContext); ok {
		return v.visitNullNotnull(ret)
	}

	return false
}

// visitDefaultColumnConstraint visits a parse tree produced by MySqlParser#defaultColumnConstraint.
func (v *visitor) visitDefaultColumnConstraint(ctx *gen.DefaultColumnConstraintContext) (value string, ok bool) {
	v.trace("VisitDefaultColumnConstraint")
	value = ctx.DefaultValue().GetText()
	text := ctx.DefaultValue().GetText()
	text = strings.Trim(text, "`")
	text = strings.Trim(text, "'")
	replacer := strings.NewReplacer("\r", "", "\n", "")
	text = replacer.Replace(text)
	if strings.HasPrefix(strings.ToUpper(text), "NULL") {
		return value, false
	}

	return value, true
}

// visitAutoIncrementColumnConstraint visits a parse tree produced by MySqlParser#autoIncrementColumnConstraint.
func (v *visitor) visitAutoIncrementColumnConstraint(_ *gen.AutoIncrementColumnConstraintContext) bool {
	v.trace("VisitAutoIncrementColumnConstraint")
	return true
}

// VisitPrimaryKeyColumnConstraint visits a parse tree produced by MySqlParser#primaryKeyColumnConstraint.
func (v *visitor) VisitPrimaryKeyColumnConstraint(ctx *gen.PrimaryKeyColumnConstraintContext) interface{} {
	v.trace("VisitPrimaryKeyColumnConstraint")
	if ctx.PRIMARY() == nil {
		var ret key
		ret = true
		return &ret
	}

	var ret primary
	ret = true
	return &ret
}

// visitUniqueKeyColumnConstraint visits a parse tree produced by MySqlParser#uniqueKeyColumnConstraint.
func (v *visitor) visitUniqueKeyColumnConstraint(_ *gen.UniqueKeyColumnConstraintContext) bool {
	v.trace("VisitUniqueKeyColumnConstraint")
	return true
}

// visitCommentColumnConstraint visits a parse tree produced by MySqlParser#commentColumnConstraint.
func (v *visitor) visitCommentColumnConstraint(ctx *gen.CommentColumnConstraintContext) string {
	v.trace("VisitCommentColumnConstraint")
	value := parseTerminalNode(
		ctx.STRING_LITERAL(),
		withTrim("`"),
		withTrim(`"`),
		withTrim(`'`),
		withReplacer(`\r`, "", `\n`, ""),
	)
	return value
}

// visitNullNotnull visits a parse tree produced by MySqlParser#nullNotnull.
func (v *visitor) visitNullNotnull(ctx *gen.NullNotnullContext) bool {
	v.trace("VisitNullNotnull")
	if ctx.NOT() != nil {
		return true
	}

	return false
}
