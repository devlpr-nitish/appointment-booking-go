package workers

import (
	"log"
	"time"

	"github.com/devlpr-nitish/appointment-booking-go/internal/handlers"
	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/devlpr-nitish/appointment-booking-go/internal/repositories"
)

func StartRequestCleaner(repo repositories.RequestRepository, hub *handlers.Hub) {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		log.Println("Started RequestCleaner background worker")
		for range ticker.C {
			// Find OPEN requests older than 10 minutes
			threshold := time.Now().Add(-10 * time.Minute)
			expiredRequests, err := repo.FindExpiredOpenRequests(threshold)
			if err != nil {
				log.Println("Error finding expired requests:", err)
				continue
			}

			for _, req := range expiredRequests {
				err := repo.UpdateStatus(req.ID, models.RequestStatusCanceled)
				if err != nil {
					log.Printf("Failed to auto-cancel request %s: %v\n", req.ID, err)
					continue
				}

				log.Printf("Auto-canceled expired request %s\n", req.ID)
				
				// Broadcast expiration to both user and subscribed experts
				if hub != nil {
					// Tell the user
					hub.BroadcastEvent("REQUEST_EXPIRED", map[string]interface{}{
						"request_id": req.ID,
					})
				}
			}
		}
	}()
}
