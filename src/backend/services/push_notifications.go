package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	webpush "github.com/SherClockHolmes/webpush-go"

	"teamitmivhs/work-order-backend/config"
	"teamitmivhs/work-order-backend/utils"
)

var (
	vapidPublicKey  string
	vapidPrivateKey string
)

func InitPushNotifications() {
	vapidPublicKey = os.Getenv("VAPID_PUBLIC_KEY")
	vapidPrivateKey = os.Getenv("VAPID_PRIVATE_KEY")
	if vapidPublicKey == "" || vapidPrivateKey == "" {
		log.Println("[WARNING] VAPID_PUBLIC_KEY/VAPID_PRIVATE_KEY not set; push notifications disabled")
		return
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

func NotifyNewWorkOrder(orderID int64, device, location, problem, priority string) {
	title, body := formatNewWorkOrderNotification(device, location, problem, priority)
	NotifyWorkOrderBroadcast(title, body, orderID, map[string]any{
		"priority": priority,
	})
}

func formatNewWorkOrderNotification(device, location, problem, priority string) (string, string) {
	priorityLabel := map[string]string{
		"low": "Rendah", "medium": "Sedang", "high": "Tinggi", "urgent": "Urgent",
	}[strings.ToLower(strings.TrimSpace(priority))]
	title := "Work Order Baru"
	if priorityLabel != "" {
		title += " • Prioritas " + priorityLabel
	}
	body := fmt.Sprintf(
		"Perangkat: %s\nLokasi: %s\nKendala: %s",
		fallback(strings.TrimSpace(device), "Belum diketahui"),
		fallback(strings.TrimSpace(location), "Belum diisi"),
		fallback(strings.TrimSpace(problem), "Belum ada deskripsi"),
	)
	return title, body
}

func NotifyWorkOrderBroadcast(title, body string, orderID int64, extra map[string]any) {
	notifyWorkOrder(title, body, orderID, extra, "")
}

func NotifyWorkOrderAdmins(title, body string, orderID int64, extra map[string]any) {
	notifyWorkOrder(title, body, orderID, extra, "admin")
}

func NotifyWorkOrderUsers(userIDs []int, title, body string, orderID int64, extra map[string]any) {
	if len(userIDs) == 0 {
		return
	}
	ids := make(map[int]bool, len(userIDs))
	for _, id := range userIDs {
		if id > 0 {
			ids[id] = true
		}
	}
	if len(ids) == 0 {
		return
	}
	notifyWorkOrder(title, body, orderID, extra, "users", ids)
}

func notifyWorkOrder(title, body string, orderID int64, extra map[string]any, target string, userIDs ...map[int]bool) {
	if vapidPublicKey == "" || vapidPrivateKey == "" {
		return
	}
	db := config.GetDB()
	if db == nil {
		return
	}

	rows, err := db.Query(`
		SELECT ps.ID, ps.UserID, ps.Endpoint, ps.P256DH, ps.Auth, COALESCE(m.Role, '')
		FROM push_subscriptions ps
		JOIN members m ON m.ID = ps.UserID
	`)
	if err != nil {
		log.Printf("[WARNING] failed to query push subscriptions: %v", err)
		return
	}
	defer rows.Close()

	payloadMap := map[string]any{
		"title":       "Work Order Baru",
		"body":        body,
		"url":         "/",
		"workOrderId": orderID,
	}
	for key, value := range extra {
		payloadMap[key] = value
	}
	if title != "" {
		payloadMap["title"] = title
	}
	payload, _ := json.Marshal(payloadMap)

	for rows.Next() {
		var id int64
		var userID int
		var endpoint, p256dh, auth, role string
		if err := rows.Scan(&id, &userID, &endpoint, &p256dh, &auth, &role); err != nil {
			log.Printf("[WARNING] failed to scan push subscription: %v", err)
			continue
		}
		if target == "admin" && !utils.IsAdminRole(role) {
			continue
		}
		if target == "users" && (len(userIDs) == 0 || !userIDs[0][userID]) {
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
