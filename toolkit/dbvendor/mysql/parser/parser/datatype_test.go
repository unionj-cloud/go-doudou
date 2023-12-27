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

func TestVisitor_VisitDataType(t *testing.T) {
	p := NewParser(WithDebugMode(true))
	accept := func(p *gen.MySqlParser, visitor *visitor) interface{} {
		return visitor.visitDataType(p.DataType())
	}

	t.Run("stringDataType", func(t *testing.T) {
		testData := map[string]int{
			`CHAR(10)`:      Char,
			`CHARACTER(10)`: Character,
			`VARCHAR(10)`:   VarChar,
			`TINYTEXT`:      TinyText,
			`TEXT`:          Text,
			`MEDIUMTEXT`:    MediumText,
			`LONGTEXT`:      LongText,
			`NCHAR(20)`:     NChar,
			`NVARCHAR(20)`:  NVarChar,
			`LONG`:          LongVarChar,
		}

		for sql, dataType := range testData {
			actual, err := p.testMysqlSyntax("test.sql", accept, sql)
			assert.Nil(t, err)
			assertTypeEqual(t, dataType, actual)
		}
	})

	t.Run("nationalStringDataType", func(t *testing.T) {
		testData := map[string]int{
			`NATIONAL VARCHAR(255)`:          NVarChar,
			`NATIONAL CHARACTER(255) BINARY`: NChar,
			`NCHAR VARCHAR(255) BINARY`:      NVarChar,
			`NCHAR VARCHAR(200)`:             NVarChar,
		}

		for sql, dataType := range testData {
			actual, err := p.testMysqlSyntax("test.sql", accept, sql)
			assert.Nil(t, err)
			assertTypeEqual(t, dataType, actual)
		}
	})

	t.Run("nationalVaryingStringDataType", func(t *testing.T) {
		testData := map[string]int{
			`NATIONAL CHAR VARYING (255)`:             NVarChar,
			`NATIONAL CHAR VARYING (255) BINARY`:      NVarChar,
			`NATIONAL CHARACTER VARYING (255)`:        NVarChar,
			`NATIONAL CHARACTER VARYING (255) BINARY`: NVarChar,
		}

		for sql, dataType := range testData {
			actual, err := p.testMysqlSyntax("test.sql", accept, sql)
			assert.Nil(t, err)
			assertTypeEqual(t, dataType, actual)
		}
	})

	t.Run("dimensionDataType", func(t *testing.T) {
		testData := map[string]int{
			`TINYINT(1)`:                       TinyInt,
			`TINYINT(1) SIGNED`:                TinyInt,
			`TINYINT(1) UNSIGNED`:              TinyInt,
			`TINYINT(1) UNSIGNED ZEROFILL`:     TinyInt,
			`SMALLINT(10)`:                     SmallInt,
			`SMALLINT(10) SIGNED`:              SmallInt,
			`SMALLINT(10) UNSIGNED`:            SmallInt,
			`SMALLINT(10) ZEROFILL`:            SmallInt,
			`MEDIUMINT(10)`:                    MediumInt,
			`MEDIUMINT(10) SIGNED`:             MediumInt,
			`MEDIUMINT(10) UNSIGNED`:           MediumInt,
			`MEDIUMINT(10) ZEROFILL`:           MediumInt,
			`INT(10)`:                          Int,
			`INT(10) SIGNED`:                   Int,
			`INT(10) UNSIGNED`:                 Int,
			`INT(10) ZEROFILL`:                 Int,
			`INTEGER(10)`:                      Integer,
			`INTEGER(10) SIGNED`:               Integer,
			`INTEGER(10) UNSIGNED`:             Integer,
			`INTEGER(10) ZEROFILL`:             Integer,
			`BIGINT(20)`:                       BigInt,
			`BIGINT(20) SIGNED`:                BigInt,
			`BIGINT(20) UNSIGNED`:              BigInt,
			`BIGINT(20) ZEROFILL`:              BigInt,
			`MIDDLEINT(20)`:                    MiddleInt,
			`MIDDLEINT(20) SIGNED`:             MiddleInt,
			`MIDDLEINT(20) UNSIGNED`:           MiddleInt,
			`MIDDLEINT(20) ZEROFILL`:           MiddleInt,
			`INT1(2)`:                          Int1,
			`INT1(2) SIGNED`:                   Int1,
			`INT1(2) UNSIGNED`:                 Int1,
			`INT1(2) ZEROFILL`:                 Int1,
			`INT2(2)`:                          Int2,
			`INT2(2) SIGNED`:                   Int2,
			`INT2(2) UNSIGNED`:                 Int2,
			`INT2(2) ZEROFILL`:                 Int2,
			`INT3(20)`:                         Int3,
			`INT3(3) SIGNED`:                   Int3,
			`INT3(3) UNSIGNED`:                 Int3,
			`INT3(3) ZEROFILL`:                 Int3,
			`INT4(4)`:                          Int4,
			`INT4(4) SIGNED`:                   Int4,
			`INT4(4) UNSIGNED`:                 Int4,
			`INT4(4) ZEROFILL`:                 Int4,
			`INT8(8)`:                          Int8,
			`INT8(8) SIGNED`:                   Int8,
			`INT8(8) UNSIGNED`:                 Int8,
			`INT8(8) ZEROFILL`:                 Int8,
			`REAL(8,10) ZEROFILL`:              Real,
			`REAL ZEROFILL`:                    Real,
			`REAL SIGNED ZEROFILL`:             Real,
			`REAL UNSIGNED ZEROFILL`:           Real,
			`DOUBLE(8,10) ZEROFILL`:            Double,
			`DOUBLE PRECISION (8,10) ZEROFILL`: Double,
			`DOUBLE ZEROFILL`:                  Double,
			`DOUBLE SIGNED ZEROFILL`:           Double,
			`DOUBLE UNSIGNED ZEROFILL`:         Double,
			`DECIMAL(8,10) ZEROFILL`:           Decimal,
			`DECIMAL ZEROFILL`:                 Decimal,
			`DECIMAL SIGNED ZEROFILL`:          Decimal,
			`DECIMAL UNSIGNED ZEROFILL`:        Decimal,
			`DEC(8,10) ZEROFILL`:               Dec,
			`DEC ZEROFILL`:                     Dec,
			`DEC SIGNED ZEROFILL`:              Dec,
			`DEC UNSIGNED ZEROFILL`:            Dec,
			`FIXED(8,10) ZEROFILL`:             Fixed,
			`FIXED ZEROFILL`:                   Fixed,
			`FIXED SIGNED ZEROFILL`:            Fixed,
			`FIXED UNSIGNED ZEROFILL`:          Fixed,
			`NUMERIC(8,10) ZEROFILL`:           Numeric,
			`NUMERIC ZEROFILL`:                 Numeric,
			`NUMERIC SIGNED ZEROFILL`:          Numeric,
			`NUMERIC UNSIGNED ZEROFILL`:        Numeric,
			`FLOAT(8,10) ZEROFILL`:             Float,
			`FLOAT ZEROFILL`:                   Float,
			`FLOAT SIGNED ZEROFILL`:            Float,
			`FLOAT UNSIGNED ZEROFILL`:          Float,
			`FLOAT4(8,10) ZEROFILL`:            Float4,
			`FLOAT4 ZEROFILL`:                  Float4,
			`FLOAT4 SIGNED ZEROFILL`:           Float4,
			`FLOAT4 UNSIGNED ZEROFILL`:         Float4,
			`FLOAT8(8,10) ZEROFILL`:            Float8,
			`FLOAT8 ZEROFILL`:                  Float8,
			`FLOAT8 SIGNED ZEROFILL`:           Float8,
			`FLOAT8 UNSIGNED ZEROFILL`:         Float8,
			`BIT`:                              Bit,
			`BIT(1)`:                           Bit,
			`TIME`:                             Time,
			`TIMESTAMP`:                        Timestamp,
			`DATETIME`:                         DateTime,
			`BINARY`:                           Binary,
			`VARBINARY`:                        VarBinary,
			`BLOB`:                             Blob,
			`YEAR`:                             Year,
		}

		for sql, dataType := range testData {
			actual, err := p.testMysqlSyntax("test.sql", accept, sql)
			assert.Nil(t, err)
			assertTypeEqual(t, dataType, actual)
		}

		testData = map[string]int{
			`TINYINT(1) UNSIGNED`: TinyInt,
			`SMALLINT UNSIGNED`:   SmallInt,
			`BIGINT UNSIGNED`:     BigInt,
		}
		for sql, dataType := range testData {
			actual, err := p.testMysqlSyntax("test.sql", accept, sql)
			assert.Nil(t, err)
			assertTypeEqual(t, dataType, actual, true)
		}
	})

	t.Run("simpleDataType", func(t *testing.T) {
		testData := map[string]int{
			`DATE`:       Date,
			`TINYBLOB`:   TinyBlob,
			`MEDIUMBLOB`: MediumBlob,
			`LONGBLOB`:   LongBlob,
			`BOOL`:       Bool,
			`BOOLEAN`:    Boolean,
			`SERIAL`:     Serial,
		}

		for sql, dataType := range testData {
			actual, err := p.testMysqlSyntax("test.sql", accept, sql)
			assert.Nil(t, err)
			assertTypeEqual(t, dataType, actual)
		}
	})

	t.Run("collectionDataType", func(t *testing.T) {
		testData := map[string]EnumSetDataType{
			`ENUM('1','2')`: {
				tp:    Enum,
				value: []string{"1", "2"},
			},
			`SET('A','B')`: {
				tp:    Set,
				value: []string{"A", "B"},
			},
			`SET('A','B') BINARY`: {
				tp:    Set,
				value: []string{"A", "B"},
			},
		}

		for sql, e := range testData {
			actual, err := p.testMysqlSyntax("test.sql", accept, sql)
			assert.Nil(t, err)
			assertEnumTypeEqual(t, e.tp, e.value, actual)
		}
	})

	t.Run("spatialDataType", func(t *testing.T) {
		testData := map[string]int{
			`GEOMETRYCOLLECTION`: GeometryCollection,
			`GEOMCOLLECTION`:     GeomCollection,
			`LINESTRING`:         LineString,
			`MULTILINESTRING`:    MultiLineString,
			`MULTIPOINT`:         MultiPoint,
			`MULTIPOLYGON`:       MultiPolygon,
			`POINT`:              Point,
			`POLYGON`:            Polygon,
			`JSON`:               Json,
			`GEOMETRY`:           Geometry,
		}

		for sql, dataType := range testData {
			actual, err := p.testMysqlSyntax("test.sql", accept, sql)
			assert.Nil(t, err)
			assertTypeEqual(t, dataType, actual)
		}
	})

	t.Run("longVarcharDataType ", func(t *testing.T) {
		testData := map[string]int{
			`LONG`:                LongVarChar,
			`LONG VARCHAR`:        LongVarChar,
			`LONG VARCHAR BINARY`: LongVarChar,
			`LONG VARCHAR BINARY CHARACTER SET 'utf8'`: LongVarChar,
			`LONG VARCHAR BINARY CHARSET 'utf8'`:       LongVarChar,
		}

		for sql, dataType := range testData {
			actual, err := p.testMysqlSyntax("test.sql", accept, sql)
			assert.Nil(t, err)
			assertTypeEqual(t, dataType, actual)
		}
	})

	t.Run("longVarbinaryDataType ", func(t *testing.T) {
		testData := map[string]int{
			`LONG VARBINARY  `: LongVarBinary,
		}

		for sql, dataType := range testData {
			actual, err := p.testMysqlSyntax("test.sql", accept, sql)
			assert.Nil(t, err)
			assertTypeEqual(t, dataType, actual)
		}
	})
}

func assertTypeEqual(t *testing.T, expected int, actual interface{}, unsigned ...bool) {
	assert.Equal(t, expected, actual.(DataType).Type())
	if len(unsigned) > 0 {
		assert.Equal(t, unsigned[0], actual.(DataType).Unsigned())
	}
}

func assertEnumTypeEqual(t *testing.T, expectedType int, values []string, actual interface{}) {
	assert.Equal(t, expectedType, actual.(DataType).Type())
	assert.Equal(t, values, actual.(DataType).Value())
}
