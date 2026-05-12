//ff:func feature=parse type=util control=iteration dimension=1
//ff:what test: TestIsSkippableLinePython
package parse

import "testing"

func TestIsSkippableLinePython(t *testing.T) {
	cases := []struct {
		line string
		want bool
	}{
		{"", true},
		{"#!/usr/bin/env python3", true},
		{"# coding: utf-8", true},
		{"# -*- coding: utf-8 -*-", true},
		{"# regular comment", true},
		{"# ff:func feature=validate", false},
		{"import os", false},
		{"def foo():", false},
	}
	for _, c := range cases {
		got := IsSkippableLinePython(c.line)
		if got != c.want {
			t.Errorf("IsSkippableLinePython(%q) = %v, want %v", c.line, got, c.want)
		}
	}
}
