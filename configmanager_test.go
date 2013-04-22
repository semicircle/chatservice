package chatservice

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

// for user + message
func TestIntegrate1(t *testing.T) {
	fmt.Println("TestIntegrate1: started")
	u := FACTORY.User()
	u.Generate().SetNick("nick").Save()
	defer u.Delete()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		msgs := u.WaitRecvMessage(time.After(5 * time.Second))
		if msgs == nil {
			t.Error("msgs == null")
		}

		if len(msgs) != 1 {
			t.Error("len(msgs) != 1, :" + strconv.Itoa(len(msgs)))
		}
		wg.Done()
	}()

	msg := FACTORY.Message()
	msg.SetAuthorId(u.Id()).SetText("text").Save()
	defer msg.Delete()
	u.WakeRecvPending(msg)

	wg.Wait()

}

// for thread + user + message
//poll on all User, then u1 post a message and see if the User can receive.
func TestIntegrate2(t *testing.T) {
	fmt.Println("TestIntegrate2: started")

	th := FACTORY.Thread()
	th.Save()
	defer th.Release()
	defer th.Delete()

	const USER_NUM = 10

	u := make([]User, USER_NUM)
	for i := 0; i < USER_NUM; i++ {
		u[i] = FACTORY.User()
		u[i].Generate().SetNick("nick").Save()
		defer u[i].Delete()
		th.Subscribe(u[i])
	}

	var wg sync.WaitGroup
	wg.Add(USER_NUM)
	for i := 0; i < USER_NUM; i++ {
		go func(inner User) {
			msgs := inner.WaitRecvMessage(time.After(5 * time.Second))
			if (msgs == nil) || (len(msgs) != 1) {
				t.Error("stage1")
			}
			wg.Done()
		}(u[i])
	}

	msg := FACTORY.Message()
	msg.SetAuthorId(u[0].Id()).SetText("text").Save()
	defer msg.Delete()
	th.Post(u[0], msg)
	wg.Wait()
}

//pressure here, everyone keeps buzz to see if any message lost.
func TestIntegrate2_2(t *testing.T) {
	fmt.Println("TestIntegrate2_2: started")

	th := FACTORY.Thread()
	th.Save()
	defer th.Release()
	defer th.Delete()

	const USER_NUM = 100
	const MAX_USER_BUZZ_TIMES = 100

	u := make([]User, USER_NUM)
	for i := 0; i < USER_NUM; i++ {
		u[i] = FACTORY.User()
		u[i].Generate().SetNick("nick").Save()
		defer u[i].Delete()
		th.Subscribe(u[i])
	}

	chMsgCnt := make(chan int)
	var factMsgCnt int

	r := rand.New(rand.NewSource(time.Now().Unix()))
	chTestStop := make(chan time.Time)
	//polling
	for i := 0; i < USER_NUM; i++ {
		go func(inner User) {
			var msgCnt int
			for {
				msgs := inner.WaitRecvMessage(chTestStop)
				if msgs == nil {
					chMsgCnt <- msgCnt
					return
				} else {
					//fmt.Println("WaitRecvMessage: msgs != null")
					msgCnt += len(msgs)
				}
			}
		}(u[i])
	}
	//buzz
	var wg sync.WaitGroup
	wg.Add(USER_NUM)
	for i := 0; i < USER_NUM; i++ {
		userbuzztimes := (r.Int() % MAX_USER_BUZZ_TIMES) + 2
		factMsgCnt += userbuzztimes
		go func(inner User) {
			deltatime := int(time.Second) / (userbuzztimes - 1)
			for i := 0; i < userbuzztimes; i++ {
				msg := FACTORY.Message()
				msg.Save()
				defer msg.Delete()
				th.Post(inner, msg)
				time.Sleep(time.Duration(deltatime))
			}
			wg.Done()
		}(u[i])
	}
	fmt.Println("TestIntegrate2_2: all go routines are launched")
	wg.Wait()
	fmt.Println("TestIntegrate2_2: buzz go routines are finished")
	time.Sleep(3 * time.Second)
	fmt.Println("TestIntegrate2_2: sending chan to stop receiver go routines")

	for i := 0; i < USER_NUM; i++ {
		chTestStop <- time.Now()
	}

	for i := 0; i < USER_NUM; i++ {
		msgCnt := <-chMsgCnt
		if msgCnt != factMsgCnt {
			t.Errorf("msgCnt: msgCnt:%d factMsgCnt:%d", msgCnt, factMsgCnt)
		}
	}
}

