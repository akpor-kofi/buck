package buckis

func (d *dict) runCommand(cmd command) {
	switch cmd.instruction {
	case SET:
		_ = d.set(DONT_SAVE, cmd.args[0].(string), cmd.args[1])

	case INCRBY:
		_, _ = d.incrBy(DONT_SAVE, cmd.args[0].(string), int(int64(cmd.args[1].(float64))))
	case HSET:
		var hashes []any

		for _, h := range cmd.args[1].([]any) {
			hashes = append(hashes, h)
		}

		_ = d.hset(DONT_SAVE, cmd.args[0].(string), hashes...)
	case HINCRBY:
		//d.hIncrBy(examples.args[0].(string), examples.args[1].(string), examples.args[2].(int))
	case SADD:
		_ = d.sadd(DONT_SAVE, cmd.args[0].(string), cmd.args[1].(string))
	case SREM:
		_ = d.srem(DONT_SAVE, cmd.args[0].(string), cmd.args[1].(string))
	case SMOVE:
		_ = d.smove(DONT_SAVE, cmd.args[0].(string), cmd.args[1].(string), cmd.args[2].(string))
	case ZADD:
		_ = d.zadd(DONT_SAVE, cmd.args[0].(string), cmd.args[1].(string), int(cmd.args[2].(float64)))
	case ZRANGESTORE:

	case ZINCRBY:
		_ = d.zincrby(DONT_SAVE, cmd.args[0].(string), cmd.args[1].(string), int(cmd.args[2].(float64)))
	case ZREM:
		_ = d.zrem(DONT_SAVE, cmd.args[0].(string), cmd.args[1].(string))
	case GADD:
		d.gAdd(DONT_SAVE, cmd.args[0].(string), cmd.args[1].(string), cmd.args[2].(string))
	case RPUSH:
		_, _ = d.rPush(DONT_SAVE, cmd.args[0].(string), cmd.args[1].(string))
	case LPUSH:
		_, _ = d.lPush(DONT_SAVE, cmd.args[0].(string), cmd.args[1].(string))
	case LPOP:
		_, _ = d.lPop(DONT_SAVE, cmd.args[0].(string))
	case RPOP:
		_, _ = d.rPop(DONT_SAVE, cmd.args[0].(string))
	case FTCREATE:
	case BFADD:
		_ = d.bfAdd(DONT_SAVE, cmd.args[0].(string), cmd.args[1].(string))
	case JSONSET:
		_ = d.jsonSet(DONT_SAVE, cmd.args[0].(string), cmd.args[1].(string), cmd.args[2])

	}

}
