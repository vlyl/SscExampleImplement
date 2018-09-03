package MultisignatureEscrowAccount

import (
	"log"
	"testing"
	"time"
)

func TestYearLater(t *testing.T) {
	now := time.Now()
	log.Print(now)
	log.Print(time.Unix(int64(YearLater(now)), 0))
}
