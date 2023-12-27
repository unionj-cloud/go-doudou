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
)

const (
	_ int = iota
	LongVarBinary
	LongVarChar
	GeometryCollection
	GeomCollection
	LineString
	MultiLineString
	MultiPoint
	MultiPolygon
	Point
	Polygon
	Json
	Geometry
	Enum
	Set
	Bit
	Time
	Timestamp
	DateTime
	Binary
	VarBinary
	Blob
	Year
	Decimal
	Dec
	Fixed
	Numeric
	Float
	Float4
	Float8
	Double
	Real
	TinyInt
	SmallInt
	MediumInt
	Int
	Integer
	BigInt
	MiddleInt
	Int1
	Int2
	Int3
	Int4
	Int8
	Date
	TinyBlob
	MediumBlob
	LongBlob
	Bool
	Boolean
	Serial
	NVarChar
	NChar
	Char
	Character
	VarChar
	TinyText
	Text
	MediumText
	LongText
)

// DataType describes the data type and value of the column in table
type DataType interface {
	Type() int
	Unsigned() bool
	// Value returns the values if the data type is Enum or Set
	Value() []string
	String() string
}

var _ DataType = (*NormalDataType)(nil)
var _ DataType = (*EnumSetDataType)(nil)

// NormalDataType describes the data type which not contains Enum and Set of column
type NormalDataType struct {
	tp       int
	unsigned bool
	text     string
}

func (n *NormalDataType) String() string {
	return n.text
}

// Unsigned returns true if the data type is unsigned.
func (n *NormalDataType) Unsigned() bool {
	return n.unsigned
}

// Type returns the data type of column
func (n *NormalDataType) Type() int {
	return n.tp
}

// Value returns nil default
func (n *NormalDataType) Value() []string {
	return nil
}

func with(tp int, unsigned bool, text string, value ...string) DataType {
	if len(value) > 0 {
		return &EnumSetDataType{
			tp:    tp,
			value: value,
		}
	}
	return &NormalDataType{tp: tp, unsigned: unsigned, text: text}
}

// EnumSetDataType describes the data type  Enum and Set of column
type EnumSetDataType struct {
	tp    int
	value []string
	text  string
}

func (e *EnumSetDataType) String() string {
	return e.text
}

// Type returns the data type of column
func (e *EnumSetDataType) Type() int {
	return e.tp
}

// Unsigned returns true if the data type is unsigned.
func (e *EnumSetDataType) Unsigned() bool {
	return false
}

// Value returns the value of data type Enum and Set
func (e *EnumSetDataType) Value() []string {
	return e.value
}

// visitDataType visits data type by switch-case
func (v *visitor) visitDataType(ctx gen.IDataTypeContext) DataType {
	v.trace("VisitDataType")
	switch t := ctx.(type) {
	case *gen.StringDataTypeContext:
		return v.visitStringDataType(t)
	case *gen.NationalStringDataTypeContext:
		return v.visitNationalStringDataType(t)
	case *gen.NationalVaryingStringDataTypeContext:
		return v.visitNationalVaryingStringDataType(t)
	case *gen.DimensionDataTypeContext:
		return v.visitDimensionDataType(t)
	case *gen.SimpleDataTypeContext:
		return v.visitSimpleDataType(t)
	case *gen.CollectionDataTypeContext:
		return v.visitCollectionDataType(t)
	case *gen.SpatialDataTypeContext:
		return v.visitSpatialDataType(t)
	case *gen.LongVarcharDataTypeContext:
		return v.visitLongVarcharDataType(t)
	case *gen.LongVarbinaryDataTypeContext:
		return v.visitLongVarbinaryDataType(t)
	}

	v.panicWithExpr(ctx.GetStart(), "invalid data type: "+ctx.GetText())
	return nil
}

