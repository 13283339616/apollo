package apollo

import (
	"os"
	"strings"
)

var initProperties = make(map[string]string)
var allProperties = make(map[string]string)

var changeChan = make(chan interface{})
var configChan = make(chan interface{})

func initAll(path string) {
	sep := string(os.PathSeparator)
	initProperties = ReadProperties(path + sep + "config.properties")

	if includeFilePath, ok := initProperties["five.include"]; ok {
		includeFilePathList := strings.Split(includeFilePath, ",")
		for _, includeFilePath := range includeFilePathList {
			includeFileProperty := ReadProperties(path + sep + includeFilePath + ".properties")
			for k, v := range includeFileProperty {
				initProperties[k] = v
			}
		}
	}
	if apolloConfig, ok := initProperties["five.apollo"]; ok {
		if apolloConfig == "on" {
			apolloEnv, envOk := initProperties["five.apollo.env"]
			if !envOk {
				panic("读取apollo配置异常")
			}
			url, urlOK := initProperties[apolloEnv+".meta"]
			if !urlOK {
				panic("读取apollo配置异常")
			}
			appId, appIdOk := initProperties["five.apollo.appId"]
			if !appIdOk {
				panic("读取apollo配置异常")
			}

			namespace, namespaceOk := initProperties["five.apollo.namespace"]
			if !namespaceOk {
				panic("读取apollo配置异常")
			}
			changeChan = Init(appId, url, namespace)
			syncConfig()
		}
	}

}

func syncConfig() {
	apolloProperMap := InitApolloConfigArr
	apolloProperties := make(map[string]string, 32)
	for key, value := range apolloProperMap {
		for k, v := range value {
			apolloProperties[key+"."+k] = v
		}
	}
	for k, v := range initProperties {
		if strings.HasPrefix(v, "${") && strings.HasSuffix(v, "}") {
			s := v[1 : len(v)-1]
			if prop, ok := apolloProperties[s]; ok {
				allProperties[k] = prop
			}
		} else {
			allProperties[k] = v
		}
	}
	for k, v := range apolloProperties {
		allProperties[k] = v
	}
}

func syncApolloTwo() {
	for {
		select {
		case <-changeChan:
			syncConfig()
			configChan <- new(interface{})
		}
	}
}

func GetAllProperties() map[string]string {
	return allProperties
}

func GetPropertiesByKey(key string) string {
	prop, ok := allProperties[key]
	if !ok {
		panic("读取配置出错")
	}
	return prop
}

func GetConfigChan() chan interface{} {
	return configChan
}
