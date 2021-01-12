/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2021 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package bytecode

// StringInterner is used to intern strings.
type StringInterner struct {
	strings map[string]string
}

// NewStringInterner creates and returns a new stringInterner.
func NewStringInterner() *StringInterner {
	return &StringInterner{
		strings: make(map[string]string),
	}
}

// Intern interns the string s and returns a string with the same contents as s,
// but that is guaranteed to be unique within si. Or maybe it's clearer this
// way: if si already contains a string with the same contents as s, it returns
// that other string: same content, but at a different memory location.
func (si *StringInterner) Intern(s string) string {
	if r, ok := si.strings[s]; ok {
		return r
	}
	si.strings[s] = s
	return s
}
