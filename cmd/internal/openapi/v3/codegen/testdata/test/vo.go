package test

// ChartReportParagraph 图表类型段落
type ChartReportParagraph struct {
}

type ChartReportParagraphStyle struct {
	// 居左，居中，居右
	Alignment *string `json:"alignment,omitempty" url:"alignment"`
	// 每个系列的前景色，图例的颜色，柱子的颜色，线条的颜色等
	ForeColor map[string]string `json:"foreColor,omitempty" url:"foreColor"`
	// word里面的图表的高，单位厘米
	Height *float64 `json:"height,omitempty" url:"height"`
	// word里面的图表的宽，单位厘米
	Width *float64 `json:"width,omitempty" url:"width"`
}

// ImageReportParagraph 图片类型段落
type ImageReportParagraph struct {
}

type ImageReportParagraphStyle struct {
	// 居左，居中，居右
	Alignment *string `json:"alignment,omitempty" url:"alignment"`

	Height *float64 `json:"height,omitempty" url:"height"`

	Separator *string `json:"separator,omitempty" url:"separator"`

	Width *float64 `json:"width,omitempty" url:"width"`
}

type MergeFieldData struct {
	Fields []string `json:"fields,omitempty" url:"fields"`

	Values []interface{} `json:"values,omitempty" url:"values"`
}

// ParagraphWrapper 报告段落包装类
type ParagraphWrapper struct {
	// 报告段落类型
	// required
	Type string `json:"type,omitempty" url:"type"`

	Value interface{} `json:"value,omitempty" url:"value"`
}

type ReportPage struct {
	// 页面下边距，单位厘米
	BottomMargin *float64 `json:"bottomMargin,omitempty" url:"bottomMargin"`
	// 页面左边距，单位厘米
	LeftMargin *float64 `json:"leftMargin,omitempty" url:"leftMargin"`
	// 页面是横向还是纵向
	OrientationType *string `json:"orientationType,omitempty" url:"orientationType"`
	// 页面纸张大小
	PaperSizeType *string `json:"paperSizeType,omitempty" url:"paperSizeType"`
	// 页面右边距，单位厘米
	RightMargin *float64 `json:"rightMargin,omitempty" url:"rightMargin"`
	// 页面上边距，单位厘米
	TopMargin *float64 `json:"topMargin,omitempty" url:"topMargin"`
}

// ReportParagraph 报告段落内容
type ReportParagraph struct {
	Bookmark *string `json:"bookmark,omitempty" url:"bookmark"`
}

type RequestPayload struct {
	ParagraphList []ParagraphWrapper `json:"paragraphList,omitempty" url:"paragraphList"`

	ReportFileName *string `json:"reportFileName,omitempty" url:"reportFileName"`

	ReportPage *ReportPage `json:"reportPage,omitempty" url:"reportPage"`

	StorageMode *string `json:"storageMode,omitempty" url:"storageMode"`

	TemplateUrl *string `json:"templateUrl,omitempty" url:"templateUrl"`
}

type ResultInteger struct {
	// 返回标记：成功标记=0，失败标记=1
	Code *int `json:"code,omitempty" url:"code"`
	// 数据
	Data *int `json:"data,omitempty" url:"data"`
	// 返回信息
	Msg *string `json:"msg,omitempty" url:"msg"`
}

type ResultListWordTemplateSubstitution struct {
	// 返回标记：成功标记=0，失败标记=1
	Code *int `json:"code,omitempty" url:"code"`
	// 数据
	Data []WordTemplateSubstitution `json:"data,omitempty" url:"data"`
	// 返回信息
	Msg *string `json:"msg,omitempty" url:"msg"`
}

type ResultString struct {
	// 返回标记：成功标记=0，失败标记=1
	Code *int `json:"code,omitempty" url:"code"`
	// 数据
	Data *string `json:"data,omitempty" url:"data"`
	// 返回信息
	Msg *string `json:"msg,omitempty" url:"msg"`
}

// TableReportParagraph 表格类型段落
type TableReportParagraph struct {
}

// TextReportParagraph 文本类型段落
type TextReportParagraph struct {
}

// TextReportParagraphFont 设置字体样式
type TextReportParagraphFont struct {
	// 中文字号，优先级高于fontSize
	ChineseFontSize *string `json:"chineseFontSize,omitempty" url:"chineseFontSize"`
	// 字体颜色，只支持黑色，红色和黄色
	FontColor *string `json:"fontColor,omitempty" url:"fontColor"`
	// 字体
	FontFamily *string `json:"fontFamily,omitempty" url:"fontFamily"`
	// 字号，默认16磅，即三号字体
	FontSize *float64 `json:"fontSize,omitempty" url:"fontSize"`
	// 是否加粗，默认false
	IsBold *bool `json:"isBold,omitempty" url:"isBold"`
}

// TextReportParagraphSentence 文本类型句子短语
type TextReportParagraphSentence struct {
	// 文本
	Content *string `json:"content,omitempty" url:"content"`

	Font *TextReportParagraphFont `json:"font,omitempty" url:"font"`
	// inline默认值为true，表示不换行
	Inline *bool `json:"inline,omitempty" url:"inline"`
}

type TextReportParagraphStyle struct {
	// 居左，居中，居右
	Alignment *string `json:"alignment,omitempty" url:"alignment"`
	// 是否清除原有段落格式
	ClearOldStyle *bool `json:"clearOldStyle,omitempty" url:"clearOldStyle"`

	Font *TextReportParagraphFont `json:"font,omitempty" url:"font"`
	// 设置首行缩进，单位磅
	Indent *float64 `json:"indent,omitempty" url:"indent"`
	// 如果inline为true，表示不换行
	Inline *bool `json:"inline,omitempty" url:"inline"`
	// 设置行距，单位磅，默认单倍行距，即12磅
	LineSpacing *float64 `json:"lineSpacing,omitempty" url:"lineSpacing"`
}

type Word2HtmlRequestPayload struct {
	DownloadUrl *string `json:"downloadUrl,omitempty" url:"downloadUrl"`
}

// WordTemplateSubstitution 模板变量
type WordTemplateSubstitution struct {
	// 完整变量字符串，用于前端展示
	Display *string `json:"display,omitempty" url:"display"`
	// 变量名
	Name *string `json:"name,omitempty" url:"name"`
	// 类型
	Type *string `json:"type,omitempty" url:"type"`
}
