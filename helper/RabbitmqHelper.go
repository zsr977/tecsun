package helper

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type channel struct {
	ch            *amqp.Channel
	notifyClose   chan *amqp.Error
	notifyConfirm chan amqp.Confirmation
}

type Rabbitmq struct {
	Url           string
	Connections   map[int]*amqp.Connection
	ConnectionNum int
	Channels      map[int]channel
	IdleChannels  []int
	BusyChannels  map[int]int
	ChannelNum    int
	Mutex         *sync.Mutex
}

func (r *Rabbitmq) lockWriteConnect(connectId int, newConn *amqp.Connection) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.Connections[connectId] = newConn
}

func (r *Rabbitmq) lockWriteChannel(channelId int, cha channel) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.Channels[channelId] = cha
}

func (r *Rabbitmq) connect() *amqp.Connection {
	conn, err := amqp.Dial(r.Url)
	if err != nil {
		logger.Errorf("连接rabbitmq失败：%v", err)
	}
	return conn
}

func (r *Rabbitmq) connectPool() {
	r.Connections = make(map[int]*amqp.Connection)
	for i := 0; i < r.ConnectionNum; i++ {
		r.Connections[i] = r.connect()
	}
}

func (r *Rabbitmq) channel(connectId int, conn *amqp.Connection) channel {
	var notifyClose = make(chan *amqp.Error)
	var notifyConfirm = make(chan amqp.Confirmation)
	cha := channel{
		notifyClose:   notifyClose,
		notifyConfirm: notifyConfirm,
	}
	if conn.IsClosed() {
		conn = r.connect()
	}
	ch, err := conn.Channel()
	if err != nil {
		ch = r.reCreateChannel(connectId, err)
	}
	err = ch.Confirm(false)
	if err != nil {
		logger.Errorf("%v", err)
	}
	ch.NotifyClose(cha.notifyClose)
	ch.NotifyPublish(cha.notifyConfirm)
	cha.ch = ch
	return cha
}

func (r *Rabbitmq) reCreateChannel(connectId int, err error) (ch *amqp.Channel) {
	if strings.Index(err.Error(), "channel/connection is not open") >= 0 || strings.Index(err.Error(), "CHANNEL_ERROR - expected 'channel.open'") >= 0 {
		var newConn *amqp.Connection
		if r.Connections[connectId].IsClosed() {
			newConn = r.connect()
		} else {
			newConn = r.Connections[connectId]
		}
		r.lockWriteConnect(connectId, newConn)
		ch, err = newConn.Channel()
		logger.Errorf("打开channel失败：%v", err)
	} else {
		logger.Errorf("打开channel失败：%v", err)
	}
	return
}

func (r *Rabbitmq) channelPool() {
	r.Channels = make(map[int]channel)
	for index, connection := range r.Connections {
		for j := 0; j < r.ChannelNum; j++ {
			key := index*r.ChannelNum + j
			r.Channels[key] = r.channel(index, connection)
			r.IdleChannels = append(r.IdleChannels, key)
		}
	}
}

func (r *Rabbitmq) getChannel() (*amqp.Channel, int) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	IdleLength := len(r.IdleChannels)
	if IdleLength > 0 {
		rand.Seed(time.Now().Unix())
		index := rand.Intn(IdleLength)
		channelId := r.IdleChannels[index]
		r.IdleChannels = append(r.IdleChannels[:index], r.IdleChannels[index+1:]...)
		r.BusyChannels[channelId] = channelId
		ch := r.Channels[channelId].ch
		return ch, channelId
	} else {
		return nil, -1
	}
}

func (r *Rabbitmq) backChannelId(channelId int) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.IdleChannels = append(r.IdleChannels, channelId)
	delete(r.BusyChannels, channelId)
}

