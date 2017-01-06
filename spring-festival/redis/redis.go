package redis

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"time"

	"github.com/astaxie/beego"

	//	"github.com/astaxie/beego/logs"
	"gopkg.in/redis.v4"
)

var client *redis.Client

func init() {
	// redis初始化
	poolsize := beego.AppConfig.DefaultInt("redis.PoolSize", 5)
	host := beego.AppConfig.Strings("redis.SentinelHost")
	name := beego.AppConfig.String("redis.SentinelName")

	beego.Info("======== begin init redis ========")
	beego.Info("SentinelHost:", host)
	beego.Info("SentinelName:", name)
	beego.Info("PoolSize:", poolsize)

	client = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    name,
		SentinelAddrs: host,
		PoolSize:      poolsize,
	})

	statusCmd := client.Ping()
	if err := statusCmd.Err(); err != nil {
		panic("init redis failed:" + err.Error())
	}
	beego.Info("======== finish redis init ========")

}

func Set(key string, object interface{}, expire int) bool {
	result, err := json.Marshal(object)
	if err != nil {
		beego.Error("redis set,json Marshal error:", err)
		return false
	}
	//	beego.Info("Marshal json:" + string(result))
	expires := time.Duration(expire) * time.Second
	err = client.Set(key, result, expires).Err()
	if err != nil {
		beego.Error("redis set error:", err)
		return false
	}
	return true
}

func Get(key string, value interface{}) bool {
	res, err := client.Get(key).Result()
	if err != nil {
		beego.Error("redis get error, key = ", key, "error:", err)
		return false
	}
	//	beego.Info("get redis result:" + res)
	err = json.Unmarshal([]byte(res), &value)
	if err != nil {
		beego.Error("redis get err,json Unmarshal error, key = ", key, ",error:", err)
		return false
	}
	return true
}

//func LPush(key string, object interface{}, expire time.Duration) int64 {
//	var value interface{}
//	switch object.(type) {
//	case utils.FeedInfo:
//		result, _ := json.Marshal(object)
//		value = string(result)
//	default:
//		value = object
//	}
//	count, err := client.LPush(key, value).Result()
//	if err != nil {
//		beego.Error("lpush error, key = %s, error: %v", key, err)
//		return 0
//	}
//	// 失效期
//	client.Expire(key, expire)
//	return count
//}

//func RPush(key string, object interface{}, expire time.Duration) int64 {
//	var value interface{}
//	switch object.(type) {
//	case utils.FeedInfo:
//		result, _ := json.Marshal(object)
//		value = string(result)
//	default:
//		value = object
//	}
//	count, err := client.RPush(key, value).Result()
//	if err != nil {
//		beego.Error("rpush error, key = %s, error: %v", key, err)
//		return 0
//	}
//	client.Expire(key, expire)
//	return count
//}

//func LRange(key string, start, stop int64) []string {
//	data, err := client.LRange(key, start, stop).Result()
//	if err != nil {
//		beego.Error("lrange error, key = %s, start = %d, stop = %d, error: %v", key, start, stop, err)
//		return nil
//	}
//	return data
//}

//func LTrim(key string, start, stop int64) bool {
//	err := client.LTrim(key, start, stop).Err()
//	if err != nil {
//		beego.Error("ltrim error, key = %s, start = %d, stop = %d, error: %v", key, start, stop, err)
//		return false
//	}
//	return true
//}

//func LRem(key string, count int64, object interface{}) bool {
//	var value interface{}
//	switch object.(type) {
//	case utils.FeedInfo:
//		result, _ := json.Marshal(object)
//		value = string(result)
//	default:
//		value = object
//	}
//	err := client.LRem(key, count, value).Err()
//	if err != nil {
//		beego.Error("lrem error, key = %s, error: %v", key, err)
//		return false
//	}
//	return true
//}

//func LLen(key string) int64 {
//	count, err := client.LLen(key).Result()
//	if err != nil {
//		beego.Error("llen error, key = %s, error: %v", key, err)
//		return 0
//	}
//	return count
//}

func Incr(key string, expire int) int64 {
	value, err := client.Incr(key).Result()
	if err != nil {
		beego.Error("incr error, key = ", key)
		return 0
	}
	expires := time.Duration(expire) * time.Second
	client.Expire(key, expires)
	return value
}

func IncrBy(key string, object int64, expire int) int64 {
	value, err := client.IncrBy(key, object).Result()
	if err != nil {
		beego.Error("incrBy error, key =", key)
		return 0
	}
	expires := time.Duration(expire) * time.Second
	client.Expire(key, expires)
	return value
}

func DecrBy(key string, object int64, expire int) int64 {
	value, err := client.DecrBy(key, object).Result()
	if err != nil {
		beego.Error("incrBy error, key =", key)
		return 0
	}
	expires := time.Duration(expire) * time.Second
	client.Expire(key, expires)
	return value
}

func Decr(key string, expire int) int64 {
	value, err := client.Decr(key).Result()
	expires := time.Duration(expire) * time.Second
	if err != nil {
		beego.Error("decr error, key =", key)
		return 0
	}
	client.Expire(key, expires)
	return value
}

func Remove(key string) bool {
	err := client.Del(key).Err()
	if err != nil {
		beego.Error("del error, key = ", key)
		return false
	}
	return true
}

// --------------------
// Encode
// 用gob进行数据编码
//
func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// -------------------
// Decode
// 用gob进行数据解码
//
func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
