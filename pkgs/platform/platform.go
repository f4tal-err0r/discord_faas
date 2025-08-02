package platform

type Platform struct {
	Name string
}

func NewPlatform(name string) *Platform {
	return &Platform{
		Name: name,
	}
}
