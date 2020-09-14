package models

// Something is a Something
type Something struct {
	Name string `json:"name"`
}

// GetSomething Gets a Something
func GetSomething() (*Something, error) {
	s := Something{Name: "Name of a Something"}
	return &s, nil
}
