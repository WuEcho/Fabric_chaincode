package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type HelloChaincode struct {
}

func (h *HelloChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	//初始化、升级链码时调用
	//_,args := GetFunctionAndParameters()
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("初始化链码时发生错误。。。,初始化参数必须为2个")
	}
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error("保存状态数据时发生错误，%s")
	}
	fmt.Println("链码初始化成功")
	return shim.Success(nil)
}

func (h *HelloChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	args := stub.GetStringArgs()
	if len(args) != 1 {
		return shim.Error("传递的参数个数错误，只能传递相应的key")
	}

	result, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("查询过程中发生错误")
	}

	if result == nil {
		return shim.Error("根据指定的key没有查询到相应的结果")
	}

	return shim.Success(result)
}

func main() {
	err := shim.Start(new(HelloChaincode))

	if err != nil {
		fmt.Printf("start chaincode faild.. %s", err)
	}
}
