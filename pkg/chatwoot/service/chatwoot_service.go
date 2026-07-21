package chatwoot_service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	chatwoot_model "github.com/evolution-foundation/evolution-go/pkg/chatwoot/model"
	instance_model "github.com/evolution-foundation/evolution-go/pkg/instance/model"
	instance_repository "github.com/evolution-foundation/evolution-go/pkg/instance/repository"
	logger_wrapper "github.com/evolution-foundation/evolution-go/pkg/logger"
	send_service "github.com/evolution-foundation/evolution-go/pkg/sendMessage/service"
)

type ChatwootService interface {
	SetChatwootConfig(instanceId string, dto chatwoot_model.ChatwootDTO, serverUrl string) (*chatwoot_model.ChatwootDTO, error)
	GetChatwootConfig(instanceId string) (*chatwoot_model.ChatwootDTO, error)
	DeleteChatwootConfig(instanceId string) error
	ProcessWhatsAppEvent(instance *instance_model.Instance, payloadMap map[string]interface{}) error
	ProcessChatwootWebhook(instanceId string, payload []byte) error
	SyncContacts(instanceId string) error
	SyncMessages(instanceId string) error
}

type chatwootService struct {
	instanceRepo  instance_repository.InstanceRepository
	sendService   send_service.SendService
	loggerWrapper *logger_wrapper.LoggerManager
}

func NewChatwootService(
	instanceRepo instance_repository.InstanceRepository,
	sendService send_service.SendService,
	loggerWrapper *logger_wrapper.LoggerManager,
) ChatwootService {
	return &chatwootService{
		instanceRepo:  instanceRepo,
		sendService:   sendService,
		loggerWrapper: loggerWrapper,
	}
}

func (s *chatwootService) SetChatwootConfig(instanceId string, dto chatwoot_model.ChatwootDTO, serverUrl string) (*chatwoot_model.ChatwootDTO, error) {
	instance, err := s.getInstance(instanceId)
	if err != nil {
		return nil, errors.New("instância não encontrada")
	}

	dto.Url = strings.TrimSuffix(dto.Url, "/")

	instance.ChatwootEnabled = dto.Enabled
	instance.ChatwootAccountId = dto.AccountId
	instance.ChatwootToken = dto.Token
	instance.ChatwootUrl = dto.Url
	instance.ChatwootInboxId = dto.InboxId
	instance.ChatwootAutoCreate = dto.AutoCreate
	instance.ChatwootSignMsg = dto.SignMsg
	instance.ChatwootReopenConversation = dto.ReopenConversation
	instance.ChatwootConversationPending = dto.ConversationPending
	instance.ChatwootImportContacts = dto.ImportContacts
	instance.ChatwootImportMessages = dto.ImportMessages
	instance.ChatwootDaysLimitImportMessages = dto.DaysLimitImportMessages

	// Se autoCreate estiver ativado e o inboxId for 0, criar caixa de entrada no Chatwoot
	if dto.Enabled && dto.AutoCreate && dto.InboxId == 0 && dto.Url != "" && dto.Token != "" && dto.AccountId != "" {
		webhookUrl := fmt.Sprintf("%s/chatwoot/webhook/%s", strings.TrimSuffix(serverUrl, "/"), instance.Name)
		if instance.Name == "" {
			webhookUrl = fmt.Sprintf("%s/chatwoot/webhook/%s", strings.TrimSuffix(serverUrl, "/"), instance.Id)
		}

		inboxId, err := s.autoCreateInbox(instance, webhookUrl)
		if err == nil && inboxId > 0 {
			instance.ChatwootInboxId = inboxId
			dto.InboxId = inboxId
		} else if err != nil {
			s.loggerWrapper.GetLogger(instanceId).LogError("[%s] Erro ao criar caixa de entrada no Chatwoot: %v", instanceId, err)
		}
	}

	err = s.instanceRepo.Update(instance)
	if err != nil {
		return nil, err
	}

	// Se importContacts estiver habilitado, disparar sincronização em background
	if dto.Enabled && dto.ImportContacts {
		go s.SyncContacts(instance.Id)
	}

	return &dto, nil
}

