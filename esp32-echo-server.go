package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	"google.golang.org/protobuf/proto"

	"github.com/Orison-Cloud/esp32-protobuf-shared/model"
)

type ProtoEvent struct {
	Data string `json:"data"`
}

// parse base64 encoded message to a protobuf model.Request
func decodeBase64(protoEvent ProtoEvent) (*model.Request, error) {

	var l int
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(protoEvent.Data)))
	l, err := base64.StdEncoding.Decode(decoded, []byte(protoEvent.Data))
	if err != nil {
		return nil, err
	}
	fmt.Printf("Decoded: %d bytes == %s\n", l, hex.EncodeToString(decoded))

	// wrap this with an envelope
	request := &model.Request{}
	err = proto.Unmarshal(decoded[:l], request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

var (
	topic    = aws.String("orison/broadcast/oem")
	endpoint = aws.String("https://a1h07teth8vuyq-ats.iot.us-east-1.amazonaws.com")
	suffix   = "-dawg"
)

func handleRequest(ctx context.Context, protoEvent ProtoEvent) (string, error) {

	fmt.Printf("Payload: %+v\n", protoEvent.Data)

	request, err := decodeBase64(protoEvent)
	if err != nil {
		return "error on decodeBase64", err
	}
	helloProto := request.GetHelloProto()

	region := os.Getenv("AWS_REGION")
	ses, err := session.NewSession()
	if err != nil {
		fmt.Printf("error")
		return "error-0", err
	}
	svc := iotdataplane.New(ses, &aws.Config{Region: aws.String(region), Endpoint: endpoint})

	// create response
	protoGoodbye := &model.GoodbyeProto{Name: helloProto.Name + suffix, Counter: helloProto.Counter + 1}
	response := &model.Response{Messagetype: model.MessageType_GOODBYEPROTO, Response: &model.Response_GoodbyeProto{protoGoodbye}}
	protoOut, err := proto.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	params := &iotdataplane.PublishInput{
		Topic:   topic,
		Payload: protoOut,
		Qos:     aws.Int64(0),
	}
	output, err := svc.Publish(params)
	if err != nil {
		return "error-2", err
	}
	fmt.Println(output)
	return "esp32-echo-server OK", nil
}

func main() {
	lambda.Start(handleRequest)
}
