package util

import "testing"

func TestAppendString(t *testing.T) {
	tests := []struct {
		name     string
		pars     []string
		expected string
	}{
		{"t1", []string{"aa", "bb", "cc"}, "aabbcc"},
		{"t2", []string{"qaz", "123", "mne"}, "qaz123mne"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := AppendString(test.pars...)
			if actual != test.expected {
				t.FailNow()
			}
		})
	}
}

func TestEncodeParams(t *testing.T) {
	tests := []struct {
		name     string
		pars     map[string]string
		expected string
	}{
		{"t1", map[string]string{"cc": "cc", "bb": "bb", "aa": "aa"}, "aa=aa&bb=bb&cc=cc"},
		{"t2", map[string]string{"ysd": "ysd", "iui": "iui", "nmm": "nmm"}, "iui=iui&nmm=nmm&ysd=ysd"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := EncodeParams(test.pars)
			if actual != test.expected {
				t.FailNow()
			}
		})
	}
}
