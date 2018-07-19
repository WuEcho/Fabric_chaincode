package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

type PaymentChaincode struct {

}

// -c '{"Args":["init","第一个账户名称","第一个账户余额",”第二个账户名称“,”第二个账户余额“]}'
func (p *PaymentChaincode)Init(stub shim.ChaincodeStubInterface)peer.Response  {

	_,args := stub.GetFunctionAndParameters()
	if len(args) != 4 {
		return shim.Error("必须指定两个账户的名称以及对应金额")
	}

	var err error
	//判断第一个账户和第二个账户余额类型
    _,err = GetArgsState(args[1])
	if err != nil{
		return shim.Error("指定的第一个账户的金额错误")
	}

    _,err = GetArgsState(args[3])
	if err != nil {
		return shim.Error("指定的第二个账户的金额错误")
	}

	//将对初始信息保存在账本中
	err = stub.PutState(args[0],[]byte(args[1]))

	if err != nil {
		return shim.Error("第一个账户信息保存失败")
	}

    err = stub.PutState(args[2],[]byte(args[3]))

	if err != nil {
		return shim.Error("第二个账户信息保存失败")
	}

	fmt.Println("初始化成功")

    return shim.Success(nil)
}


//从指定的账户转账给指定的目标账户指定金额
func (p *PaymentChaincode)Invoke(stub shim.ChaincodeStubInterface) peer.Response  {
	fun,args := stub.GetFunctionAndParameters()
	if fun == "query" {
		return query(stub,args)
	}else if fun == "invoke" {
		return invoke(stub,args)
	}else if fun == "set" {
		return set(stub,args)
	}else if fun == "get"{
		return get(stub,args)
	}
	return shim.Success(nil)
}

//根据指定账户查询
func query(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("必须且只能指定一个要查询的账户名称")
	}	
	
	result,err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据指定的账户名称查询对应状态数据时发生错误")
	}

	if result == nil {
		return shim.Error("根据指定的账户名称没有查询到相应的数据")
	}
	
	return shim.Success(result)
}

//转账
//-c '{"Args"：【"invoke","原账户","目标账户","转账金额"】}'
func invoke(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	//1.判断参数的长度是否符合要求
	if len(args) != 3 {
		return shim.Error("参数错误，必须指定原账户，目标账户，转账金额")
	}
	
	//判断转账金额类型是否正确
	v,err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("转账金额错误，请重新设置")
	}
	
	//查询原账户里面的余额
	v1,err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据指定的原账户名称查询出现错误")
	}

	if v1 == nil {
		return shim.Error("根据指定的原账户名称没有查询到数据")
	}
	
	//将查询的数据转换类型
	v2,err := strconv.Atoi(string(v1))
	if err != nil {
		return shim.Error("转换查询到的原账户金额类型的时候出错")
	}
	
	if v2 < v {
		//如果转账金额大于原账户中有的金额就会报错
		return shim.Error("转账的金额大于原账户中的余额，请重新操作")
	}else {
		v2 = v2 - v
	}
	
	//查询目标账户的余额
	var tarv int
	bv,err := stub.GetState(args[1])
	if err != nil {
		return shim.Error("根据指定的目标账户名称查询出现错误")
	}

	if bv == nil {
		err = stub.PutState(args[1],[]byte(args[2]))

		if err != nil{
			return shim.Error("更新目标账户的状态发生错误")
		}

	}else {
		tarv,err = strconv.Atoi(string(bv))
		if err != nil {
			return shim.Error("转换金额类型出错")
		}
		
		tarv = tarv + v
	}

	err = stub.PutState(args[1],[]byte(strconv.Itoa(tarv)))

	if err != nil{
		return shim.Error("更新目标账户的状态发生错误")
	}


	//更新账户的余额
 	err	= stub.PutState(args[0],[]byte(strconv.Itoa(v2)))
	if err != nil {
		return shim.Error("更新原账户的状态发生错误")
	}
	
	
	return shim.Success([]byte("转账成功"))
}

func GetArgsState(value string) (int,error) {
    v,err := strconv.Atoi(value)
	if err != nil {
		return 0,err
	}

	return v,err
}

//实现向指定的账户存钱
// -c '{"Args":["set","目标账户名称","金额"]}'
func set(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("必须指定在存款的目标账户及金额")
	}
	
	//判断金额是否合法
	v,err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("指定存储的目标账户金额出错")
	}
	
	//根据账户查询金额
	result,err := getValueWithkey(stub,args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	
	//将查询到的结果进行转换int
	v1,err := transByteToInt(result)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	v = v + v1
	
	//将新的余额保存
  	err	= stub.PutState(args[0],[]byte(string(v)))
	if err != nil{
		return shim.Error(err.Error())
	}
	
	return shim.Success([]byte("保存成功"))
}

//实现向指定账户中取钱
//-c '{"Args":["get","目标账户名称","金额"]}'
func get(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	result,err := getValueWithkey(stub,args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	
	v1,err := transByteToInt([]byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}
	
	v,err := transByteToInt(result)
	if err != nil{
		return shim.Error(err.Error())
	}

	if v < v1 {
		//账户余额不足
		return shim.Error("余额不足")
	}else {
		
		v = v - v1
	}
	
	//保存状态
	err = stub.PutState(args[0],[]byte(string(v)))
	if err != nil {
		return shim.Error("保存状态错误")
	}
	
	return shim.Success(nil)
}


func getValueWithkey(stub shim.ChaincodeStubInterface,key string) ([]byte,error) {
	//对目标账户进行查询
	result,err := stub.GetState(key)
	if err != nil{
		return nil,fmt.Errorf("查询目标账户出现错误")
	}

	if result == nil {
		return nil,fmt.Errorf("根据目标账户查询没有查询到结果")
	}
	
	return result,nil
}

//将[]byte类型转换成int
func transByteToInt(value []byte) (int,error) {
	v,err := strconv.Atoi(string(value))
	if err != nil {
		return 0,fmt.Errorf("转换格式失败")
	}
	return v,nil
}

func main()  {

	err := shim.Start(new(PaymentChaincode))
	if err != nil {
		fmt.Println("启动链码失败")
	}
}
