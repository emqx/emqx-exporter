package prober

import (
	"context"
	"emqx-exporter/config"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type MQTTProbe struct {
	Client  mqtt.Client
	MsgChan <-chan mqtt.Message
}

type mqttProbeManager struct {
	probes map[string]*MQTTProbe
	sync.RWMutex
}

var manager mqttProbeManager

func init() {
	manager = mqttProbeManager{
		probes: make(map[string]*MQTTProbe),
	}
	go func() {
		for {
			manager.Lock()
			defer manager.Unlock()
			for target, probe := range manager.probes {
				if probe == nil {
					delete(manager.probes, target)
					continue
				}
			}
			manager.Unlock()

			select {
			case <-context.Background().Done():
				return
			case <-time.After(60 * time.Second):
			}
		}
	}()
}

func initMQTTProbe(probe config.Probe, logger log.Logger) (*MQTTProbe, error) {
	var msgChan = make(chan mqtt.Message)

	opt := mqtt.NewClientOptions().AddBroker(probe.Scheme + "://" + probe.Target)
	opt.SetClientID(probe.ClientID)
	opt.SetUsername(probe.Username)
	opt.SetPassword(probe.Password)
	opt.SetKeepAlive(time.Duration(probe.KeepAlive) * time.Second)
	if probe.TLSClientConfig != nil {
		opt.SetTLSConfig(probe.TLSClientConfig.ToTLSConfig())
	}
	opt.SetOnConnectHandler(func(c mqtt.Client) {
		level.Info(logger).Log("msg", "Connected to MQTT broker", "target", probe.Target)
		token := c.Subscribe(probe.Topic, probe.QoS, func(c mqtt.Client, m mqtt.Message) {
			msgChan <- m
		})
		token.Wait()
		if token.Error() != nil {
			level.Error(logger).Log("msg", "Failed to subscribe to MQTT topic", "target", probe.Target, "topic", probe.Topic, "qos", probe.QoS, "err", token.Error())
			return
		}
		level.Info(logger).Log("msg", "Subscribed to MQTT topic", "target", probe.Target, "topic", probe.Topic, "qos", probe.QoS)
	})
	opt.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		level.Error(logger).Log("msg", "Lost connection to MQTT broker", "target", probe.Target, "err", err)
	})
	c := mqtt.NewClient(opt)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		level.Error(logger).Log("msg", "Failed to connect to MQTT broker", "target", probe.Target, "err", token.Error())
		return nil, token.Error()
	}

	return &MQTTProbe{
		Client:  c,
		MsgChan: msgChan,
	}, nil
}

func ProbeMQTT(probe config.Probe, logger log.Logger) bool {
	mqttProbe, ok := manager.probes[probe.Target]
	if !ok {
		var err error
		if mqttProbe, err = initMQTTProbe(probe, logger); err != nil {
			return false
		}
		manager.Lock()
		defer manager.Unlock()
		manager.probes[probe.Target] = mqttProbe
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
	case <-time.After(time.Duration(probe.KeepAlive) * time.Second):
		level.Info(logger).Log("msg", "MQTT probe receive message timeout", "target", probe.Target)
		return false
	}

	return true
}
