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
	"github.com/unionj-cloud/go-doudou/v2/toolkit/dbvendor/mysql/parser/gen"
	"strings"
)

type AlterTable struct {
	// Name describes the literal of table, this name can be specified as db_name.tbl_name,
	// https://dev.mysql.com/doc/refman/8.0/en/create-table.html#create-table-name
	Name    string
	Columns []*ColumnDeclaration
}

// VisitAlterTable visits a parse tree produced by MySqlParser#alterTable.
func (v *visitor) VisitAlterTable(ctx *gen.AlterTableContext) interface{} {
	v.trace("VisitAlterTable")
	var ret AlterTable
	tableName := ctx.TableName().GetText()
	tableName = strings.Trim(tableName, "`")
	tableName = strings.Trim(tableName, "'")
	replacer := strings.NewReplacer("\r", "", "\n", "")
	tableName = replacer.Replace(tableName)
	ret.Name = tableName
	for _, item := range ctx.AllAlterSpecification() {
		data := v.visitAlterSpecification(item)
		if data == nil {
			continue
		}

		switch r := data.(type) {
		case *ColumnDeclaration:
			ret.Columns = append(ret.Columns, r)
		}
	}

	return &ret
}

// VisitCreateDefinition visits a parse tree produced by MySqlParser#createDefinition.
func (v *visitor) visitAlterSpecification(ctx gen.IAlterSpecificationContext) interface{} {
	v.trace("VisitAlterSpecification")
	switch tx := ctx.(type) {
	case *gen.AlterByModifyColumnContext:
		ret := v.getColumnDeclaration(tx)
		ret.ColumnDefinition.Type = ModifyColumn
		return &ret
	case *gen.AlterByAddColumnContext:
		ret := v.getColumnDeclaration(tx)
		ret.ColumnDefinition.Type = AddColumn
		return &ret
	case *gen.AlterByDropColumnContext:
		v.trace("*gen.AlterByDropColumnContext")
	default:
		v.trace("Not support " + tx.GetText())
	}
	return nil
}

type iAlterTableColumnContext interface {
	Uid(i int) gen.IUidContext
	ColumnDefinition() gen.IColumnDefinitionContext
}

func (v *visitor) getColumnDeclaration(ctx iAlterTableColumnContext) ColumnDeclaration {
	var ret ColumnDeclaration
	ret.Name = v.visitUid(ctx.Uid(0))
	iDefinitionContext := ctx.ColumnDefinition()
	definitionContext, ok := iDefinitionContext.(*gen.ColumnDefinitionContext)
	if ok {
		out := v.VisitColumnDefinition(definitionContext)
		if cd, ok := out.(*ColumnDefinition); ok {
			ret.ColumnDefinition = cd
		}
	}
	return ret
}
