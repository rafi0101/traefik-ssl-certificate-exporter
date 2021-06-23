package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"rafi0101/traefik-ssl-certificate-exporter/models"
	"reflect"
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

			// ------------- [START] Write privateKey to destination path ------------- //
			privateKeyFilePath := certPath + "privkey.pem"

			//Open or Create privateKey is not exists
			privateKeyFile, err := os.OpenFile(privateKeyFilePath, os.O_RDWR|os.O_CREATE, 0600)
			if err != nil {
				fmt.Println("Failed to create privateKey file for ", certificate.Domain.Main, " :", err)
				os.Exit(1)
			}

			//Read old privateKey
			privateKeyOld, err := ioutil.ReadFile(privateKeyFilePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			//Compate old and new privateKey. If they are identical the file won't be touched
			if !reflect.DeepEqual(privateKeyOld, privateKey) {
				_, err = privateKeyFile.Write(privateKey)
				if err != nil {
					fmt.Println("Failed to write privateKey for ", certificate.Domain.Main, " :", err)
					os.Exit(1)
				}
				privateKeyFile.Chown(ownerId, groupId)
			}
			// ------------- [END] Write privateKey to destination path ------------- //

			// ------------- [START] Write fullChain to destination path ------------- //
			fullChainFilePath := certPath + "fullchain.pem"

			//Open or Create fullChain is not exists
			fullChainFile, err := os.OpenFile(fullChainFilePath, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				fmt.Println("Failed to create fullChain file for ", certificate.Domain.Main, " :", err)
				os.Exit(1)
			}

			//Read old fullChain
			fullChainOld, err := ioutil.ReadFile(fullChainFilePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			//Compate old and new fullChain. If they are identical the file won't be touched
			if !reflect.DeepEqual(fullChainOld, fullChain) {
				_, err = fullChainFile.Write(fullChain)
				if err != nil {
					fmt.Println("Failed to write fullChain for ", certificate.Domain.Main, " :", err)
					os.Exit(1)
				}
				fullChainFile.Chown(ownerId, groupId)
			}
			// ------------- [END] Write fullChain to destination path ------------- //

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

			// ------------- [START] Write cert to destination path ------------- //
			certFilePath := certPath + "cert.pem"

			//Open or Create cert is not exists
			certFile, err := os.OpenFile(certFilePath, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				fmt.Println("Failed to create cert file for ", certificate.Domain.Main, " :", err)
				os.Exit(1)
			}

			//Read old cert
			certOld, err := ioutil.ReadFile(certFilePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			//Compate old and new cert. If they are identical the file won't be touched
			if !reflect.DeepEqual(string(certOld), cert) {
				_, err = certFile.WriteString(cert)
				if err != nil {
					fmt.Println("Failed to write cert for ", certificate.Domain.Main, " :", err)
					os.Exit(1)
				}
				certFile.Chown(ownerId, groupId)
			}
			// ------------- [END] Write cert to destination path ------------- //

			// ------------- [START] Write chain to destination path ------------- //
			chainFilePath := certPath + "chain.pem"

			//Open or Create chain is not exists
			chainFile, err := os.OpenFile(chainFilePath, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				fmt.Println("Failed to create chain file for ", certificate.Domain.Main, " :", err)
				os.Exit(1)
			}

			//Read old chain
			chainOld, err := ioutil.ReadFile(chainFilePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			//Compate old and new chain. If they are identical the file won't be touched
			if !reflect.DeepEqual(string(chainOld), chain) {
				_, err = chainFile.WriteString(chain)
				if err != nil {
					fmt.Println("Failed to write chain for ", certificate.Domain.Main, " :", err)
					os.Exit(1)
				}
				chainFile.Chown(ownerId, groupId)
			}
			// ------------- [END] Write chain to destination path ------------- //

		}
	}

}
