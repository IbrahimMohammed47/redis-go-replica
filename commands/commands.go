package commands

import (
	"strconv"
	"strings"
	"time"

	"github.com/IbrahimMohammed47/codecrafters-redis-go/resp"
)

var database = make(map[string]resp.Resp)

func HandleCommand(respArr resp.Resp) resp.Resp {
	arr := respArr.(*resp.Array)
	if len(arr.Value) == 0 {
		return resp.NewErrorWithString("first argument must be a command")
	}
	cmdToken := arr.Value[0]
	if cmdToken.Type().String() != "<bulkbytes>" {
		return resp.NewErrorWithString("first argument must be a command")
	} else {
		var res resp.Resp
		cmd := strings.ToUpper(string(cmdToken.(*resp.BulkBytes).Value))
		switch cmd {
		default:
			res = resp.NewErrorWithString("first argument must be a command")
		case "ECHO":
			res = EchoCommand(arr.Value[1:])
		case "PING":
			res = PingCommand()
		case "SET":
			res = SetCommand(database, arr.Value[1:])
		case "GET":
			res = GetCommand(database, arr.Value[1:])
		}
		return res
	}
}

func EchoCommand(respArr []resp.Resp) resp.Resp {
	return respArr[0]
}

func PingCommand() resp.Resp {
	return resp.NewString("PONG")
}

func SetCommand(db map[string]resp.Resp, args []resp.Resp) resp.Resp {
	if len(args) != 2 && len(args) != 4 {
		return resp.NewErrorWithString("invalid arguments count for SET")
	}
	var key string
	switch v := args[0].(type) {
	default:
		return resp.NewErrorWithString("invalid key type for SET")
	case *resp.String:
		key = v.Value
	case *resp.BulkBytes:
		key = string(v.Value)
	}
	if len(args) == 4 {
		px, ok := args[2].(*resp.BulkBytes)
		if !ok || strings.ToUpper(string(px.Value)) != "PX" {
			return resp.NewErrorWithString("invalid option for SET")
		}
		ttlStr, ok := args[3].(*resp.BulkBytes)
		if !ok {
			return resp.NewErrorWithString("invalid ttl value for PX option")
		}
		ttl, _ := strconv.Atoi(string(ttlStr.Value))
		go func() {
			<-time.After(time.Duration(ttl) * time.Millisecond)
			delete(database, key)
		}()
	}
	db[key] = args[1]
	return resp.NewString("OK")
}

func GetCommand(db map[string]resp.Resp, args []resp.Resp) resp.Resp {
	if len(args) != 1 {
		return resp.NewErrorWithString("invalid arguments count for GET")
	}
	var key string
	switch v := args[0].(type) {
	default:
		return resp.NewErrorWithString("invalid key type for GET")
	case *resp.String:
		key = v.Value
	case *resp.BulkBytes:
		key = string(v.Value)
	}
	val, ok := db[key]
	if !ok {
		return resp.NewBulkBytes(nil)
	}
	return val
}
