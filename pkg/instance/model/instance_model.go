package instance_model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Instance struct {
	Id               string    `json:"id" gorm:"type:uuid;primaryKey"`
	Name             string    `json:"name"`
	Token            string    `json:"token" gorm:"unique"`
	Webhook          string    `json:"webhook"`
	RabbitmqEnable   string    `json:"rabbitmqEnable"`
	WebSocketEnable  string    `json:"websocketEnable"`
	NatsEnable       string    `json:"natsEnable"`
	Jid              string    `json:"jid" gorm:"column:jid"`
	Qrcode           string    `json:"qrcode" gorm:"type:text"`
	Connected        bool      `json:"connected"`
	Expiration       int64     `json:"expiration"`
	DisconnectReason string    `json:"disconnect_reason"`
	Events           string    `json:"events"`
	OsName           string    `json:"os_name"`
	Proxy            string    `json:"proxy"`
	ClientName       string    `json:"client_name"`
	CreatedAt        time.Time `json:"createdAt" gorm:"autoCreateTime"`

	// Advanced Settings
	AlwaysOnline  bool   `json:"alwaysOnline" gorm:"default:false"`
	RejectCall    bool   `json:"rejectCall" gorm:"default:false"`
	MsgRejectCall string `json:"msgRejectCall" gorm:"default:''"`
	ReadMessages  bool   `json:"readMessages" gorm:"default:false"`
	IgnoreGroups  bool   `json:"ignoreGroups" gorm:"default:false"`
	IgnoreStatus  bool   `json:"ignoreStatus" gorm:"default:false"`

	// Chatwoot Integration
	ChatwootEnabled                   bool   `json:"chatwootEnabled" gorm:"default:false"`
	ChatwootAccountId                 string `json:"chatwootAccountId" gorm:"default:''"`
	ChatwootToken                     string `json:"chatwootToken" gorm:"default:''"`
	ChatwootUrl                       string `json:"chatwootUrl" gorm:"default:''"`
	ChatwootInboxId                   int    `json:"chatwootInboxId" gorm:"default:0"`
	ChatwootAutoCreate                bool   `json:"chatwootAutoCreate" gorm:"default:false"`
	ChatwootSignMsg                   bool   `json:"chatwootSignMsg" gorm:"default:false"`
	ChatwootReopenConversation        bool   `json:"chatwootReopenConversation" gorm:"default:false"`
	ChatwootConversationPending       bool   `json:"chatwootConversationPending" gorm:"default:false"`
	ChatwootImportContacts            bool   `json:"chatwootImportContacts" gorm:"default:false"`
	ChatwootImportMessages            bool   `json:"chatwootImportMessages" gorm:"default:false"`
	ChatwootDaysLimitImportMessages   int    `json:"chatwootDaysLimitImportMessages" gorm:"default:0"`
}

// AdvancedSettings representa as configurações avançadas de uma instância
type AdvancedSettings struct {
	AlwaysOnline  bool   `json:"alwaysOnline"`
	RejectCall    bool   `json:"rejectCall"`
	MsgRejectCall string `json:"msgRejectCall"`
	ReadMessages  bool   `json:"readMessages"`
	IgnoreGroups  bool   `json:"ignoreGroups"`
	IgnoreStatus  bool   `json:"ignoreStatus"`
}

func (m *Instance) BeforeCreate(tx *gorm.DB) (err error) {
	if m.Id == "" {
		m.Id = uuid.New().String()
	}
	return
}
