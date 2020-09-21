package dataprovider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBooks(t *testing.T) {
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
	books, err := parseBooks(testInput)
	assert.Nil(t, err)

	expectedBooks := []Book{
		{
			title:       "Shakespeare's Tragedy of Romeo and Juliet",
			author:      "William Shakespeare",
			bookLinkref: "/ebooks/47960",
		}, {
			title:       "Dramas de Guillermo Shakspeare [vol. 1] (Spanish)",
			author:      "William Shakespeare",
			bookLinkref: "/ebooks/53207",
		},
	}

	assert.Equal(t, expectedBooks, books)
}
