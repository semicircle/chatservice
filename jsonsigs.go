package chatservice

import (
	"time"
)

// This file defines the json formatted message signals between client and server.

// Description:
// 1. sigHead is included by all sigs.
// 2. user/message/thread/group are inclued by other sigs.
// 3. All response sigs contains a 'Req' which is the corrsponding request sig.
// 4. All response sigs has a Result field, 
//    when the request is successed, Result is 'successed',
//    otherwise, it repersent the fail reason

// all the json blobs have 'SigType'
type sigHead struct {
	SigType string
}

//inner type user
type user struct {
	Id         idType  `json:",string,omitempty"`
	Nick       string  `json:",omitempty"`
	PrivateKey keyType `json:"PrivateKey,string,omitempty"`
}

//inner type msg
type message struct {
	Id     idType `json:",string"`
	Author user   `json:",omitempty"` //Not make sense when post request
	Thread thread `json:",omitempty"`
	Text   string
	Level  MessageCode `json:",int"`
	Date   time.Time
}

//inner type thread
type thread struct {
	Id    idType `json:",string"`
	Title string `json:",omitempty"`
}

//inner type group
type group struct {
	Id    idType `json:",string"`
	Title string `json:",omitempty"`
}

//
// Requests about User
//

//MsgType: userRegisterReq
type userRegisterReq struct {
	SigType string
	User    user
}

//MsgType: userRegisterRsp
type userRegisterRsp struct {
	SigType string
	Result  string
	User    user
	Req     userRegisterReq
}

//
// Requests about Thread
//

//MsgType: threadSubscribeReq
type threadSubscribeReq struct {
	SigType string
	Thread  thread
}

//MsgType: threadSubscribeReq
type threadSubscribeRsp struct {
	SigType string
	Result  string
	Req     threadSubscribeReq
}

//MsgType: threadPostMessageReq
type threadPostMessageReq struct {
	SigType string
	Thread  thread
	Message message
}

//MsgType: threadPostMessageRsp 
type threadPostMessageRsp struct {
	SigType string
	Result  string
	Req     threadPostMessageReq
}

//MsgType: threadCreateReq 
type threadCreateReq struct {
	SigType string
	Thread  thread
	Group   group
}

//MsgType: threadCreateRsp
type threadCreateRsp struct {
	SigType string
	Result  string
	Req     threadCreateReq
}

//MsgType: threadAchieveReq
type threadAchieveReq struct {
	SigType string
	Thread  thread
}

//MsgType: threadAchieveRsp
type threadAchieveRsp struct {
	SigType string
	Result  string
	Achieve []message
	Req     threadAchieveReq
}

//
// Server sigs
//

//MsgType: serverMessages
type serverMessages struct {
	SigType  string
	Messages []message
}
