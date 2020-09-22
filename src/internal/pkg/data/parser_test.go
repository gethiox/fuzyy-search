package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindBooks(t *testing.T) {
	testInput := `
(...)

<li class="booklink">
<a class="link" href="/ebooks/47960" accesskey="8">
<span class="cell leftcell with-cover">
<img class="cover-thumb" src="/cache/epub/47960/pg47960.cover.small.jpg" alt="">
</span>
<span class="cell content">
<span class="title">Shakespeare's Tragedy of Romeo and Juliet</span>
<span class="subtitle">William Shakespeare</span>
<span class="extra">224 downloads</span>
</span>
<span class="hstrut"></span>
</a>
</li>

<li class="booklink">
<a class="link" href="/ebooks/53207" accesskey="9">
<span class="cell leftcell with-cover">
<img class="cover-thumb" src="/cache/epub/53207/pg53207.cover.small.jpg" alt="">
</span>
<span class="cell content">
<span class="title">Dramas de Guillermo Shakspeare [vol. 1] (Spanish)</span>
<span class="subtitle">William Shakespeare</span>
<span class="extra">103 downloads</span>
</span>
<span class="hstrut"></span>
</a>
</li>

(...)
`
	books, err := findBooks(testInput)
	assert.Nil(t, err)

	expectedBooks := []Book{
		{
			Title:       "Shakespeare's Tragedy of Romeo and Juliet",
			Author:      "William Shakespeare",
			bookLinkref: "/ebooks/47960",
		}, {
			Title:       "Dramas de Guillermo Shakspeare [vol. 1] (Spanish)",
			Author:      "William Shakespeare",
			bookLinkref: "/ebooks/53207",
		},
	}

	assert.Equal(t, expectedBooks, books)
}

func TestFindTxtLinkref(t *testing.T) {
	testInput := `
(...)

<td class="noprint">
<a href="/ebooks/send/msdrive/32571.kindle.noimages" title="Send to OneDrive." rel="nofollow"><span class="icon icon_msdrive"></span></a>
</td>
</tr><tr class="even" about="https://www.gutenberg.org/files/32571/32571-0.txt" typeof="pgterms:file">
<td><span class="icon icon_book"></span></td>
<td property="dcterms:format" content="text/plain; charset=utf-8" datatype="dcterms:IMT" class="unpadded icon_save"><a href="/files/32571/32571-0.txt" type="text/plain; charset=utf-8" class="link" title="Download">Plain Text UTF-8</a></td>
<td class="noscreen">https://www.gutenberg.org/files/32571/32571-0.txt</td>
<td class="right" property="dcterms:extent" content="315190">308 kB</td>
<td class="noprint">
</td>

(...)
`
	linkref, err := findTxtLinkref(testInput)
	assert.Nil(t, err)

	assert.Equal(t, "/files/32571/32571-0.txt", linkref)

}
