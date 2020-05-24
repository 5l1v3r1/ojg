// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"strings"
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestNodeParseString(t *testing.T) {
	for i, d := range []data{
		{src: "null", value: nil},
		{src: "true", value: true},
		{src: "false", value: false},
		{src: "123", value: 123},
		{src: "-321", value: -321},
		{src: "12.3", value: 12.3},
		{src: `12345678901234567890`, value: oj.Big("12345678901234567890")},
		{src: `9223372036854775807`, value: 9223372036854775807},             // max int
		{src: `9223372036854775808`, value: oj.Big("9223372036854775808")},   // max int + 1
		{src: `-9223372036854775807`, value: -9223372036854775807},           // min int
		{src: `-9223372036854775808`, value: oj.Big("-9223372036854775808")}, // min int -1
		{src: `0.9223372036854775808`, value: oj.Big("0.9223372036854775808")},
		{src: `123456789012345678901234567890`, value: oj.Big("123456789012345678901234567890")},
		{src: `0.123456789012345678901234567890`, value: oj.Big("0.123456789012345678901234567890")},
		{src: `0.1e20000`, value: oj.Big("0.1e20000")},
		{src: `1.2e1025`, value: oj.Big("1.2e1025")},
		{src: `-1.2e-1025`, value: oj.Big("-1.2e-1025")},

		{src: `"xyz"`, value: "xyz"},

		{src: "[]", value: []interface{}{}},
		{src: "[true]", value: []interface{}{true}},
		{src: "[true,false]", value: []interface{}{true, false}},
		{src: "[[]]", value: []interface{}{[]interface{}{}}},
		{src: "[[true]]", value: []interface{}{[]interface{}{true}}},

		{src: "{}", value: map[string]interface{}{}},
		{src: `{"abc":true}`, value: map[string]interface{}{"abc": true}},
		{src: `{"abc":{"def":3}}`, value: map[string]interface{}{"abc": map[string]interface{}{"def": 3}}},

		{src: `{"abc": [{"x": {"y": [{"b": true}]},"z": 7}]}`,
			value: map[string]interface{}{
				"abc": []interface{}{
					map[string]interface{}{
						"x": map[string]interface{}{
							"y": []interface{}{
								map[string]interface{}{
									"b": true,
								},
							},
						},
						"z": 7,
					},
				},
			}},
	} {
		var err error
		var v interface{}
		if d.onlyOne || d.noComment {
			p := oj.NodeParser{NoComment: d.noComment}
			v, err = p.Parse([]byte(d.src))
		} else {
			var p oj.NodeParser
			v, err = p.Parse([]byte(d.src))
		}
		if 0 < len(d.expect) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.expect, err.Error(), i, ": ", d.src)
		} else {
			tt.Nil(t, err, d.src)
			tt.Equal(t, d.value, v, ": ", d.src)
		}
	}
}

func TestNodeParseCallback(t *testing.T) {
	var results []byte
	cb := func(n oj.Node) bool {
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, n.String()...)
		return false
	}
	var p oj.NodeParser
	v, err := p.Parse([]byte(callbackJSON), cb)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, `1 [2] {"x":3} true false 123`, string(results))
}

func TestNodeParseReaderCallback(t *testing.T) {
	var results []byte
	cb := func(n oj.Node) bool {
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, n.String()...)
		return false
	}
	var p oj.NodeParser
	v, err := p.ParseReader(strings.NewReader(callbackJSON), cb)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, `1 [2] {"x":3} true false 123`, string(results))
}