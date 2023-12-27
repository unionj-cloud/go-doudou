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
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type CaseChangingStream struct {
	antlr.CharStream

	upper bool
}

// newCaseChangingStream returns a new CaseChangingStream that forces
// all tokens read from the underlying stream to be either upper case
// or lower case based on the upper argument.
func newCaseChangingStream(in antlr.CharStream, upper bool) *CaseChangingStream {
	return &CaseChangingStream{in, upper}
}

// LA gets the value of the symbol at offset from the current position
// from the underlying CharStream and converts it to either upper case
// or lower case.
func (is *CaseChangingStream) LA(offset int) int {
	in := is.CharStream.LA(offset)
	if in < 0 {
		// Such as antlr.TokenEOF which is -1
		return in
	}
	if is.upper {
		return int(unicode.ToUpper(rune(in)))
	}
	return int(unicode.ToLower(rune(in)))
}
