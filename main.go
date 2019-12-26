package main

import (
	"fmt"
	block "github.com/corgi-kx/blockchain_golang/blc"
	"github.com/corgi-kx/blockchain_golang/cli"
	"github.com/corgi-kx/blockchain_golang/database"
	"github.com/corgi-kx/blockchain_golang/network"
	log "github.com/corgi-kx/logcustom"
	"github.com/spf13/viper"
	"os"
)

//初始化系统,读取config.yaml里面的配置信息并进行赋值
func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	logPath := viper.GetString("blockchain.log_path")
	listenHost := viper.GetString("network.listen_host")
	listenPort := viper.GetString("network.listen_port")
	rendezvousString := viper.GetString("network.rendezvous_string")
	protocolID := viper.GetString("network.protocol_id")
	tokenRewardNum := viper.GetInt("blockchain.token_reward_num")
	tradePoolLength := viper.GetInt("blockchain.trade_pool_length")
	mineDifficultyValue := viper.GetInt("blockchain.mine_difficulty_value")
	chineseMnwordPath := viper.GetString("blockchain.chinese_mnemonic_path")

	network.TradePoolLength = tradePoolLength
	network.ListenHost = listenHost
	network.RendezvousString = rendezvousString
	network.ProtocolID = protocolID
	network.ListenPort = listenPort
	database.ListenPort = listenPort
	block.ListenPort = listenPort
	block.TokenRewardNum = tokenRewardNum
	block.TargetBits = uint(mineDifficultyValue)
	block.ChineseMnwordPath = chineseMnwordPath

	//将日志输出到指定文件
	file, err := os.OpenFile(fmt.Sprintf("%slog%s.txt", logPath, listenPort), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Error(err)
	}
	log.SetOutputAll(file)
}

func main() {
	//引入命令行信息,并进行调用
	c := cli.New()
	c.Run()
}
