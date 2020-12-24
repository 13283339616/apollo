package apollo

func GetSnowflakeId() (snowId string) {
	node, err := SnowFlakeNewWorker(1)
	if err != nil {
		panic(err)
	}
	snowId = string(node.SnowFlakeGetId())
	return
}
