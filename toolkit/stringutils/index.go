// SPDX-License-Identifier: MIT OR Unlicense

package stringutils

import (
	"math"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"
)

// IndexAll extracts all of the locations of a string inside another string
// up-to the defined limit and does so without regular expressions
// which makes it faster than FindAllIndex in most situations while
// not being any slower. It performs worst when working against random
// data.
//
// Some benchmark results to illustrate the point (find more in index_benchmark_test.go)
//
// BenchmarkFindAllIndex-8                         2458844	       480.0 ns/op
// BenchmarkIndexAll-8                            14819680	        79.6 ns/op
//
// For pure literal searches IE no regular expression logic this method
// is a drop in replacement for re.FindAllIndex but generally much faster.
//
// Similar to how FindAllIndex the limit option can be passed -1
// to get all matches.
//
// Note that this method is explicitly case sensitive in its matching.
// A return value of nil indicates no match.
func IndexAll(haystack string, needle string, limit int) [][]int {
	// The below needed to avoid timeout crash found using go-fuzz
	if len(haystack) == 0 || len(needle) == 0 {
		return nil
	}

	// Return contains a slice of slices where index 0 is the location of the match in bytes
	// and index 1 contains the end location in bytes of the match
	var locs [][]int

	// Perform the first search outside the main loop to make the method
	// easier to understand
	searchText := haystack
	offSet := 0
	loc := strings.Index(searchText, needle)

	if limit <= -1 {
		// Similar to how regex FindAllString works
		// if we have -1 as the limit set to max to
		//  try to get everything
		limit = math.MaxInt32
	} else {
		// Increment by one because we do count++ at the start of the loop
		// and as such there is a off by 1 error in the return otherwise
		limit++
	}

	var count int
	for loc != -1 {
		count++

		if count == limit {
			break
		}

		// trim off the portion we already searched, and look from there
		searchText = searchText[loc+len(needle):]
		locs = append(locs, []int{loc + offSet, loc + offSet + len(needle)})

		// We need to keep the offset of the match so we continue searching
		offSet += loc + len(needle)

		// strings.Index does checks of if the string is empty so we don't need
		// to explicitly do it ourselves
		loc = strings.Index(searchText, needle)
	}

	// Retain compatibility with FindAllIndex method
	if len(locs) == 0 {
		return nil
	}

	return locs
}

// if the IndexAllIgnoreCase method is called frequently with the same patterns
// (which is a common case) this is here to speed up the case permutations
// it is limited to a size of 10 so it never gets that large but really
// allows things to run faster
var _permuteCache = map[string][]string{}
var _permuteCacheLock = sync.Mutex{}

// CacheSize this is public so it can be modified depending on project needs
// you can increase this value to cache more of the case permutations which
// can improve performance if doing the same searches over and over
var CacheSize = 10

