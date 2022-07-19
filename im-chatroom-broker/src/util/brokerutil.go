package util

func GetBrokerIp() string {
	ip, e := ExternalIP()
	if e != nil {
		Panic(e)
	}

	brokerAddress := ip.String()
	return brokerAddress
}
