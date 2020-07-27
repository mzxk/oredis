package oredis

import "github.com/gomodule/redigo/redis"

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
