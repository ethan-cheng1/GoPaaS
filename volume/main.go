package main


import (
	"flag"
	"fmt"
	"git.imooc.com/coding-535/common"
	"git.imooc.com/coding-535/volume/domain/repository"
	"path/filepath" 
    
	//"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	service2 "git.imooc.com/coding-535/volume/domain/service"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v3"
	opentracing2 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
    "github.com/asim/go-micro/v3/server"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/opentracing/opentracing-go"
	"git.imooc.com/coding-535/volume/handler"
	//hystrix2 "git.imooc.com/coding-535/volume/plugin/hystrix"
	"strconv"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	volume "git.imooc.com/coding-535/volume/proto/volume"

)

var (
    // Service address
	hostIp = "192.168.0.108"
    // Service address
    serviceHost = hostIp
    // Service port
	servicePort = "8081"
     
	// Registry configuration
	consulHost  = hostIp
	consulPort int64 = 8500
	// Distributed tracing
	tracerHost = hostIp
	tracerPort = 6831
	// Circuit breaker port, each service must be unique
	//hystrixPort = 9092
	// Monitoring port, each service must be unique
	prometheusPort = 9192
)

func main() {
    // Need to start mysql, consul middleware services locally
	//1. Service registry
	consul:=consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			consulHost+":"+strconv.FormatInt(consulPort,10),
		}
	})
	//2. Configuration center for frequently changing variables
	consulConfig,err := common.GetConsulConfig(consulHost,consulPort,"/micro/config")
	if err !=nil {
		common.Error(err)
	}
	//3. Use configuration center to connect to mysql
	mysqlInfo := common.GetMysqlFromConsul(consulConfig,"mysql")
	// Initialize database
	db,err := gorm.Open("mysql",mysqlInfo.User+":"+mysqlInfo.Pwd+"@("+mysqlInfo.Host+":3306)/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err !=nil {
        // Command line output for easy error viewing
		fmt.Println(err)
		common.Fatal(err)
	}
	defer db.Close()
	// Disable plural table names
	db.SingularTable(true)

	//4. Add distributed tracing
	t,io,err := common.NewTracer("go.micro.service.volume",tracerHost+":"+strconv.Itoa(tracerPort))
	if err !=nil {
		common.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// Add circuit breaker, enable as client
	//hystrixStreamHandler := hystrix.NewStreamHandler()
	//hystrixStreamHandler.Start()

	// Add logging center
	//1) Program logs need to be written to log files
	//2) Add filebeat.yml file to the program
	//3) Start filebeat with command: ./filebeat -e -c filebeat.yml
	fmt.Println("Logs are unified in the root directory micro.log file, please click to view logs!")

	// Start monitoring program
	//go func() {
	//	//http://192.168.0.112:9092/turbine/turbine.stream
	//	// Dashboard access address: http://127.0.0.1:9002/hystrix, URL must end with /hystrix
	//	err = http.ListenAndServe(net.JoinHostPort("0.0.0.0",strconv.Itoa(hystrixPort)),hystrixStreamHandler)
	//	if err !=nil {
	//		common.Error(err)
	//	}
	//}()

	//5. Add monitoring
	common.PrometheusBoot(prometheusPort)


	// Download kubectl: https://kubernetes.io/docs/tasks/tools/#tabset-2
	// macOS:
	// 1.curl -LO "https://dl.k8s.io/release/v1.21.0/bin/darwin/amd64/kubectl"
	// 2.chmod +x ./kubectl
	// 3.sudo mv ./kubectl /usr/local/bin/kubectl
	//   sudo chown root: /usr/local/bin/kubectl
	// 5.kubectl version --client
	// 6. In cluster mode, directly copy the server ~/.kube/config file to local ~/.kube/config
	//    Note: - Domain names in config must resolve correctly
	//        - Production environment can create another certificate
	// 7.kubectl get ns to check if normal
	//
	//6. Create k8s connection
	// Connect from outside the cluster
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		common.Fatal(err.Error())
	}

	// Configuration outside the cluster
	//config, err := rest.InClusterConfig()
	//if err != nil {
	//	panic(err.Error())
	//}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		common.Fatal(err.Error())
	}

	//7. Create service
	service := micro.NewService(
		// Custom service address, must be written before other parameters
		micro.Server(server.NewServer(func(options *server.Options) {
			options.Advertise =serviceHost+":"+servicePort
		})),
		micro.Name("go.micro.service.volume"),
		micro.Version("latest"),
		// Specify service port
		micro.Address(":"+servicePort),
		// Add service registry
		micro.Registry(consul),
		// Add distributed tracing
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// Only works when acting as a client, if there are calls to others, don't actively call in principle
		//micro.WrapClient(hystrix2.NewClientHystrixWrapper()),
		// Add rate limiting
		micro.WrapHandler(ratelimit.NewHandlerWrapper(1000)),
	)
 
	service.Init()

	// Can only execute once
	//err = repository.NewVolumeRepository(db).InitTable()
	//if err != nil {
	//	common.Fatal(err)
	//}

	// Register handler for quick access to developed services
	volumeDataService:=service2.NewVolumeDataService(repository.NewVolumeRepository(db),clientset)
	volume.RegisterVolumeHandler(service.Server(), &handler.VolumeHandler{ VolumeDataService:volumeDataService})

	// Start service
	if err := service.Run(); err != nil {
        // Output startup failure information
		common.Fatal(err)
	}
}

