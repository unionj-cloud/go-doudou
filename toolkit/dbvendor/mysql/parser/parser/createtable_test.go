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

func TestVisitor_VisitCreateTable(t *testing.T) {
	p := NewParser(WithDebugMode(true))
	accept := func(p *gen.MySqlParser, visitor *visitor) interface{} {
		ctx := p.CreateTable()
		return visitor.visitCreateTable(ctx)
	}

	t.Run("copyCreateTableContext", func(t *testing.T) {
		_, err := p.testMysqlSyntax("test.sql", accept,
			`create table new_t  (like t1);`)
		assert.Error(t, err)
	})

	t.Run("queryCreateTable", func(t *testing.T) {
		_, err := p.testMysqlSyntax("test.sql", accept,
			`CREATE TABLE test (a INT NOT NULL AUTO_INCREMENT,PRIMARY KEY (a),
				KEY(b))ENGINE=InnoDB SELECT b,c FROM test2;`)
		assert.Error(t, err)
	})

	t.Run("columnCreateTable_normal_case", func(t *testing.T) {
		v, err := p.testMysqlSyntax("test.sql", accept,
			"CREATE TABLE `user` (\n  "+
				"`id` bigint NOT NULL AUTO_INCREMENT,\n  "+
				"`number` varchar(255) NOT NULL DEFAULT '' COMMENT '学号',\n  "+
				"`name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '用户名称',\n "+
				" `password` varchar(255) NOT NULL DEFAULT '' COMMENT '用户密码',\n "+
				" `gender` char(5) NOT NULL COMMENT '男｜女｜未公开',\n  "+
				"`create_time` timestamp NULL DEFAULT NULL,\n  "+
				"`update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n  "+
				"PRIMARY KEY (`id`),\n  "+
				"UNIQUE KEY `number_unique` (`number`) USING BTREE,\n  "+
				"UNIQUE KEY `number_unique2` (`number`"+
				") USING BTREE\n) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
		assert.Nil(t, err)
		table, ok := v.(*CreateTable)
		assert.True(t, ok)
		expected := &CreateTable{
			Name: "user",
			Columns: []*ColumnDeclaration{
				{
					Name: "id",
					ColumnDefinition: &ColumnDefinition{
						DataType: &NormalDataType{tp: BigInt},
						ColumnConstraint: &ColumnConstraint{
							NotNull:       true,
							AutoIncrement: true,
						},
					},
				},
				{
					Name: "number",
					ColumnDefinition: &ColumnDefinition{
						DataType: &NormalDataType{tp: VarChar},
						ColumnConstraint: &ColumnConstraint{
							NotNull:         true,
							Comment:         "学号",
							HasDefaultValue: true,
						},
					},
				},
				{
					Name: "name",
					ColumnDefinition: &ColumnDefinition{
						DataType: &NormalDataType{tp: VarChar},
						ColumnConstraint: &ColumnConstraint{
							Comment: "用户名称",
						},
					},
				},
				{
					Name: "password",
					ColumnDefinition: &ColumnDefinition{
						DataType: &NormalDataType{tp: VarChar},
						ColumnConstraint: &ColumnConstraint{
							NotNull:         true,
							Comment:         "用户密码",
							HasDefaultValue: true,
						},
					},
				},
				{
					Name: "gender",
					ColumnDefinition: &ColumnDefinition{
						DataType: &NormalDataType{tp: Char},
						ColumnConstraint: &ColumnConstraint{
							NotNull: true,
							Comment: "男｜女｜未公开",
						},
					},
				},
				{
					Name: "create_time",
					ColumnDefinition: &ColumnDefinition{
						DataType:         &NormalDataType{tp: Timestamp},
						ColumnConstraint: &ColumnConstraint{},
					},
				},
				{
					Name: "update_time",
					ColumnDefinition: &ColumnDefinition{
						DataType: &NormalDataType{tp: Timestamp},
						ColumnConstraint: &ColumnConstraint{
							HasDefaultValue: true,
						},
					},
				},
			},
			Constraints: []*TableConstraint{
				{
					ColumnPrimaryKey: []string{"id"},
				},
				{
					ColumnUniqueKey: []string{"number"},
				},
				{
					ColumnUniqueKey: []string{"number"},
				},
			},
		}
		assertCreateTableEqual(t, expected, table)
	})

	t.Run("columnCreateTable_every_case", func(t *testing.T) {
		v, err := p.testMysqlSyntax("test.sql", accept,
			`create table if not exists foo (
					id bigint(20) not null primary key auto_increment default 0 comment 'id',
					class_id varchar not null default '' comment '班级id',
					name char(10) not null key default '' comment '姓名',
					mobile varchar(15) not null unique default '' comment '手机号',
					gender enum('男','女') not null default '男' comment '性别',
					flag boolean not null default 'false' comment '标志位',
					document JSON NOT NULL,
					location POINT comment '地理位置',
					primary key ('id'),
					key 'name_idx' ('name'),
					unique key 'class_mobile_uni' ('class_id','mobile')
				)`)
		assert.Nil(t, err)
		table, ok := v.(*CreateTable)
		assert.True(t, ok)
		assertCreateTableEqual(t, &CreateTable{
			Name: "foo",
			Columns: []*ColumnDeclaration{
				{
					Name: "id",
					ColumnDefinition: &ColumnDefinition{
						DataType: &NormalDataType{tp: BigInt},
						ColumnConstraint: &ColumnConstraint{
							NotNull:         true,
							HasDefaultValue: true,
							AutoIncrement:   true,
							Primary:         true,
							Comment:         "id",
						},
					},
				},
				{
					Name: "class_id",
					ColumnDefinition: &ColumnDefinition{
						DataType: &NormalDataType{tp: VarChar},
						ColumnConstraint: &ColumnConstraint{
							NotNull:         true,
							Comment:         "班级id",
							HasDefaultValue: true,
						},
					},
				},
				{
					Name: "name",
					ColumnDefinition: &ColumnDefinition{
						DataType: &NormalDataType{tp: Char},
						ColumnConstraint: &ColumnConstraint{
							NotNull:         true,
							Key:             true,
							Comment:         "姓名",
							HasDefaultValue: true,
						},
					},
				},
				{
					Name: "mobile",
					ColumnDefinition: &ColumnDefinition{
						DataType: &NormalDataType{tp: VarChar},
						ColumnConstraint: &ColumnConstraint{
							NotNull:         true,
							Unique:          true,
							Comment:         "手机号",
							HasDefaultValue: true,
						},
					},
				},
				{
					Name: "gender",
					ColumnDefinition: &ColumnDefinition{
						DataType: &EnumSetDataType{
							tp:    Enum,
							value: []string{"男", "女"},
						},
						ColumnConstraint: &ColumnConstraint{
							NotNull:         true,
							HasDefaultValue: true,
							Comment:         "性别",
						},
					},
				},
				{
					Name: "flag",
					ColumnDefinition: &ColumnDefinition{
						DataType: &EnumSetDataType{
							tp: Boolean,
						},
						ColumnConstraint: &ColumnConstraint{
							NotNull:         true,
							HasDefaultValue: true,
							Comment:         "标志位",
						},
					},
				},
				{
					Name: "document",
					ColumnDefinition: &ColumnDefinition{
						DataType: &EnumSetDataType{
							tp: Json,
						},
						ColumnConstraint: &ColumnConstraint{
							NotNull: true,
						},
					},
				},
				{
					Name: "location",
					ColumnDefinition: &ColumnDefinition{
						DataType: &EnumSetDataType{
							tp: Point,
						},
						ColumnConstraint: &ColumnConstraint{
							Comment: "地理位置",
						},
					},
				},
			},
			Constraints: []*TableConstraint{
				{
					ColumnPrimaryKey: []string{"id"},
				},
				{
					ColumnUniqueKey: []string{"class_id", "mobile"},
				},
			},
		}, table)
	})
}

func TestGetTableFromCreateTable(t *testing.T) {
	p := NewParser(WithDebugMode(true))
	accept := func(p *gen.MySqlParser, visitor *visitor) interface{} {
		ctx := p.CreateTable()
		return visitor.visitCreateTable(ctx)
	}
	v, err := p.testMysqlSyntax("test.sql", accept,
		"CREATE TABLE `foo`.`bar` (\n  "+
			"`id` bigint NOT NULL AUTO_INCREMENT,\n  "+
			"PRIMARY KEY (`id`)\n"+
			") USING BTREE\n) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	assert.Nil(t, err)
	createTable, ok := v.(*CreateTable)
	assert.True(t, ok)
	assert.Equal(t, "foo`.`bar", createTable.Name)
	table := createTable.Convert()
	assert.Equal(t, "bar", table.Name)
}

func assertCreateTableEqual(t *testing.T, expected, actual *CreateTable) {
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, len(expected.Columns), len(actual.Columns))
	sort.SliceStable(expected.Columns, func(i, j int) bool {
		return expected.Columns[i].Name < expected.Columns[j].Name
	})
	sort.SliceStable(actual.Columns, func(i, j int) bool {
		return actual.Columns[i].Name < actual.Columns[j].Name
	})

	for i, expectedColumn := range expected.Columns {
		actualColumn := actual.Columns[i]
		assert.Equal(t, expectedColumn.Name, actualColumn.Name)
		assertColumnDefinition(t, expectedColumn.ColumnDefinition, actualColumn.ColumnDefinition)
	}

	for i, expectedConstraint := range expected.Constraints {
		actualConstraint := actual.Constraints[i]
		assert.Equal(t, *expectedConstraint, *actualConstraint)
		assert.Equal(t, expectedConstraint, actualConstraint)
	}
}

func assertColumnDefinition(t *testing.T, expected, actual *ColumnDefinition) {
	assertDataType(t, expected.DataType, actual.DataType)
	assert.Equal(t, *expected.ColumnConstraint, *actual.ColumnConstraint)
}

func assertDataType(t *testing.T, expected, actual DataType) {
	assert.Equal(t, expected.Type(), actual.Type())
	assert.Equal(t, expected.Value(), actual.Value())
}