func (s *chatwootService) GetChatwootConfig(instanceId string) (*chatwoot_model.ChatwootDTO, error) {
	instance, err := s.getInstance(instanceId)
	if err != nil {
		return nil, errors.New("instância não encontrada")
	}

	return &chatwoot_model.ChatwootDTO{
		Enabled:                 instance.ChatwootEnabled,
		AccountId:               instance.ChatwootAccountId,
		Token:                   instance.ChatwootToken,
		Url:                     instance.ChatwootUrl,
		InboxId:                 instance.ChatwootInboxId,
		AutoCreate:              instance.ChatwootAutoCreate,
		SignMsg:                 instance.ChatwootSignMsg,
		ReopenConversation:      instance.ChatwootReopenConversation,
		ConversationPending:     instance.ChatwootConversationPending,
		ImportContacts:          instance.ChatwootImportContacts,
		ImportMessages:          instance.ChatwootImportMessages,
		DaysLimitImportMessages: instance.ChatwootDaysLimitImportMessages,
	}, nil
}

func (s *chatwootService) DeleteChatwootConfig(instanceId string) error {
	instance, err := s.getInstance(instanceId)
	if err != nil {
		return errors.New("instância não encontrada")
	}

	instance.ChatwootEnabled = false
	instance.ChatwootAccountId = ""
	instance.ChatwootToken = ""
	instance.ChatwootUrl = ""
	instance.ChatwootInboxId = 0
	instance.ChatwootAutoCreate = false
	instance.ChatwootSignMsg = false
	instance.ChatwootReopenConversation = false
	instance.ChatwootConversationPending = false
	instance.ChatwootImportContacts = false
	instance.ChatwootImportMessages = false
	instance.ChatwootDaysLimitImportMessages = 0

	return s.instanceRepo.Update(instance)
}

func (s *chatwootService) autoCreateInbox(instance *instance_model.Instance, webhookUrl string) (int, error) {
	inboxName := fmt.Sprintf("WhatsApp - %s", instance.Name)
	if instance.Name == "" {
		inboxName = fmt.Sprintf("WhatsApp - %s", instance.Id)
	}

	reqBody := chatwoot_model.ChatwootCreateInboxReq{
		Name: inboxName,
		Channel: chatwoot_model.ChatwootCreateInboxChannel{
			Type:       "api",
			WebhookUrl: webhookUrl,
		},
	}

	bodyBytes, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("%s/api/v1/accounts/%s/inboxes", instance.ChatwootUrl, instance.ChatwootAccountId)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_access_token", instance.ChatwootToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("chatwoot API retornou status %d: %s", resp.StatusCode, string(respBody))
	}

	var res chatwoot_model.ChatwootCreateInboxResp
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return 0, err
	}

	return res.ID, nil
}

func (s *chatwootService) ProcessWhatsAppEvent(instance *instance_model.Instance, payloadMap map[string]interface{}) error {
	latestInstance, err := s.instanceRepo.GetInstanceByID(instance.Id)
	if err == nil && latestInstance != nil {
		instance = latestInstance
	}

	if !instance.ChatwootEnabled || instance.ChatwootUrl == "" || instance.ChatwootToken == "" || instance.ChatwootAccountId == "" || instance.ChatwootInboxId == 0 {
		return nil
	}

	event, _ := payloadMap["event"].(string)
	if event != "Message" && event != "SendMessage" && event != "messages.upsert" {
		return nil
	}

	data, ok := payloadMap["data"].(map[string]interface{})
	if !ok {
		return nil
	}

	key, _ := data["key"].(map[string]interface{})
	if key == nil {
		// Tenta extrair key do formato Info se existir
		if info, ok := data["Info"].(map[string]interface{}); ok {
			remoteJid, _ := info["Chat"].(string)
			fromMe, _ := info["IsFromMe"].(bool)
			id, _ := info["ID"].(string)
			key = map[string]interface{}{
				"remoteJid": remoteJid,
				"fromMe":    fromMe,
				"id":        id,
			}
		} else {
			return nil
		}
	}

	remoteJid, _ := key["remoteJid"].(string)

	// Resolver JID real de telefone caso remoteJid seja um LID (ex: 260313823891664@lid -> 556399400537@s.whatsapp.net)
	if chat, ok := data["Chat"].(string); ok && chat != "" && !strings.HasSuffix(chat, "@lid") {
		remoteJid = chat
	} else if sender, ok := data["Sender"].(string); ok && sender != "" && !strings.HasSuffix(sender, "@lid") {
		remoteJid = sender
	} else if senderAlt, ok := data["SenderAlt"].(string); ok && senderAlt != "" && !strings.HasSuffix(senderAlt, "@lid") {
		remoteJid = senderAlt
	}

	if remoteJid == "" || strings.HasSuffix(remoteJid, "@broadcast") {
		return nil
	}

	if instance.IgnoreGroups && strings.HasSuffix(remoteJid, "@g.us") {
		return nil
	}

	fromMe, _ := key["fromMe"].(bool)
	pushName, _ := data["pushName"].(string)
	if pushName == "" {
		if pn, ok := data["PushName"].(string); ok && pn != "" {
			pushName = pn
		} else {
			pushName = strings.Split(remoteJid, "@")[0]
		}
	}

	phoneNumber := "+" + strings.Split(remoteJid, "@")[0]

	// 1. Criar ou Buscar Contato no Chatwoot
	contactId, err := s.findOrCreateContact(instance, pushName, phoneNumber, remoteJid)
	if err != nil {
		s.loggerWrapper.GetLogger(instance.Id).LogError("[%s] Erro ao buscar/criar contato no Chatwoot: %v", instance.Id, err)
		return err
	}

	// 2. Criar ou Buscar Conversa no Chatwoot
	conversationId, err := s.findOrCreateConversation(instance, contactId, remoteJid)
	if err != nil {
		s.loggerWrapper.GetLogger(instance.Id).LogError("[%s] Erro ao buscar/criar conversa no Chatwoot: %v", instance.Id, err)
		return err
	}

	// 3. Extrair Conteúdo da Mensagem
	content := extractTextFromMessageData(data)
	if content == "" {
		content = "[Mídia / Mensagem não suportada]"
	}

	msgType := "incoming"
	if fromMe {
		msgType = "outgoing"
	}

	msgId := ""
	if id, ok := key["id"].(string); ok {
		msgId = id
	}

	// 4. Enviar Mensagem para o Chatwoot
	return s.postMessageToChatwoot(instance, conversationId, content, msgType, msgId)
}

