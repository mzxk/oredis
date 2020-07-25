package oredis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type Oredis struct {
	pool   *redis.Pool
	script map[string]string
}

//Script 这个函数保存对应的名字和lua脚本到redis服务器，并缓存hash在本地
func (t *Oredis) Script(name, lua string) error {
	c := t.Get()
	defer c.Close()
	hash, err := redis.String(c.Do("script", "load", lua))
	if err != nil {
		return err
	}
	t.script[name] = hash
	return nil
}

//Eval 这个函数将调用本地缓存的对应名字的hash来调用lua脚本
func (t *Oredis) Eval(name string, keys []interface{}, args ...interface{}) (interface{}, error) {
	ags := []interface{}{
		t.script[name],
		len(keys),
	}
	ags = append(ags, keys...)
	ags = append(ags, args...)
	c := t.Get()
	defer c.Close()
	return c.Do("evalsha", ags...)
}
func (t *Oredis) Get() redis.Conn {
	conn := t.pool.Get()
	return conn
}
func (t *Oredis) GetDB(i int) (redis.Conn, error) {
	conn := t.pool.Get()
	_, err := conn.Do("SELECT", i)
	return conn, err
}

//New return redis.pool
func New(add, pwd string) *Oredis {
	p := &redis.Pool{
		MaxIdle:   10,
		MaxActive: 200,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", add,
				redis.DialConnectTimeout(1*time.Second),
				redis.DialReadTimeout(10*time.Second),
				redis.DialWriteTimeout(10*time.Second),
				redis.DialPassword(pwd),
			)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		// Use the TestOnBorrow function to check the health of an idle connection
		// before the connection is returned to the application.
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			// if time.Since(t) < time.Minute {
			// 	return nil
			// }
			_, err := c.Do("PING")
			return err
		},
		IdleTimeout: 300 * time.Second,
		// If Wait is true and the pool is at the MaxActive limit,
		// then Get() waits for a connection to be returned to the pool before returning
		Wait: true,
	}
	t := &Oredis{pool: p, script: map[string]string{}}
	return t
}