func (r *Rabbitmq) reDeclareExchange(exchangeName string, exchangeKind string, channelId int, err error) (ch *amqp.Channel) {
	var connectionId int
	if strings.Index(err.Error(), "channel/connection is not open") >= 0 {
		if channelId == -1 {
			connectionId = 0
		} else {
			connectionId = int(channelId / r.ChannelNum)
		}
		cha := r.channel(connectionId, r.Connections[connectionId])
		r.lockWriteChannel(channelId, cha)
		if err := cha.ch.ExchangeDeclare(
			exchangeName,
			exchangeKind,
			true,
			false,
			false,
			false,
			nil); err != nil {
			logger.Errorf("%v", err)
		}
		return cha.ch
	} else {
		logger.Errorf("%v", err)
		return nil
	}
}

func (r *Rabbitmq) dataForm(notice interface{}) []byte {
	temp := make(map[string]interface{})
	for k, v := range notice.(map[interface{}]interface{}) {
		temp[k.(string)] = v
	}
	if body, err := json.Marshal(temp); err != nil {
		logger.Errorf("格式化数据失败：%v", err)
		return nil
	} else {
		return body
	}
}

func (r *Rabbitmq) publish(channelId int, ch *amqp.Channel, exchangeName string, exchangeKind string, routeKey string, data []byte) (err error) {
	err = ch.Publish(
		exchangeName,
		routeKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         data,
		})

	if err != nil {
		if strings.Index(err.Error(), "channel/connection is not open") >= 0 {
			err = r.rePublish(exchangeName, exchangeKind, channelId, routeKey, data, err)
		}
	}
	return
}

func (r *Rabbitmq) rePublish(exchangeName string, exchangeKind string, channelId int, routeKey string, data []byte, errMsg error) (err error) {
	ch := r.reDeclareExchange(exchangeName, exchangeKind, channelId, errMsg)
	err = ch.Publish(
		exchangeName,
		routeKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         data,
		})
	return
}

func (r *Rabbitmq) DeclareExchange(exchangeName string, exchangeKind string) error {
	defer func() {
		msg := recover()
		if msg != nil {
			errMsg, _ := msg.(string)
			_ = errors.New(errMsg)
			return
		}
	}()
	ch, _ := r.getChannel()
	cha := channel{}
	if ch == nil {
		cha = r.channel(0, r.Connections[0])
		defer func(ch *amqp.Channel) {
			err := ch.Close()
			if err != nil {
				logger.Errorf("%v", err)
			}
		}(cha.ch)
		ch = cha.ch
	}
	return ch.ExchangeDeclare(
		exchangeName,
		exchangeKind,
		true,
		false,
		false,
		false,
		nil)
}

func (r *Rabbitmq) DestroyExchange(exchangeName string) error {
	defer func() {
		msg := recover()
		if msg != nil {
			errMsg, _ := msg.(string)
			_ = errors.New(errMsg)
			return
		}
	}()
	ch, _ := r.getChannel()
	cha := channel{}
	if ch == nil {
		cha = r.channel(0, r.Connections[0])
		defer func(ch *amqp.Channel) {
			err := ch.Close()
			if err != nil {
				logger.Errorf("%v", err)
			}
		}(cha.ch)
		ch = cha.ch
	}
	return ch.ExchangeDelete(exchangeName, false, false)
}

func (r *Rabbitmq) BindQueue(name, key, exchange string) error {
	defer func() {
		msg := recover()
		if msg != nil {
			errMsg, _ := msg.(string)
			_ = errors.New(errMsg)
			return
		}
	}()
	ch, _ := r.getChannel()
	cha := channel{}
	if ch == nil {
		cha = r.channel(0, r.Connections[0])
		defer func(ch *amqp.Channel) {
			err := ch.Close()
			if err != nil {
				logger.Errorf("%v", err)
			}
		}(cha.ch)
		ch = cha.ch
	}
	q, _ := ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil)
	return ch.QueueBind(
		q.Name,
		key,
		exchange,
		false,
		nil)
}

