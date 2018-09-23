package tseg

import (
	"fmt"
	"testing"
)

type test struct {
	str string
	seg []string
}

var tests = []test{
	{"charactersetencoding", []string{"character", "set", "encoding"}},
	{"niceday", []string{"nice", "day"}},
	{"niceweather", []string{"nice", "weather"}},
	{"input", []string{"input"}},
	{"expertsexchange", []string{"experts", "exchange"}},
}

func TestGetTextSegmentation(t *testing.T) {
	sr := Segmentator{DictPath: "dict.txt", TextPath: "text.txt"}
	for _, val := range tests {
		seg, err := sr.GetSegmentation(val.str)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(seg) != len(val.seg) {
			t.Errorf("Строка на входе: %v\nРазбиение на выходе: %v\nПравильный ответ: %v\n", val.str, seg, val.seg)
			return
		}
		for j := range seg {
			if seg[j] != val.seg[j] {
				t.Errorf("Строка на входе: %v\nРазбиение на выходе: %v\nПравильный ответ: %v\n", val.str, seg, val.seg)
				return
			}
		}
	}
}
