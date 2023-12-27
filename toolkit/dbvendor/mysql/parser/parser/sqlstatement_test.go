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

func TestVisitor_VisitSqlStatements(t *testing.T) {
	p := NewParser(WithDebugMode(true))
	accept := func(p *gen.MySqlParser, visitor *visitor) interface{} {
		root := p.Root()
		return root.Accept(visitor)
	}

	t.Run("modifyColumn", func(t *testing.T) {
		ret, err := p.testMysqlSyntax("test.sql", accept, `alter table my_test_table
          
            MODIFY REVIEW_COMMENTS TEXT(4000) null default '123' comment '审核意见'
         , 
            MODIFY CREATE_USER VARCHAR(255) not null comment '创建人'`)
		assert.Nil(t, err)
		assert.NotNil(t, ret)
	})

	t.Run("addColumn", func(t *testing.T) {
		_, err := p.testMysqlSyntax("test.sql", accept, `alter table my_test_table
          
            ADD CUSTOM_FIELD6 VARCHAR(255) null default null comment 'custom_field6'
         , 
            ADD CUSTOM_FIELD5 VARCHAR(255) null default 'bbc' comment 'custom_field5'`)
		assert.Nil(t, err)
	})

	t.Run("empty", func(t *testing.T) {
		_, err := p.testMysqlSyntax("test.sql", accept, ``)
		assert.Nil(t, err)
	})

	t.Run("createDatabase", func(t *testing.T) {
		ret, err := p.testMysqlSyntax("test.sql", accept, "create database user")
		assert.Nil(t, err)
		assert.Equal(t, []*CreateTable(nil), ret)
	})

	t.Run("createSingleTable", func(t *testing.T) {
		ret, err := p.testMysqlSyntax("test.sql", accept, `
			create table if not exists user(
				id bigint(11) primary key not null default 0 comment '主键ID'
			)
		`)
		tables, ok := ret.([]*CreateTable)
		assert.True(t, ok)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(tables))
		assertCreateTableEqual(t, &CreateTable{
			Name: "user",
			Columns: []*ColumnDeclaration{
				{
					Name: "id",
					ColumnDefinition: &ColumnDefinition{
						DataType: &NormalDataType{tp: BigInt},
						ColumnConstraint: &ColumnConstraint{
							NotNull:         true,
							HasDefaultValue: true,
							AutoIncrement:   false,
							Primary:         true,
							Comment:         "主键ID",
						},
					},
				},
			},
		}, tables[0])
	})

	t.Run("ddlWithOtherSql", func(t *testing.T) {
		ret, err := p.testMysqlSyntax("test.sql", accept, `
			-- ddl create table
			create table if not exists user(
				id bigint(11) primary key not null comment 'id'
			)
			-- ddl create database
			create database foo;
			-- dml select
			select * from bar;
			-- dml update
			update foo set bar = 'test';
			-- dml insert
			insert into foo ('id','name') values ('1','bar');
		`)
		assert.Nil(t, err)
		assert.NotNil(t, ret)
		tables, ok := ret.([]*CreateTable)
		assert.True(t, ok)
		assert.Equal(t, 1, len(tables))
		assert.Equal(t, "user", tables[0].Name)
	})
}