//more pressure and message filled with values, and values will be checked.
//fork from 2_2
func TestIntegrate2_3(t *testing.T) {
	fmt.Println("TestIntegrate2_3: started")

	th := FACTORY.Thread()
	th.Save()
	defer th.Release()
	defer th.Delete()

	const USER_NUM = 300
	const MAX_USER_BUZZ_TIMES = 5

	u := make([]User, USER_NUM)
	for i := 0; i < USER_NUM; i++ {
		u[i] = FACTORY.User()
		u[i].Generate().SetNick("nick").Save()
		defer u[i].Delete()
		th.Subscribe(u[i])
	}

	chMsgCnt := make(chan int)
	chMsgText := make(chan int)
	chMsgAuthorId := make(chan int)
	chMsgDate := make(chan int)
	chMsgDisplyAuthor := make(chan int)
	chMsgThreadId := make(chan int)
	var factMsgCnt, factMsgText, factMsgAuthorId, factMsgThreadId, factMsgDate, factMsgDisplyAuthor int

	r := rand.New(rand.NewSource(time.Now().Unix()))
	chTestStop := make(chan time.Time)
	//polling
	for i := 0; i < USER_NUM; i++ {
		go func(inner User) {
			var msgCnt, msgText, msgAuthorId, msgThreadId, msgDate, msgDisplyAuthor int
			for {
				msgs := inner.WaitRecvMessage(chTestStop)
				if msgs == nil {
					chMsgCnt <- msgCnt
					chMsgText <- msgText
					chMsgAuthorId <- msgAuthorId
					chMsgDate <- msgDate
					chMsgDisplyAuthor <- msgDisplyAuthor
					chMsgThreadId <- msgThreadId
					return
				} else {
					//fmt.Println("WaitRecvMessage: msgs != null")
					msgCnt += len(msgs)
					for _, msg := range msgs {
						msgTextInc, _ := strconv.Atoi(msg.Text())
						msgText += msgTextInc
						msgAuthorId += int(msg.AuthorId())
						msgThreadId += int(msg.ThreadId())
						msgDate += int(msg.Date().UnixNano()) % 1000
						msgDisplyAuthorInc, _ := strconv.Atoi(msg.DisplyAuthor())
						msgDisplyAuthor += msgDisplyAuthorInc
					}
				}
			}
		}(u[i])
	}
	//buzz
	var wg sync.WaitGroup
	wg.Add(USER_NUM)
	for i := 0; i < USER_NUM; i++ {
		userbuzztimes := (r.Int() % MAX_USER_BUZZ_TIMES) + 2
		factMsgCnt += userbuzztimes
		go func(inner User) {
			deltatime := int(time.Second) / (userbuzztimes - 1)
			for i := 0; i < userbuzztimes; i++ {
				msg := FACTORY.Message()
				// a little bit boring...
				msgtext := r.Int() % 1000
				msgauthorid := inner.Id()
				msgthreadid := th.Id()
				msgdate := time.Now()
				msgdisplyauthor := r.Int() % 1000
				factMsgText += msgtext
				factMsgAuthorId += int(msgauthorid)
				factMsgDate += int(msgdate.UnixNano()) % 1000
				factMsgDisplyAuthor += msgdisplyauthor
				factMsgThreadId += int(msgthreadid)
				msg.SetText(strconv.Itoa(msgtext))
				msg.SetAuthorId(msgauthorid)
				msg.SetThreadId(msgthreadid)
				msg.SetDate(msgdate)
				msg.SetDisplyAuthor(strconv.Itoa(msgdisplyauthor))
				msg.Save()

				defer msg.Delete()
				th.Post(inner, msg)
				time.Sleep(time.Duration(deltatime))
			}
			wg.Done()
		}(u[i])
	}
	fmt.Println("TestIntegrate2_3: all go routines are launched")
	wg.Wait()
	fmt.Println("TestIntegrate2_3: buzz go routines are finished")
	time.Sleep(2 * time.Second)

	fmt.Println("TestIntegrate2_3: sending chan to stop receiver go routines")
	for i := 0; i < USER_NUM; i++ {
		chTestStop <- time.Now()
	}

	for i := 0; i < USER_NUM; i++ {
		msgCnt := <-chMsgCnt
		msgText := <-chMsgText
		msgAuthorId := <-chMsgAuthorId
		msgDate := <-chMsgDate
		msgDisplyAuthor := <-chMsgDisplyAuthor
		msgThreadId := <-chMsgThreadId
		if msgCnt != factMsgCnt {
			t.Errorf("msgCnt: msgCnt:%d factMsgCnt:%d", msgCnt, factMsgCnt)
		}
		if msgText != factMsgText {
			t.Errorf("msgText: msgText:%d factMsgText:%d", msgText, factMsgText)
		}
		if msgAuthorId != factMsgAuthorId {
			t.Errorf("msgAuthorId: msgAuthorId:%d factMsgAuthorId:%d", msgAuthorId, factMsgAuthorId)
		}
		if msgDate != factMsgDate {
			t.Errorf("msgDate: msgDate:%d factMsgDate:%d", msgDate, factMsgDate)
		}
		if msgDisplyAuthor != factMsgDisplyAuthor {
			t.Errorf("msgDisplyAuthor: msgDisplyAuthor:%d factMsgDisplyAuthor:%d", msgDisplyAuthor, factMsgDisplyAuthor)
		}
		if msgThreadId != factMsgThreadId {
			t.Errorf("msgThreadId: msgThreadId:%d factMsgThreadId:%d", msgThreadId, factMsgThreadId)
		}
	}
}