func (s *chatwootService) findOrCreateContact(instance *instance_model.Instance, name, phoneNumber, identifier string) (int, error) {
	searchUrl := fmt.Sprintf("%s/api/v1/accounts/%s/contacts/search?q=%s", instance.ChatwootUrl, instance.ChatwootAccountId, phoneNumber)
	req, _ := http.NewRequest("GET", searchUrl, nil)
	req.Header.Set("api_access_token", instance.ChatwootToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		var searchRes chatwoot_model.ChatwootSearchContactResp
		if json.NewDecoder(resp.Body).Decode(&searchRes) == nil && len(searchRes.Payload) > 0 {
			resp.Body.Close()
			return searchRes.Payload[0].ID, nil
		}
		resp.Body.Close()
	}

	createUrl := fmt.Sprintf("%s/api/v1/accounts/%s/contacts", instance.ChatwootUrl, instance.ChatwootAccountId)
	contactReq := chatwoot_model.ChatwootContactReq{
		InboxId:     instance.ChatwootInboxId,
		Name:        name,
		PhoneNumber: phoneNumber,
		Identifier:  identifier,
	}

	bodyBytes, _ := json.Marshal(contactReq)
	req, err = http.NewRequest("POST", createUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_access_token", instance.ChatwootToken)

	resp, err = client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("erro ao criar contato status %d: %s", resp.StatusCode, string(b))
	}

	var contactResp chatwoot_model.ChatwootContactResp
	if err := json.NewDecoder(resp.Body).Decode(&contactResp); err != nil {
		return 0, err
	}

	if contactResp.Payload.Contact.ID > 0 {
		return contactResp.Payload.Contact.ID, nil
	}
	return contactResp.ID, nil
}

func (s *chatwootService) findOrCreateConversation(instance *instance_model.Instance, contactId int, sourceId string) (int, error) {
	// 1. Tentar buscar conversas ativas do contato no Chatwoot para agrupamento na mesma thread
	getConvUrl := fmt.Sprintf("%s/api/v1/accounts/%s/contacts/%d/conversations", instance.ChatwootUrl, instance.ChatwootAccountId, contactId)
	req, err := http.NewRequest("GET", getConvUrl, nil)
	if err == nil {
		req.Header.Set("api_access_token", instance.ChatwootToken)
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err == nil {
			if resp.StatusCode == 200 {
				var convSearchRes struct {
					Payload []chatwoot_model.ChatwootConversationResp `json:"payload"`
				}
				if json.NewDecoder(resp.Body).Decode(&convSearchRes) == nil {
					resp.Body.Close()
					for _, conv := range convSearchRes.Payload {
						if conv.InboxID == instance.ChatwootInboxId {
							if conv.Status == "open" || conv.Status == "pending" || conv.Status == "snoozed" {
								return conv.ID, nil
							}
							if conv.Status == "resolved" && instance.ChatwootReopenConversation {
								_ = s.toggleConversationStatus(instance, conv.ID, "open")
								return conv.ID, nil
							}
						}
					}
				} else {
					resp.Body.Close()
				}
			} else {
				resp.Body.Close()
			}
		}
	}

	// 2. Se não houver conversa aberta para este contato na caixa de entrada, criar uma nova
	url := fmt.Sprintf("%s/api/v1/accounts/%s/conversations", instance.ChatwootUrl, instance.ChatwootAccountId)
	convReq := chatwoot_model.ChatwootConversationReq{
		SourceId:  sourceId,
		InboxId:   instance.ChatwootInboxId,
		ContactId: contactId,
	}
	if instance.ChatwootConversationPending {
		convReq.Status = "pending"
	}

	bodyBytes, _ := json.Marshal(convReq)
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_access_token", instance.ChatwootToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var convResp chatwoot_model.ChatwootConversationResp
	if err := json.NewDecoder(resp.Body).Decode(&convResp); err != nil {
		return 0, err
	}

	return convResp.ID, nil
}

