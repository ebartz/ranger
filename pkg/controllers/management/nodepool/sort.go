package nodepool

import (
	v3 "github.com/ranger/ranger/pkg/generated/norman/management.cattle.io/v3"
)

//implement natural sort of a slice of nodeHostnames

type byHostname []*v3.Node

func (n byHostname) Len() int      { return len(n) }
func (n byHostname) Swap(i, j int) { n[i], n[j] = n[j], n[i] }

func (n byHostname) Less(i, j int) bool {
	s := n[i].Spec.RequestedHostname
	t := n[j].Spec.RequestedHostname
	return NaturalLess(s, t)
}

// from https://github.com/fvbommel/util/blob/efcd4e0f97874370259c7d93e12aad57911dea81/sortorder/natsort.go
func isdigit(b byte) bool { return '0' <= b && b <= '9' }

// NaturalLess compares two strings using natural ordering. This means that e.g.
// "abc2" < "abc12".
//
// Non-digit sequences and numbers are compared separately. The former are
// compared bytewise, while the latter are compared numerically (except that
// the number of leading zeros is used as a tie-breaker, so e.g. "2" < "02")
//
// Limitation: only ASCII digits (0-9) are considered.
func NaturalLess(str1, str2 string) bool {
	idx1, idx2 := 0, 0
	for idx1 < len(str1) && idx2 < len(str2) {
		c1, c2 := str1[idx1], str2[idx2]
		dig1, dig2 := isdigit(c1), isdigit(c2)
		switch {
		case dig1 != dig2: // Digits before other characters.
			return dig1 // True if LHS is a digit, false if the RHS is one.
		case !dig1: // && !dig2, because dig1 == dig2
			// UTF-8 compares bytewise-lexicographically, no need to decode
			// codepoints.
			if c1 != c2 {
				return c1 < c2
			}
			idx1++
			idx2++
		default: // Digits
			// Eat zeros.
			for ; idx1 < len(str1) && str1[idx1] == '0'; idx1++ {
			}
			for ; idx2 < len(str2) && str2[idx2] == '0'; idx2++ {
			}
			// Eat all digits.
			nonZero1, nonZero2 := idx1, idx2
			for ; idx1 < len(str1) && isdigit(str1[idx1]); idx1++ {
			}
			for ; idx2 < len(str2) && isdigit(str2[idx2]); idx2++ {
			}
			// If lengths of numbers with non-zero prefix differ, the shorter
			// one is less.
			if len1, len2 := idx1-nonZero1, idx2-nonZero2; len1 != len2 {
				return len1 < len2
			}
			// If they're equal, string comparison is correct.
			if nr1, nr2 := str1[nonZero1:idx1], str2[nonZero2:idx2]; nr1 != nr2 {
				return nr1 < nr2
			}
			// Otherwise, the one with less zeros is less.
			// Because everything up to the number is equal, comparing the index
			// after the zeros is sufficient.
			if nonZero1 != nonZero2 {
				return nonZero1 < nonZero2
			}
		}
		// They're identical so far, so continue comparing.
	}
	// So far they are identical. At least one is ended. If the other continues,
	// it sorts last.
	return len(str1) < len(str2)
}
