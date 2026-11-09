package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"rafi0101/traefik-ssl-certificate-exporter/models"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {

	// pflag.String("certResolver", "dode", "Acme certresolver to extract these certs (requiered)")
	pflag.String("source", "traefik/acme.json", "path to traefik acme.json")
	pflag.String("dest", "certs/", "path to destination where to store certificates")
	pflag.Int("owner", 0, "owner for the extracted cert/keys")
	pflag.Int("group", 0, "group for the extracted cert/keys")

	pflag.String("config", "", "Path to config file")

	pflag.Usage = func() {
		fmt.Println("traefik-ssl-certificate-exporter exports the ssl certificates from the provided traefik acme.json")
		pflag.PrintDefaults()
		os.Exit(0)
	}
	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)
	viper.SetConfigFile(viper.GetString("config"))

	viper.ReadInConfig()

	ownerId := viper.GetInt("owner")
	groupId := viper.GetInt("group")

	acmejson, err := ioutil.ReadFile(viper.GetString("source"))
	if err != nil {
		fmt.Println("Failed to read traefik acme.json", err)
		os.Exit(1)
	}

	// //Unmarshal traefik acme.json
	acme := new(models.ProviderMdl)
	if err := json.Unmarshal(acmejson, acme); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//For each certificateprovider
	for _, certifcateProviderValue := range *acme {

		//For each certificate within one provider
		for _, certificate := range certifcateProviderValue.Certificates {

			//Get cert destination folder path, and replace * (Wildcard) with "_"
			certPath := viper.GetString("dest") + strings.Replace(certificate.Domain.Main, "*", "_", -1) + "/"
			os.MkdirAll(certPath, 0755)
			os.Chown(certPath, ownerId, groupId)

			//Decode private key
			privateKey, err := base64.StdEncoding.DecodeString(certificate.Key)
			if err != nil {
				fmt.Println("Failed to decode private key for ", certificate.Domain.Main, " :", err)
				os.Exit(1)
			}
			//Decode fullChain cert
			fullChain, err := base64.StdEncoding.DecodeString(certificate.Certificate)
			if err != nil {
				fmt.Println("Failed to decode fullchain for ", certificate.Domain.Main, " :", err)
				os.Exit(1)
			}

			//Convert fullChain []byte to string
			fullChainString := string(fullChain)

			//Get Index where cert ends and chain begins
			fullChainIndex := strings.Index(fullChainString, "\n-----BEGIN CERTIFICATE-----")

			if fullChainIndex < 0 {
				fmt.Println("Could not read fullchain cert")
				os.Exit(1)
			}
			//Cert we need is on top of the file. chain is everything else
			cert := fullChainString[:fullChainIndex]
			chain := fullChainString[fullChainIndex+1:]

			files := []struct {
				name string
				data []byte
				perm os.FileMode
			}{
				{"privkey.pem", privateKey, 0600},
				{"fullchain.pem", fullChain, 0644},
				{"cert.pem", []byte(cert), 0644},
				{"chain.pem", []byte(chain), 0644},
			}

			for _, file := range files {
				if err := writeIfChanged(certPath+file.name, file.data, file.perm, ownerId, groupId); err != nil {
					fmt.Println("Failed to write ", file.name, " for ", certificate.Domain.Main, " :", err)
					os.Exit(1)
				}
			}

		}
	}

}

// writeIfChanged writes data to path, replacing any previous content.
// If the file already holds identical data it won't be touched.
func writeIfChanged(path string, data []byte, perm os.FileMode, ownerId, groupId int) error {
	old, err := ioutil.ReadFile(path)
	if err == nil && bytes.Equal(old, data) {
		return nil
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := ioutil.WriteFile(path, data, perm); err != nil {
		return err
	}
	os.Chown(path, ownerId, groupId)
	return nil
}
