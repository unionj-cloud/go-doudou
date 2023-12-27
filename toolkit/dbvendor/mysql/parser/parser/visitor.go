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
	"fmt"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/dbvendor/mysql/parser/console"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/dbvendor/mysql/parser/gen"
)

type visitor struct {
	gen.BaseMySqlParserVisitor
	prefix string
	debug  bool
	logger console.Console
}

func (v *visitor) trace(msg ...interface{}) {
	if v.debug {
		v.logger.Debug("Visit Trace: " + fmt.Sprint(msg...))
	}
}

func (v *visitor) panicWithExpr(expr Token, msg string) {
	if len(v.prefix) == 0 {
		err := fmt.Errorf("%v:%v %s", expr.GetLine(), expr.GetColumn(), msg)
		if v.debug {
			v.logger.Error(err)
		}

		panic(err)
		return
	}

	err := fmt.Errorf("%v line %v:%v %s", v.prefix, expr.GetLine(), expr.GetColumn(), msg)
	if v.debug {
		v.logger.Error(err)
	}

	panic(err)
}
