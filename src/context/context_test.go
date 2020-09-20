package context

import (
	"testing"
)

func TestExpectedContext(t *testing.T) {
	content := `   Rom. She speakes.
Oh speake againe bright Angell, for thou art
As glorious to this night being ore my head,
As is a winged messenger of heauen
Vnto the white vpturned wondring eyes
Of mortalls that fall backe to gaze on him,
When he bestrides the lazie puffing Cloudes,
And sailes vpon the bosome of the ayre

   Iul. O Romeo, Romeo, wherefore art thou Romeo?
Denie thy Father and refuse thy name:
Or if thou wilt not, be but sworne to my Loue,
And Ile no longer be a Capulet

   Rom. Shall I heare more, or shall I speake at this?
  Iu. 'Tis but thy name that is my Enemy:
Thou art thy selfe, though not a Mountague,
What's Mountague? it is nor hand nor foote,
Nor arme, nor face, O be some other name
Belonging to a man.
What? in a names that which we call a Rose,
By any other word would smell as sweete,
So Romeo would, were he not Romeo cal'd,
Retaine that deare perfection which he owes,
Without that title Romeo, doffe thy name,
And for thy name which is no part of thee,
Take all my selfe`

	expectedContext := "O Romeo, Romeo, wherefore art thou Romeo?\n" +
		"Denie thy Father and refuse thy name:\n" +
		"Or if thou wilt not, be but sworne to my Loue,\n" +
		"And Ile no longer be a Capulet"

	provider := NewProvider()
	context, err := provider.ProvideContext(content, 321, 335)
	if err != nil {
		t.Fatal("Failed to gather context: ", err)
	}

	if context != expectedContext {
		t.Logf("Returned context is different than expected\n")
		t.Logf("Wanted: %#v\n", expectedContext)
		t.Logf("  Have: %#v\n", context)
		t.Fail()
	}
}