func (s *chatwootService) toggleConversationStatus(instance *instance_model.Instance, conversationId int, status string) error {
	url := fmt.Sprintf("%s/api/v1/accounts/%s/conversations/%d/toggle_status", instance.ChatwootUrl, instance.ChatwootAccountId, conversationId)
	bodyMap := map[string]string{"status": status}
	bodyBytes, _ := json.Marshal(bodyMap)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_access_token", instance.ChatwootToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (s *chatwootService) postMessageToChatwoot(instance *instance_model.Instance, conversationId int, content, messageType, sourceId string) error {
	url := fmt.Sprintf("%s/api/v1/accounts/%s/conversations/%d/messages", instance.ChatwootUrl, instance.ChatwootAccountId, conversationId)
	msgReq := chatwoot_model.ChatwootMessageReq{
		Content:     content,
		MessageType: messageType,
		Private:     false,
		SourceID:    sourceId,
	}

	bodyBytes, _ := json.Marshal(msgReq)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_access_token", instance.ChatwootToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao postar mensagem no chatwoot status %d: %s", resp.StatusCode, string(b))
	}

	return nil
}

func (s *chatwootService) ProcessChatwootWebhook(instanceIdOrName string, payload []byte) error {
	var webhookPayload chatwoot_model.ChatwootWebhookPayload
	err := json.Unmarshal(payload, &webhookPayload)
	if err != nil {
		return fmt.Errorf("payload de webhook inválido: %v", err)
	}

	// Filtra apenas mensagens de saída enviadas por agentes
	if webhookPayload.Event != "message_created" || webhookPayload.MessageType != "outgoing" || webhookPayload.Private {
		return nil
	}

	instance, err := s.getInstance(instanceIdOrName)
	if err != nil {
		return fmt.Errorf("instância %s não encontrada", instanceIdOrName)
	}

	if !instance.Connected {
		return errors.New("instância do WhatsApp desconectada")
	}

	phoneNumber := webhookPayload.Contact.PhoneNumber
	if phoneNumber == "" {
		phoneNumber = webhookPayload.Contact.Identifier
	}

	phoneNumber = strings.TrimPrefix(phoneNumber, "+")
	if phoneNumber == "" {
		return errors.New("número do destinatário não encontrado no payload do Chatwoot")
	}

	text := webhookPayload.Content
	if instance.ChatwootSignMsg && webhookPayload.Sender.Name != "" && webhookPayload.Sender.Type == "user" {
		text = fmt.Sprintf("*%s*: %s", webhookPayload.Sender.Name, text)
	}

	// Se houver anexos de mídia (imagem, vídeo, áudio, arquivo)
	if len(webhookPayload.Attachments) > 0 {
		for _, att := range webhookPayload.Attachments {
			mediaData := &send_service.MediaStruct{
				Number:   phoneNumber,
				Url:      att.DataUrl,
				Caption:  text,
				Type:     getMediaType(att.FileType),
				Filename: fmt.Sprintf("attachment_%d", att.ID),
			}
			_, err := s.sendService.SendMediaUrl(mediaData, instance)
			if err != nil {
				s.loggerWrapper.GetLogger(instance.Id).LogError("[%s] Erro ao enviar mídia via WhatsApp: %v", instance.Id, err)
			}
		}
		return nil
	}

	// Envio de mensagem de texto simples
	if text != "" {
		textData := &send_service.TextStruct{
			Number: phoneNumber,
			Text:   text,
		}
		_, err := s.sendService.SendText(textData, instance)
		if err != nil {
			s.loggerWrapper.GetLogger(instance.Id).LogError("[%s] Erro ao enviar texto via WhatsApp: %v", instance.Id, err)
			return err
		}
	}

	return nil
}

