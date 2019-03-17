# Private Medical Care

This is Blockchain Project made by Hyperledger fabric. Made by Jungkyu Choi.

I made PMC for patients and docotrs to share information of Patient Medical Chart.

Then, they know each other about their patient medical prescriptions.

So, they can know which drug have side effect to their patient. The each medical information of patients is secure and fast

on Hyperledger fabric. Orderer type: solo, 1orderer, 2orgs, 2peers in each org

If you have question, please don't hesitate to contact me.

Git: https://github.com/sizzflyer/PMC 

Email: sizzflyer@gmail.com

I made MSP, Chaincode.

To launch this,

You need Prerequisites, Hyperledger fabric, fabric-ca, fabric-samples, Golang, Docker etc.

1. Chaincode test

Now, then open 3 Terminal,

Terminal1

> cd $GOPATH/src/github.com/hyperledger/fabric-samples/chaincode-docker-devmode

> docker-compose -f docker-compose-simple.yaml up

Terminal2

> docker exec -it chaincode bash

<!!! You must copy pmc.go file to hyperledger/fabric-samples/chaincode/pmc/pmc.go!!>

> cd chaincode/pmc/

> go build -o pmcjk

> CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0 ./pmcjk

Terminal3

>docker exec -it cli bash

> peer chaincode install -p chaincodedev/chaincode/pmcjk/ -n mycc -v 0

> peer chaincode instantiate -n mycc -v 0 -c '{"Args":["initChart"]}' -C myc

//Invoke

> peer chaincode invoke -n mycc -c '{"Args":["initChart"]}' -C myc

Create new chart. Next type of Struct is MedicalChart. When you invoke args, these args put into world state. I made auto indexing chart Algorithm in pmc2.go. If you invoke "createChart", params input into chain network as key-value. By auto indexing chart Algorithm, the chart index is automatically counted(key:"CHART1","CHART2"..)

	type MedicalChart struct {
  		Name       string `json:"name"`       //환자 이름	  
 		BirthDate  string `json:"birthDate"`  //생년월일	  
  		Diagnosis  string `json:"diagnosis"`  // 진단명(병명)	  
 		Medication string `json:"medication"` //복용 중인 약물 등	  
  		ASTN       string `json:"astn"`       //복용시 부작용이 있는 약물들 After Skin Test Negative	  
 		MCP        string `json:"mcp"`        //앞으로의 진료 계획 Medical Care Plan
	}

'{"Args":["function","Name","BirthDate","Diagnosis","Medication","ASTN","MCP"]}'

> peer chaincode invoke -n mycc -c '{"Args":["createChart","JK","890101","Cold","Aspirin","Antibiotic","Rest at home"]}' -C myc

> peer chaincode invoke -n mycc -c '{"Args":["createChart","SasKe","890101","Headache","Aspirin","None","Come after 3 days"]}' -C myc

//Query

Single chart query. you want CHART1?, just input arg "1". or CHART2?, input arg "2".

> peer chaincode query -n mycc -c '{"Args":["queryChart", "1"]}' -C myc

All chart query.

> peer chaincode query -n mycc -c '{"Args":["queryAllCharts"]}' -C myc

2. End-to-End Test

Query and execute transaction test. I made MSP of PMC. doctororg, patientorg, ordererjk.

> cd PrivateMedicalCare/first-network/

> ./byfn.sh generate

> ./byfn.sh up

If you see the mesage " END ", it is a success.
