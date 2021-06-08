# esp32-AWS-echo-server

## Simple echo-server that has a 
* IoT rule directing MQTT payload, listening on `orison/oem-10/oem/telemetry` channel
* decodes protobuf `model.Request` payload and extracts HelloProto message
* extracts name, counter, appends `-dawg` to `name` , and increments `counter` +1
* creates a `model.Response protobuf`, and loads a `model.GoodByeProto` struct inside, with info on last line above 
* rewrites out to `/orison/broadcast/oem` MQTT topic

# Build It
```
go build lambdaTest
```

# Upload it to AWS and Version it
```
./pushAWS.sh
```

# Good to Know

In AWS, an IoT rule was configured to send directly to lambda, but must base64 encode the binary protobuf, else breaks the lambda rule. 
Several attempts at bypassing this were made for future reference

* a new lambdaHandler was implemented that took []byte argument and overrode Invoke. 
  This worked great as long as you read direct from MQTT topic, but then won't fire the lambda serverless automatically
* setup a rule to have MQTT topic directly to kinesis stream. 
  This also worked, but the abundance of meta-data surrounding the []byte payload was HUGE, **and** AWS base64 encoded it automatically!!

I just had the IoT rule base64 encode it and then decoded it inside the lambda. DONE