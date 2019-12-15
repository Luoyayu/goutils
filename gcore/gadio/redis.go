package gadio

import (
	"github.com/pkg/errors"
)

func (r *Redis) Stores(key string, entity *RadioData, ownerId int64, ownerName string, msgId int, chatId int64) (err error) {
	defer func() {
		err = errors.Wrap(err, "Stores ->")
	}()
	if err = r.HMSet(key, map[string]interface{}{
		"attr:title":         entity.Attrs.Title,
		"attr:desc":          entity.Attrs.Desc,
		"attr:cover":         entity.Attrs.Cover,
		"attr:published-at":  entity.Attrs.PublishedAt,
		"channel:owner-id":   ownerId,
		"channel:owner-name": ownerName,
		"channel:msg-id":     msgId,
		"channel:chat-id":    chatId,
	}).Err(); err != nil {
		return
	}

	return
}
