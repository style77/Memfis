package managers

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/thanhpk/randstr"
)

const CHUNK_SIZE = 65536 // <- for debugging reasons 1048576

var firstCreation = true

func GenerateSecretKey(params ...int) string {
	var length int
	if len(params) == 0 {
		length = 20
	} else {
		length = params[0]
	}

	key := randstr.String(length)

	if firstCreation {
		log.Println("Generating app secret key.")
		config, err := LoadConfig()
		if err != nil {
			log.Fatal("Cannot load config: ", err)
		}

		if config.SaveSecret {
			WriteToConfig("SECRET_KEY", key)
		}
		firstCreation = false
	}

	return key
}

var SecretKey string = GenerateSecretKey()

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Encrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(SecretKey))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func Decrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func (dm *dataModel) FormatToString() string {
	return fmt.Sprintf("%s;%s;%v", dm.Type, dm.Name, dm.Value)
}

func (dm *dataModel) clear() bool {
	dm.Name = ""
	dm.Value = ""
	dm.Type = ""
	debug.FreeOSMemory()
	return true
}

func (dm *dataModel) update(value any) (bool, dataModel) {
	modelType := defineModelsType(value)

	dm.Value = value
	dm.Type = modelType
	return true, *dm
}

func (ch *chunkModel) Encode() (*bytes.Buffer, error) {
	var dataToWrite []string
	for _, s := range ch.data {
		if s != (&dataModel{}) {
			strf := s.FormatToString()
			strf, _ = Encrypt(strf)
			dataToWrite = append(dataToWrite, strf)
		}
	}

	buf := new(bytes.Buffer)

	data := strings.Join(dataToWrite[:], "\n")
	_, err := buf.WriteString(data)
	if err != nil {
		log.Fatal(err)
	}
	return buf, err
}

func (ch *chunkModel) Decode(encodedString string) ([]dataModel, error) {
	dm := []dataModel{}
	encodedData := strings.Split(encodedString, "\n")
	for _, s := range encodedData {
		decodedData, err := Decrypt(s)
		if err != nil {
			return dm, err
		}
		strf := strings.Split(decodedData, ";")
		newDm := dataModel{Type: strf[0], Name: strf[1], Value: strf[2]}
		dm = append(dm, newDm)
	}
	return dm, nil
}

func (ch *chunkModel) SaveDataChunk() bool {
	// It's eventually in case of some issue that'll raise panic()

	log.Println("Trying to save data.")

	f, err := os.Create("C:\\Users\\style\\Downloads\\Memfis\\tmp\\dat" + strconv.Itoa(ch.chunkId)) // TODO change this path
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	f.Sync()

	w := bufio.NewWriter(f)
	data, _ := ch.Encode()

	_, err = w.WriteString(data.String())
	if err != nil {
		log.Fatal(err)
	}

	w.Flush()

	success := true

	return success
}

func (ch *chunkModel) contains(data dataModel) bool {
	for _, a := range ch.data {
		if a.Name == data.Name {
			return true
		}
	}
	return false
}

func (ch *chunkModel) AddDataToChunk(data dataModel) []*dataModel {
	if !ch.contains(data) {
		log.Println("Adding data to chunk.")
		ch.data = append(ch.data, &data)
		return ch.data
	}
	log.Println("Can't add data to chunk. Data is already in.")
	return nil
}
