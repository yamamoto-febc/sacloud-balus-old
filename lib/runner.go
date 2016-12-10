package lib

import (
	"fmt"
	"github.com/yamamoto-febc/sakura-iot-go"
	"log"
	"net/http"
	"os"
)

func Start(option *Option) error {

	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(os.Stdout)
	out := log.Printf

	balusHandler := NewWebhookHandler(option, out)

	handler := &sakura_iot_go.WebhookHandler{
		Secret:     option.Secret,
		HandleFunc: balusHandler.HandleRequest,
		Debug:      option.Debug,
	}

	addr := fmt.Sprintf("%s:%d", "", option.Port)

	out("[INFO] start ListenAndServe. addr:[%s] path:[%s] secret:[%s]\n", addr, option.Path, option.Secret)
	http.Handle(option.Path, handler)
	return http.ListenAndServe(addr, nil)
}
