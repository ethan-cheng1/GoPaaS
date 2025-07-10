package main

import (
	"fmt"
	"git.imooc.com/coding-535/common"
	go_micro_service_pod "git.imooc.com/coding-535/pod/proto/pod"
	"git.imooc.com/coding-535/podApi/handler"
	hystrix2 "git.imooc.com/coding-535/podApi/plugin/hystrix"
	"git.imooc.com/coding-535/podApi/proto/podApi"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v3"
	"github.com/asim/go-micro/plugins/wrapper/select/roundrobin/v3"
	opentracing2 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"

	"github.com/asim/go-micro/v3/server"
	"github.com/opentracing/opentracing-go"


	"net"
	"net/http"
	"strconv"
)

var (
	// Service address
	hostIp = "192.168.0.105"
	// Service address
	serviceHost = hostIp
	// Service port
	servicePort = "8082"
	// Registry configuration
	consulHost  = hostIp
	consulPort int64 = 8500
	// Tracing
	tracerHost = hostIp
	tracerPort = 6831
	// Circuit breaker port, each service must be unique
	hystrixPort = 9092
	// Monitoring port, each service must be unique
	prometheusPort = 9192
)

func main()  {
	//1. Registry
	consul:=consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			consulHost+":"+strconv.FormatInt(consulPort,10),
		}
	})

	//2. Add tracing
	t,io,err := common.NewTracer("go.micro.api.podApi",tracerHost+":"+strconv.Itoa(tracerPort))
	if err != nil {
		common.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//3. Add circuit breaker
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	//4. Add logging
	//1) Program logs need to be written to log files
	//2) Add filebeat.yml file to the program
	//3) Start filebeat with command: ./filebeat -e -c filebeat.yml
	fmt.Println("Logs are recorded in the root directory micro.log file, please check the logs!")

	//6. Start circuit breaker monitoring
	go func() {
		//http://192.168.0.108:9092/turbine/turbine.stream
		// Dashboard access address http://127.0.0.1:9002/hystrix, URL must end with /hystrix
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0",strconv.Itoa(hystrixPort)),hystrixStreamHandler)
		if err != nil {
			common.Error(err)
		}
	}()

	//7. Add monitoring collection address
	common.PrometheusBoot(prometheusPort)

	//8. Create service
	service := micro.NewService(
		// Custom service address, must be written before other parameters
		micro.Server(server.NewServer(func(options *server.Options) {
			options.Advertise = serviceHost+":"+servicePort
		})),
		micro.Name("go.micro.api.podApi"),
		micro.Version("latest"),
		// Specify service port
		micro.Address(":"+servicePort),
		// Add registry
		micro.Registry(consul),
		// Add tracing
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// Start circuit breaker as client scope
		micro.WrapClient(hystrix2.NewClientHystrixWrapper()),
		// Add rate limiting
		micro.WrapHandler(ratelimit.NewHandlerWrapper(1000)),
		// Add load balancing
		micro.WrapClient(roundrobin.NewClientWrapper()),
		)

	service.Init()

	podService := go_micro_service_pod.NewPodService("go.micro.service.pod",service.Client())
	// Register controller
	if err := podApi.RegisterPodApiHandler(service.Server(),&handler.PodApi{PodService:podService});err !=nil {
		common.Error(err)
	}
	// Start service
	if err := service.Run(); err!=nil {
		common.Fatal(err)
	}
}
