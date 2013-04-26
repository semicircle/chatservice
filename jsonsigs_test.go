package chatservice

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestJsonSig(t *testing.T) {

	object := new(message)

	object.Date = time.Now()
	object.Id = idType(1234)
	object.Author.Id = idType(4321)
	object.Author.Nick = "nick"
	object.Author.PrivateKey = keyType(0)
	object.Level = MessageCode(1)
	object.Text = "UTF8测试 , 中文可不可以?"

	var jsondata []byte
	var err error
	if jsondata, err = json.Marshal(object); err != nil {
		t.Error(err)
	} else {
		fmt.Println(string(jsondata))
	}

	var v message
	if err = json.Unmarshal(jsondata, &v); err != nil {
		t.Error(err)
	} else {
		fmt.Println(v)
	}

	if *object != v {
		t.Error("not equal??")
	}

}