// visitStringDataType visits a parse tree produced by MySqlParser#stringDataType.
func (v *visitor) visitStringDataType(ctx *gen.StringDataTypeContext) DataType {
	v.trace(`VisitStringDataType`)
	text := parseToken(ctx.GetTypeName(), withUpperCase(), withTrim("`"))
	switch text {
	case `CHAR`:
		return with(Char, false, text)
	case `CHARACTER`:
		return with(Character, false, text)
	case `VARCHAR`:
		return with(VarChar, false, text)
	case `TINYTEXT`:
		return with(TinyText, false, text)
	case `TEXT`:
		return with(Text, false, text)
	case `MEDIUMTEXT`:
		return with(MediumText, false, text)
	case `LONGTEXT`:
		return with(LongText, false, text)
	case `NCHAR`:
		return with(NChar, false, text)
	case `NVARCHAR`:
		return with(NVarChar, false, text)
	case `LONG`:
		return with(LongVarChar, false, text)
	}

	v.panicWithExpr(ctx.GetTypeName(), "invalid data type: "+text)
	return nil
}

// visitNationalStringDataType visits a parse tree produced by MySqlParser#nationalVaryingStringDataType.
func (v *visitor) visitNationalStringDataType(ctx *gen.NationalStringDataTypeContext) DataType {
	v.trace(`VisitNationalStringDataType`)
	text := parseToken(ctx.GetTypeName(), withUpperCase(), withTrim("`"))
	switch text {
	case `VARCHAR`:
		return with(NVarChar, false, text)
	case `CHARACTER`:
		return with(NChar, false, text)
	}

	v.panicWithExpr(ctx.GetTypeName(), "invalid data type: "+text)
	return nil
}

// visitNationalVaryingStringDataType visits a parse tree produced by MySqlParser#nationalVaryingStringDataType.
func (v *visitor) visitNationalVaryingStringDataType(_ *gen.NationalVaryingStringDataTypeContext) DataType {
	v.trace("VisitNationalVaryingStringDataType")
	return with(NVarChar, false, "")
}

// visitDimensionDataType visits a parse tree produced by MySqlParser#dimensionDataType.
func (v *visitor) visitDimensionDataType(ctx *gen.DimensionDataTypeContext) DataType {
	v.trace("VisitDimensionDataType")
	text := parseToken(ctx.GetTypeName(), withUpperCase(), withTrim("`"))
	unsigned := ctx.UNSIGNED() != nil
	switch text {
	case `BIT`:
		return with(Bit, unsigned, text)
	case `TIME`:
		return with(Time, unsigned, text)
	case `TIMESTAMP`:
		return with(Timestamp, unsigned, text)
	case `DATETIME`:
		return with(DateTime, unsigned, text)
	case `BINARY`:
		return with(Binary, unsigned, text)
	case `VARBINARY`:
		return with(VarBinary, unsigned, text)
	case `BLOB`:
		return with(Blob, unsigned, text)
	case `YEAR`:
		return with(Year, unsigned, text)
	case `DECIMAL`:
		return with(Decimal, unsigned, text)
	case `DEC`:
		return with(Dec, unsigned, text)
	case `FIXED`:
		return with(Fixed, unsigned, text)
	case `NUMERIC`:
		return with(Numeric, unsigned, text)
	case `FLOAT`:
		return with(Float, unsigned, text)
	case `FLOAT4`:
		return with(Float4, unsigned, text)
	case `FLOAT8`:
		return with(Float8, unsigned, text)
	case `DOUBLE`:
		return with(Double, unsigned, text)
	case `REAL`:
		return with(Real, unsigned, text)
	case `TINYINT`:
		return with(TinyInt, unsigned, text)
	case `SMALLINT`:
		return with(SmallInt, unsigned, text)
	case `MEDIUMINT`:
		return with(MediumInt, unsigned, text)
	case `INT`:
		return with(Int, unsigned, text)
	case `INTEGER`:
		return with(Integer, unsigned, text)
	case `BIGINT`:
		return with(BigInt, unsigned, text)
	case `MIDDLEINT`:
		return with(MiddleInt, unsigned, text)
	case `INT1`:
		return with(Int1, unsigned, text)
	case `INT2`:
		return with(Int2, unsigned, text)
	case `INT3`:
		return with(Int3, unsigned, text)
	case `INT4`:
		return with(Int4, unsigned, text)
	case `INT8`:
		return with(Int8, unsigned, text)
	}

	v.panicWithExpr(ctx.GetTypeName(), "invalid data type: "+text)
	return nil
}

