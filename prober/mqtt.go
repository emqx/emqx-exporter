package prober

import (
	"emqx-exporter/config"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type MQTTProbe struct {
	Client  mqtt.Client
	MsgChan <-chan mqtt.Message
}

var mqttProbe *MQTTProbe

func initMQTTProbe(probe config.Probe, logger log.Logger) (*MQTTProbe, error) {
	opt := mqtt.NewClientOptions().AddBroker(probe.Scheme + "://" + probe.Target).SetClientID(probe.ClientID).SetUsername(probe.Username).SetPassword(probe.Password)
	opt.SetOnConnectHandler(func(c mqtt.Client) {
		level.Info(logger).Log("msg", "Connected to MQTT broker")
	})
	opt.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		level.Error(logger).Log("msg", "Lost connection to MQTT broker", "err", err)
	})
	c := mqtt.NewClient(opt)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		level.Error(logger).Log("msg", "Failed to connect to MQTT broker", "err", token.Error())
		return nil, token.Error()
	}

	var msgChan = make(chan mqtt.Message)
	if token := c.Subscribe(probe.Topic, probe.QoS, func(c mqtt.Client, m mqtt.Message) {
		msgChan <- m
	}); token.Wait() && token.Error() != nil {
		level.Error(logger).Log("msg", "Failed to subscribe to MQTT topic", "err", token.Error())
		return nil, token.Error()
	}

	return &MQTTProbe{
		Client:  c,
		MsgChan: msgChan,
	}, nil
}

func ProbeMQTT(probe config.Probe, logger log.Logger) bool {
	if mqttProbe == nil {
		var err error
		if mqttProbe, err = initMQTTProbe(probe, logger); err != nil {
			return false
		}
	}

	if !mqttProbe.Client.IsConnected() {
		return false
	}

	if token := mqttProbe.Client.Publish(probe.Topic, probe.QoS, false, "hello world"); token.Wait() && token.Error() != nil {
		return false
	}

	select {
	case msg := <-mqttProbe.MsgChan:
		if msg == nil {
			return false
		}
	case <-time.After(5 * time.Second):
		return false
	}

	return true
}
