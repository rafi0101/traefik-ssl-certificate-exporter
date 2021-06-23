package main

import (
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

			//Write privateKey to destination path
			privateKeyFile, err := os.Create(certPath + "privkey.pem")
			if err != nil {
				fmt.Println("Failed to create privateKey file for ", certificate.Domain.Main, " :", err)
				os.Exit(1)
			}
			_, err = privateKeyFile.Write(privateKey)
			if err != nil {
				fmt.Println("Failed to write privateKey for ", certificate.Domain.Main, " :", err)
				os.Exit(1)
			}
			privateKeyFile.Chmod(0600)
			privateKeyFile.Chown(ownerId, groupId)

			//Write fullChain to destination path
			fullChainFile, err := os.Create(certPath + "fullchain.pem")
			if err != nil {
				fmt.Println("Failed to create fullchain file for ", certificate.Domain.Main, " :", err)
				os.Exit(1)
			}
			_, err = fullChainFile.Write(fullChain)
			if err != nil {
				fmt.Println("Failed to write fullchain for ", certificate.Domain.Main, " :", err)
				os.Exit(1)
			}
			fullChainFile.Chmod(0644)
			fullChainFile.Chown(ownerId, groupId)

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

			//Write cert to destination path
			certFile, err := os.Create(certPath + "cert.pem")
			if err != nil {
				fmt.Println("Failed to create cert file for ", certificate.Domain.Main, " :", err)
				os.Exit(1)
			}
			_, err = certFile.WriteString(cert)
			if err != nil {
				fmt.Println("Failed to write cert for ", certificate.Domain.Main, " :", err)
				os.Exit(1)
			}
			certFile.Chmod(0644)
			certFile.Chown(ownerId, groupId)

			//Write chain to destination path
			chainFile, err := os.Create(certPath + "chain.pem")
			if err != nil {
				fmt.Println("Failed to create chain file for ", certificate.Domain.Main, " :", err)
				os.Exit(1)
			}
			_, err = chainFile.WriteString(chain)
			if err != nil {
				fmt.Println("Failed to write chain for ", certificate.Domain.Main, " :", err)
				os.Exit(1)
			}
			chainFile.Chmod(0644)
			chainFile.Chown(ownerId, groupId)

		}
	}

}
