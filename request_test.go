package curl

import (
	"fmt"
	"log"
	"testing"
	"time"
)

// TestHttpGet 測試 Http Get
func TestHttpGet(t *testing.T) {
	// var err error

	url := "http://dev.cqgame.api:8100/dev/promoweb/v1/util/gameList"

	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiI1YTRiM2JlODY1ZTM5MmIzMmY2MTNmNTYiLCJhY2NvdW50IjoicHJvbW9zaXRlIiwibmlja25hbWUiOiJwcm9tb3NpdGUiLCJ0eXBlIjoiT1AiLCJyb2xlIjoic2l0ZSIsImp0aSI6Ijk3NTU4OTU3NiIsImlhdCI6MTUxNDg3OTk3NSwiaXNzIjoiQ3lwcmVzcyIsInN1YiI6Ik9QVG9rZW4ifQ.kM7NZHHX4G1aBiVSYbokVeDGHakL1VTnqmkusTFq7Ek",
	}

	// 初始化
	req := NewRequest(
		SetTimeOut(time.Second*30),
		SetSkipTLSVerify(true),
		SetURL(url),
		SetHeader(headers),
	)

	_, err := req.Get()

	if err != nil {
		log.Fatal("Error Msg ", err)
	}
}

// TestHttpPostRawData 測試透過 raw data 呼叫 Http Post
func TestHttpPostRawData(t *testing.T) {
	url := "http://cherry.local.com/api/game/create"

	param := make(map[string]interface{})
	param["game_id"] = "AD1"
	param["guide"] = make(map[string]string)
	param["guide_banner"] = make(map[string]string)
	param["background"] = make(map[string]string)
	param["icon"] = make(map[string]string)
	param["guide_related"] = []string{}
	param["inside"] = 1
	param["respin"] = false
	param["free_game"] = false
	param["edited_by"] = "Neil"
	param["symbol"] = []string{}
	param["type"] = 1

	req := NewRequest(
		SetURL(url),
		SetRawData(param),
	)

	resp, err := req.Post()
	if err != nil {
		log.Fatal("Error Mssage ", err)
	}

	fmt.Println(string(resp.Body))
}

// TestHttpPutFromData 測試透過 from data 呼叫 Http Put
func TestHttpPutFromData(t *testing.T) {
	url := "http://dev.cqgame.api:8100/dev/promoweb/marquee/set"

	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiI1YTRiM2JlODY1ZTM5MmIzMmY2MTNmNTYiLCJhY2NvdW50IjoicHJvbW9zaXRlIiwibmlja25hbWUiOiJwcm9tb3NpdGUiLCJ0eXBlIjoiT1AiLCJyb2xlIjoic2l0ZSIsImp0aSI6Ijk3NTU4OTU3NiIsImlhdCI6MTUxNDg3OTk3NSwiaXNzIjoiQ3lwcmVzcyIsInN1YiI6Ik9QVG9rZW4ifQ.kM7NZHHX4G1aBiVSYbokVeDGHakL1VTnqmkusTFq7Ek",
	}

	param := make(map[string]interface{})

	req := NewRequest(
		SetURL(url),
		SetHeader(headers),
		SetFormData(param),
	)
	resp, err := req.Put()

	if err != nil {
		log.Fatal("Error Mssage ", err)
	}

	fmt.Println(string(resp.Body))
}

// TestHttpSendFile 測試上傳檔案 呼叫 Http Post
func TestHttpSendFile(t *testing.T) {
	url := "http://rd3-dev-imgcenter.guardians.one/api/file/upload_image"

	headers := map[string]string{
		"Token": "88e1a19120b059c33f6c9eb66fb755fb",
	}

	param := map[string]interface{}{}
	param["encrypt"] = true
	param["path"] = "links"

	filename := "./1.png"

	req := NewRequest(
		SetURL(url),
		SetHeader(headers),
	)

	if err := req.SetFileData(param, filename); err != nil {
		log.Fatal("Error Msg ", err)
	}

	resp, err := req.Post()
	if err != nil {
		log.Fatal("Error Mssage ", err)
	}

	fmt.Println(string(resp.Body))
}
