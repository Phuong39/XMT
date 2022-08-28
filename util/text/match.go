// Copyright (C) 2020 - 2022 iDigitalFlame
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package text

import (
	"sync"

	"github.com/iDigitalFlame/xmt/util"
)

var (
	// MatchAny is a Matcher that can be used to match any string or byte array.
	MatchAny anyRegexp
	// MatchNone is a Matcher that can be used to NOT match any string or byte
	// array.
	MatchNone falseRegexp

	builders = sync.Pool{
		New: func() any {
			return new(util.Builder)
		},
	}
)

// String is a wrapper for strings to support the Stringer interface.
type String string

// Matcher is an alias of a string that can contain specific variable
// instructions to be replaced when calling the 'String' function. This alias
// provides the 'Match' function, which returns a Regexp struct that will match
// any value generated.
//
// # Matcher Verb Guide
//
// <N> indicates a non-negative number value that is REQUIRED.
// [N] indicated a non-negative number value that is OPTIONAL.
//     - If omitted, a random positive 8 to 32bit value will be used.
//
// If the number is followed bfy a 'f', this will FORCE the count and will use it
// directly instead of a range to N. Otherwise, a [1, N] (inclusive) value will
// be generated to be used instead.
//
// In both cases, the '<', '[', ']', and '>' are only used to indicate usage,
// and not actually used in the string value.
//
//  | Verb   | Description                                         | RegEx            |
//  | ====== | =================================================== | ================ |
//  | %<N>n  | 1 to N count of random single-digit numbers         | [0-9]{1,N}       |
//  | %<N>fn | N count of random single-digit numbers              | [0-9]{N}         |
//  | %<N>c  | 1 to N count of random ASCII non-number characters  | [a-zA-Z]{1,N}    |
//  | %<N>fc | N count of random ASCII non-number characters       | [a-zA-Z]{N}      |
//  | %<N>u  | 1 to N count of random ASCII uppercase characters   | [A-Z]{1,N}       |
//  | %<N>fu | N count of random ASCII uppercase characters        | [A-Z]{N}         |
//  | %<N>l  | 1 to N count of random ASCII lowercase characters   | [a-z]{1,N}       |
//  | %<N>fl | N count of random ASCII lowercase characters        | [a-z]{n}         |
//  | %s     | Random 1 to 256 count of random ASCII characters    | ([a-zA-Z0-9]+)   |
//  | %[N]s  | 1 to N (or random) count of random ASCII characters | [a-zA-Z0-9]{1,N} |
//  | %<N>fs | N count of random ASCII characters                  | [a-zA-Z0-9]{N}   |
//  | %d     | String literal number 0 to 4,294,967,296            | ([0-9]+)         |
//  | %<N>d  | String literal number 0 to N                        | ([0-9]+)         |
//  | %<N>fd | String literal number N                             | ([0-9]+)         |
//  | %h     | Hex string literal number 0 to 4,294,967,296        | ([a-fA-F0-9]+)   |
//  | %<N>h  | Hex string literal number 0 to N                    | ([a-fA-F0-9]+)   |
//  | %<N>fh | Hex string literal number N                         | ([a-fA-F0-9]+)   |
//
// All other values are ignored and directly added to the resulting string value.
type Matcher string
type anyRegexp struct{}

// Regexp is a compatibility interface that is used to allow for UnMatch values
// to be used lacking the ability of RE2 used in Golang to do negative lookahead.
type Regexp interface {
	String() string
	Match([]byte) bool
	MatchString(string) bool
}
type falseRegexp struct{}

// Raw returns the raw string value of this Matcher without preforming any
// replacements.
func (s Matcher) Raw() string {
	return string(s)
}

// String returns the string value of itself.
func (s String) String() string {
	return string(s)
}

// Match returns a valid Regexp struct that is guaranteed to match any string
// generated by the  Matcher's 'String' function.
func (s Matcher) Match() Regexp {
	return s.MatchEx(true)
}
func (anyRegexp) String() string {
	return "*"
}

// UnMatch returns a valid Regexp struct that is guaranteed to not match any
// string generated by the Matcher's 'String' function.
func (s Matcher) UnMatch() Regexp {
	return s.MatchEx(false)
}
func (falseRegexp) String() string {
	return "false"
}
func (anyRegexp) Match(_ []byte) bool {
	return true
}
func (falseRegexp) Match(_ []byte) bool {
	return false
}
func (anyRegexp) MatchString(_ string) bool {
	return true
}
func (falseRegexp) MatchString(_ string) bool {
	return false
}
