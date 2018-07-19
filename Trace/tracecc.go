package main

import (
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	Bank_Flag_Loan      = 1
	Bank_Flag_Repayment = 2
)

//保存状态
func putAccount(stub shim.ChaincodeStubInterface, account Account) bool {

	accBytes, err := json.Marshal(account)
	if err != nil {
		return false
	}

	err = stub.PutState(account.CardNo, accBytes)
	if err != nil {
		return false
	}
	return true
}

//获取状态
func getAccount(stub shim.ChaincodeStubInterface, cardNo string) (Account, bool) {

	var account Account
	b, err := stub.GetState(cardNo)

	if err != nil {
		return account, false
	}

	json.Unmarshal(b, &account)

	return account, true
}

//贷款
//-c '{"Args":["loan","身份账户id","银行名称","金额"]}'
func loan(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("指定的所需参数错误")
	}

	v, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("")
	}

	bank := Bank{
		BankName:  args[1],
		Amount:    v,
		Flag:      Bank_Flag_Loan,
		StartTime: "2011-01-10",
		EndTime:   "2021-01-09",
	}

	account := Account{
		CardNo: args[0],
		Aname:  "jack",
		Gender: "男",
		Mobile: "1509923812",
		Bank:   bank,
	}

	b := putAccount(stub, account)
	if !b {
		return shim.Error("保存贷款数据失败")
	}
	return shim.Success([]byte("保存贷款数据成功"))
}

//还款
//-c '{"Args:["repayment","身份账号ID","银行名称","金额"]"}'
func repayment(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("指定的所需参数错误")
	}

	v, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("指定的还款金额错误")
	}

	bank := Bank{
		BankName:  args[1],
		Amount:    v,
		Flag:      Bank_Flag_Repayment,
		StartTime: "2011-09-10",
		EndTime:   "2012-07-09",
	}

	account := Account{
		CardNo: args[0],
		Aname:  "jack",
		Gender: "男",
		Mobile: "7633276",
		Bank:   bank,
		//历史数据无需考虑 由fabric自行维护
	}

	b := putAccount(stub, account)
	if !b {
		return shim.Error("存还款数据失败")
	}
	return shim.Success([]byte("保存还款数据成功"))
}

//查询历史信息
//-c '{"Args":["queryAccountByCardNo","身份证号码"]}'
func queryAccountByCardNo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("queryAccountByCardNo")
	}

	acc, bl := getAccount(stub, args[0])
	if !bl {
		return shim.Error("根据指定的身份证号码查询对应数据时错误")
	}

	//查询历史数据
	accIterator, err := stub.GetHistoryForKey(acc.CardNo)
	if err != nil {
		return shim.Error("获取历史数据时发生错误")
	}

	defer accIterator.Close()

	var historys []HistoryItem
	var account Account
	for accIterator.HasNext() {
		hisData, err := accIterator.Next()
		if err != nil {
			return shim.Error("处理迭代器数据时发生错误")
		}
		var hisItem HistoryItem
		//获取当前的交易编号
		hisItem.TxID = hisData.TxId
		//将交易数据转换
		err = json.Unmarshal(hisData.Value, &account)
		if err != nil {
			return shim.Error("反序列化历史数据时发生错误")
		}

		if hisData.Value == nil {
			var empty Account
			hisItem.Account = empty
		} else {
			hisItem.Account = account
		}

		historys = append(historys, hisItem)
	}

	acc.Historys = historys

	accByte, err := json.Marshal(acc)
	if err != nil {
		return shim.Error("序列化数据时发生错误")
	}

	return shim.Success(accByte)
}
