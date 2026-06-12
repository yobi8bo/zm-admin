package cache

import "context"

type CaptchaStore struct{}

func (CaptchaStore) Set(id, value string) error {
	return Set(context.Background(), CaptchaKey(id), value, CaptchaTTL())
}

func (CaptchaStore) Get(id string, clear bool) string {
	ctx := context.Background()
	var value string
	var err error
	if clear {
		value, err = GetDelString(ctx, CaptchaKey(id))
	} else {
		value, err = GetString(ctx, CaptchaKey(id))
	}
	if err != nil {
		return ""
	}
	return value
}

func (s CaptchaStore) Verify(id, answer string, clear bool) bool {
	return s.Get(id, clear) == answer
}
