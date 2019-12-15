package d4j

import (
	"testing"
)

func TestHandleRSS(t *testing.T) {
	//assert.NoError(t, GetBookInfoFromRss())

	/*next, book, err := GetBookInfo("14796")
	  _ = StoreBook(redisC, book)
	  if err != nil {
	      log.Println(err)
	  } else {
	      log.Printf("%+v\t%+v\n", book, book.ShareLink)
	      log.Println("next: ", next)
	  }*/

	// 4593 -> 200
	// 4594 -> 301 -> 4593
	// 4592 -> 404

	//CrawlAll(-1)
}
