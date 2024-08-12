package semver

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// SEMVERRE is the semantic versioning regexp
const SEMVERRE = `(\d+)\.(\d+)\.(\d+)`
const SEMVERRE_PRE_RE = `(.*)` + SEMVERRE
const SEMVERRE_SUF_RE = SEMVERRE + `(.*)`

var ErrParsingError = errors.New("parsing error")

// SemVer is a struct representing a semantic version
type SemVer struct {
	major  uint
	minor  uint
	patch  uint
	prefix string
	presep string
	suffix string
	sufsep string
}

// ByVersionDescending sorts versions in descending order. Suffixes are parsed according to semver
type ByVersionDescending []SemVer

func (a ByVersionDescending) Len() int      { return len(a) }
func (a ByVersionDescending) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByVersionDescending) Less(i, j int) bool {
	return Max(a[i], a[j]) == a[i]
}

// MaxSlice returns the highest version in a list
func MaxSlice(v []SemVer) SemVer {
	if len(v) == 0 {
		return SemVer{}
	}
	sort.Sort(ByVersionDescending(v))
	return v[0]
}

func Max(a, b SemVer) SemVer {
	if a.major > b.major || a.minor > b.minor || a.patch > b.patch {
		return a
	}
	if a.major < b.major || a.minor < b.minor || a.patch < b.patch {
		return b
	}
	if MaxLabel(a.suffix, b.suffix) == a.suffix {
		return a
	}
	return b
}

// Version returns the semantic version, without any prefixes or suffixes
func (s SemVer) Version() string {
	return fmt.Sprintf("%d.%d.%d", s.major, s.minor, s.patch)
}

// PreSuffix returns the prefix and suffix of s
func (s SemVer) PreSuffix() (string, string) {
	return s.prefix, s.suffix
}

// String returns the complete string semantic version
func (s SemVer) String() string {
	v := strings.Builder{}
	if s.prefix != "" {
		v.WriteString(s.prefix)
		v.WriteString(s.presep)
	}
	v.WriteString(s.Version())
	if s.suffix != "" {
		v.WriteString(s.sufsep)
		v.WriteString(s.suffix)
	}
	return v.String()
}

// IncrementMajor increments the major version and sets minor and patch to zero
func (s *SemVer) IncrementMajor() {
	s.major++
	s.minor, s.patch = 0, 0
}

// IncrementMinor increments the minor version and sets patch to zero
func (s *SemVer) IncrementMinor() {
	s.minor++
	s.patch = 0
}

// IncrementPatch increments the patch version
func (s *SemVer) IncrementPatch() {
	s.patch++
}

func (s *SemVer) Suffix(suffix string) {
	s.suffix = suffix
}

func (s *SemVer) Prefix(prefix string) {
	s.prefix = prefix
}

func (s *SemVer) Presep(pre string) {
	s.presep = pre
}

func (s *SemVer) Sufsep(suf string) {
	s.sufsep = suf
}

// Parse returns a new SemVer, or an error is a parsing error occurs
func Parse(s string) (SemVer, error) {
	var sv SemVer
	var err error
	sv.major, sv.minor, sv.patch, err = version(s)
	if err != nil {
		return SemVer{}, err
	}
	sv.prefix, err = prefix(s)
	if err != nil {
		return SemVer{}, err
	}
	sv.suffix, err = suffix(s)
	if err != nil {
		return SemVer{}, err
	}
	return sv, nil
}

func ParseSeparated(s, prefixSeparator, suffixSeparator string) (SemVer, error) {
	sv, err := Parse(s)
	if err != nil {
		return SemVer{}, err
	}
	sv.separators(prefixSeparator, suffixSeparator)
	return sv, nil
}

// Separators will handle using specific separator characters to separate prefixes, version strings and suffixes
func (s *SemVer) separators(pre, suf string) {
	s.presep = pre
	s.sufsep = suf

	func() {
		if pre != "" && len(s.prefix) < 1 {
			return
		}
		// handle if there's no prefix, only a prefix character. Which would be weird, but hey (ex: prefix separator: "-", version string: "-1.2.3")
		if s.presep == s.prefix {
			s.prefix = ""
			return
		}
		if s.prefix[len(s.prefix)-len(s.presep):] != s.presep {
			return
		}
		s.prefix = s.prefix[0 : len(s.prefix)-len(s.presep)]
	}()

	func() {
		if suf != "" && len(s.suffix) < 1 {
			return
		}
		if s.sufsep == s.suffix {
			s.suffix = ""
			return
		}
		if s.suffix[:len(s.sufsep)] != s.sufsep {
			return
		}
		s.suffix = s.suffix[len(s.sufsep):]
	}()

}

// version extracts major,minor and patch versions from a semantic version string
func version(s string) (uint, uint, uint, error) {
	re := regexp.MustCompile(SEMVERRE)
	matches := re.FindStringSubmatch(s)
	if len(matches) != 4 {
		return 0, 0, 0, ErrParsingError
	}
	rv := make([]uint, 3)
	for i, v := range matches[1:] {
		p, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return 0, 0, 0, ErrParsingError
		}
		rv[i] = uint(p)
	}
	return rv[0], rv[1], rv[2], nil
}

// prefix extracts a given prefix from a semver string
func prefix(s string) (string, error) {
	re := regexp.MustCompile(SEMVERRE_PRE_RE)
	matches := re.FindStringSubmatch(s)
	if len(matches) < 5 {
		return "", ErrParsingError
	}
	return matches[1], nil
}

// suffix returns any suffix from a semver string
func suffix(s string) (string, error) {
	re := regexp.MustCompile(SEMVERRE_SUF_RE)
	matches := re.FindStringSubmatch(s)
	if len(matches) < 5 {
		return "", ErrParsingError
	}
	return matches[len(matches)-1], nil
}
