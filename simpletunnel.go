package chatservice

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// A simple tunnel which meet the requirements.
type SimpleTunnel struct {
	onlineUsers map[idType]*onlineUserStru
	receiveChan chan recvStru
}

type onlineUserStru struct {
	jsonblob chan []byte
	u        User
}

type recvStru struct {
	jsonblob []byte
	u        User
}

func NewSimpleTunnel(baseUrl string) (ret *SimpleTunnel) {
	ret = &SimpleTunnel{make(map[idType]*onlineUserStru), make(chan recvStru, 10)}
	http.Handle(baseUrl+"/", ret)
	http.ListenAndServe(":8090", nil)
	return
}

func (st *SimpleTunnel) makeUserOnline(userid idType) (u *onlineUserStru) {
	var ok bool
	if u, ok = st.onlineUsers[userid]; !ok {
		user := FACTORY.User()
		if err := user.Load(userid); err != nil {
			return nil
		}
		u = &onlineUserStru{make(chan []byte, 10), user}
		st.onlineUsers[userid] = u
	}
	return
}

func (st *SimpleTunnel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ok bool
	var ou *onlineUserStru
	var v []string
	if strings.HasSuffix(r.URL.Path, "request") {
		r.ParseForm()
		var useridstr, jsonblob string

		if v, ok = r.Form["userid"]; !ok {
			io.WriteString(w, "no userid?")
			return
		}
		useridstr = strings.Join(v, "")
		if v, ok = r.Form["data"]; !ok {
			io.WriteString(w, "no data?")
			return
		}
		jsonblob = strings.Join(v, "")

		userid, _ := strconv.Atoi(useridstr)
		if ou = st.makeUserOnline(idType(userid)); ou == nil {
			io.WriteString(w, "invalid userid?")
			return
		}
		st.receiveChan <- recvStru{[]byte(jsonblob), ou.u}

	} else if strings.HasSuffix(r.URL.Path, "longpollingjsonp") {
		r.ParseForm()
		var callback, useridstr string

		if v, ok = r.Form["callback"]; !ok {
			io.WriteString(w, "no callback? emmm...only jsonp supported.")
			return
		}
		callback = strings.Join(v, "")
		if v, ok = r.Form["userid"]; !ok {
			io.WriteString(w, "no userid?")
			return
		}
		useridstr = strings.Join(v, "")

		userid, _ := strconv.Atoi(useridstr)
		if ou = st.makeUserOnline(idType(userid)); ou == nil {
			io.WriteString(w, "invalid userid?")
			return
		}

		select {
		case <-time.After(5 * time.Minute):
			return
		case item := <-ou.jsonblob:
			w.Header().Add("Content-Type", "text/javascript")
			io.WriteString(w, callback+"('")
			io.WriteString(w, string(item))
			io.WriteString(w, "');")
			return
		}
	}
}

func (st *SimpleTunnel) RecvFromClient(timeout chan time.Time) (User, []byte, error) {
	select {
	case <-timeout:
		return nil, nil, nil
	case item := <-st.receiveChan:
		return item.u, item.jsonblob, nil
	}
	return nil, nil, nil
}

func (st *SimpleTunnel) SendToClient(u User, jsonblob []byte, timeout <-chan time.Time) error {
	// is user exist(online)?
	if ou, ok := st.onlineUsers[u.Id()]; !ok {
		select {
		case <-timeout:
			return nil
		case ou.jsonblob <- jsonblob:
			return nil
		}
	} else {
		return errors.New("User is offlline")
	}
	return nil
}
