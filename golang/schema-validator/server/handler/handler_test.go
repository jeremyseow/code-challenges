package handler

import (
	"schema-validator/event"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func BenchmarkUnmarshalProtoJson(b *testing.B) {
	jsonPayload := `{
		"device_model":"xiaomi",
		"os_type":"android",
		"events":[
			{
				"event_name":"event1",
				"event_timestamp":1234567890,
				"params_struct":{
					"param1":{
						"string_value":"hello"
					},
					"param2":{
						"int_value":123
					},
					"param3":{
						"bool_value":true
					},
					"param4":{
						"double_value":123.45
					},
					"param5":{
						"string_array_value": {
							"string_values":["hello", "world"]
						}
					},
					"param6":{
						"int_array_value": {
							"int_values":[1, 2, 3]
						}
					},
					"param7":{
						"bool_array_value": {
							"bool_values":[true, false, true]
						}
					},
					"param8":{
						"double_array_value": {
							"double_values":[1.1, 2.2, 3.3]
						}
					}
				}
			}
		]
	}`

	b.ResetTimer()

	// BenchmarkUnmarshalProtoJson-16    	   65067	     16743 ns/op	    4536 B/op	     144 allocs/op
	for i := 0; i < b.N; i++ {
		var clientRequest event.ClientRequest
		err := protojson.Unmarshal([]byte(jsonPayload), &clientRequest)
		if err != nil {
			b.Fatal(err)
		}
		// fmt.Println(&clientRequest)
	}

}

func BenchmarkUnmarshalJson(b *testing.B) {
	jsonPayload := `{
		"device_model": "xiaomi",
		"os_type": "android",
		"events": [
			{
				"event_name": "event1",
				"event_timestamp": 1234567890,
				"params_struct":{
					"param1":{
						"string_value":"hello"
					},
					"param2":{
						"int_value":123
					},
					"param3":{
						"bool_value":true
					},
					"param4":{
						"double_value":123.45
					},
					"param5":{
						"string_array_value": {
							"string_values":["hello", "world"]
						}
					},
					"param6":{
						"int_array_value": {
							"int_values":[1, 2, 3]
						}
					},
					"param7":{
						"bool_array_value": {
							"bool_values":[true, false, true]
						}
					},
					"param8":{
						"double_array_value": {
							"double_values":[1.1, 2.2, 3.3]
						}
					}
				}
			}
		]
	}`

	b.ResetTimer()

	// BenchmarkUnmarshalJson-16    	  307692	      3805 ns/op	    2865 B/op	      64 allocs/op
	for i := 0; i < b.N; i++ {
		var clientRequest event.ClientRequest
		err := jsoniter.Unmarshal([]byte(jsonPayload), &clientRequest)
		if err != nil {
			b.Fatal(err)
		}
		// fmt.Println(&clientRequest)
	}

}

func BenchmarkUnmarshalProtoOneOf(b *testing.B) {
	samplePayload := event.ClientRequest{
		DeviceModel: "xiaomi",
		OsType:      "android",
		Events: []*event.Event{
			{
				EventName:      "event1",
				EventTimestamp: 1234567890,
				ParamsOneof: map[string]*event.DataValueOneOf{
					"param1": {
						Kind: &event.DataValueOneOf_StringValue{StringValue: "hello"},
					},
					"param2": {
						Kind: &event.DataValueOneOf_IntValue{IntValue: 123},
					},
					"param3": {
						Kind: &event.DataValueOneOf_BoolValue{BoolValue: true},
					},
					"param4": {
						Kind: &event.DataValueOneOf_DoubleValue{DoubleValue: 123.45},
					},
					"param5": {
						Kind: &event.DataValueOneOf_StringArrayValue{StringArrayValue: &event.StringArray{StringValues: []string{"hello", "world"}}},
					},
					"param6": {
						Kind: &event.DataValueOneOf_IntArrayValue{IntArrayValue: &event.IntArray{IntValues: []int64{1, 2, 3}}},
					},
					"param7": {
						Kind: &event.DataValueOneOf_BoolArrayValue{BoolArrayValue: &event.BoolArray{BoolValues: []bool{true, false, true}}},
					},
					"param8": {
						Kind: &event.DataValueOneOf_DoubleArrayValue{DoubleArrayValue: &event.DoubleArray{DoubleValues: []float64{1.1, 2.2, 3.3}}},
					},
				},
			},
			{
				EventName:      "event2",
				EventTimestamp: 1234567890,
				ParamsOneof: map[string]*event.DataValueOneOf{
					"param1": {
						Kind: &event.DataValueOneOf_StringValue{StringValue: "hello"},
					},
					"param2": {
						Kind: &event.DataValueOneOf_IntValue{IntValue: 123},
					},
					"param3": {
						Kind: &event.DataValueOneOf_BoolValue{BoolValue: true},
					},
					"param4": {
						Kind: &event.DataValueOneOf_DoubleValue{DoubleValue: 123.45},
					},
					"param5": {
						Kind: &event.DataValueOneOf_StringArrayValue{StringArrayValue: &event.StringArray{StringValues: []string{"hello", "world"}}},
					},
					"param6": {
						Kind: &event.DataValueOneOf_IntArrayValue{IntArrayValue: &event.IntArray{IntValues: []int64{1, 2, 3}}},
					},
					"param7": {
						Kind: &event.DataValueOneOf_BoolArrayValue{BoolArrayValue: &event.BoolArray{BoolValues: []bool{true, false, true}}},
					},
					"param8": {
						Kind: &event.DataValueOneOf_DoubleArrayValue{DoubleArrayValue: &event.DoubleArray{DoubleValues: []float64{1.1, 2.2, 3.3}}},
					},
				},
			},
		},
	}

	protoPayload, err := proto.Marshal(&samplePayload)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	// BenchmarkUnmarshalProtoOneOf-16    	  132045	      8860 ns/op	    3352 B/op	     117 allocs/op
	for i := 0; i < b.N; i++ {
		var clientRequest event.ClientRequest
		err := proto.Unmarshal(protoPayload, &clientRequest)
		if err != nil {
			b.Fatal(err)
		}
		// fmt.Println(&clientRequest)
	}
}

func BenchmarkUnmarshalProtoStruct(b *testing.B) {
	samplePayload := event.ClientRequest{
		DeviceModel: "xiaomi",
		OsType:      "android",
		Events: []*event.Event{
			{
				EventName:      "event1",
				EventTimestamp: 1234567890,
				ParamsStruct: map[string]*event.DataValueStruct{
					"param1": {
						StringValue: proto.String("hello"),
					},
					"param2": {
						IntValue: proto.Int64(123),
					},
					"param3": {
						BoolValue: proto.Bool(true),
					},
					"param4": {
						DoubleValue: proto.Float64(123.45),
					},
					"param5": {
						StringArrayValue: &event.StringArray{StringValues: []string{"hello", "world"}},
					},
					"param6": {
						IntArrayValue: &event.IntArray{IntValues: []int64{1, 2, 3}},
					},
					"param7": {
						BoolArrayValue: &event.BoolArray{BoolValues: []bool{true, false, true}},
					},
					"param8": {
						DoubleArrayValue: &event.DoubleArray{DoubleValues: []float64{1.1, 2.2, 3.3}},
					},
				},
			},
			{
				EventName:      "event2",
				EventTimestamp: 1234567890,
				ParamsStruct: map[string]*event.DataValueStruct{
					"param1": {
						StringValue: proto.String("hello"),
					},
					"param2": {
						IntValue: proto.Int64(123),
					},
					"param3": {
						BoolValue: proto.Bool(true),
					},
					"param4": {
						DoubleValue: proto.Float64(123.45),
					},
					"param5": {
						StringArrayValue: &event.StringArray{StringValues: []string{"hello", "world"}},
					},
					"param6": {
						IntArrayValue: &event.IntArray{IntValues: []int64{1, 2, 3}},
					},
					"param7": {
						BoolArrayValue: &event.BoolArray{BoolValues: []bool{true, false, true}},
					},
					"param8": {
						DoubleArrayValue: &event.DoubleArray{DoubleValues: []float64{1.1, 2.2, 3.3}},
					},
				},
			},
		},
	}

	protoPayload, err := proto.Marshal(&samplePayload)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	// BenchmarkUnmarshalProtoStruct-16    	  165014	      7145 ns/op	    4056 B/op	     109 allocs/op
	for i := 0; i < b.N; i++ {
		var clientRequest event.ClientRequest
		err := proto.Unmarshal(protoPayload, &clientRequest)
		if err != nil {
			b.Fatal(err)
		}
		// fmt.Println(&clientRequest)
	}
}

func BenchmarkUnmarshalGoAny(b *testing.B) {
	jsonPayload := `{
		"device_model":"xiaomi",
		"os_type":"android",
		"events":[
			{
				"event_name":"event1",
				"event_timestamp":12314567890,
				"params_any":{
					"param1": "hello",
					"param2": 123,
					"param3": true,
					"param4": 123.45,
					"param5": ["hello", "world"],
					"param6": [1, 2, 3],
					"param7": [true, false, true],
					"param8": [1.1, 2.2, 3.3]
				}
			},
			{
				"event_name":"event2",
				"event_timestamp":12314567890,
				"params_any":{
					"param1": "hello",
					"param2": 123,
					"param3": true,
					"param4": 123.45,
					"param5": ["hello", "world"],
					"param6": [1, 2, 3],
					"param7": [true, false, true],
					"param8": [1.1, 2.2, 3.3]
				}
			}
		]
	}`

	b.ResetTimer()

	// BenchmarkUnmarshalGoAny-16    	  131049	      8652 ns/op	    5066 B/op	     167 allocs/op
	for i := 0; i < b.N; i++ {
		var clientRequestAny event.ClientRequestAny
		err := jsoniter.Unmarshal([]byte(jsonPayload), &clientRequestAny)
		if err != nil {
			b.Fatal(err)
		}

		var clientRequest event.ClientRequest
		clientRequest.DeviceModel = clientRequestAny.DeviceModel
		clientRequest.OsType = clientRequestAny.OsType
		clientRequest.Events = make([]*event.Event, len(clientRequestAny.Events))
		for idx, evt := range clientRequestAny.Events {
			clientRequest.Events[idx] = &event.Event{
				EventName:      evt.EventName,
				EventTimestamp: evt.EventTimestamp,
				ParamsOneof:    make(map[string]*event.DataValueOneOf, len(evt.ParamsAny)),
			}
			for paramKey, paramValue := range evt.ParamsAny {
				switch castedValue := paramValue.(type) {
				case string:
					clientRequest.Events[idx].ParamsOneof[paramKey] = &event.DataValueOneOf{
						Kind: &event.DataValueOneOf_StringValue{StringValue: castedValue},
					}
				case int64:
					clientRequest.Events[idx].ParamsOneof[paramKey] = &event.DataValueOneOf{
						Kind: &event.DataValueOneOf_IntValue{IntValue: castedValue},
					}
				case bool:
					clientRequest.Events[idx].ParamsOneof[paramKey] = &event.DataValueOneOf{
						Kind: &event.DataValueOneOf_BoolValue{BoolValue: castedValue},
					}
				case float64:
					clientRequest.Events[idx].ParamsOneof[paramKey] = &event.DataValueOneOf{
						Kind: &event.DataValueOneOf_DoubleValue{DoubleValue: castedValue},
					}
				case []string:
					clientRequest.Events[idx].ParamsOneof[paramKey] = &event.DataValueOneOf{
						Kind: &event.DataValueOneOf_StringArrayValue{StringArrayValue: &event.StringArray{StringValues: castedValue}},
					}
				case []int64:
					clientRequest.Events[idx].ParamsOneof[paramKey] = &event.DataValueOneOf{
						Kind: &event.DataValueOneOf_IntArrayValue{IntArrayValue: &event.IntArray{IntValues: castedValue}},
					}
				case []bool:
					clientRequest.Events[idx].ParamsOneof[paramKey] = &event.DataValueOneOf{
						Kind: &event.DataValueOneOf_BoolArrayValue{BoolArrayValue: &event.BoolArray{BoolValues: castedValue}},
					}
				case []float64:
					clientRequest.Events[idx].ParamsOneof[paramKey] = &event.DataValueOneOf{
						Kind: &event.DataValueOneOf_DoubleArrayValue{DoubleArrayValue: &event.DoubleArray{DoubleValues: castedValue}},
					}
				default:
					// fmt.Printf("Unsupported type: %T\n", paramValue)
				}
			}
		}

		// fmt.Println(&clientRequest)
	}

}
