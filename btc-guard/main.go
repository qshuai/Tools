package btc_guard

import (
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/btcsuite/btcd/rpcclient"
)

func main() {
	logs.Info("program start")

	conf, err := config.NewConfig("ini", "app.conf")
	if err != nil {
		panic(err)
	}

	logs.Info("create rpc client")
	client := newRPCClient(conf)
	logs.Info("create rpc client successfully!")

	i := 1
	retChan := make(chan interface{})

	for {
		logs.Info("exec %d times!", i)

		go func() {
			retChan <- getInfo(client)
		}()

		select {
		case ret := <-retChan:
			if !isRunning(conf) {
				if err = start(conf); err != nil {
					logs.Error(err)
					continue
				} else {
					// stop extra ten minute because of the possible reindex incoming
					time.Sleep(60 * time.Minute)
					logs.Info("start omnicored process successfully")
				}
			}
			logs.Info("rpc result:", ret)
		case <-time.After(30 * time.Second):
			logs.Error("rpc request timeout...")
			if err := stop(conf); err != nil {
				logs.Error(err)
				continue
			} else {
				logs.Info("stop omnicored process successfully")
				time.Sleep(1 * time.Minute)
				if err = start(conf); err != nil {
					logs.Error(err)
					continue
				} else {
					// stop extra ten minute because of the possible reindex incoming
					time.Sleep(60 * time.Minute)
					logs.Info("start omnicored process successfully")
				}
			}
		}

		time.Sleep(5 * time.Minute)
	}
}

func getInfo(client *rpcclient.Client) chan interface{} {
	logs.Info("start exec getinfo rpc command")
	retChan := make(chan interface{})

	ret, err := client.GetInfo()
	if err != nil {
		return nil
	}

	retChan <- ret

	return retChan
}

func isRunning(conf config.Configer) bool {
	cmd := exec.Command("/bin/bash", conf.String("cmd::isrunning"))
	r, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		logs.Error("start omnicored command failed:", err)
		return false
	}

	ret, _ := ioutil.ReadAll(r)
	if string(ret) == "" {
		return false
	}

	return false
}

func start(conf config.Configer) error {
	cmd := exec.Command("/bin/bash", conf.String("cmd::start"))
	err := cmd.Start()
	if err != nil {
		logs.Error("start omnicored command failed:", err)
		return err
	}

	logs.Info("program started!")
	return nil
}

func stop(conf config.Configer) error {
	cmd := exec.Command("/bin/bash", conf.String("cmd::stop"))
	err := cmd.Start()
	if err != nil {
		logs.Error("stop omnicored command failed:", err)
		return err
	}

	logs.Info("program stopped!")
	return nil
}

func newRPCClient(conf config.Configer) *rpcclient.Client {
	// acquire configure item
	link := conf.String("rpc::url") + ":" + conf.String("rpc::port")
	user := conf.String("rpc::user")
	passwd := conf.String("rpc::passwd")

	// rpc client instance
	connCfg := &rpcclient.ConnConfig{
		Host:         link,
		User:         user,
		Pass:         passwd,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		panic(err)
	}

	return client
}
