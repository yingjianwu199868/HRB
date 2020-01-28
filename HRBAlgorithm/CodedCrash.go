package HRBAlgorithm

import (
	"fmt"
	"strconv"
)

var codeCounter map[string] int
var codeElements map[string] []string

func InitCrash() {
	codeCounter = make(map[string] int)
	codeElements= make(map[string] []string)
}

func CrashECBroadCast(s string, round int) {
	//need to make sure that coded element > f
	fmt.Println("MyID" + MyID)
	var shards[][] byte
	if faulty == 0 {
		shards = Encode(s, total - 1, 1)
	} else {
		shards = Encode(s, faulty + 1, total - (faulty + 1))
	}
	fmt.Println("Shards are ", shards)


	for i := 0; i < total; i++ {
		code := ConvertBytesToString(shards[i])
		m := MSGStruct{Header:MSG, Id:MyID, SenderId:MyID, Data: code, Round: round}
		sendReq := PrepareSend{M: m, SendTo: serverList[i]}
		SendReqChan <- sendReq
	}
}

func identifierCreate(id string, round int) string{
	return id + ":" + strconv.Itoa(round)
}

func crashRecMsg(m Message) {
	identifier := identifierCreate(m.GetId(), m.GetRound())
	count, exist :=codeCounter[identifier]

	if exist {
		codeCounter[identifier] = count + 1
	} else {
		codeCounter[identifier] = 1
		codeElements[identifier] = make([]string, total)
	}

	intMyId, _ := serverMap[MyID]
	codeElements[identifier][intMyId] = m.GetData()

	code := m.GetData()
	id := m.GetId();
	round := m.GetRound()
	//Send Echo
	for i := 0; i < total; i++ {
		message := ECHOStruct{Header:ECHO, Id:id, SenderId:MyID, Data: code, Round: round}
		sendReq := PrepareSend{M: message, SendTo: serverList[i]}
		SendReqChan <- sendReq
	}
}

func listToShards(list []string) [][]byte{
	shards := make([][]byte, total)
	for i:=0; i < len(list); i++ {
		if len(list[i]) != 0 {
			shards[i], _ = ConvertStringToBytes(list[i])
		}
	}
	fmt.Println(shards)
	return shards
}

func crashRecEcho(m Message) {
	if m.GetSenderId() != MyID {
		identifier := identifierCreate(m.GetId(), m.GetRound())
		count, exist :=codeCounter[identifier]

		if exist {
			codeCounter[identifier] = count + 1
		} else {
			codeCounter[identifier] = 1
			codeElements[identifier] = make([]string, total)
		}

		senderId,_ := serverMap[m.GetSenderId()]
		codeElements[identifier][senderId] = m.GetData()

		var data string
		shards := listToShards(codeElements[identifier])
		id := m.GetId()
		round := m.GetRound()

		if faulty == 0 {
			if count + 1 == total - 1 {
				data, _ = Decode(shards, total - 1, 1)
				for i := 0; i < total; i++ {
					m := ACCStruct{Header:ACC, Id:id, SenderId:MyID, HashData: data, Round: round}
					sendReq := PrepareSend{M: m, SendTo: serverList[i]}
					SendReqChan <- sendReq
				}
			}

		} else if count + 1 == faulty + 1 {
			//decode elements back
			data, _ = Decode(shards, faulty + 1, total - (faulty + 1 ))

			for i := 0; i < total; i++ {
				m := ACCStruct{Header:ACC, Id:id, SenderId:MyID, HashData: data, Round: round}
				sendReq := PrepareSend{M: m, SendTo: serverList[i]}
				SendReqChan <- sendReq
			}
		}
	}
}

func crashRecAcc(m Message) {
	fmt.Println("Receive " + m.GetHashData() + " " + identifierCreate(m.GetId(), m.GetRound()))
}