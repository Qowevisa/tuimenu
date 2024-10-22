package simple

type SimpleMenuOption func(conf *MenuConfig)

func WithCustomTitle(title string) SimpleMenuOption {
	return func(conf *MenuConfig) {
		conf.Title = title
	}
}

func WithCustomBackKey(back string) SimpleMenuOption {
	return func(conf *MenuConfig) {
		conf.BackKey = back
	}
}

func WithUsageOfEscapeCodes() SimpleMenuOption {
	return func(conf *MenuConfig) {
		conf.UsingEscapeCodes = true
	}
}

type InterruptOption func(intr *MenuInterrupt)

type IntrOptStatus uint32

const (
	IntrOptStatus_ClearAfterFinish IntrOptStatus = 1 << iota
)

func ClearAfterFinishingFunc() InterruptOption {
	return func(intr *MenuInterrupt) {
		intr.Status |= IntrOptStatus_ClearAfterFinish
	}
}
