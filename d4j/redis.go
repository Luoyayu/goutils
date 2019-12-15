package d4j

/* HashMap
* 1)   key: d4j:bookId
       keyExample: d4j:14659
           - field: share_url
           - field: share_key
           - field: name
           - field: author
           - field: category
 * 2)   key: d4j:bookId:tags
 * 3)   key: d4j:bookId:keywords
*/

func (r *Redis) Stores(key string, book *Struct) error {
	kv := map[string]interface{}{
		"name":      book.Name,
		"author":    book.Author,
		"share_url": book.ShareLink.Url,
		"share_key": book.ShareLink.Key,
		"category":  book.Category,
		"title":     book.Title,
	}
	stat := r.HMSet(key, kv)
	r.SAdd("d4j:"+book.ID+":tags", book.Tags)
	r.SAdd("d4j:"+book.ID+":keywords", book.KeyWords)
	return stat.Err()
}
