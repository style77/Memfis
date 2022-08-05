package managers

type dataModel struct {
	Type  string      `json:"type"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type chunkModel struct {
	chunkId int
	data    []*dataModel
}

type config struct {
	ServerAddress   string `mapstructure:"SERVER_ADDRESS"`
	SaveSecret      bool   `mapstructure:"SAVE_SECRET"`
	SecretKey       string `mapstructure:"SECRET_KEY"`
	Username        string `mapstructure:"USERNAME"`
	Password        string `mapstructure:"PASSWORD"`
	TokenExpireTime int    `mapstructure:"API_TOKEN_EXPIRE_TIME"`
}
