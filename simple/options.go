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