func (r *Rabbitmq) UnbindQueue(name, key, exchange string) error {
	defer func() {
		msg := recover()
		if msg != nil {
			errMsg, _ := msg.(string)
			_ = errors.New(errMsg)
			return
		}
	}()
	ch, _ := r.getChannel()
	cha := channel{}
	if ch == nil {
		cha = r.channel(0, r.Connections[0])
		defer func(ch *amqp.Channel) {
			err := ch.Close()
			if err != nil {
				logger.Errorf("%v", err)
			}
		}(cha.ch)
		ch = cha.ch
	}
	_ = ch.QueueUnbind(name, key, exchange, nil)
	_, err := ch.QueueDelete(name, false, false, false)
	return err
}

func (r *Rabbitmq) initRabbitmq() error {
	if r.Url == "" {
		return errors.New("未定义rabbitmq连接地址")
	}
	if r.ConnectionNum == 0 {
		r.ConnectionNum = 3
	}
	if r.ChannelNum == 0 {
		r.ChannelNum = 5
	}
	r.BusyChannels = make(map[int]int)
	r.Mutex = new(sync.Mutex)
	r.connectPool()
	r.channelPool()
	return nil
}

var ServerPool = map[string]Rabbitmq{}

func RegisterRabbitmq(user, password, vhost string) Rabbitmq {
	if ServerPool[user+vhost].Url == "" {
		server := Rabbitmq{Url: "amqp://" + user + ":" + password + "@" + os.Getenv("RABBITMQ_AMQP") + "/" + vhost}
		_ = server.initRabbitmq()
		ServerPool[user+vhost] = server
	}
	return ServerPool[user+vhost]
}

func (r *Rabbitmq) PutIntoQueue(exchangeName string, exchangeKind string, routeKey string, notice interface{}) (message interface{}, pubErr error) {
	defer func() {
		msg := recover()
		if msg != nil {
			errMsg, _ := msg.(string)
			pubErr = errors.New(errMsg)
			return
		}
	}()
	ch, channelId := r.getChannel()
	cha := channel{}
	if ch == nil {
		cha = r.channel(0, r.Connections[0])
		defer func(ch *amqp.Channel) {
			err := ch.Close()
			if err != nil {
				logger.Errorf("%v", err)
			}
		}(cha.ch)
		ch = cha.ch
	}
	data, ok := notice.([]byte)
	if !ok {
		data = r.dataForm(notice)
	}
	var tryTime = 1
	for {
		pubErr = r.publish(channelId, ch, exchangeName, exchangeKind, routeKey, data)
		if pubErr != nil {
			if tryTime <= 5 {
				tryTime++
				continue
			} else {
				r.backChannelId(channelId)
				return notice, pubErr
			}
		}
		select {
		case confirm := <-r.Channels[channelId].notifyConfirm:
			if confirm.Ack {
				r.backChannelId(channelId)
				return notice, nil
			}
			return notice, errors.New("ack failed")
		case chaConfirm := <-cha.notifyConfirm:
			if chaConfirm.Ack {
				return notice, nil
			}
			return notice, errors.New("ack failed")
		case <-time.After(3 * time.Second):
			r.backChannelId(channelId)
			confirmErr := errors.New("Can not receive the confirm . ")
			return notice, confirmErr
		}
	}
}

func (r *Rabbitmq) GetFromQueue(name string, c chan map[string]interface{}) {
	defer func() {
		msg := recover()
		if msg != nil {
			errMsg, _ := msg.(string)
			_ = errors.New(errMsg)
			return
		}
	}()
	ch, _ := r.getChannel()
	cha := channel{}
	if ch == nil {
		cha = r.channel(0, r.Connections[0])
		defer func(ch *amqp.Channel) {
			err := ch.Close()
			if err != nil {
				logger.Errorf("%v", err)
			}
		}(cha.ch)
		ch = cha.ch
	}
	qu, _ := ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil)
	_ = ch.Qos(
		1,
		0,
		false,
	)
	// 消费任务
	msg, err := ch.Consume(
		qu.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	logger.Error(err, "接收消息失败")
	go func() {
		for d := range msg {
			var v map[string]interface{}
			_ = json.Unmarshal(d.Body, &v)
			c <- v
			err := d.Ack(false)
			logger.Error(err, "手动确认消息投递失败")
		}
	}()
}
