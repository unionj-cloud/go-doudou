package test

type RecognizeCharacterResultVO struct {
	Probability float32 `json:"probability" url:"probability"`

	Text string `json:"text" url:"text"`

	TextRectangles TextRectanglesVO `json:"textRectangles" url:"textRectangles"`
}

type RecognizePdfResultVO struct {
	Angle int64 `json:"angle" url:"angle"`

	Height int64 `json:"height" url:"height"`

	OrgHeight int64 `json:"orgHeight" url:"orgHeight"`

	OrgWidth int64 `json:"orgWidth" url:"orgWidth"`

	PageIndex int64 `json:"pageIndex" url:"pageIndex"`

	Width int64 `json:"width" url:"width"`

	WordsInfo []RecognizePdfWordsInfoVO `json:"wordsInfo" url:"wordsInfo"`
}

type RecognizePdfWordsInfoPositionsVO struct {
	X int64 `json:"x" url:"x"`

	Y int64 `json:"y" url:"y"`
}

type RecognizePdfWordsInfoVO struct {
	Angle int64 `json:"angle" url:"angle"`

	Height int64 `json:"height" url:"height"`

	Positions []RecognizePdfWordsInfoPositionsVO `json:"positions" url:"positions"`

	Width int64 `json:"width" url:"width"`

	Word string `json:"word" url:"word"`

	X int64 `json:"x" url:"x"`

	Y int64 `json:"y" url:"y"`
}

type ResultListRecognizeCharacterResultVO struct {
	Code int `json:"code" url:"code"`

	Data []RecognizeCharacterResultVO `json:"data" url:"data"`

	Msg string `json:"msg" url:"msg"`
}

type ResultRecognizePdfResultVO struct {
	Code int `json:"code" url:"code"`

	Data RecognizePdfResultVO `json:"data" url:"data"`

	Msg string `json:"msg" url:"msg"`
}

type Resultstring struct {
	Code int `json:"code" url:"code"`

	Data string `json:"data" url:"data"`

	Msg string `json:"msg" url:"msg"`
}

type TextRectanglesVO struct {
	Angle int `json:"angle" url:"angle"`

	Height int `json:"height" url:"height"`

	Left int `json:"left" url:"left"`

	Top int `json:"top" url:"top"`

	Width int `json:"width" url:"width"`
}
