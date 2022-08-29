package cmd

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/st2projects/ssh-sentinel-client/config"
	"github.com/st2projects/ssh-sentinel-client/helpers"
	"github.com/st2projects/ssh-sentinel-core/model"
	"io"
	"net/http"
	"os"
)

var configPath string

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign a new key",
	Run: func(cmd *cobra.Command, args []string) {
		config.MakeConfig(configPath)

		conf := config.Config

		certValid, expDate := helpers.IsCertValid(conf.GetCertFile())
		if certValid {
			log.Infof("Existing cert valid until %s", expDate)
		} else {
			log.Info("Creating new cert")
			signNewKey(conf)
		}
	},
}

func signNewKey(conf *config.ConfigType) {
	pubKey := conf.GetPublicKey()

	if !helpers.PathExists(pubKey) {
		panic("Key " + conf.PublicKey + " does not exist")
	}

	key, err := os.ReadFile(pubKey)

	if err != nil {
		panic(err)
	}

	signReq := &model.KeySignRequest{
		Username:   conf.Username,
		APIKey:     conf.APIKey,
		Principals: conf.Principals,
		Key:        string(key),
	}

	signReqBytes, err := json.Marshal(signReq)

	if err != nil {
		panic(err)
	}

	resp, err := http.Post(conf.EndPoint, "application/json", bytes.NewBuffer(signReqBytes))

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	signResp := &model.KeySignResponse{}

	json.Unmarshal(body, signResp)

	if !signResp.Success {
		log.Errorf("Sign request failed with err: %s", signResp.Message)
	} else {
		err := os.WriteFile(helpers.ExpandPath(conf.GetCertFile()), []byte(signResp.SignedKey), os.FileMode(0600))
		if err != nil {
			log.Errorf("Failed to write new cert %s", err.Error())
		} else {
			log.Info("Signed new key")
		}
	}
}

func init() {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().StringVarP(&configPath, "config", "c", "", "Config file")

	signCmd.MarkFlagRequired("config")
}
