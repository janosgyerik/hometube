package hometube

import "testing"

func TestParseURL(t *testing.T) {
	var tests = []struct {
		name, url, want, errstr string
	}{
		{
			"video url",
			"https://www.youtube.com/watch?v=H0FcOPb-9rE",
			"https://www.youtube.com/watch?v=H0FcOPb-9rE",
			"",
		},
		{
			"video url with extra args",
			"https://www.youtube.com/watch?v=H0FcOPb-9rE&foo=bar",
			"https://www.youtube.com/watch?v=H0FcOPb-9rE",
			"",
		},
		{
			"list url",
			"https://www.youtube.com/watch?list=OLAK5uy_kioHHyij8sGCFNG3aOh5C9nPGHExtNT74",
			"https://www.youtube.com/watch?list=OLAK5uy_kioHHyij8sGCFNG3aOh5C9nPGHExtNT74",
			"",
		},
		{
			"list url with extra args",
			"https://www.youtube.com/watch?list=OLAK5uy_kioHHyij8sGCFNG3aOh5C9nPGHExtNT74&foo=bar",
			"https://www.youtube.com/watch?list=OLAK5uy_kioHHyij8sGCFNG3aOh5C9nPGHExtNT74",
			"",
		},
		{
			"list url with video",
			"https://www.youtube.com/watch?v=Yk4x4CP_-Lk&list=OLAK5uy_kioHHyij8sGCFNG3aOh5C9nPGHExtNT74",
			"https://www.youtube.com/watch?list=OLAK5uy_kioHHyij8sGCFNG3aOh5C9nPGHExtNT74",
			"",
		},
		{
			"url without video or list",
			"https://www.youtube.com/watch?foo=bar",
			"",
			"url must have v or list param",
		},
		{
			"garbage url",
			"foo",
			"",
			"url must have v or list param",
		},
	}

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			ans, err := parseURL(tt.url)
			if ans != nil && ans.sanitized != tt.want {
				t.Errorf("got: %s\nwant: %s", ans, tt.want)
			}
			if err != nil && err.Error() != tt.errstr {
				t.Errorf("got: %s\nwant: %s", err.Error(), tt.want)
			}
		})
	}
}
