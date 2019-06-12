package main

import (
	"fmt"
	"io/ioutil"

	"bytes"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/miguelmota/go-ethereum-hdwallet"
)

var (
	conf Config
)

type Config struct {
	Mnemo       []string `toml:"mnemo"`
	PrivateKeys []string `toml:"private_key"`
}

func main() {
	var (
		pk     []string
		buffer bytes.Buffer
	)
	ConfFileName := "getprivatekey.toml"
	if _, err := toml.DecodeFile(ConfFileName, &conf); err != nil {
		fmt.Println("Ошибка загрузки файла конфигурации:", err.Error())
		return
	}
	for _, mn := range conf.Mnemo {

		wallet, err := hdwallet.NewFromMnemonic(mn)
		if err != nil {
			panic(err)
		}
		path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
		account, err := wallet.Derive(path, false)
		if err != nil {
			panic(err)
		}

		strAdrs := account.Address.String()
		addrss := fmt.Sprintf("M%s", strings.ToLower(strAdrs[1:len(strAdrs)]))
		fmt.Println(addrss)
		privKeyStr, err := wallet.PrivateKeyHex(account)
		if err != nil {
			panic(err)
		}
		pk = append(pk, privKeyStr)

	}
	conf.PrivateKeys = pk

	encoder := toml.NewEncoder(&buffer)
	if err := encoder.Encode(conf); err != nil {
		fmt.Errorf("Ошибка сохранения файла: %s", err.Error())
		return
	}
	ioutil.WriteFile(ConfFileName, []byte(buffer.String()), 0644)

}
