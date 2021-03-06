// Copyright the Hyperledger Fabric contributors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/hyperledgendary/fabric-chaincode-wasm/wasmruntime"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/wapc/wapc-go"
)

// ChaincodeConfig is used to configure the chaincode server. See chaincode.env.example
type ChaincodeConfig struct {
	CCID    string
	Address string
	WasmCC  string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	log.Printf("[host] Wasm Contract runtime..")
	ctx := context.Background()

	config := ChaincodeConfig{
		CCID:    os.Getenv("CHAINCODE_ID"),
		Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
		WasmCC:  os.Getenv("CHAINCODE_WASM_FILE"),
	}

	wasmBytes, err := ioutil.ReadFile(config.WasmCC)
	check(err)

	module, err := wapc.New(consoleLog, wasmBytes, hostCall)
	if err != nil {
		panic(err)
	}
	defer module.Close()

	instance, err := module.Instantiate()
	if err != nil {
		panic(err)
	}
	defer instance.Close()

	// result, err := instance.Invoke(ctx, "hello", []byte(name))
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(result))

	wrt := wasmruntime.NewRuntime(ctx, *instance)

	server := &shim.ChaincodeServer{
		CCID:    config.CCID,
		Address: config.Address,
		CC:      wrt,
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}

	err = server.Start()
	check(err)

	return
}

func consoleLog(msg string) {
	fmt.Println(msg)
}

func hostCall(ctx context.Context, binding, namespace, operation string, payload []byte) ([]byte, error) {
	// Route the payload to any custom functionality accordingly.
	// You can even route to other waPC modules!!!
	log.Printf("bd %s ns %s op %s payload length %d\n", binding, namespace, operation, len(payload))

	// Todo add default cases?
	switch namespace {
	case "LedgerService":
		switch operation {
		case "CreateState":
			return []byte("CreateState ftw!"), nil
		case "ReadState":
			return []byte("ReadState ftw!"), nil
		}
	}
	return []byte("Hello from Go"), nil
}