// visitSimpleDataType visits a parse tree produced by MySqlParser#simpleDataType.
func (v *visitor) visitSimpleDataType(ctx *gen.SimpleDataTypeContext) DataType {
	v.trace("VisitSimpleDataType")
	text := parseToken(
		ctx.GetTypeName(),
		withUpperCase(),
		withTrim("`"),
	)

	switch text {
	case `DATE`:
		return with(Date, false, text)
	case `TINYBLOB`:
		return with(TinyBlob, false, text)
	case `MEDIUMBLOB`:
		return with(MediumBlob, false, text)
	case `LONGBLOB`:
		return with(LongBlob, false, text)
	case `BOOL`:
		return with(Bool, false, text)
	case `BOOLEAN`:
		return with(Boolean, false, text)
	case `SERIAL`:
		return with(Serial, false, text)
	}

	v.panicWithExpr(ctx.GetTypeName(), "invalid data type: "+text)
	return nil
}

// visitCollectionDataType visits a parse tree produced by MySqlParser#collectionDataType.
func (v *visitor) visitCollectionDataType(ctx *gen.CollectionDataTypeContext) DataType {
	v.trace("VisitCollectionDataType")
	text := parseToken(
		ctx.GetTypeName(),
		withUpperCase(),
		withTrim("`"),
	)

	var values []string
	if ctx.CollectionOptions() != nil {
		optionsCtx, ok := ctx.CollectionOptions().(*gen.CollectionOptionsContext)
		if ok {
			for _, e := range optionsCtx.AllSTRING_LITERAL() {
				value := parseTerminalNode(
					e, withTrim("`"),
					withTrim(`"`),
					withTrim(`'`),
				)
				values = append(values, value)
			}
		}
	}

	switch text {
	case `ENUM`:
		return with(Enum, false, text, values...)
	case `SET`:
		return with(Set, false, text, values...)
	}

	v.panicWithExpr(ctx.GetTypeName(), "invalid data type: "+text)
	return nil
}

// visitSpatialDataType visits a parse tree produced by MySqlParser#spatialDataType.
func (v *visitor) visitSpatialDataType(ctx *gen.SpatialDataTypeContext) DataType {
	v.trace("VisitSpatialDataType")
	text := parseToken(
		ctx.GetTypeName(),
		withUpperCase(),
		withTrim("`"),
	)

	switch text {
	case `GEOMETRYCOLLECTION`:
		return with(GeometryCollection, false, text)
	case `GEOMCOLLECTION`:
		return with(GeomCollection, false, text)
	case `LINESTRING`:
		return with(LineString, false, text)
	case `MULTILINESTRING`:
		return with(MultiLineString, false, text)
	case `MULTIPOINT`:
		return with(MultiPoint, false, text)
	case `MULTIPOLYGON`:
		return with(MultiPolygon, false, text)
	case `POINT`:
		return with(Point, false, text)
	case `POLYGON`:
		return with(Polygon, false, text)
	case `JSON`:
		return with(Json, false, text)
	case `GEOMETRY`:
		return with(Geometry, false, text)
	}

	v.panicWithExpr(ctx.GetTypeName(), "invalid data type: "+text)
	return nil
}

// visitLongVarcharDataType visits a parse tree produced by MySqlParser#longVarcharDataType.
func (v *visitor) visitLongVarcharDataType(_ *gen.LongVarcharDataTypeContext) DataType {
	v.trace("VisitLongVarcharDataType")
	return with(LongVarChar, false, "")
}

// visitLongVarbinaryDataType visits a parse tree produced by MySqlParser#longVarbinaryDataType.
func (v *visitor) visitLongVarbinaryDataType(_ *gen.LongVarbinaryDataTypeContext) DataType {
	v.trace("VisitLongVarbinaryDataType")
	return with(LongVarBinary, false, "")
}
