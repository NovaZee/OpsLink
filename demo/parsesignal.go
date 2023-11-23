package main

import (
	"encoding/hex"
	"fmt"
	"github.com/denovo/permission/protoc/signal"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
	"time"
)

func main() {
	ping()
	renewal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2ODYzMDIyMTYyNTE0MDQwMjU3LCJ1c2VyX25hbWUiOiJ6ZjYiLCJleHAiOjE3MDA3MTU4ODgsImlzcyI6IjM4Mzg0LVNlYXJjaEVuZ2luZSJ9.rho2sYEHyQNOYSWRZ2MSe_I8CclUNklEKpx5WZYVR-8")
	//token, _ := util.ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2ODYzMDIyMTYyNTE0MDQwMjU3LCJ1c2VyX25hbWUiOiJ6ZjYiLCJleHAiOjE3MDA2NDc1NTgsImlzcyI6IjM4Mzg0LVNlYXJjaEVuZ2luZSJ9.MF21fwd-0VtO5zRP-VrCo2kk3YrE9OxPURrV1YaFsMw")
	//println(token.ExpiresAt)
}
func renewal(token string) {
	request := signal.SignalRequest{
		Message: &signal.SignalRequest_Renewal{
			Renewal: &signal.RefreshToken{
				Token: token,
			},
		},
	}

	data, err := proto.Marshal(&request)
	if err != nil {
		fmt.Println("序列化出错:", err)
		return
	}

	//转base64 方便 postman调用调试
	//encodedData := base64.StdEncoding.EncodeToString(data)
	//转成16进制
	hexString := hex.EncodeToString(data)
	// 打印二进制数据
	fmt.Printf("二进制数据: \n", hexString)
	marshal := protojson.Format(&request)
	fmt.Println("json数据", marshal)
}

func ping() {
	request := signal.SignalRequest{
		Message: &signal.SignalRequest_Ping{Ping: time.Now().Unix()},
	}

	data, err := proto.Marshal(&request)
	if err != nil {
		fmt.Println("序列化出错:", err)
		return
	}

	//转base64 方便 postman调用调试
	//encodedData := base64.StdEncoding.EncodeToString(data)
	//转成6进制
	hexString := hex.EncodeToString(data)
	// 打印二进制数据
	fmt.Printf("二进制数据: %v\n", hexString)
	marshal := protojson.Format(&request)
	fmt.Println("json数据", marshal)
}

func refreshToken(id int64, token string) {
	request := signal.SignalRequest{
		Message: &signal.SignalRequest_Renewal{
			Renewal: &signal.RefreshToken{
				Token: token,
			},
		},
	}

	data, err := proto.Marshal(&request)
	if err != nil {
		fmt.Println("序列化出错:", err)
		return
	}

	//转base64 方便 postman调用调试
	//encodedData := base64.StdEncoding.EncodeToString(data)
	//转成6进制
	hexString := hex.EncodeToString(data)
	// 打印二进制数据
	fmt.Printf("二进制数据: %v\n", hexString)
	marshal := protojson.Format(&request)
	fmt.Println("json数据", marshal)
}

func parseRefreshToken(hexString string) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		fmt.Println("解码出错:", err)
		return
	}
	response := &signal.SignalResponse{}
	err = proto.Unmarshal(bytes, response)
	fmt.Println("解析后结构体:", response)
}
