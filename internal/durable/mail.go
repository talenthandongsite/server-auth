package durable

import (
	"context"
	"fmt"
	"log"

	"github.com/talenthandongsite/server-auth/internal/variable"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type MailClient struct {
	service *gmail.Service
}

func InitMailClient(ctx context.Context) (*MailClient, error) {
	apiKey := variable.GetEnv(ctx, variable.GOOGLE_API_KEY)

	srv, err := gmail.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &MailClient{srv}, nil
}

func (m *MailClient) Test() {
	user := "me"
	r, err := m.service.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}
}
