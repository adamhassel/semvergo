package semver

import (
	"cmp"
	"strconv"
	"strings"
)

// pre is the entire pre-release version pre. It is a dot-separated list of identifiers
type pre string

// identifier is a single component of the prerelease label
type identifier string

// identifiers is the parsed list of individual pre label components
type identifiers []identifier

func (id identifiers) String() string {
	sl := make([]string, 0, len(id))
	for _, s := range id {
		sl = append(sl, string(s))
	}
	return strings.Join(sl, ".")
}

type label struct {
	raw string
	t   identifiers
}

func (l label) String() string {
	return l.raw
}

func max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// idMax  returns the largest (most significant) of two identifiers
func idMax(a, b identifier) identifier {
	// is one or both pre numeric?
	ia, na := a.numeric()
	ib, nb := b.numeric()

	switch {
	case na && nb:
		return identifier(strconv.Itoa(max(ia, ib)))
	case na:
		return b
	case nb:
		return a
	}
	return max(a, b)
}

// idsMax returns the largest (most significant) of two sets of identifiers
func idsMax(a, b identifiers) identifiers {
	// nothing is always higher priority than something. It's called a 'pre-release label' after all
	if len(a) == 0 || len(b) == 0 {
		return identifiers{}
	}

	sm := a // smaller (element count)
	lg := b // larger (element count

	if len(a) > len(b) {
		sm = b
		lg = a
	}

	// compare each element of the smaller set to the equivalent in the larger set
	for i := range sm {
		if a[i] == b[i] {
			continue
		}
		if idMax(a[i], b[i]) == a[i] {
			return a
		}
		if idMax(a[i], b[i]) == b[i] {
			return b
		}
	}
	// a and b are equivalent. Return the longest
	return lg
}

func (t pre) components() identifiers {
	tmp := strings.Split(string(t), ".")

	// if t is the empty string, strings.Split will return a slice of length 1 with the empty string as the only member. But we'd rather have an empty list
	if len(tmp) == 1 && tmp[0] == "" {
		return identifiers{}
	}
	ids := make(identifiers, 0, len(tmp))
	for _, id := range tmp {
		ids = append(ids, identifier(id))
	}
	return ids
}

func (i identifier) numeric() (int, bool) {
	n, err := strconv.Atoi(string(i))
	return n, err == nil
}

// MaxSlice returns the max (meaning most significant, higher precedence) of p and q, according to Semver
func (p pre) Max(q pre) string {
	a := p.components()
	b := q.components()
	m := idsMax(a, b)
	return m.String()
}

func MaxLabel(a, b string) string {
	return pre(a).Max(pre(b))
}
