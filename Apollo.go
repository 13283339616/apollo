package apollo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var InitApolloConfigArr = make(map[string]map[string]string)
var initApolloNamespaceVersion = make([]apolloNameSpace, 0, 8)

type apolloNameSpace struct {
	NamespaceName  string `json:"namespaceName"`
	NotificationId int    `json:"notificationId"`
}
type apolloConfig struct {
	AppId          string            `json:"appId"`
	Cluster        string            `json:"cluster"`
	NamespaceName  string            `json:"namespaceName"`
	Configurations map[string]string `json:"configurations"`
}

func syncApollo(apolloUrl, appId string, changeChan chan interface{}) {
	for {
		time.Sleep(time.Minute)
		fmt.Print("apollo开始同步远程")
		getApolloNamespaceVersion(apolloUrl, appId)
		var wgSyncApollo sync.WaitGroup
		for _, v := range initApolloNamespaceVersion {
			wgSyncApollo.Add(1)
			go syncNamespaceApollo(v.NamespaceName, apolloUrl, appId, &wgSyncApollo)
		}
		wgSyncApollo.Wait()
		changeChan <- new(interface{})
		fmt.Print("apollo同步成功")
	}
}
func syncNamespaceApollo(namespace, apolloUrl, appId string, wgNameApollo *sync.WaitGroup) {
	defer wgNameApollo.Done()
	resp, err := http.Get(apolloUrl + "/configs/" + appId + "/default/" + namespace)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var apolloConfig = new(apolloConfig)
	json.Unmarshal(bytes, &apolloConfig)
	var apolloConfigMap = make(map[string]string, 32)
	for k, v := range apolloConfig.Configurations {
		apolloConfigMap[k] = v
	}
	InitApolloConfigArr[apolloConfig.NamespaceName] = apolloConfigMap
}

func getApolloNamespaceVersion(apolloUrl, appId string) {
	marshal, err := json.Marshal(initApolloNamespaceVersion)
	if err != nil {
		panic(err)
	}
	urlObj, err := url.ParseRequestURI(apolloUrl + "/notifications/v2")
	if err != nil {
		panic(err)
	}
	data := url.Values{}
	data.Add("appId", appId)
	data.Add("cluster", "default")
	data.Add("notifications", string(marshal))
	encodeData := data.Encode()
	urlObj.RawQuery = encodeData
	request, err := http.NewRequest("GET", urlObj.String(), nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &initApolloNamespaceVersion)
}

func Init(appId, url, namespace string) (changeChan chan interface{}) {

	var wgApollo sync.WaitGroup
	for _, s := range strings.Split(namespace, ",") {
		apolloNameSpace := apolloNameSpace{
			NamespaceName:  s,
			NotificationId: -1,
		}
		initApolloNamespaceVersion = append(initApolloNamespaceVersion, apolloNameSpace)
	}
	getApolloNamespaceVersion(url, appId)
	for _, v := range initApolloNamespaceVersion {
		wgApollo.Add(1)
		go syncNamespaceApollo(v.NamespaceName, url, appId, &wgApollo)
	}
	wgApollo.Wait()

	go syncApollo(url, appId, changeChan)
	return
}
