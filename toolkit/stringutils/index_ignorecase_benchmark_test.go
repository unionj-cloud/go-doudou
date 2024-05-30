// SPDX-License-Identifier: MIT OR Unlicense

package stringutils

import (
	"regexp"
	"testing"
)

func BenchmarkFindAllIndexCaseInsensitive(b *testing.B) {
	r := regexp.MustCompile(`(?i)test`)
	haystack := []byte(testMatchEndCase)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 1 {
			b.Error("Expected single match")
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseCaseInsensitive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(testMatchEndCase, "test", -1)

		if len(matches) != 1 {
			b.Error("Expected single match")
		}
	}
}

func BenchmarkFindAllIndexLargeCaseInsensitive(b *testing.B) {
	r := regexp.MustCompile(`(?i)test`)
	haystack := []byte(testMatchEndCaseLarge)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 1 {
			b.Error("Expected single match")
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseLargeCaseInsensitive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(testMatchEndCaseLarge, "test", -1)
		if len(matches) != 1 {
			b.Error("Expected single match")
		}
	}
}

func BenchmarkFindAllIndexUnicodeCaseInsensitive(b *testing.B) {
	r := regexp.MustCompile(`(?i)test`)
	haystack := []byte(testUnicodeMatchEndCase)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 1 {
			b.Error("Expected single match")
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseUnicodeCaseInsensitive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(testUnicodeMatchEndCase, "test", -1)
		if len(matches) != 1 {
			b.Error("Expected single match")
		}
	}
}

func BenchmarkFindAllIndexUnicodeLargeCaseInsensitive(b *testing.B) {
	r := regexp.MustCompile(`(?i)test`)
	haystack := []byte(testUnicodeMatchEndCaseLarge)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 1 {
			b.Error("Expected single match")
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseUnicodeLargeCaseInsensitive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(testUnicodeMatchEndCaseLarge, "test", -1)
		if len(matches) != 1 {
			b.Error("Expected single match")
		}
	}
}

// This benchmark simulates a bad case of there being many
// partial matches where the first character in the needle
// can be found throughout the haystack
func BenchmarkFindAllIndexManyPartialMatchesCaseInsensitive(b *testing.B) {
	r := regexp.MustCompile(`(?i)1test`)
	haystack := []byte(testMatchEndCase)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 1 {
			b.Error("Expected single match")
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseManyPartialMatchesCaseInsensitive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(testMatchEndCase, "1test", -1)
		if len(matches) != 1 {
			b.Error("Expected single match")
		}
	}
}

// This benchmark simulates a bad case of there being many
// partial matches where the first character in the needle
// can be found throughout the haystack
func BenchmarkFindAllIndexUnicodeManyPartialMatchesCaseInsensitive(b *testing.B) {
	r := regexp.MustCompile(`(?i)Ⱥtest`)
	haystack := []byte(testUnicodeMatchEndCase)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 1 {
			b.Error("Expected single match")
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseUnicodeManyPartialMatchesCaseInsensitive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(testUnicodeMatchEndCase, "Ⱥtest", -1)
		if len(matches) != 1 {
			b.Error("Expected single match")
		}
	}
}

func BenchmarkFindAllIndexUnicodeCaseInsensitiveVeryLarge(b *testing.B) {
	var large string
	for i := 0; i <= 100; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	r := regexp.MustCompile(`(?i)Ⱥtest`)
	haystack := []byte(large)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 101 {
			b.Error("Expected single match got", len(matches))
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseUnicodeCaseInsensitiveVeryLarge(b *testing.B) {
	var large string
	for i := 0; i <= 100; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(large, "Ⱥtest", -1)
		if len(matches) != 101 {
			b.Error("Expected single match got", len(matches))
		}
	}
}

func BenchmarkFindAllIndexFoldingCaseInsensitiveVeryLarge(b *testing.B) {
	var large string
	for i := 0; i <= 100; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	r := regexp.MustCompile(`(?i)ſ`)
	haystack := []byte(large)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 101 {
			b.Error("Expected single match got", len(matches))
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseFoldingCaseInsensitiveVeryLarge(b *testing.B) {
	var large string
	for i := 0; i <= 100; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(large, "ſ", -1)
		if len(matches) != 101 {
			b.Error("Expected single match got", len(matches))
		}
	}
}

func BenchmarkFindAllIndexFoldingCaseInsensitiveNeedle1(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	r := regexp.MustCompile(`(?i)a`)
	haystack := []byte(large)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseFoldingCaseInsensitiveNeedle1(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(large, "a", -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkFindAllIndexFoldingCaseInsensitiveNeedle2(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	r := regexp.MustCompile(`(?i)aa`)
	haystack := []byte(large)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseFoldingCaseInsensitiveNeedle2(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(large, "aa", -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkFindAllIndexFoldingCaseInsensitiveNeedle3(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	r := regexp.MustCompile(`(?i)aaa`)
	haystack := []byte(large)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseFoldingCaseInsensitiveNeedle3(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(large, "aaa", -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkFindAllIndexFoldingCaseInsensitiveNeedle4(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	r := regexp.MustCompile(`(?i)aaaa`)
	haystack := []byte(large)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseFoldingCaseInsensitiveNeedle4(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(large, "aaaa", -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkFindAllIndexFoldingCaseInsensitiveNeedle5(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	r := regexp.MustCompile(`(?i)aaaaa`)
	haystack := []byte(large)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseFoldingCaseInsensitiveNeedle5(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(large, "aaaaa", -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkFindAllIndexFoldingCaseInsensitiveNeedle6(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	r := regexp.MustCompile(`(?i)aaaaaa`)
	haystack := []byte(large)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseFoldingCaseInsensitiveNeedle6(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(large, "aaaaaa", -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkFindAllIndexFoldingCaseInsensitiveNeedle7(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	r := regexp.MustCompile(`(?i)aaaaaaa`)
	haystack := []byte(large)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseFoldingCaseInsensitiveNeedle7(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(large, "aaaaaaa", -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkFindAllIndexFoldingCaseInsensitiveNeedle8(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	r := regexp.MustCompile(`(?i)aaaaaaaa`)
	haystack := []byte(large)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseFoldingCaseInsensitiveNeedle8(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(large, "aaaaaaaa", -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkFindAllIndexFoldingCaseInsensitiveNeedle9(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	r := regexp.MustCompile(`(?i)aaaaaaaaa`)
	haystack := []byte(large)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseFoldingCaseInsensitiveNeedle9(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(large, "aaaaaaaaa", -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkFindAllIndexFoldingCaseInsensitiveNeedle10(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	r := regexp.MustCompile(`(?i)aaaaaaaaaa`)
	haystack := []byte(large)

	for i := 0; i < b.N; i++ {
		matches := r.FindAllIndex(haystack, -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}

func BenchmarkIndexesAllIgnoreCaseFoldingCaseInsensitiveNeedle10(b *testing.B) {
	var large string
	for i := 0; i <= 10; i++ {
		large += testUnicodeMatchEndCaseLarge
	}

	for i := 0; i < b.N; i++ {
		matches := IndexAllIgnoreCase(large, "aaaaaaaaaa", -1)
		if len(matches) != 0 {
			b.Error("Expected no matches got", len(matches))
		}
	}
}
