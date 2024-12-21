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

func newMQTTProbe(probe config.Probe, logger log.Logger) *MQTTProbe {
	var isReady = make(chan struct{})
	var msgChan = make(chan mqtt.Message)

	opt := mqtt.NewClientOptions().AddBroker(probe.Scheme + "://" + probe.Target)
	opt.SetCleanSession(true)
	opt.SetClientID(probe.ClientID)
	opt.SetUsername(probe.Username)
	opt.SetPassword(probe.Password)
	opt.SetKeepAlive(time.Duration(probe.KeepAlive) * time.Second)
	opt.SetPingTimeout(time.Duration(probe.PingTimeout) * time.Second)
	opt.SetConnectTimeout(time.Duration(probe.ConnectTimeout) * time.Second)
	if probe.TLSClientConfig != nil {
		opt.SetTLSConfig(probe.TLSClientConfig.ToTLSConfig())
	}
	opt.SetOnConnectHandler(func(c mqtt.Client) {
		optReader := c.OptionsReader()
		level.Debug(logger).Log("msg", "Connected to MQTT broker", "target", probe.Target, "client_id", optReader.ClientID())
		token := c.Subscribe(probe.Topic, probe.QoS, func(c mqtt.Client, m mqtt.Message) {
			msgChan <- m
		})
		token.Wait()
		if token.Error() != nil {
			level.Error(logger).Log("msg", "Failed to subscribe to MQTT topic", "target", probe.Target, "topic", probe.Topic, "qos", probe.QoS, "err", token.Error())
			return
		}
		isReady <- struct{}{}
		level.Debug(logger).Log("msg", "Subscribed to MQTT topic", "target", probe.Target, "topic", probe.Topic, "qos", probe.QoS)
	})
	opt.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		level.Error(logger).Log("msg", "Lost connection to MQTT broker", "target", probe.Target, "err", err)
	})
	c := mqtt.NewClient(opt)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		level.Error(logger).Log("msg", "Failed to connect to MQTT broker", "target", probe.Target, "err", token.Error())
		return nil
	}

	select {
	case <-isReady:
	case <-time.After(time.Duration(probe.KeepAlive) * time.Second):
		level.Error(logger).Log("msg", "MQTT probe connect timeout", "target", probe.Target)
		return nil
	}

	return &MQTTProbe{
		Client:  c,
		MsgChan: msgChan,
	}
}

func (mp *MQTTProbe) Probe(probe config.Probe, logger log.Logger) bool {
	defer mp.Client.Disconnect(0)

	if !mp.Client.IsConnected() {
		level.Error(logger).Log("msg", "MQTT client is not connected", "target", probe.Target)
		return false
	}

	level.Debug(logger).Log("msg", "Publishing MQTT message", "target", probe.Target, "topic", probe.Topic, "qos", probe.QoS)

	message := "from emqx-exporter MQTT probe"
	if token := mp.Client.Publish(probe.Topic, probe.QoS, false, message); token.Wait() && token.Error() != nil {
		level.Error(logger).Log("msg", "Failed to publish MQTT message", "target", probe.Target, "topic", probe.Topic, "qos", probe.QoS, "err", token.Error())
		return false
	}

	select {
	case msg := <-mp.MsgChan:
		if msg != nil && string(msg.Payload()) == message {
			level.Debug(logger).Log("msg", "MQTT probe receive message success", "target", probe.Target, "topic", probe.Topic, "qos", probe.QoS)
			return true
		}
		level.Error(logger).Log("msg", "MQTT probe receive message failed", "target", probe.Target, "topic", probe.Topic, "qos", probe.QoS)
		return false
	case <-time.After(time.Duration(probe.KeepAlive) * time.Second):
		level.Error(logger).Log("msg", "MQTT probe receive message timeout", "target", probe.Target)
		return false
	}
}