func (s *chatwootService) SyncContacts(instanceId string) error {
	instance, err := s.getInstance(instanceId)
	if err != nil {
		return errors.New("instância não encontrada")
	}

	if !instance.ChatwootEnabled || instance.ChatwootUrl == "" || instance.ChatwootToken == "" {
		return errors.New("chatwoot não está ativado para esta instância")
	}

	s.loggerWrapper.GetLogger(instance.Id).LogInfo("[%s] Iniciando sincronização de contatos com Chatwoot...", instance.Id)
	return nil
}

func (s *chatwootService) SyncMessages(instanceId string) error {
	instance, err := s.getInstance(instanceId)
	if err != nil {
		return errors.New("instância não encontrada")
	}

	if !instance.ChatwootEnabled || instance.ChatwootUrl == "" || instance.ChatwootToken == "" {
		return errors.New("chatwoot não está ativado para esta instância")
	}

	s.loggerWrapper.GetLogger(instance.Id).LogInfo("[%s] Iniciando sincronização de mensagens com Chatwoot...", instance.Id)
	return nil
}

func (s *chatwootService) getInstance(instanceIdOrName string) (*instance_model.Instance, error) {
	instance, err := s.instanceRepo.GetInstanceByID(instanceIdOrName)
	if err != nil || instance == nil {
		instance, err = s.instanceRepo.GetInstanceByName(instanceIdOrName)
	}
	return instance, err
}

func extractTextFromMessageData(data map[string]interface{}) string {
	msg, ok := data["message"].(map[string]interface{})
	if !ok || msg == nil {
		msg, ok = data["Message"].(map[string]interface{})
		if !ok || msg == nil {
			if text, ok := data["text"].(string); ok && text != "" {
				return text
			}
			return ""
		}
	}

	if conversation, ok := msg["conversation"].(string); ok && conversation != "" {
		return conversation
	}

	if extendedMsg, ok := msg["extendedTextMessage"].(map[string]interface{}); ok {
		if text, ok := extendedMsg["text"].(string); ok && text != "" {
			return text
		}
	}

	if imageMsg, ok := msg["imageMessage"].(map[string]interface{}); ok {
		if caption, ok := imageMsg["caption"].(string); ok && caption != "" {
			return fmt.Sprintf("[Imagem]: %s", caption)
		}
		return "[Imagem]"
	}

	if videoMsg, ok := msg["videoMessage"].(map[string]interface{}); ok {
		if caption, ok := videoMsg["caption"].(string); ok && caption != "" {
			return fmt.Sprintf("[Vídeo]: %s", caption)
		}
		return "[Vídeo]"
	}

	if audioMsg, ok := msg["audioMessage"].(map[string]interface{}); ok {
		_ = audioMsg
		return "[Áudio / Nota de Voz]"
	}

	if documentMsg, ok := msg["documentMessage"].(map[string]interface{}); ok {
		if title, ok := documentMsg["title"].(string); ok && title != "" {
			if caption, ok := documentMsg["caption"].(string); ok && caption != "" {
				return fmt.Sprintf("[Documento - %s]: %s", title, caption)
			}
			return fmt.Sprintf("[Documento]: %s", title)
		}
		if caption, ok := documentMsg["caption"].(string); ok && caption != "" {
			return fmt.Sprintf("[Documento]: %s", caption)
		}
		return "[Documento]"
	}

	if stickerMsg, ok := msg["stickerMessage"].(map[string]interface{}); ok {
		_ = stickerMsg
		return "[Sticker / Figurinha]"
	}

	if contactMsg, ok := msg["contactMessage"].(map[string]interface{}); ok {
		if displayName, ok := contactMsg["displayName"].(string); ok && displayName != "" {
			return fmt.Sprintf("[Contato]: %s", displayName)
		}
		return "[Contato]"
	}

	if locationMsg, ok := msg["locationMessage"].(map[string]interface{}); ok {
		if name, ok := locationMsg["name"].(string); ok && name != "" {
			return fmt.Sprintf("[Localização]: %s", name)
		}
		return "[Localização]"
	}

	if text, ok := msg["text"].(string); ok && text != "" {
		return text
	}

	return ""
}

func getMediaType(fileType string) string {
	switch strings.ToLower(fileType) {
	case "image":
		return "image"
	case "audio":
		return "audio"
	case "video":
		return "video"
	default:
		return "document"
	}
}
