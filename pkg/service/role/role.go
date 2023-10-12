package role

type Role struct {
	Id       int64  `json:"id" yaml:"id"`
	Name     string `json:"name" yaml:"name"`
	Password string `json:"password" yaml:"password"`
	Token    string `json:"token" yaml:"token"`
}
