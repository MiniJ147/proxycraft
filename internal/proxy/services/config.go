package services

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/minij147/proxycraft/pkg/consts"
)

func logic(hostIp string, config consts.ServerConfig) {
	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(config)
		if err != nil {
			log.Println("Config-Server: failed json marshal", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte{})
			return
		}

		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write(data)
	})

	log.Println("starting config server...")
	log.Fatal(http.ListenAndServe(hostIp, nil))

}

func ConfigRun(hostIp string, config consts.ServerConfig) {
	go logic(hostIp, config)
}
