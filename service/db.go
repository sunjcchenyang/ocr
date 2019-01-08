package service

import (
	"go.uber.org/zap"
	mgo "gopkg.in/mgo.v2"
)

// Newmgo connect the mgdb
func Newmgo() *mgo.Session {
	Addr := "127.0.0.1"
	session, err := mgo.Dial(Addr) //连接数据库
	if err != nil {
		Logger.Error("connect the db is error", zap.String("the message is", err.Error()))
		return nil
	}
	session.SetPoolLimit(100)
	Logger.Info("Connect", zap.String("[Connect the mgdb]", "Success connect the mgdb"))
	return session
}
