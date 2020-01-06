package pointSign

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/guregu/dynamo"
	"google.golang.org/grpc/grpclog"
	"strconv"
	"time"
)

type BuddySign struct {
	Uid string

	LastSignTs int64 `dynamo:"ts"`
}

func createSignKey(path string, uuid string, uid string, db *dynamo.DB) (*errors, string) {
	var ts int64 = 0
	table := dynamo.Table{}
	if logcfg.IsTest == 1 {
		table = db.Table("buddySignTest")
	} else {
		table = db.Table("buddySign")
	}

	var result BuddySign
	var isExpired bool = false
	errDB := table.Get("Uid", uid).One(&result)

	if errDB != nil && errDB != dynamo.ErrNotFound {
		return errDB, ""
	} else if errDB == dynamo.ErrNotFound {
		result.Uid = uid
		ts = time.Now().Unix() +(60*60*24*21)
		result.LastSignTs = ts
		isExpired = true
	} else {
		//已经过期了
		if result.LastSignTs  <= ts - 120 {
			result.LastSignTs = ts
			isExpired = true
		}
	}

	tss := strconv.FormatInt(result.LastSignTs, 10)
	sstring := path + "-" + tss + "-" + uuid + "-" + uid + "-" + "ddhjxjsbgsn5"
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(sstring))
	cipherStr := md5Ctx.Sum(nil)
	if isExpired {
		errDB = table.Put(result).Run()
		if errDB != nil {
			return errDB, ""

		}
	}
	return nil, tss + "-" + uuid + "-" + uid + "-" + hex.EncodeToString(cipherStr)

}
