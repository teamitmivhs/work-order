package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	webpush "github.com/SherClockHolmes/webpush-go"

	"teamitmivhs/work-order-backend/config"
)

var (
	vapidPublicKey  string
	vapidPrivateKey string
)

func InitPushNotifications() {
	vapidPublicKey = os.Getenv("VAPID_PUBLIC_KEY")
	vapidPrivateKey = os.Getenv("VAPID_PRIVATE_KEY")
	if vapidPublicKey == "" || vapidPrivateKey == "" {
		privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
		if err != nil {
			log.Printf("[WARNING] failed to generate VAPID keys: %v", err)
			return
		}
		vapidPrivateKey = privateKey
		vapidPublicKey = publicKey
		log.Println("[WARNING] VAPID_PUBLIC_KEY/VAPID_PRIVATE_KEY not set; generated temporary push keys")
	}

	if err := ensurePushSubscriptionTable(); err != nil {
		log.Printf("[WARNING] failed to ensure push subscription table: %v", err)
	}
}

func PushPublicKey() string {
	return vapidPublicKey
}

func ensurePushSubscriptionTable() error {
	db := config.GetDB()
	if db == nil {
		return nil
	}
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS push_subscriptions (
			ID int NOT NULL AUTO_INCREMENT,
			UserID int NOT NULL,
			Endpoint text NOT NULL,
			EndpointHash char(64) NOT NULL,
			P256DH varchar(255) NOT NULL,
			Auth varchar(255) NOT NULL,
			UserAgent varchar(255) DEFAULT NULL,
			CreatedAt datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UpdatedAt datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (ID),
			UNIQUE KEY uq_push_endpoint_hash (EndpointHash),
			KEY idx_push_user (UserID),
			CONSTRAINT fk_push_user FOREIGN KEY (UserID) REFERENCES members (ID) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
	`)
	return err
}

func SavePushSubscription(userID int, endpoint, p256dh, auth, userAgent string) error {
	if endpoint == "" || p256dh == "" || auth == "" {
		return fmt.Errorf("push subscription is incomplete")
	}
	db := config.GetDB()
	if db == nil {
		return fmt.Errorf("database is not available")
	}
	_, err := db.Exec(`
		INSERT INTO push_subscriptions (UserID, Endpoint, EndpointHash, P256DH, Auth, UserAgent)
		VALUES (?, ?, SHA2(?, 256), ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			UserID = VALUES(UserID),
			Endpoint = VALUES(Endpoint),
			P256DH = VALUES(P256DH),
			Auth = VALUES(Auth),
			UserAgent = VALUES(UserAgent),
			UpdatedAt = CURRENT_TIMESTAMP
	`, userID, endpoint, endpoint, p256dh, auth, nullableString(userAgent))
	return err
}

func DeletePushSubscription(endpoint string) error {
	if endpoint == "" {
		return nil
	}
	db := config.GetDB()
	if db == nil {
		return nil
	}
	_, err := db.Exec("DELETE FROM push_subscriptions WHERE EndpointHash = SHA2(?, 256)", endpoint)
	return err
}

func NotifyNewWorkOrder(orderID int64, device, location, priority string) {
	if vapidPublicKey == "" || vapidPrivateKey == "" {
		return
	}
	db := config.GetDB()
	if db == nil {
		return
	}

	rows, err := db.Query("SELECT ID, Endpoint, P256DH, Auth FROM push_subscriptions")
	if err != nil {
		log.Printf("[WARNING] failed to query push subscriptions: %v", err)
		return
	}
	defer rows.Close()

	body := fmt.Sprintf("#%d %s - %s", orderID, fallback(device, "Device"), fallback(location, "Lokasi belum diisi"))
	payload, _ := json.Marshal(map[string]any{
		"title":       "Work order baru masuk",
		"body":        body,
		"url":         "/index",
		"workOrderId": orderID,
		"priority":    priority,
	})

	for rows.Next() {
		var id int64
		var endpoint, p256dh, auth string
		if err := rows.Scan(&id, &endpoint, &p256dh, &auth); err != nil {
			log.Printf("[WARNING] failed to scan push subscription: %v", err)
			continue
		}
		sendPush(id, endpoint, p256dh, auth, payload)
	}
}

func sendPush(subscriptionID int64, endpoint, p256dh, auth string, payload []byte) {
	resp, err := webpush.SendNotification(payload, &webpush.Subscription{
		Endpoint: endpoint,
		Keys: webpush.Keys{
			P256dh: p256dh,
			Auth:   auth,
		},
	}, &webpush.Options{
		Subscriber:      "mailto:admin@localhost",
		TTL:             86400,
		VAPIDPublicKey:  vapidPublicKey,
		VAPIDPrivateKey: vapidPrivateKey,
	})
	if err != nil {
		log.Printf("[WARNING] failed to send push notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusGone || resp.StatusCode == http.StatusNotFound {
		db := config.GetDB()
		if db != nil {
			_, _ = db.Exec("DELETE FROM push_subscriptions WHERE ID = ?", subscriptionID)
		}
	}
}

func fallback(value, fallbackValue string) string {
	if value == "" {
		return fallbackValue
	}
	return value
}

func nullableString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}
