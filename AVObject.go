package lean

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// "picture": {
//          "id": "543cbaede4b07db196f50f3c",
//          "__type": "File"
//  }

type LeanFile struct {
	File
	Id string
}

//{
//  "__type": "Bytes",
//  "base64": "5b6I5aSa55So5oi36KGo56S65b6I5Zac5qyi5oiR5Lus55qE5paH5qGj6aOO5qC877yM5oiR5Lus5bey5bCGIExlYW5DbG91ZCDmiYDmnInmlofmoaPnmoQgTWFya2Rvd24g5qC85byP55qE5rqQ56CB5byA5pS+5Ye65p2l44CC"
//}
type LeanByte struct {
	bytes []byte
}

type LeanTime struct {
	time.Time
}

//{
//  "__type": "Pointer",
//  "className": "Post",
//  "objectId": "55a39634e4b0ed48f0c1845c"
//}
type LeanPointer struct {
	ClassName string
	ObjectID  string
}

//relation can only modify by AddRelation or RemoveRelation.
//and the api return only this struct
type LeanRelation struct {
	ClassName string `json:'className'`
}

func (t *LeanPointer) typeName() string {
	return "Pointer"
}

func (t *LeanPointer) fillByMap(m map[string]string) error {
	t.ClassName = m["className"]
	t.ObjectID = m["objectId"]
	return nil
}

type AVObject interface {
	typeName() string
	fillByMap(map[string]string) error
}

func NewLeanByte(in []byte) LeanByte {
	return LeanByte{
		bytes: in,
	}
}

func (this *LeanByte) GetBytes() *[]byte {
	return &(this.bytes)
}

func (this *LeanByte) typeName() string {
	return "Bytes"
}

func (this *LeanByte) fillByMap(in map[string]string) error {
	if bytes, err := base64.StdEncoding.DecodeString(in["base64"]); nil != err {
		return err
	} else {
		this.bytes = bytes
		return nil
	}
}

func NewLeanTime(t time.Time) LeanTime {
	return LeanTime{t}
}

//func (t *LeanPointer) UnmarshalJSON(i []byte) error {
//	if converErr := bytes2AvObject(i, t); nil != converErr {
//		return converErr
//	} else {
//		return nil
//	}
//}

func (t LeanFile) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf(`{
		"__type": "File",
		"id":"%s"
	}`, t.Id)
	return []byte(str), nil
}

func (t *LeanFile) UnmarshalJSON(i []byte) error {
	//do your serializing here
	if converErr := bytes2AvObject(i, t); nil != converErr {
		return converErr
	} else {
		return nil
	}
}

func (t LeanPointer) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf(`{
		"__type": "Pointer",
		"className": "%s",
		"objectId":"%s"
	}`, t.ClassName, t.ObjectID)
	return []byte(str), nil
}

//convet json into TimeStamp, we take only the unix timestamp seconds
func (t *LeanByte) UnmarshalJSON(i []byte) error {
	//do your serializing here
	if converErr := bytes2AvObject(i, t); nil != converErr {
		return converErr
	} else {
		return nil
	}
}

func (t LeanByte) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf(`{
		"__type": "Bytes",
		"base64": "%x"
	}`, base64.StdEncoding.EncodeToString(t.bytes))
	return []byte(str), nil
}

//2015-07-14T02:31:50.100Z
//convet TimeStamp into json, we take only the unix timestamp seconds
func (t LeanTime) MarshalJSON() ([]byte, error) {
	stamp := t.UTC().Format("2006-01-02T15:04:05.000Z")
	str := fmt.Sprintf(`{
		"__type": "Date",
		"iso": "%s"
	}`, stamp)
	return []byte(str), nil
}

func (t *LeanFile) typeName() string {
	return "File"
}

func (t *LeanFile) fillByMap(m map[string]string) error {
	if id, exist := m["id"]; false == exist {
		return errors.New("wrong data type")
	} else {
		t.Id = id
	}
	return nil
}
func (t *LeanTime) typeName() string {
	return "Date"
}

func (t *LeanTime) fillByMap(m map[string]string) error {
	timeStr, exist := m["iso"]
	if false == exist {
		return errors.New("wrong data type")
	}
	if ret, err := str2Date(timeStr); nil != err {
		return err
	} else {
		t = ret
	}
	return nil
}

//convet json into TimeStamp, we take only the unix timestamp seconds
func (t *LeanTime) UnmarshalJSON(i []byte) error {
	//do your serializing here
	var timeStr string
	if err := json.Unmarshal(i, &timeStr); err != nil {
		if converErr := bytes2AvObject(i, t); nil != converErr {
			println("also can not convert :" + timeStr)
			return converErr
		} else {
			return nil
		}
		return err
	} else {
		if ret, err := str2Date(timeStr); err != nil {
		} else {
			t = ret
		}
		return nil
	}
}

func str2Date(str string) (*LeanTime, error) {
	if ret, err := time.Parse("2006-01-02T15:04:05.000Z", str); nil != err {
		return nil, err
	} else {
		return &LeanTime{ret}, nil
	}
}

func bytes2AvObject(str []byte, obj AVObject) error {
	object := map[string]string{}
	if err := json.Unmarshal([]byte(str), &object); nil != err {
		return err
	} else {
		if object["__type"] != obj.typeName() {
			return errors.New("type wrong! name is :" + obj.typeName())
		}
		if err := obj.fillByMap(object); nil != err {
			return err
		}
		return nil
	}
}
