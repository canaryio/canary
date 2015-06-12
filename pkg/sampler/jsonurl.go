package sampler

import "net/url"

// The JsonURL type allows us to parse the URL from a JSON document at load
// time, instead of when it's used (which could be multiple times)
type JsonURL struct {
	*url.URL
}

func NewJsonURL(str string) (u *JsonURL, err error) {
	_url, err := url.Parse(str)
	if err == nil {
		u = &JsonURL{_url}
	}

	return
}

func (self *JsonURL) UnmarshalJSON(data []byte) (err error) {
	// strip off leading and trailing double-quotes
	tmp, err := NewJsonURL(string(data[1:len(data) - 1]))
	
	if err == nil {
		*self = *tmp
	}

	return
}

