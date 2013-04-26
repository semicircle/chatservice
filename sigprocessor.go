package chatservice

import (
	"encoding/json"
	"time"
)

type SigProcessor struct {
}

func NewSigProcessor() *SigProcessor {
	return &SigProcessor{}
}

func (h *SigProcessor) Handle(u User, jsondata []byte) ([]byte, error) {
	var head sigHead
	if err := json.Unmarshal(jsondata, &head); err != nil {
		return nil, err
	}

	switch head.SigType {
	case "userRegisterReq":
		return h.userRegisterReqHandleFunc(u, jsondata)
	case "threadSubscribeReq":
		return h.threadSubscribeReqHandleFunc(u, jsondata)
	case "threadPostMessageReq":
		return h.threadPostMessageReqHandleFunc(u, jsondata)
	}

	return nil, nil
}

func (h *SigProcessor) WaitServerSigs(u User, timeout chan time.Time) ([]byte, error) {
	var rsp serverMessages
	rsp.SigType = "serverMessages"
	if msgs := u.WaitRecvMessage(timeout); msgs != nil {
		rsp.Messages = make([]message, len(msgs))
		for i := 0; i < len(msgs); i++ {
			rsp.Messages[i].Date = msgs[i].Date()
			rsp.Messages[i].Text = msgs[i].Text()
			rsp.Messages[i].Level = msgs[i].Type()
			rsp.Messages[i].Thread.Id = msgs[i].ThreadId()
			rsp.Messages[i].Author.Id = msgs[i].AuthorId()
			rsp.Messages[i].Author.Nick = msgs[i].DisplyAuthor()
		}
	}

	if ret, err := json.Marshal(rsp); err == nil {
		return ret, nil
	} else {
		return nil, err
	}
	return nil, nil
}

func (h *SigProcessor) userRegisterReqHandleFunc(u User, jsondata []byte) ([]byte, error) {
	var rsp userRegisterRsp
	rsp.SigType = "userRegisterRsp"
	if err := json.Unmarshal(jsondata, &rsp.Req); err != nil {
		return nil, err
	}

	newuser := FACTORY.User()
	newuser.Generate()
	newuser.SetNick(rsp.Req.User.Nick)
	if err := newuser.Save(); err != nil {
		rsp.Result = "failed: " + err.Error()
	} else {
		rsp.Result = "successed"
		rsp.User.Id = newuser.Id()
		rsp.User.Nick = newuser.Nick()
		rsp.User.PrivateKey = newuser.PrivateKey()
	}

	if ret, err := json.Marshal(rsp); err == nil {
		return ret, nil
	} else {
		return nil, err
	}
	return nil, nil
}

func (h *SigProcessor) threadSubscribeReqHandleFunc(u User, jsondata []byte) ([]byte, error) {
	var rsp threadSubscribeRsp
	rsp.SigType = "threadSubscribeRsp"
	if err := json.Unmarshal(jsondata, &rsp.Req); err != nil {
		return nil, err
	}

	thread := FACTORY.Thread()
	if err := thread.Load(rsp.Req.Thread.Id); err != nil {
		rsp.Result = "failed: " + err.Error()
	} else if err := thread.Subscribe(u); err != nil {
		rsp.Result = "failed: " + err.Error()
	} else {
		rsp.Result = "successed"
	}

	if ret, err := json.Marshal(rsp); err == nil {
		return ret, nil
	} else {
		return nil, err
	}
	return nil, nil
}

func (h *SigProcessor) threadPostMessageReqHandleFunc(u User, jsondata []byte) ([]byte, error) {
	var rsp threadPostMessageRsp
	rsp.SigType = "threadPostMessageRsp"
	if err := json.Unmarshal(jsondata, &rsp.Req); err != nil {
		return nil, err
	}

	thread := FACTORY.Thread()
	message := FACTORY.Message()
	if err := thread.Load(rsp.Req.Thread.Id); err != nil {
		rsp.Result = "failed thread.Load: " + err.Error()
		goto enough
	}
	message.SetDate(rsp.Req.Message.Date)
	message.SetText(rsp.Req.Message.Text)
	message.SetType(MESSAGE_TYPE_NORMAL)
	message.SetAuthorId(u.Id())
	message.SetDisplyAuthor(u.Nick())
	if err := message.Save(); err != nil {
		rsp.Result = "failed message.Save(): " + err.Error()
		goto enough
	}

	if err := thread.Post(u, message); err != nil {
		rsp.Result = "failed thread.Post: " + err.Error()
		goto enough
	}

	rsp.Result = "successed"

enough:
	if ret, err := json.Marshal(rsp); err == nil {
		return ret, nil
	} else {
		return nil, err
	}
	return nil, nil
}
