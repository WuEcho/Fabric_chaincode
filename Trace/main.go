package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type TraceChaincode struct {
}

func (t *TraceChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {

	//初始化
	initTest(stub)
	return shim.Success(nil)
}

/**
 * loan: 贷款
 * repayemnt: 还款
 * queryAccountByCardNo: 根据账户身份证号码查询相应数据(包含历史数据)
 */
func (t *TraceChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fun, args := stub.GetFunctionAndParameters()
	if fun == "loan" {
		return loan(stub, args)
	} else if fun == "repayment" {
		return repayment(stub, args)
	} else if fun == "queryAccountByCardNo" {
		return queryAccountByCardNo(stub, args)
	}

	return shim.Error("调用参数名称错误")
}

func initTest(stub shim.ChaincodeStubInterface) peer.Response {

	bank := Bank{
		BankName:  "ccic",
		Amount:    6000,
		Flag:      1,
		StartTime: "2010-01-10",
		EndTime:   "2011-01-09",
	}

	bank1 := Bank{
		BankName:  "ccic",
		Amount:    1000,
		Flag:      2,
		StartTime: "2010-02-10",
		EndTime:   "2011-02-09",
	}

	account := Account{
		CardNo:   "12321421",
		Aname:    "jack",
		Gender:   "男",
		Mobile:   "1509923812",
		Bank:     bank,
		Historys: nil,
	}

	account2 := Account{
		CardNo:   "12321421",
		Aname:    "jack",
		Gender:   "男",
		Mobile:   "1509923812",
		Bank:     bank1,
		Historys: nil,
	}

	//将对象进行存储或在网络上进行传输必须要对其序列化
	accBtyes, err := json.Marshal(account)
	if err != nil {
		return shim.Error("序列化账户对象失败")
	}

	accBtyes1, err := json.Marshal(account2)
	if err != nil {
		return shim.Error("序列化账户2对象失败")
	}

	//保存
	err = stub.PutState(account.CardNo, accBtyes)
	if err != nil {
		return shim.Error("保存账户1到账本失败")
	}

	err = stub.PutState(account2.CardNo, accBtyes1)
	if err != nil {
		return shim.Error("保存账户2到账本失败")
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(TraceChaincode))
	if err != nil {
		fmt.Printf("启动链码时发生错误")
	}
}
