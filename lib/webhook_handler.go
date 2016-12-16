package lib

import (
	"encoding/hex"
	"fmt"
	"github.com/cryptix/wav"
	sakura "github.com/yamamoto-febc/sakura-iot-go"
	"io/ioutil"
	"os"
	"sort"
)

const (
	BUF_SIZE        = 1024
	DATA_CHANNEL    = 0
	CONTROL_CHANNEL = 1
	START_REQUEST   = 0
	END_REQUEST     = 1
)

type WebhookHandler struct {
	channelBuffer channelSortWrapper
	option        *Option
	out           func(string, ...interface{})
}

func NewWebhookHandler(option *Option, logFunc func(string, ...interface{})) *WebhookHandler {
	w := &WebhookHandler{
		option: option,
		out:    logFunc,
	}
	if w.out == nil {
		w.out = func(f string, v ...interface{}) {}
	}
	w.initBuf()
	return w
}

func (w *WebhookHandler) infolog(f string, v ...interface{}) {
	if len(v) == 0 {
		w.out(f)
	} else {
		w.out(f, v)
	}
}
func (w *WebhookHandler) debuglog(f string, v ...interface{}) {
	if w.option.Debug {
		if len(v) == 0 {
			w.out(f)
		} else {
			w.out(f, v)
		}
	}

}

func (w *WebhookHandler) HandleRequest(p sakura.Payload) {
	// TODO 同時リクエストや複数のモジュール対応は後で
	
	w.debuglog("[DEBUG] Start handle request\n")
	// 開始リクエストを受けたらバッファをクリアする
	if w.hasStartRequest(p) {
		w.debuglog("[DEBUG] Receive a start request\n")
		w.initBuf()
	}

	// バッファリング
	w.debuglog("[DEBUG] Buffuring received data\n")
	w.pushToBuffer(p)

	// 終了リクエストを受けたらバッファをソートしてwav変換を行う
	if w.hasEndRequest(p) {
		w.debuglog("[DEBUG] Received a end request \n")

		// ソートしたバイト列を取得
		rawData, err := w.channelBuffer.GetSortedBytes()
		if err != nil {
			w.infolog("[ERROR] Failed on decoding channel data : %s", err)
			w.initBuf()
			return
		}

		// waveファイル作成
		path, err := w.createWave(rawData)
		if err != nil {
			w.infolog("[ERROR] Failed on creating wave file : %s", err)
			w.initBuf()
			return
		}

		file, err := os.Open(path)
		if err != nil {
			w.infolog("[ERROR] Failed on creating wave file : %s", err)
			w.initBuf()
			return
		}
		defer file.Close()
		defer os.Remove(path)

		// wav変換後、Azureで音声判定を行う
		speechToTextWorker := NewSpeechToTextWorker(w.option)

		result, err := speechToTextWorker.HasMagicalSpel(file)

		if err != nil {
			w.infolog("[ERROR] Failed on Recognizing from wave file: %s", err)
			w.initBuf()
			return
		}

		// TODO   判定OKならバルス！を実行してIncommingWebhookへポストする
		if result {
			w.infolog("[INFO] Spell 'バルス' was casted. Accept.")
		} else {
			w.infolog("[INFO] Spell 'バルス' was not casted. ")
		}
	}

}

func (w *WebhookHandler) initBuf() {
	w.channelBuffer = channelSortWrapper{channels: []sakura.Channel{}}
}

func (w *WebhookHandler) pushToBuffer(p sakura.Payload) {
	for _, c := range p.Payload.Channels {
		if c.Channel != DATA_CHANNEL {
			continue
		}
		w.channelBuffer.Add(c)
	}
}

func (w *WebhookHandler) hasStartRequest(p sakura.Payload) bool {
	return w.hasRequestWithControlValue(p, START_REQUEST)
}

func (w *WebhookHandler) hasEndRequest(p sakura.Payload) bool {
	return w.hasRequestWithControlValue(p, END_REQUEST)
}

func (w *WebhookHandler) hasRequestWithControlValue(p sakura.Payload, v int) bool {
	for _, c := range p.Payload.Channels {
		if c.Channel == CONTROL_CHANNEL {
			value, err := c.GetInt()
			if err != nil {
				return false
			}
			if value == int32(v) {
				return true
			}
		}
	}
	return false
}

func (w *WebhookHandler) createWave(data []byte) (string, error) {
	var wavFormat = wav.File{
		SampleRate:      1024 * 2,
		Channels:        1,
		SignificantBits: 8,
	}

	f, _ := ioutil.TempFile(".", "sacloud_balus_")
	defer f.Close()
	writer, err := wavFormat.NewWriter(f)
	if err != nil {
		return "", err
	}
	defer writer.Close()

	str := ""
	for i, b := range data {
		if i%8 == 0 {
			w.out("[DUMP] %s", str)
			str = ""
		}
		str = str + fmt.Sprintf("%X", b)

		err = writer.WriteSample([]byte{b})
		if err != nil {
			return "", err
		}
	}

	return f.Name(), nil
}

// ----------------------------------------------------------------------------

type channelSortWrapper struct {
	channels []sakura.Channel
}

func (c *channelSortWrapper) Add(channel sakura.Channel) {
	c.channels = append(c.channels, channel)
	//c = &channelSortWrapper(channels)
}

func (c *channelSortWrapper) Len() int {
	return len(c.channels)
}

func (c *channelSortWrapper) Less(i, j int) bool {
	return c.channels[i].Datetime.UnixNano() < c.channels[j].Datetime.UnixNano()
}

func (c *channelSortWrapper) Swap(i, j int) {
	c.channels[i], c.channels[j] = c.channels[j], c.channels[i]
}

func (c *channelSortWrapper) GetSortedBytes() ([]byte, error) {
	sort.Sort(c)

	res := []byte{}

	for _, ch := range c.channels {
		hexString, err := ch.GetHexString()
		if err != nil {
			return res, err
		}
		bytes, err := hex.DecodeString(hexString)
		if err != nil {
			return res, err
		}
		res = append(res, bytes...)
	}

	return res, nil
}
