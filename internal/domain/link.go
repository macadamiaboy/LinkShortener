package domain

import "errors"

type Link struct {
	ID        int64  `json:"id"`
	ShortCode string `json:"short_code"`
	LongUrl   string `json:"long_url"`
	Clicks    int32  `json:"clicks"`
}

func (l *Link) Validate() error {
	if l.LongUrl == "" {
		return errors.New("url required")
	}

	if l.ShortCode == "" {
		return errors.New("code required")
	}

	return nil
}
