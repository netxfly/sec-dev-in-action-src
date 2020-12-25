package util

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"sec-dev-in-action-src/waf/x-waf-admin/setting"
)

type Message struct {
	Status  int
	Message string
}

func MakeMd5(srcStr string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(srcStr)))
}

func EncryptPass(src string) string {
	return fmt.Sprintf("%s", MakeMd5(MakeMd5(src)[5:10]))
}

func WriteNginxConf(proxyConfig string, siteName string, vhostPath string) (err error) {
	proxyConfigFile := path.Join(vhostPath, fmt.Sprintf("%v.conf", siteName))
	log.Println(proxyConfigFile)
	fileConfig, err := os.Create(proxyConfigFile)
	log.Println(fileConfig, err)
	defer fileConfig.Close()
	proxyConfig = strings.Replace(proxyConfig, "\r\n", "\n", -1)
	_, err = fileConfig.WriteString(proxyConfig)

	return err
}

func ReloadNginx() (err error) {
	log.Println("start to Reload nginx")
	ret, err := exec.Command(setting.NginxBin, "-t").Output()
	log.Println(ret, err)
	if err == nil {
		ret1, err := exec.Command(setting.NginxBin, "-s", "reload").Output()
		log.Println(ret1, err)
	}
	return err
}
