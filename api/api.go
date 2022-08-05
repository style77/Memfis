package api

// First I need to clarify that:
// Token - API token that is used to perform every action
// Key - Key that is used to obtain API token (1 use)

import (
	managers "Memfis/managers"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (t *token) formatString() string {
	t.Formatted = managers.GenerateSecretKey(8) + "." + strconv.FormatInt(t.Expiring, 10) + "." + managers.Encode([]byte(t.user.Username))
	return t.Formatted
}

func (t *token) isExpired() bool {
	return !(time.Now().UTC().Unix() < t.Expiring)
}

var AuthorizeKeys = []key{}
var tokens = []token{}

func contains(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Returns one-time-use token that is necessary to use to get API KEY
func authorize(c *gin.Context) {
	config, _ := managers.LoadConfig()

	var u userModel

	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Wrong parameters passed."})
	}

	if u.Username == config.Username && u.Password == config.Password {

		generatedKey := managers.GenerateSecretKey()

		AuthorizeKeys = append(AuthorizeKeys, key{Key: generatedKey, user: u})

		c.JSON(http.StatusOK, gin.H{
			"token": generatedKey,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong credentials."})
	}
}

func keyInAuthorizedKeys(k key_req) (bool, key) {
	for _, d := range AuthorizeKeys {
		if d.Key == k.Token && !d.used {
			return true, d
		}
	}
	return false, key{}
}

func getToken(c *gin.Context) {

	var k key_req

	if err := c.BindJSON(&k); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Wrong parameters passed."})
	}

	z, kt := keyInAuthorizedKeys(k)
	kt.used = true

	if !z {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong token passed."})
	} else {

		config, err := managers.LoadConfig()
		if err != nil {
			log.Fatal(err)
		}

		interval := int32(config.TokenExpireTime)

		token := token{}
		token.Key = managers.GenerateSecretKey(8)
		token.Expiring = time.Now().Add(time.Hour * time.Duration(interval)).Unix()
		token.user = kt.user
		token.formatString()

		tokens = append(tokens, token)
		c.JSON(http.StatusOK, gin.H{"token": token.Formatted, "expiring": token.Expiring, "user": token.user.Username})
	}
}

func authorizeToken(tokenToCheck string) (bool, token) {
	for _, authorizedToken := range tokens {
		if tokenToCheck == authorizedToken.Formatted && !authorizedToken.isExpired() {
			return true, authorizedToken
		}
	}
	return false, token{}
}

func getData(c *gin.Context) {
	var k reqDataModel

	if err := c.BindJSON(&k); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Wrong parameters passed."})
	}

	z, _ := authorizeToken(k.Token)

	if !z {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong token passed or token has expired."})
	} else {
		dm, exists := managers.FindData(k.Name)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"message": "Couldn't find data with name \"" + k.Name + "\"."})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"formatted": dm.FormatToString(),
				"type":      dm.Type,
				"name":      dm.Name,
				"value":     dm.Value,
			})
		}
	}
}

func setData(c *gin.Context) {
	var k reqDataModel

	if err := c.BindJSON(&k); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Wrong parameters passed."})
	}

	z, _ := authorizeToken(k.Token)

	if !z {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong token passed or token has expired."})
	} else {
		dm, z := managers.CreateDataModel(k.Name, k.Value)
		if !z {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Data with name \"" + k.Name + "\" already exists."})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"formatted": dm.FormatToString(),
				"type":      dm.Type,
				"name":      dm.Name,
				"value":     dm.Value,
			})
		}
	}
}

func updateData(c *gin.Context) {
	var k reqDataModel

	if err := c.BindJSON(&k); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Wrong parameters passed."})
	}

	z, _ := authorizeToken(k.Token)

	if !z {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong token passed or token has expired."})
	} else {
		err, dm := managers.UpdateData(k.Name, k.Value)
		if !err {
			c.JSON(http.StatusNotFound, gin.H{"message": "Couldn't find data with name \"" + k.Name + "\"."})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"formatted": dm.FormatToString(),
				"type":      dm.Type,
				"name":      dm.Name,
				"value":     dm.Value,
				"updated":   time.Now().UTC().Unix(),
			})
		}
	}
}

func deleteData(c *gin.Context) {
	var k reqDataModel

	if err := c.BindJSON(&k); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Wrong parameters passed."})
	}

	z, _ := authorizeToken(k.Token)

	if !z {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong token passed or token has expired."})
	} else {
		z, dm := managers.RemoveData(k.Name)
		if !z {
			c.JSON(http.StatusNotFound, gin.H{"message": "Couldn't find data with name \"" + k.Name + "\"."})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"formatted": dm.FormatToString(),
				"type":      dm.Type,
				"name":      dm.Name,
				"value":     dm.Value,
				"updated":   time.Now().UTC().Unix(),
			})
		}
	}
}

func Run(serverAddress string) {
	log.Println("Setuping gin.")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.TrustedPlatform = gin.PlatformGoogleAppEngine

	r.POST("/authorize", authorize)
	r.POST("/token", getToken)

	r.GET("/data", getData)
	r.POST("/data", setData)
	r.PATCH("/data", updateData)
	r.DELETE("/data", deleteData)

	log.Printf("Running server on: %s.\n", serverAddress)
	r.Run(serverAddress)
}
