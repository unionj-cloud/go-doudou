package columnenum

type ColumnType string

const (
	BitType        ColumnType = "BIT"
	TextType       ColumnType = "TEXT"
	BlobType       ColumnType = "BLOB"
	DateType       ColumnType = "DATE"
	DatetimeType   ColumnType = "DATETIME"
	DecimalType    ColumnType = "DECIMAL"
	DoubleType     ColumnType = "DOUBLE"
	EnumType       ColumnType = "ENUM"
	FloatType      ColumnType = "FLOAT"
	GeometryType   ColumnType = "GEOMETRY"
	MediumintType  ColumnType = "MEDIUMINT"
	JsonType       ColumnType = "JSON"
	IntType        ColumnType = "INT"
	LongtextType   ColumnType = "LONGTEXT"
	LongblobType   ColumnType = "LONGBLOB"
	BigintType     ColumnType = "BIGINT"
	MediumtextType ColumnType = "MEDIUMTEXT"
	MediumblobType ColumnType = "MEDIUMBLOB"
	SmallintType   ColumnType = "SMALLINT"
	TinyintType    ColumnType = "TINYINT"
	VarcharType    ColumnType = "VARCHAR(255)"
)
