package trigger

import (
	"testing"
)

const trigger1 = `
{
   "TriggerId":100,
   "TriggerName":"Basic Filter 1",
   "TriggerType":"WatchTransactions",
   "CreatorId":223,
   "CreationDate":"2019-03-24 17:45:12",
   "ContractABI":"",
   "Filters":[
      {
         "FilterType":"BasicFilter",
         "ParameterName":"To",
         "ParameterType":"Address",
         "Condition":{
            "Predicate":"Is",
            "Attribute":"0x174bfa6600bf90c885c7c01c7031389ed1461ab9"
         }
      },
      {
         "FilterType":"BasicFilter",
         "ParameterName":"Nonce",
         "ParameterType":"Int",
         "Condition":{
            "Predicate":"BiggerThan",
            "Attribute": "1000"
         }
      }
   ]
}`

const trigger2 = `
{
   "TriggerId":101,
   "TriggerName":"Basic + Function filters",
   "TriggerType":"WatchTransactions",
   "CreatorId":223,
   "CreationDate":"2019-03-24 17:45:12",
   "ContractABI":"",
   "Filters":[
      {
         "FilterType":"BasicFilter",
         "ParameterName":"To",
         "ParameterType":"Address",
         "Condition":{
            "Predicate":"Is",
            "Attribute":"0xe8663a64a96169ff4d95b4299e7ae9a76b905b31"
         }
      },
      {
         "FilterType":"CheckFunctionParameter",
   		 "ToContract":"0xe8663a64a96169ff4d95b4299e7ae9a76b905b31",
         "FunctionName":"depositToken",
         "ParameterName":"_to",
         "ParameterType":"Address",
         "Condition":{
            "Predicate":"Is",
            "Attribute":"0000000000000000000000007abe49749989a53b8d9e584b0ee93bb773ca0b9e"
         }
      },
      {
         "FilterType":"CheckFunctionParameter",
   		 "ToContract":"0xe8663a64a96169ff4d95b4299e7ae9a76b905b31",
         "FunctionName":"depositToken",
         "ParameterName":"_not_there",
         "ParameterType":"Address",
         "Condition":{
            "Predicate":"Is",
            "Attribute":"0000000000000000000000007abe49749989a53b8d9e584b0ee93bb773ca0b9e"
         }
      }
   ]
}`

func TestNewTriggerJson(t *testing.T) {
	_, err := NewTriggerJson(trigger2)
	if err != nil {
		t.Error(err)
	}
}

func TestTriggerJson_ToTrigger(t *testing.T) {
	tjs, _ := NewTriggerJson(trigger1)
	trig, err := tjs.ToTrigger()
	if err != nil {
		t.Error(err)
	}
	_, ok := trig.Filters[0].Condition.(ConditionTo)
	if ok != true {
		t.Error("Expected type ConditionTo")
	}
}

func TestValidateFilter(t *testing.T) {

	block := getBlockFromFile("../resources/block.json")
	trigger := getTriggerFromJson(trigger2)
	abi := `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"pausedPublic","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_value","type":"uint256"}],"name":"burn","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"pausedOwnerAdmin","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_subtractedValue","type":"uint256"}],"name":"decreaseApproval","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_value","type":"uint256"}],"name":"burnFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newAdmin","type":"address"}],"name":"changeAdmin","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_addedValue","type":"uint256"}],"name":"increaseApproval","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"token","type":"address"},{"name":"amount","type":"uint256"}],"name":"emergencyERC20Drain","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newPausedPublic","type":"bool"},{"name":"newPausedOwnerAdmin","type":"bool"}],"name":"pause","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"admin","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_admin","type":"address"},{"name":"_totalTokenAmount","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_burner","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Burn","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"previousAdmin","type":"address"},{"indexed":true,"name":"newAdmin","type":"address"}],"name":"AdminTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"newState","type":"bool"}],"name":"PausePublic","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"newState","type":"bool"}],"name":"PauseOwnerAdmin","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"previousOwner","type":"address"},{"indexed":true,"name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`

	// Testing one filter VS one transaction

	// BasicFilter / To

	if ValidateFilter(block.Transactions[0], trigger.Filters[0], abi) != true {
		t.Error("Basic Filter / To should match")
	}
	if ValidateFilter(block.Transactions[1], trigger.Filters[0], abi) != false {
		t.Error("Basic Filter / To should NOT match")
	}

	// FunctionParameter / Address

	if ValidateFilter(block.Transactions[0], trigger.Filters[1], abi) != true {
		t.Error("FuncParam should match")
	}

	if ValidateFilter(block.Transactions[1], trigger.Filters[1], abi) != false {
		t.Error("FuncParam should NOT match")
	}

}
