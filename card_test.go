package gmopg

import (
	"testing"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken("4111111111111111", "202512", "1234", "GMO PG")
	if err != nil {
		t.Errorf("GenerateToken returns error: %v", err)
	}
	// want := "Uaq9/EaY8vY7VycbubHK3qRTQOVGpXODNuJEggHxKoLz2npoRFiTYxlCYedBQyujHFbGkPXsDbaBO+6cvchMcV+z0xcKnZ8WBcgf1iEbdt/Y6mKigOGKJFyF98UTb6aT0h2IzbO/1J1WvGhfJ37gIlGSz24MK0PdKZW+r40Jj9dfYWqMm0hilsUUecqBMFwISft1A4t2/K27fLcMhPt3mpKdmoUpyPSG7S/cqfGnd42VUVesXt8RQqmDqLY8X9pHQG8RucT/km1Kk+WjOaw6uHl4fjirWm8ExeW2psnpK9/iqgzjRlwx56kgv1Ul0IbYmOujJLIju2uTV8VYTEljVQ=="
	if len(token) == 0 {
		t.Errorf("token is empty: %s", token)
	}
}
