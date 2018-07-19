package main

type Bank struct {
	/*银行名称*/
	BankName string `json:"BankName"`
	/*金额*/
	Amount int `json:"Amount"`
	/*标记*/
	Flag int `json:"Flag"` //1.贷款  2.还款
	/*起始时间*/
	StartTime string `json:"StartTime"`
	/*结束时间*/
	EndTime string `json:"EndTime"`
}

type Account struct {
	/*证件号*/
	CardNo string `json:"CardNo"`
	/*用户名*/
	Aname string `json:"Aname"`
	/*性别*/
	Gender string `json:"Gender"`
	/*电话*/
	Mobile string `json:"Mobile"`
	/*银行*/
	Bank Bank `json:"Bank"`

	Historys []HistoryItem
}

type HistoryItem struct {
	TxID    string
	Account Account
}
