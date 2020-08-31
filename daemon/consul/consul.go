package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"github.com/vietnamz/cli-common/daemon"
	"strconv"
	"strings"
	"time"
)

const (
	defaultConsulHost = "localhost"
	defaultConsulPort = "8500"
)

var DefaultMysql = Service{
	Name: "mysql",
	Port: "3300",
	Address: "localhost",
}

var DefaultMq = Service{
	Name: "rabbitmq",
	Port: "5671",
	Address: "localhost",
}

type ServiceConfig struct {
	Host string
	Port string
	Scheme string
}

func DefaultServiceConfig() *ServiceConfig{
	return &ServiceConfig{
		Host: defaultConsulHost,
		Port: defaultConsulPort,
		Scheme: "http",
	}
}

type ServicesMap map[string]Service

type ConsulMonitor struct {

	name string

	// monitor worker configuration.
	Config *ServiceConfig

	// keep track the services.
	Services ServicesMap

	// consul client
	Client *consulapi.Client

	// to stop the timer.
	Done chan bool

	// already initialize
	WasInit bool

	// node name, take from consul
	Node string

}

func NewConsulMonitor(config *ServiceConfig )  *ConsulMonitor {
	var consul ConsulMonitor
	consul.Config = config
	consul.name = "trading"
	consul.init()
	return &consul
}

func (m *ConsulMonitor) initServiceMap()  {
	m.Services  = map[string]Service {
		"mysql": DefaultMysql,
		"rabbitmq": DefaultMq,
	}
}
type Service struct {
	Name string
	Port string
	Address string
}

func (c *ConsulMonitor) check(done chan bool, ticker *time.Ticker)  {
	logrus.Debugf("Start tick")
	for {
		select {
		case <-done:
			logrus.Debugf("done")
			return
		case t := <-ticker.C:
			logrus.Debugln("Tick at ", t)
			options := consulapi.QueryOptions {

			}
			srv,_, _ := c.Client.Catalog().Services(&options)

			logrus.Debugf("start checkout the services\n")
			for k, _ := range srv {
				catalogs, _, _ := c.Client.Catalog().Service(k, "", &options)

				for i := range catalogs {
					// TODO: we have to go with tags or service id.
					// Since address is attached to node instaed of service itself.
					addr := catalogs[i].Address
					// split the service id to take the address
					parts := strings.Split(catalogs[i].ServiceID, ":")
					if len(parts) < 2 {
						logrus.Debugf("The serivce id invalid format\n")
						// hack, since consul is a service registry. we can ignore.
						if catalogs[i].ServiceName == "consul" {
							continue
						}
					} else {
						// re-assign the addr
						addr = parts[1]
					}
					logrus.Debugf("Fist is %s\n", parts[0])
					logrus.Debugf("Second is %s\n", parts[1])
					newServ := Service{
						Name:    catalogs[i].ServiceName,
						Address: addr,
						Port:    strconv.Itoa(catalogs[i].ServicePort),
					}
					c.checkService(newServ)
				}
			}
		}
	}

}

func  (c *ConsulMonitor) Start()  {
	logrus.Debugf("start")
	// start a tick every 5 second.
	ticker := time.NewTicker(5 * time.Second)
	c.Done = make(chan bool)
	go c.check(c.Done, ticker)
}


// register the service to monitor.
func (c *ConsulMonitor) Register()  {
	/*
	// AgentCheck represents a check known to the agent
	type AgentCheck struct {
		Node        string
		CheckID     string
		Name        string
		Status      string
		Notes       string
		Output      string
		ServiceID   string
		ServiceName string
	}
	 */
	// TODO. this params should be passed from command line.


	hostAddress := "host.docker.internal"
	port := ":" + strconv.Itoa(8080)
	scheme := "http://"
	pingPath := "/ping"
	host := strings.Join([]string{ scheme, hostAddress, port, pingPath}, "")
	logrus.Debugln(host)
	serviceId := strings.Join([]string{c.name, daemon.GetLocalIpAddresses()[0]}, ":")
	agentService := consulapi.AgentService{
		ID: serviceId,
		Service: serviceId,
		Address: hostAddress,
		Port: 8081,
	}
	logrus.Debugf("Time duration 3s is %s", time.Duration(3).String())
	/*healthCheckDefinition := consulapi.HealthCheckDefinition {
		HTTP: host,
		Method: "GET",
		IntervalDuration: time.Duration(3 * 1000),
	}*/
	/*consulapi.AgentCheck{
		Name: c.name,
		CheckID: c.name +"testing",
		Node: c.Node ,
		ServiceID: serviceId,
		ServiceName: c.name,
		Notes: "testing",
		Type: "http",
		Definition: healthCheckDefinition,
		Output: "trading is alive",
	}*/
	register := consulapi.CatalogRegistration {
		Datacenter: "",
		Address: daemon.GetLocalIpAddresses()[0],
		Service: &agentService,
		Node: c.Node,
		//Check: &agentCheck,
	}






	// register health check.
	_, err := c.Client.Catalog().Register(&register, nil)
	if err != nil {
		logrus.Errorf("Failed to register with consul %s", err)
	}
	check := consulapi.AgentServiceCheck{
		HTTP: host,
		Method: "GET",
		Interval: "3s",
	}
	checkRegistration := consulapi.AgentCheckRegistration{
		ID: c.name,
		Name: c.name,
		ServiceID: serviceId,
		AgentServiceCheck: check,
	}
	err = c.Client.Agent().CheckRegister(&checkRegistration)
	if err != nil {
		logrus.Errorf("%s", err)
	}
}

// stop the monitor
func (c *ConsulMonitor) Stop()  {
	c.Done <- true
	logrus.Debugf("Timer is stopped")
}

func (c *ConsulMonitor) init() error {
	c.initServiceMap()
	conf := consulapi.Config{
		Address: "localhost:8500",
		Scheme: "http",
	}
	cli, err := consulapi.NewClient(&conf)
	c.Client = cli
	if err != nil {
		logrus.Errorf("failed to create new consul client %s", err)
	}
	node, err := c.Client.Agent().NodeName()
	if err != nil {
		logrus.Errorf("failed to get node name %s", err)
	}
	c.Node = node
	c.WasInit = true
	return nil
}

// Initialized the monitor.
// We should implement ReInit instead of init here.
// make a bool flag to distinguish two state.
func (c *ConsulMonitor) ReInit() error {
	if !c.WasInit {
		panic("Invalid call")
		return nil
	}
	c.init()
	return nil
}

func (c *ConsulMonitor) checkService( service Service) {
	if _,ok := c.Services[service.Name]; ok{
		if service.Address != c.Services[service.Name].Address {
			fmt.Printf("service %s address changes old %s, new %s\n",
				service.Name,
				c.Services[service.Name].Address,
				service.Address)
		}
		if service.Port != c.Services[service.Name].Port {
			fmt.Printf("service %s port changes old %s, new %s\n",
				service.Name,
				c.Services[service.Name].Port,
				service.Port)
		}
		c.Services[service.Name] = service
	}
	c.Services[service.Name] = service
}