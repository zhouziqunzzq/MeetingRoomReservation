package lockcontroller

import (
	"encoding/hex"
	"github.com/zhouziqunzzq/MeetingRoomReservation/config"
	"github.com/yanzay/log"
)

func Unlock(ip string) (err error) {
	cmd, err := hex.DecodeString(config.GlobalConfig.LOCK_CMD)
	if err != nil {
		log.Error("Failed while parsing lock cmd")
		log.Error(err.Error())
		return
	}
	log.Println(ip + ":" + config.GlobalConfig.LOCK_PORT)
	log.Println(cmd)
	return
}
