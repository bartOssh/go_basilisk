package services

import "testing"

// CollisionsNumToCheck is a number of tokens to generate and check against each other
const CollisionsNumToCheck = 10000

// TestRandToken is a dump token generation algorithm test
func TestRandToken(t *testing.T) {
	zeroToken, err := RandToken(128)
	if err != nil {
		t.Errorf("unknown error: %s", err)
	}
	for i := 0; i < CollisionsNumToCheck; i++ {
		firstToken, err := RandToken(128)
		if err != nil {
			t.Errorf("unknown error: %s", err)
		}
		if err != nil {
			t.Errorf("unknown error: %s", err)
		}
		secondToken, err := RandToken(128)
		if firstToken == secondToken {
			t.Errorf("token are in collision, two tokens generated one after another cannot be equal\n %v == %v\n", firstToken, secondToken)
		}
		if zeroToken == firstToken || zeroToken == secondToken {
			t.Errorf("tokens are in collision, pattern is to obvious and generation is not random\n repeated token : %v\n", zeroToken)
		}
	}
}
