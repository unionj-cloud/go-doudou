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

package console

import "testing"

func TestConsole(t *testing.T) {
	console := NewColorConsole()
	console.Info("info", "info")
	console.InfoF("%v", "info")

	console.Debug("debug", "debug")
	console.DebugF("%v", "debug")

	console.Warning("warning", "warning")
	console.WarningF("%v", "warning")

	console.Error("error")
	console.ErrorF("%v", "error")
}
