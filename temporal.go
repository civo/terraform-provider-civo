// +build ignore

package main

//import (
//	"fmt"
//	"github.com/civo/civogo"
//)
//
//const (
//	apiKey = "nTF2YPA5E1i0f4cB3jx6KkWDLgd8lvIJtwCpUerMNXHmGsQhzO"
//)
//
//func main() {
//	client, _ := civogo.NewClient(apiKey)
//
//	config := &civogo.NetworkConfig{Label: "Temporal-terraform"}
//
//	//test := civogo.Network{}
//	//
//	resp, err := client.RenameNetwork(config)
//	if err != nil {
//		fmt.Errorf("failed to create a new config: %s", err)
//	}
//	//
//	//for _, network := range resp {
//	//	if network.ID == "d2579bec-68f4-4923-83af-ddd0abe9eaf8"{
//	//		test = network
//	//	}
//	//}
//
//	fmt.Println(test.CIDR)
//}

import (
	"fmt"
	"github.com/civo/civogo"
)

const apiKey = "nTF2YPA5E1i0f4cB3jx6KkWDLgd8lvIJtwCpUerMNXHmGsQhzO"

func main() {

	client, _ := civogo.NewClient(apiKey)

	//resp, _ := client.GetDNSRecord("a5042dea-0494-4ff0-9f3d-6b0dc72b1468", "499a028f-1859-46c5-ae21-3227e0898c36")
	//
	//config := &civogo.LoadBalancerConfig{Hostname: "pepe.domain.com", Port: "80", Protocol: "http", Backends:[]civogo.LoadBalancerBackendConfig{{InstanceID: "cf7ff95a-0e88-4e57-9202-b3476e8451fb", Port: 80, Protocol: "http"}}}
	//_, _ = client.UpdateDNSRecord(resp, config)

	//resp, err := client.CreateLoadBalancer(config)
	resp, err := client.FindFirewallRule("446c091e-b0fb-4469-8e52-3306994f3604", "8ff")
	if err != nil {
		fmt.Printf("(%s)", err)
	}

	fmt.Println(resp)
}
