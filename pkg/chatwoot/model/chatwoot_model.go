package chatwoot_model

// ChatwootDTO representa o payload de configuração enviado e retornado pela API
type ChatwootDTO struct {
	Enabled                   bool   `json:"enabled"`
	AccountId                 string `json:"accountId"`
	Token                     string `json:"token"`
	Url                       string `json:"url"`
	InboxId                   int    `json:"inboxId"`
	AutoCreate                bool   `json:"autoCreate"`
	SignMsg                   bool   `json:"signMsg"`
	ReopenConversation        bool   `json:"reopenConversation"`
	ConversationPending       bool   `json:"conversationPending"`
	ImportContacts            bool   `json:"importContacts"`
	ImportMessages            bool   `json:"importMessages"`
	DaysLimitImportMessages   int    `json:"daysLimitImportMessages"`
}

// Structs para criação de Inbox no Chatwoot
type ChatwootCreateInboxChannel struct {
	Type       string `json:"type"`
	WebhookUrl string `json:"webhook_url"`
}

type ChatwootCreateInboxReq struct {
	Name    string                     `json:"name"`
	Channel ChatwootCreateInboxChannel `json:"channel"`
}

type ChatwootCreateInboxResp struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ChannelID int    `json:"channel_id"`
}

// Structs para Contatos no Chatwoot
type ChatwootContactReq struct {
	InboxId     int    `json:"inbox_id,omitempty"`
	Name        string `json:"name"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Identifier  string `json:"identifier,omitempty"`
}

type ChatwootContact struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Identifier  string `json:"identifier"`
	Email       string `json:"email"`
}

type ChatwootContactResp struct {
	ID      int `json:"id"`
	Payload struct {
		Contact ChatwootContact `json:"contact"`
	} `json:"payload"`
}

type ChatwootSearchContactResp struct {
	Payload []ChatwootContact `json:"payload"`
}

// Structs para Conversas no Chatwoot
type ChatwootConversationReq struct {
	SourceId  string `json:"source_id"`
	InboxId   int    `json:"inbox_id"`
	ContactId int    `json:"contact_id"`
	Status    string `json:"status,omitempty"`
}

type ChatwootConversationResp struct {
	ID        int    `json:"id"`
	AccountID int    `json:"account_id"`
	InboxID   int    `json:"inbox_id"`
	Status    string `json:"status"`
}

type ChatwootConversationSearchResp struct {
	Payload []ChatwootConversationResp `json:"payload"`
}

// Structs para Envio de Mensagens no Chatwoot
type ChatwootMessageReq struct {
	Content           string                 `json:"content"`
	MessageType       string                 `json:"message_type"` // "incoming" ou "outgoing"
	Private           bool                   `json:"private"`
	SourceID          string                 `json:"source_id,omitempty"`
	ContentAttributes map[string]interface{} `json:"content_attributes,omitempty"`
}

// Structs para Webhook do Chatwoot (Chatwoot -> EvolutionGO)
type ChatwootWebhookPayload struct {
	ID                int                         `json:"id,omitempty"`
	Event             string                      `json:"event"`
	MessageType       string                      `json:"message_type"`
	Content           string                      `json:"content"`
	Private           bool                        `json:"private"`
	SourceID          string                      `json:"source_id,omitempty"`
	ContentAttributes map[string]interface{}    `json:"content_attributes,omitempty"`
	InReplyTo         interface{}                 `json:"in_reply_to,omitempty"`
	Sender            ChatwootWebhookSender       `json:"sender"`
	Contact           ChatwootWebhookContact      `json:"contact"`
	Conversation      ChatwootWebhookConversation `json:"conversation"`
	Attachments       []ChatwootWebhookAttachment `json:"attachments"`
}

type ChatwootWebhookSender struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"` // "user" (agente), "contact", "bot"
}

type ChatwootWebhookContact struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Identifier  string `json:"identifier"`
}

type ChatwootWebhookConversation struct {
	ID      int    `json:"id"`
	InboxID int    `json:"inbox_id"`
	Status  string `json:"status"`
}

type ChatwootWebhookAttachment struct {
	ID       int    `json:"id"`
	FileType string `json:"file_type"` // "image", "audio", "video", "file"
	DataUrl  string `json:"data_url"`
}
