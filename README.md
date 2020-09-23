## Build and run (Docker)

Make sure your Docker is installed and running, then simply run (in root repository directory):

```shell
./build/docker-build.sh && ./build/docker-run.sh
```

## Build and run (locally)

Make sure your Go is up to date and execute (in root repository directory):

```shell
./build/local-build.sh && ./build/local-run.sh
```

### Making requests

App is expecting JSON-encoded data on /search endpoint with these fields:
- title
- phrase

HTTPie example:
```shell
echo '{"title": "Romeo & Juliet", "phrase": "oh romeo romeo"}'  | http "http://localhost:8000/search" 
```

### Example app logs

```
2020/09/23 11:48:11 Starting http server
[stdin] echo '{"title": "Romeo and Juliet", "phrase": "oh romeo romeo"}'  | http "http://localhost:8000/search"
2020/09/23 11:48:12 Searching books with "Romeo and Juliet" title
2020/09/23 11:48:13 Read 21 book positions from external source
2020/09/23 11:48:13 Run job monitor
2020/09/23 11:48:13 Run search engine workers
2020/09/23 11:48:13 search worker 0 started
2020/09/23 11:48:13 search worker 3 started
2020/09/23 11:48:13 search worker 6 started
2020/09/23 11:48:13 search worker 1 started
2020/09/23 11:48:13 search worker 2 started
2020/09/23 11:48:13 Waiting to process all jobs
2020/09/23 11:48:13 search worker 7 started
2020/09/23 11:48:13 search worker 4 started
2020/09/23 11:48:13 search worker 5 started
2020/09/23 11:48:17 no result for book ("The Tragedy of Romeo and Juliet" - William Shakespeare [/ebooks/1112]): pattern not found
2020/09/23 11:48:21 no result for book ("Romeo and Juliet" - William Shakespeare [/ebooks/1513]): pattern not found
2020/09/23 11:48:25 no result for book ("Beautiful Stories from Shakespeare" - William Shakespeare and E. Nesbit [/ebooks/1430]): pattern not found
2020/09/23 11:48:30 no result for book ("Tales from Shakespeare" - Charles Lamb and Mary Lamb [/ebooks/573]): pattern not found
2020/09/23 11:48:34 no result for book ("Romeo and Juliet" - William Shakespeare [/ebooks/1777]): pattern not found
2020/09/23 11:48:40 no result for book ("Shakespeare's Tragedy of Romeo and Juliet" - William Shakespeare [/ebooks/47960]): pattern not found
2020/09/23 11:48:46 no result for book ("Dramas de Guillermo Shakspeare [vol. 1] (Spanish)" - William Shakespeare [/ebooks/53207]): pattern not found
2020/09/23 11:48:50 no result for book ("Romeo and Juliet. French (French)" - William Shakespeare [/ebooks/18143]): pattern not found
2020/09/23 11:48:54 result found
[stdout] HTTP/1.1 200 OK
[stdout] Content-Length: 201
[stdout] Content-Type: text/plain; charset=utf-8
[stdout] Date: Wed, 23 Sep 2020 09:48:54 GMT
[stdout] 
[stdout] Oh gentle Romeo,
[stdout]      If thou dost love, pronounce it faithfully;
[stdout]      Or if thou think I am too quickly won,
[stdout]      I'll frown and be perverse, and say thee nay,
[stdout]      So thou wilt woo: but else not
2020/09/23 11:48:54 Search worker 7 closed
2020/09/23 11:48:54 Search worker 4 closed
2020/09/23 11:48:54 Search worker 1 closed
2020/09/23 11:48:54 Search worker 5 closed
2020/09/23 11:48:54 Search worker 2 closed
2020/09/23 11:48:54 Search worker 0 closed
2020/09/23 11:48:54 Search worker 6 closed
2020/09/23 11:48:54 Search worker 3 closed
2020/09/23 11:48:58 Downloading interrupted
[stdin] echo '{"title": "Romeo and Juliet", "phrase": "oh romeo romeo"}'  | http "http://localhost:8000/search" ### the same request
2020/09/23 11:49:51 found cached query result
[stdout] HTTP/1.1 200 OK
[stdout] Content-Length: 201
[stdout] Content-Type: text/plain; charset=utf-8
[stdout] Date: Wed, 23 Sep 2020 09:49:51 GMT
[stdout] 
[stdout] Oh gentle Romeo,
[stdout]      If thou dost love, pronounce it faithfully;
[stdout]      Or if thou think I am too quickly won,
[stdout]      I'll frown and be perverse, and say thee nay,
[stdout]      So thou wilt woo: but else not

[stdin] echo '{"title": "Romeo & Juliet", "phrase": "o my romeo"}'  | http "http://localhost:8000/search" 
2020/09/23 11:51:11 Searching books with "Romeo & Juliet" title
2020/09/23 11:51:12 Read 21 book positions from external source
2020/09/23 11:51:12 Run job monitor
2020/09/23 11:51:12 Run search engine workers
2020/09/23 11:51:12 Waiting to process all jobs
2020/09/23 11:51:12 Load book from cache ("The Tragedy of Romeo and Juliet" - William Shakespeare)
2020/09/23 11:51:12 search worker 7 started
2020/09/23 11:51:12 search worker 1 started
2020/09/23 11:51:12 search worker 3 started
2020/09/23 11:51:12 search worker 0 started
2020/09/23 11:51:12 search worker 5 started
2020/09/23 11:51:12 search worker 2 started
2020/09/23 11:51:12 search worker 4 started
2020/09/23 11:51:12 search worker 6 started
2020/09/23 11:51:12 Load book from cache ("Romeo and Juliet" - William Shakespeare)
2020/09/23 11:51:12 Load book from cache ("Beautiful Stories from Shakespeare" - William Shakespeare and E. Nesbit)
2020/09/23 11:51:12 Load book from cache ("Tales from Shakespeare" - Charles Lamb and Mary Lamb)
2020/09/23 11:51:12 Load book from cache ("Romeo and Juliet" - William Shakespeare)
2020/09/23 11:51:12 Load book from cache ("Shakespeare's Tragedy of Romeo and Juliet" - William Shakespeare)
2020/09/23 11:51:12 Load book from cache ("Dramas de Guillermo Shakspeare [vol. 1] (Spanish)" - William Shakespeare)
2020/09/23 11:51:12 Load book from cache ("Romeo and Juliet. French (French)" - William Shakespeare)
2020/09/23 11:51:12 Load book from cache ("Characters of Shakespeare's Plays" - William Hazlitt)
2020/09/23 11:51:12 Load book from cache ("Romeo and Juliet. Polish (Polish)" - William Shakespeare)
2020/09/23 11:51:12 result found
2020/09/23 11:51:12 Search worker 7 closed
2020/09/23 11:51:12 Search worker 5 closed
2020/09/23 11:51:12 no result for book ("Romeo and Juliet. French (French)" - William Shakespeare [/ebooks/18143]): pattern not found
2020/09/23 11:51:12 Search worker 3 closed
2020/09/23 11:51:12 Search worker 6 closed
2020/09/23 11:51:12 Search worker 2 closed
2020/09/23 11:51:12 Search worker 4 closed
2020/09/23 11:51:12 Search worker 0 closed
2020/09/23 11:51:12 Search worker 1 closed
2020/09/23 11:51:14 text book not available ("Romeo and Juliet" - William Shakespeare [/ebooks/26268]): failed to get txt linkref: finding txt linkref failed: txt linkref is not available
[stdout] HTTP/1.1 200 OK
[stdout] Content-Length: 201
[stdout] Content-Type: text/plain; charset=utf-8
[stdout] Date: Wed, 23 Sep 2020 09:51:12 GMT
[stdout] 
[stdout] O Romeo, Romeo, wherefore art thou Romeo?
[stdout] Deny thy father and refuse thy name.
[stdout] Or if thou wilt not, be but sworn my love,
[stdout] And I’ll no longer be a Capulet.
[stdout] 
[stdout] ROMEO.
[stdout] [_Aside._] Shall I hear more,
(...)
```