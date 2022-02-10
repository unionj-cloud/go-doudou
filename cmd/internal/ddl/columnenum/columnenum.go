package columnenum

// ColumnType database data types
type ColumnType string

const (
	// BitType bit
	BitType ColumnType = "BIT"
	// TextType text
	TextType ColumnType = "TEXT"
	// BlobType blob
	BlobType ColumnType = "BLOB"
	// DateType date
	DateType ColumnType = "DATE"
	// DatetimeType datatime
	DatetimeType ColumnType = "DATETIME"
	// DecimalType decimal
	DecimalType ColumnType = "DECIMAL"
	// DoubleType double
	DoubleType ColumnType = "DOUBLE"
	// EnumType enum
	EnumType ColumnType = "ENUM"
	// FloatType float
	FloatType ColumnType = "FLOAT"
	// GeometryType geometry
	GeometryType ColumnType = "GEOMETRY"
	// MediumintType medium int
	MediumintType ColumnType = "MEDIUMINT"
	// JSONType json
	JSONType ColumnType = "JSON"
	// IntType int
	IntType ColumnType = "INT"
	// LongtextType long text
	LongtextType ColumnType = "LONGTEXT"
	// LongblobType long blob
	LongblobType ColumnType = "LONGBLOB"
	// BigintType big int
	BigintType ColumnType = "BIGINT"
	// MediumtextType medium text
	MediumtextType ColumnType = "MEDIUMTEXT"
	// MediumblobType medium blob
	MediumblobType ColumnType = "MEDIUMBLOB"
	// SmallintType small int
	SmallintType ColumnType = "SMALLINT"
	// TinyintType tiny int
	TinyintType ColumnType = "TINYINT"
	// VarcharType varchar
	VarcharType ColumnType = "VARCHAR(255)"
)