// IndexAllIgnoreCase extracts all of the locations of a string inside another string
// up-to the defined limit. It is designed to be faster than uses of FindAllIndex with
// case insensitive matching enabled, by looking for string literals first and then
// checking for exact matches. It also does so in a unicode aware way such that a search
// for S will search for S s and ſ which a simple strings.ToLower over the haystack
// and the needle will not.
//
// The result is the ability to search for literals without hitting the regex engine
// which can at times be horribly slow. This by contrast is much faster. See
// index_ignorecase_benchmark_test.go for some head to head results. Generally
// so long as we aren't dealing with random data this method should be considerably
// faster (in some cases thousands of times) or just as fast. Of course it cannot
// do regular expressions, but that's fine.
//
// For pure literal searches IE no regular expression logic this method
// is a drop in replacement for re.FindAllIndex but generally much faster.
func IndexAllIgnoreCase(haystack string, needle string, limit int) [][]int {
	// The below needed to avoid timeout crash found using go-fuzz
	if len(haystack) == 0 || len(needle) == 0 {
		return nil
	}

	// One of the problems with finding locations ignoring case is that
	// the different case representations can have different byte counts
	// which means the locations using strings or bytes Index can be off
	// if you apply strings.ToLower to your haystack then use strings.Index.
	//
	// This can be overcome using regular expressions but suffers the penalty
	// of hitting the regex engine and paying the price of case
	// insensitive match there.
	//
	// This method tries something else which is used by some regex engines
	// such as the one in Rust where given a str literal if you get
	// all the case options of that such as turning foo into foo Foo fOo FOo foO FoO fOO FOO
	// and then use Boyer-Moore or some such for those. Of course using something
	// like Aho-Corasick or Rabin-Karp to get multi match would be a better idea so you
	// can match all of the input in one pass.
	//
	// If the needle is over some amount of characters long you chop off the first few
	// and then search for those. However this means you are not finding actual matches and as such
	// you the need to validate a potential match after you have found one.
	// The confirmation match is done in a loop because for some literals regular expression
	// is still to slow, although for most its a valid option.
	var locs [][]int

	// Char limit is the cut-off where we switch from all case permutations
	// to just the first 3 and then check for an actual match
	// in my tests 3 speeds things up the most against test data
	// of many famous books concatenated together and large
	// amounts of data from /dev/urandom
	var charLimit = 3

	if utf8.RuneCountInString(needle) <= charLimit {
		// We are below the limit we set, so get all the search
		// terms and search for that

		// Generally speaking I am against caches inside libraries but in this case...
		// when the IndexAllIgnoreCase method is called repeatedly it quite often
		// ends up performing case folding on the same thing over and over again which
		// can become the most expensive operation. So we keep a VERY small cache
		// to avoid that being an issue.
		_permuteCacheLock.Lock()
		searchTerms, ok := _permuteCache[needle]
		if !ok {
			if len(_permuteCache) > CacheSize {
				_permuteCache = map[string][]string{}
			}
			searchTerms = PermuteCaseFolding(needle)
			_permuteCache[needle] = searchTerms
		}
		_permuteCacheLock.Unlock()

		// This is using IndexAll in a loop which was faster than
		// any implementation of Aho-Corasick or Boyer-Moore I tried
		// but in theory Aho-Corasick / Rabin-Karp or even a modified
		// version of Boyer-Moore should be faster than this.
		// Especially since they should be able to do multiple comparisons
		// at the same time.
		// However after some investigation it turns out that this turns
		// into a fancy  vector instruction on AMD64 (which is all we care about)
		// and as such its pretty hard to beat.
		for _, term := range searchTerms {
			locs = append(locs, IndexAll(haystack, term, limit)...)
		}

		// if the limit is not -1 we need to sort and return the first X results so we maintain compatibility with how
		// FindAllIndex would work
		if limit > 0 && len(locs) > limit {

			// now sort the results to we can get the first X results
			// Now rank based on which ones are the best and sort them on that rank
			// then get the top amount and the surrounding lines
			sort.Slice(locs, func(i, j int) bool {
				return locs[i][0] < locs[j][0]
			})

			return locs[:limit]
		}
	} else {
		// Over the character limit so look for potential matches and only then check to find real ones

		// Note that we have to use runes here to avoid cutting bytes off so
		// cast things around to ensure it works
		needleRune := []rune(needle)

		// Generally speaking I am against caches inside libraries but in this case...
		// when the IndexAllIgnoreCase method is called repeatedly it quite often
		// ends up performing case folding on the same thing over and over again which
		// can become the most expensive operation. So we keep a VERY small cache
		// to avoid that being an issue.
		_permuteCacheLock.Lock()
		searchTerms, ok := _permuteCache[string(needleRune[:charLimit])]
		if !ok {
			if len(_permuteCache) > CacheSize {
				_permuteCache = map[string][]string{}
			}
			searchTerms = PermuteCaseFolding(string(needleRune[:charLimit]))
			_permuteCache[string(needleRune[:charLimit])] = searchTerms
		}
		_permuteCacheLock.Unlock()

		// This is using IndexAll in a loop which was faster than
		// any implementation of Aho-Corasick or Boyer-Moore I tried
		// but in theory Aho-Corasick / Rabin-Karp or even a modified
		// version of Boyer-Moore should be faster than this.
		// Especially since they should be able to do multiple comparisons
		// at the same time.
		// However after some investigation it turns out that this turns
		// into a fancy  vector instruction on AMD64 (which is all we care about)
		// and as such its pretty hard to beat.
		haystackRune := []rune(haystack)

		for _, term := range searchTerms {
			potentialMatches := IndexAll(haystack, term, -1)

			for _, match := range potentialMatches {
				// We have a potential match, so now see if it actually matches
				// by getting the actual value out of our haystack
				if len(haystackRune) < match[0]+len(needleRune) {
					continue
				}

				// Because the length of the needle might be different to what we just found as a match
				// based on byte size we add enough extra on the end to deal with the difference
				e := len(needle) + len(needle) - 1
				for match[0]+e > len(haystack) {
					e--
				}

				// Cut off the number at the end to the number we need which is the length of the needle runes
				toMatchRune := []rune(haystack[match[0] : match[0]+e])
				toMatchEnd := len(needleRune)
				if len(toMatchRune) < len(needleRune) {
					toMatchEnd = len(toMatchRune)
				}

				toMatch := toMatchRune[:toMatchEnd]

				// old logic here
				//toMatch = []rune(haystack[match[0] : match[0]+e])[:len(needleRune)]

				// what we need to do is iterate the runes of the haystack portion we are trying to
				// match and confirm that the same rune position is a actual match or case fold match
				// if they are keep looking, if they are not bail out as its not a real match
				isMatch := false
				for i := 0; i < len(toMatch); i++ {
					isMatch = false

					// Check against the actual term and if that's a match we can avoid folding
					// and doing those comparisons to hopefully save some CPU time
					if toMatch[i] == needleRune[i] {
						isMatch = true
					} else {
						// Not a match so case fold to actually check
						for _, j := range AllSimpleFold(toMatch[i]) {
							if j == needleRune[i] {
								isMatch = true
							}
						}
					}

					// Bail out as there is no point to continue checking at this point
					// as we found no match and there is no point burning more CPU checking
					if !isMatch {
						break
					}
				}

				if isMatch {
					// When we have confirmed a match we add it to our total
					// but adjust the positions to the match and the length of the
					// needle to ensure the byte count lines up
					locs = append(locs, []int{match[0], match[0] + len(string(toMatch))})
				}

			}
		}

		// if the limit is not -1 we need to sort and return the first X results so we maintain compatibility with how
		// FindAllIndex would work
		if limit > 0 && len(locs) > limit {

			// now sort the results to we can get the first X results
			// Now rank based on which ones are the best and sort them on that rank
			// then get the top amount and the surrounding lines
			sort.Slice(locs, func(i, j int) bool {
				return locs[i][0] < locs[j][0]
			})

			return locs[:limit]
		}
	}

	// Retain compatibility with FindAllIndex method
	if len(locs) == 0 {
		return nil
	}

	return locs
}