//fork from 2_2
//bad guy test: one User subscribed but never call WaitRecvMessage, see if the whole system stuck.
func TestIntegrate2_4(t *testing.T) {
	fmt.Println("TestIntegrate2_4: started")

	th := FACTORY.Thread()
	th.Save()
	defer th.Release()
	defer th.Delete()

	const USER_NUM = 10
	const MAX_USER_BUZZ_TIMES = 100000

	u := make([]User, USER_NUM+1)     // +1 bad guy
	for i := 0; i < USER_NUM+1; i++ { // +1 bad guy
		u[i] = FACTORY.User()
		u[i].Generate().SetNick("nick").Save()
		defer u[i].Delete()
		th.Subscribe(u[i])
	}

	chMsgCnt := make(chan int)
	var factMsgCnt int

	r := rand.New(rand.NewSource(time.Now().Unix()))
	chTestStop := make(chan time.Time)
	//polling
	for i := 0; i < USER_NUM; i++ {
		go func(inner User) {
			var msgCnt int
			for {
				msgs := inner.WaitRecvMessage(chTestStop)
				if msgs == nil {
					chMsgCnt <- msgCnt
					return
				} else {
					//fmt.Println("WaitRecvMessage: msgs != null")
					msgCnt += len(msgs)
				}
			}
		}(u[i])
	}
	//buzz
	var wg sync.WaitGroup
	wg.Add(USER_NUM)
	for i := 0; i < USER_NUM; i++ {
		userbuzztimes := (r.Int() % MAX_USER_BUZZ_TIMES) + 2
		factMsgCnt += userbuzztimes
		go func(inner User) {
			deltatime := int(time.Second) / (userbuzztimes - 1)
			for i := 0; i < userbuzztimes; i++ {
				msg := FACTORY.Message()
				msg.Save()
				defer msg.Delete()
				th.Post(inner, msg)
				time.Sleep(time.Duration(deltatime))
			}
			wg.Done()
		}(u[i])
	}
	//fmt.Println("before wg.Wait()")
	wg.Wait()
	//fmt.Println("after wg.Wait()")
	time.Sleep(0 * time.Second)

	for i := 0; i < USER_NUM; i++ { // <-------- here also be changed.
		chTestStop <- time.Now()
	}

	for i := 0; i < USER_NUM; i++ { // <-------- here also be changed.
		msgCnt := <-chMsgCnt
		if msgCnt != factMsgCnt {
			return //success.
		}
	}

}
