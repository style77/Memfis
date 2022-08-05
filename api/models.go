package api

type userModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type key struct {
	Key  string
	user userModel
	used bool
}

type token struct {
	Key       string
	Expiring  int64 // timestamp
	user      userModel
	Formatted string
}

type key_req struct {
	Token string `json:"token"`
}

type reqDataModel struct {
	Token string      `json:"token"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
