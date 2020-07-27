package oredis

import (
	"errors"
)

//Do 封装了conn.Do
func (t *Oredis) Do(commandName string, args ...interface{}) (interface{}, error) {
	conn := t.pool.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
}

//DoDB 封装了conn.Do
func (t *Oredis) DoDB(db int, commandName string, args ...interface{}) (interface{}, error) {
	conn := t.pool.Get()
	defer conn.Close()
	_ = conn.Send("multi")
	_ = conn.Send("select", db)
	_ = conn.Send(commandName, args...)
	return conn.Do("exec")
}

//Multi 减少调用代码，已经自行添加了multi和exec，只要在输入里增加命令就可以，注意，每一个数组的第一个参数必须是字符串
func (t *Oredis) Multi(command [][]interface{}) (interface{}, error) {
	conn := t.pool.Get()
	defer conn.Close()
	_ = conn.Send("multi")
	for _, v := range command {
		commandName, ok := v[0].(string)
		if !ok {
			return nil, errors.New("commandNameMustString")
		}
		_ = conn.Send(commandName, v[1:]...)
	}
	return conn.Do("exec")
}
