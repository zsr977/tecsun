package helper

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/config_center"
	"gopkg.in/yaml.v2"
	"os"
	"sync"
)

type ConfigCenter struct {
	wg     sync.WaitGroup
	dc     config_center.DynamicConfiguration
	dataId string
	group  string
}

func (c *ConfigCenter) Init(dataId string, arg ...string) {
	group := "DEFAULT_GROUP"
	if len(arg) > 0 && arg[0] != "" {
		group = arg[0]
	}
	c.dataId = dataId
	c.group = group
	dc, err := config.NewConfigCenterConfigBuilder().
		SetProtocol(config.GetRootConfig().ConfigCenter.Protocol).
		SetAddress(config.GetRootConfig().ConfigCenter.Address).
		Build().GetDynamicConfiguration()
	if err != nil {
		panic(err)
	}
	c.dc = dc
	c.dc.AddListener(dataId, &ConfigCenter{})
	yamlFile, err := c.dc.GetProperties(c.dataId, config_center.WithGroup(c.group))
	if err != nil {
		panic(err)
	}
	if conf, err := c.Parse(yamlFile); err != nil {
		panic(err)
	} else {
		c.Listener(conf)
	}
}

func (c *ConfigCenter) Process(configType *config_center.ConfigChangeEvent) {
	c.wg.Add(1)
	if conf, err := c.Parse(configType.Value.(string)); err != nil {
		logger.Errorf("%s %s,But %v", configType.Key, configType.ConfigType, err)
	} else {
		c.Listener(conf)
	}
}

func (c *ConfigCenter) Parse(value string) (conf map[string]interface{}, err error) {
	err = yaml.Unmarshal([]byte(value), &conf)
	return
}

func (c *ConfigCenter) Listener(conf map[string]interface{}) {
	mysqlConf := conf["mysql"].(map[interface{}]interface{})
	redisConf := conf["redis"].(map[interface{}]interface{})
	rabbitmqConf := conf["rabbitmq"].(map[interface{}]interface{})
	minioConf := conf["minio"].(map[interface{}]interface{})
	_ = os.Setenv("MYSQL", mysqlConf["username"].(string)+":"+
		mysqlConf["password"].(string)+"@tcp("+
		mysqlConf["host"].(string)+":"+
		mysqlConf["port"].(string)+")/"+
		mysqlConf["database"].(string)+
		"?charset=utf8")
	_ = os.Setenv("REDIS", redisConf["host"].(string)+":"+redisConf["port"].(string))
	_ = os.Setenv("RABBITMQ_AMQP", rabbitmqConf["amqp"].(string))
	_ = os.Setenv("MINIO_ENDPOINT", minioConf["endpoint"].(string))
	_ = os.Setenv("MINIO_ACCESS_KEY", minioConf["accessKey"].(string))
	_ = os.Setenv("MINIO_SECRET_KEY", minioConf["secretKey"].(string))
	_ = os.Setenv("RABBITMQ_AMQP", rabbitmqConf["amqp"].(string))
	_ = os.Setenv("RABBITMQ_HTTP", rabbitmqConf["http"].(string))
	_ = os.Setenv("RABBITMQ_ADMIN", "guest")
	_ = os.Setenv("RABBITMQ_SECRET", rabbitmqConf["guestPassword"].(string))
	_ = os.Setenv("RABBITMQ_USER", rabbitmqConf["taskUser"].(string))
	_ = os.Setenv("RABBITMQ_PASS", rabbitmqConf["taskPassword"].(string))
}

func RegisterHelper() {
	// 注册配置服务
	var c ConfigCenter
	c.Init("runtime")
}