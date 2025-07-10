package main

import (
	"fmt"
	"git.imooc.com/coding-535/base/handler"
	hystrix2 "git.imooc.com/coding-535/base/plugin/hystrix"
	base "git.imooc.com/coding-535/base/proto/base"
	"git.imooc.com/coding-535/common"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v3"
	opentracing2 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	log "github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/registry"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/opentracing/opentracing-go"
	"net"
	"net/http"
)

func main() {
	//1. Service registry
	consul:=consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			// Can also be service name if deployed in docker-compose
			"localhost:8500",
		}
	})
	//2. Configuration center for frequently changing variables
	consulConfig,err := common.GetConsulConfig("localhost",8500,"/micro/config")
	if err !=nil {
		common.Error(err)
	}

	//3. Use configuration center to connect to mysql
	mysqlInfo := common.GetMysqlFromConsul(consulConfig,"mysql")
	// Initialize database
	db,err := gorm.Open("mysql",mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err !=nil {
		common.Fatal(err)
	}
	fmt.Println("MySQL connection successful")
	common.Info("MySQL connection successful")
	defer db.Close()
	// Disable plural table names
	db.SingularTable(true)

	//4. Add distributed tracing
	t,io,err := common.NewTracer("base","localhost:6831")
	if err !=nil {
		common.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//5. Add circuit breaker
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	//6. Add logging center
	//1) Program logs need to be written to log files
	//2) Add filebeat.yml file to the program
	//3) Start filebeat with command: ./filebeat -e -c filebeat.yml
	common.Info("Logging system added!")

	// Start monitoring program
	go func() {
		//http://192.168.0.112:9092/turbine/turbine.stream
		// Dashboard access address: http://127.0.0.1:9002/hystrix, URL must end with /hystrix
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0","9092"),hystrixStreamHandler)
		fmt.Println("333")
		if err !=nil {
			fmt.Println(err)
		}
	}()

	//7. Add monitoring
	common.PrometheusBoot(9192)

	// Create service
	service := micro.NewService(
		micro.Name("base-cap"),
		micro.Version("latest"),
		// Add service registry
		micro.Registry(consul),
		// Add distributed tracing
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// Only works when acting as a client
		micro.WrapClient(hystrix2.NewClientHystrixWrapper()),
		// Add rate limiting
		micro.WrapHandler(ratelimit.NewHandlerWrapper(1000)),
	)

	// Initialize service
	service.Init()

	// Register handler for quick access to developed services
	base.RegisterBaseHandler(service.Server(), new(handler.Base))

	// Start service
	if err := service.Run(); err != nil {
        // Output startup failure information
		log.Fatal(err)
	}
}
