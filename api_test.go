package memberclicks

import (
	"log"
	"os"

	"golang.org/x/net/context"

	_ "github.com/joho/godotenv/autoload"
)

var (
	ctx = context.Background()
	mc  = New(os.Getenv("MEMBERCLICKS_ORG_ID"), os.Getenv("MEMBERCLICKS_CLIENT_ID"), os.Getenv("MEMBERCLICKS_CLIENT_SECRET"))
)

func init() {
	if err := mc.Auth(ctx); err != nil {
		log.Fatalf("Could not authorize with MemberClicks: %v", err)
	}
}
