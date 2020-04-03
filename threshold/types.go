package threshold

type Threshold interface {
	Attempt(key string) error
	IsAllow(key string, failAction ...func()) bool
}
