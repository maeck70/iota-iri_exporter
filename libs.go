/*
MIT License

Copyright (c) 2018 Marcel van Eck

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import "strconv"

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func btof(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func actify(b float64) string {
	if b == 1 {
		return "Active"
	}
	return "Inactive"
}

func stoi(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func stof(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// BytesToString will convert a Bytes stream to a String datatype.
// Input: []byte
// Output: string
func BytesToString(data []byte) string {
	return string(data[:])
}
