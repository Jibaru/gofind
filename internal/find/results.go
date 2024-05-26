package find

type Result struct {
	Name string `json:"name"`
}

func (f Result) String() string {
	return f.Name
}
