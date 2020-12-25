package util

import "C"
import (
	"bytes"
	"encoding/xml"
	"regexp"

	"exchange_zero_trust_api/logger"

	. "github.com/magicmonty/activesync-go/activesync"
	. "github.com/magicmonty/activesync-go/activesync/base"
	. "github.com/magicmonty/wbxml-go/wbxml"
)

/*
```
<?xml version="1.0" encoding="utf-8"?>
<O:Provision xmlns:O="Provision" xmlns:S="Settings">
    <S:DeviceInformation>
        <S:Set>
            <S:Model>MIX 2</S:Model>
            <S:IMEI>888833336669999</S:IMEI>
            <S:FriendlyName>MIX 2</S:FriendlyName>
            <S:OS>Android 8.0.0</S:OS>
            <S:PhoneNumber>+8618599999999</S:PhoneNumber>
            <S:UserAgent>Android/8.0.0-EAS-1.3</S:UserAgent>
            <S:MobileOperator>中国联通 (46001)</S:MobileOperator>
        </S:Set>
    </S:DeviceInformation>
    <O:Policies>
        <O:Policy>
            <O:PolicyType>MS-EAS-Provisioning-WBXML</O:PolicyType>
        </O:Policy>
    </O:Policies>
</O:Provision>
```
*/

type (
	Provision struct {
		XMLName           xml.Name `xml:"Provision"`
		Text              string   `xml:",chardata"`
		O                 string   `xml:"O,attr"`
		S                 string   `xml:"S,attr"`
		DeviceInformation struct {
			Text string `xml:",chardata"`
			Set  struct {
				Text           string `xml:",chardata"`
				Model          string `xml:"Model"`
				IMEI           string `xml:"IMEI"`
				FriendlyName   string `xml:"FriendlyName"`
				OS             string `xml:"OS"`
				PhoneNumber    string `xml:"PhoneNumber"`
				UserAgent      string `xml:"UserAgent"`
				MobileOperator string `xml:"MobileOperator"`
			} `xml:"Set"`
		} `xml:"DeviceInformation"`
		Policies struct {
			Text   string `xml:",chardata"`
			Policy struct {
				Text       string `xml:",chardata"`
				PolicyType string `xml:"PolicyType"`
			} `xml:"Policy"`
		} `xml:"Policies"`
	}

	DeviceInfo struct {
		Model          string `json:"model"`
		IMEI           string `json:"imei"`
		FriendlyName   string `json:"friendly_name"`
		PhoneNumber    string `json:"phone_number"`
		MobileOperator string `json:"mobile_operator"`
	}
)

func removeInvalidChars(b []byte) []byte {
	re := regexp.MustCompile("[^\x09\x0A\x0D\x20-\uD7FF\uE000-\uFFFD\u10000-\u10FFFF]")
	return re.ReplaceAll(b, []byte{})
}

func EncodeXML(xmlString []byte) {
	xmlString = removeInvalidChars([]byte(xmlString))
	w := bytes.NewBuffer(make([]byte, 0))
	e := NewEncoder(
		MakeCodeBook(PROTOCOL_VERSION_14_1),
		string(xmlString),
		w)
	err := e.Encode()
	if err != nil {
		logger.Log.Println(err.Error())
	} else {
		logger.Log.Println(w)
	}
}

func getDecodeResult(data ...byte) string {
	var result string
	result, _ = Decode(bytes.NewBuffer(data), MakeCodeBook(PROTOCOL_VERSION_14_1))
	return result
}

func Parse(data string) (DeviceInfo, error) {
	result := DeviceInfo{}
	xmlData := getDecodeResult([]byte(data)...)

	out := Provision{}
	err := xml.Unmarshal([]byte(xmlData), &out)
	if err != nil {
		return result, err
	}

	logger.Log.Printf("Model: %v\n", out.DeviceInformation.Set.Model)
	logger.Log.Printf("IMEI: %v\n", out.DeviceInformation.Set.IMEI)
	logger.Log.Printf("FriendlyName: %v\n", out.DeviceInformation.Set.FriendlyName)
	logger.Log.Printf("PhoneNumber: %v\n", out.DeviceInformation.Set.PhoneNumber)
	logger.Log.Printf("MobileOperator: %v\n", out.DeviceInformation.Set.MobileOperator)

	result.Model = out.DeviceInformation.Set.Model
	result.IMEI = out.DeviceInformation.Set.IMEI
	result.FriendlyName = out.DeviceInformation.Set.FriendlyName
	result.PhoneNumber = out.DeviceInformation.Set.PhoneNumber
	result.MobileOperator = out.DeviceInformation.Set.MobileOperator

	return result, err
}
