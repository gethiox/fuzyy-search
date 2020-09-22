### Example app logs

demonstrates working caching mechanism
```text
2020/09/22 17:20:03 Starting http server
[first request]
2020/09/22 17:20:09 Searching books with "The Tragedy of Romeo and Juliet" title
2020/09/22 17:20:10 Read 13 book positions from external source
2020/09/22 17:20:10 Run job monitor
2020/09/22 17:20:10 Run search engine workers
2020/09/22 17:20:10 Waiting to process all jobs
2020/09/22 17:20:15 no result for book ("The Tragedy of Romeo and Juliet" - William Shakespeare [/ebooks/1112]): pattern not found
2020/09/22 17:20:20 no result for book ("Romeo and Juliet" - William Shakespeare [/ebooks/1513]): pattern not found
2020/09/22 17:20:27 no result for book ("Romeo and Juliet" - William Shakespeare [/ebooks/1777]): pattern not found
2020/09/22 17:20:33 no result for book ("Shakespeare's Tragedy of Romeo and Juliet" - William Shakespeare [/ebooks/47960]): pattern not found
2020/09/22 17:20:40 no result for book ("Romeo and Juliet. French (French)" - William Shakespeare [/ebooks/18143]): pattern not found
2020/09/22 17:20:48 no result for book ("Romeo and Juliet. Polish (Polish)" - William Shakespeare [/ebooks/27062]): pattern not found
2020/09/22 17:20:50 text book not available ("Romeo and Juliet" - William Shakespeare [/ebooks/26268]): failed to get txt linkref: finding txt linkref failed: txt linkref is not available
2020/09/22 17:20:50 no result for book ("Romeo and Juliet" - William Shakespeare [/ebooks/26268]): pattern not found
2020/09/22 17:20:57 no result for book ("Romeo and Juliet" - William Shakespeare [/ebooks/2261]): pattern not found
2020/09/22 17:21:04 no result for book ("Romeo and Juliet. German (German)" - William Shakespeare [/ebooks/6996]): pattern not found
2020/09/22 17:21:10 no result for book ("The Works of William Shakespeare [Cambridge Edition] [Vol. 7 of 9]" - William Shakespeare [/ebooks/47715]): pattern not found
2020/09/22 17:21:16 no result for book ("Romeo and Juliet. Finnish (Finnish)" - William Shakespeare [/ebooks/15643]): pattern not found
2020/09/22 17:21:22 no result for book ("Romeo and Juliet. Dutch (Dutch)" - William Shakespeare [/ebooks/49880]): pattern not found
2020/09/22 17:21:27 no result for book ("Romeo and Juliet. Greek (Greek)" - William Shakespeare [/ebooks/31808]): pattern not found
2020/09/22 17:21:27 All jobs processed
[second request, query the same as previous]
2020/09/22 17:21:37 Searching books with "The Tragedy of Romeo and Juliet" title
2020/09/22 17:21:37 Read 13 book positions from cache
2020/09/22 17:21:37 Run job monitor
2020/09/22 17:21:37 Run search engine workers
2020/09/22 17:21:37 Load book from cache ("The Tragedy of Romeo and Juliet" - William Shakespeare)
2020/09/22 17:21:37 Load book from cache ("Romeo and Juliet" - William Shakespeare)
2020/09/22 17:21:37 Load book from cache ("Romeo and Juliet" - William Shakespeare)
2020/09/22 17:21:37 Load book from cache ("Shakespeare's Tragedy of Romeo and Juliet" - William Shakespeare)
2020/09/22 17:21:37 Load book from cache ("Romeo and Juliet. French (French)" - William Shakespeare)
2020/09/22 17:21:37 Load book from cache ("Romeo and Juliet. Polish (Polish)" - William Shakespeare)
2020/09/22 17:21:37 Load book from cache ("Romeo and Juliet" - William Shakespeare)
2020/09/22 17:21:37 Load book from cache ("Romeo and Juliet" - William Shakespeare)
2020/09/22 17:21:37 Load book from cache ("Romeo and Juliet. German (German)" - William Shakespeare)
2020/09/22 17:21:37 Load book from cache ("The Works of William Shakespeare [Cambridge Edition] [Vol. 7 of 9]" - William Shakespeare)
2020/09/22 17:21:37 Load book from cache ("Romeo and Juliet. Finnish (Finnish)" - William Shakespeare)
2020/09/22 17:21:37 Load book from cache ("Romeo and Juliet. Dutch (Dutch)" - William Shakespeare)
2020/09/22 17:21:37 Load book from cache ("Romeo and Juliet. Greek (Greek)" - William Shakespeare)
2020/09/22 17:21:37 Waiting to process all jobs
2020/09/22 17:21:37 no result for book ("The Tragedy of Romeo and Juliet" - William Shakespeare [/ebooks/1112]): pattern not found
2020/09/22 17:21:37 no result for book ("Romeo and Juliet" - William Shakespeare [/ebooks/1513]): pattern not found
2020/09/22 17:21:37 no result for book ("Shakespeare's Tragedy of Romeo and Juliet" - William Shakespeare [/ebooks/47960]): pattern not found
2020/09/22 17:21:37 no result for book ("Romeo and Juliet" - William Shakespeare [/ebooks/1777]): pattern not found
2020/09/22 17:21:37 no result for book ("Romeo and Juliet" - William Shakespeare [/ebooks/26268]): pattern not found
2020/09/22 17:21:37 no result for book ("Romeo and Juliet. French (French)" - William Shakespeare [/ebooks/18143]): pattern not found
2020/09/22 17:21:37 no result for book ("Romeo and Juliet" - William Shakespeare [/ebooks/2261]): pattern not found
2020/09/22 17:21:37 no result for book ("Romeo and Juliet. German (German)" - William Shakespeare [/ebooks/6996]): pattern not found
2020/09/22 17:21:37 no result for book ("Romeo and Juliet. Finnish (Finnish)" - William Shakespeare [/ebooks/15643]): pattern not found
2020/09/22 17:21:37 no result for book ("Romeo and Juliet. Polish (Polish)" - William Shakespeare [/ebooks/27062]): pattern not found
2020/09/22 17:21:37 no result for book ("Romeo and Juliet. Greek (Greek)" - William Shakespeare [/ebooks/31808]): pattern not found
2020/09/22 17:21:37 no result for book ("The Works of William Shakespeare [Cambridge Edition] [Vol. 7 of 9]" - William Shakespeare [/ebooks/47715]): pattern not found
2020/09/22 17:21:37 no result for book ("Romeo and Juliet. Dutch (Dutch)" - William Shakespeare [/ebooks/49880]): pattern not found
2020/09/22 17:21:37 All jobs processed
[third request, changed phrase]
2020/09/22 17:21:57 Searching books with "The Tragedy of Romeo and Juliet" title
2020/09/22 17:21:57 Read 13 book positions from cache
2020/09/22 17:21:57 Run job monitor
2020/09/22 17:21:57 Run search engine workers
2020/09/22 17:21:57 Load book from cache ("The Tragedy of Romeo and Juliet" - William Shakespeare)
2020/09/22 17:21:57 Load book from cache ("Romeo and Juliet" - William Shakespeare)
2020/09/22 17:21:57 Load book from cache ("Romeo and Juliet" - William Shakespeare)
2020/09/22 17:21:57 Load book from cache ("Shakespeare's Tragedy of Romeo and Juliet" - William Shakespeare)
2020/09/22 17:21:57 Load book from cache ("Romeo and Juliet. French (French)" - William Shakespeare)
2020/09/22 17:21:57 Load book from cache ("Romeo and Juliet. Polish (Polish)" - William Shakespeare)
2020/09/22 17:21:57 Load book from cache ("Romeo and Juliet" - William Shakespeare)
2020/09/22 17:21:57 Load book from cache ("Romeo and Juliet" - William Shakespeare)
2020/09/22 17:21:57 Load book from cache ("Romeo and Juliet. German (German)" - William Shakespeare)
2020/09/22 17:21:57 Load book from cache ("The Works of William Shakespeare [Cambridge Edition] [Vol. 7 of 9]" - William Shakespeare)
2020/09/22 17:21:57 Load book from cache ("Romeo and Juliet. Finnish (Finnish)" - William Shakespeare)
2020/09/22 17:21:57 Load book from cache ("Romeo and Juliet. Dutch (Dutch)" - William Shakespeare)
2020/09/22 17:21:57 Load book from cache ("Romeo and Juliet. Greek (Greek)" - William Shakespeare)
2020/09/22 17:21:57 Waiting to process all jobs
[phrase found]
2020/09/22 17:21:57 result ok 
2020/09/22 17:21:57 no result for book ("Romeo and Juliet. Polish (Polish)" - William Shakespeare [/ebooks/27062]): pattern not found
2020/09/22 17:21:57 no result for book ("Romeo and Juliet" - William Shakespeare [/ebooks/26268]): pattern not found
2020/09/22 17:21:57 no result for book ("Romeo and Juliet. French (French)" - William Shakespeare [/ebooks/18143]): pattern not found
2020/09/22 17:21:57 no result for book ("Romeo and Juliet. Finnish (Finnish)" - William Shakespeare [/ebooks/15643]): pattern not found
[fourth request, same query as previous one, answer directly returned from cache]
2020/09/22 17:22:21 found cached query result```