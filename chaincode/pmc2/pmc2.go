// Coding by Jungkyu Choi
// Github: https://github.com/sizzflyer

// Terminal1
// docker-compose -f docker-compose-simple.yaml up
// Terminal2
// docker exec -it chaincode bash
// cd pmc2/
// go build -o pmc2
// CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0 ./pmc2
// Terminal3
// docker exec -it cli bash
// peer chaincode install -p chaincodedev/chaincode/pmc2/ -n mycc -v 0
// peer chaincode instantiate -n mycc -v 0 -c '{"Args":["initChart"]}' -C myc
// peer chaincode invoke -n mycc -c '{"Args":["initChart"]}' -C myc
// peer chaincode invoke -n mycc -c '{"Args":["createChart","Naruto","890101","Sick","Aspirin","juice","rest at home"]}' -C myc
// peer chaincode invoke -n mycc -c '{"Args":["createChart","SasKe","890101","Sick","Aspirin","juice","rest at home"]}' -C myc
// peer chaincode query -n mycc -c '{"Args":["queryChart", "1"]}' -C myc
// peer chaincode query -n mycc -c '{"Args":["queryAllCharts"]}' -C myc

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SmartContractJK struct {
}

// 진료 차트
type MedicalChart struct {
	Name       string `json:"name"`       //환자 이름
	BirthDate  string `json:"birthDate"`  //생년월일
	Diagnosis  string `json:"diagnosis"`  // 진단명(병명)
	Medication string `json:"medication"` //복용 중인 약물 등
	ASTN       string `json:"astn"`       //복용시 부작용이 있는 약물들 After Skin Test Negative
	MCP        string `json:"mcp"`        //앞으로의 진료 계획 Medical Care Plan
}

// 차트 넘버 자동 뽑기
func (s *SmartContractJK) chartNumber(APIstub shim.ChaincodeStubInterface) int {
	startKey := "CHART0"
	endKey := "CHART999"
	var key string

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		fmt.Println(err)
	}
	firstChart, _ := APIstub.GetState(key)
	if firstChart == nil {
		key = startKey
	}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			fmt.Println(err)
		}
		key = queryResponse.Key
	}
	rst := pickNumber(key) + 1
	return rst
}

// 문자에서 숫자 뽑기
func pickNumber(s string) int {
	re := regexp.MustCompile("[0-9]+")
	rstStr := re.FindAllString(s, -1)
	rstNu, err := strconv.Atoi(rstStr[0])
	if err != nil {
		fmt.Println(err)
	}
	return rstNu
}

// Init
func (s *SmartContractJK) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (s *SmartContractJK) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
	function, args := APIstub.GetFunctionAndParameters()
	// 차트 초기화
	// 차트 생성
	// 진료 기록 열람
	// 모든 진료 기록 열람
	if function == "initChart" {
		return s.initChart(APIstub)
	} else if function == "createChart" {
		return s.createChart(APIstub, args)
	} else if function == "queryChart" {
		return s.queryChart(APIstub, args)
	} else if function == "queryAllCharts" {
		return s.queryAllCharts(APIstub)
	}
	return shim.Error("There is no command of these.: initChart, createChart, queryChart, queryAllCharts")
}

// 차트 초기화 (instantiate chaincode on network)
func (s *SmartContractJK) initChart(APIstub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("<initChart 시작>")

	charts := []MedicalChart{
		MedicalChart{Name: "JungkyuChoi", BirthDate: "870115", Diagnosis: "cold", Medication: "drug", ASTN: "drug", MCP: "See u later"},
		MedicalChart{Name: "GildongHong", BirthDate: "450815", Diagnosis: "cold", Medication: "drug", ASTN: "drug", MCP: "Cured"},
	}

	i := 0
	for i < len(charts) {
		chartAsBytes, err := json.Marshal(charts[i])
		if err != nil {
			return shim.Error(err.Error())
		}
		APIstub.PutState("CHART"+strconv.Itoa(i), chartAsBytes)
		fmt.Println("Added CHART"+strconv.Itoa(i), string(chartAsBytes))
		i = i + 1
	}
	fmt.Println("<initChart 종료>")
	return shim.Success(nil)
}

// 차트 생성
func (s *SmartContractJK) createChart(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("<createChart 시작>")

	if len(args) != 6 {
		return shim.Error("Expecting Args: 6")
	}

	var chart = MedicalChart{Name: args[0], BirthDate: args[1], Diagnosis: args[2], Medication: args[3], ASTN: args[4], MCP: args[5]}

	chartAsBytes, err := json.Marshal(chart)
	if err != nil {
		fmt.Println("json화 에러:", err)
	}
	cn := s.chartNumber(APIstub)
	APIstub.PutState("CHART"+strconv.Itoa(cn), chartAsBytes)

	fmt.Println("CHART"+strconv.Itoa(cn), string(chartAsBytes))
	fmt.Println("<createChart 종료>")

	return shim.Success(chartAsBytes)
}

// 진료 차트 열람
func (s *SmartContractJK) queryChart(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("<queryChart 시작>")

	if len(args) > 1 {
		fmt.Println(" 항목 1개만 입력하세요. ")
		return shim.Error("Expecting Arg: 1")
	}

	if _, err := strconv.Atoi(args[0]); err != nil {
		fmt.Println(" 숫자만 입력하세요. ")
		return shim.Error("Only input Number")
	}
	cn := "CHART" + args[0]
	chartAsBytes, err := APIstub.GetState(cn)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println(cn, string(chartAsBytes))
	fmt.Println("<queryChart 종료>")
	return shim.Success(chartAsBytes)
}

// 모든 진료 차트 열람
func (s *SmartContractJK) queryAllCharts(APIstub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("<queryAllCharts 시작>")
	startKey := "CHART0"
	endKey := "CHART999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[\n")

	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}\n")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllCharts:\n%s\n", buffer.String())
	fmt.Println("<queryAllCharts 종료>")

	return shim.Success(buffer.Bytes())
}

func main() {
	err := shim.Start(new(SmartContractJK))
	if err != nil {
		fmt.Printf("스마트 컨트랙트 실행 에러:%s", err)
	}
}
