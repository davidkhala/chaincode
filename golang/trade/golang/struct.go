package golang

import "strings"

type StringList struct {
	Strings []string
}

func (s *StringList) construct() {
	if s.Strings == nil {
		s.Strings = []string{}
	}
}
func (s *StringList) Has(value string) bool {
	s.construct();
	for _, e := range s.Strings {
		if e == value {
			return true
		}
	}
	return false
}
func (s *StringList) String() string {
	return strings.Join(s.Strings,",")
}
func (s *StringList) Put(value string) {
	s.construct();
	s.Strings = append(s.Strings,value )
}
