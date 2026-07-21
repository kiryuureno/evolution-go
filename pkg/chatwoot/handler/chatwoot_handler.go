package chatwoot_handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	chatwoot_model "github.com/evolution-foundation/evolution-go/pkg/chatwoot/model"
	chatwoot_service "github.com/evolution-foundation/evolution-go/pkg/chatwoot/service"
)

type ChatwootHandler interface {
	SetChatwootConfig(ctx *gin.Context)
	FindChatwootConfig(ctx *gin.Context)
	DeleteChatwootConfig(ctx *gin.Context)
	WebhookHandler(ctx *gin.Context)
	SyncContactsHandler(ctx *gin.Context)
	SyncMessagesHandler(ctx *gin.Context)
}

type chatwootHandler struct {
	chatwootService chatwoot_service.ChatwootService
}

func NewChatwootHandler(chatwootService chatwoot_service.ChatwootService) ChatwootHandler {
	return &chatwootHandler{
		chatwootService: chatwootService,
	}
}

// SetChatwootConfig godoc
// @Summary Configurar integração do Chatwoot
// @Description Define e habilita as configurações do Chatwoot para uma instância
// @Tags chatwoot
// @Accept json
// @Produce json
// @Param instanceId path string true "ID ou Nome da Instância"
// @Param body body chatwoot_model.ChatwootDTO true "Payload de Configuração do Chatwoot"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /chatwoot/set/{instanceId} [post]
func (h *chatwootHandler) SetChatwootConfig(ctx *gin.Context) {
	instanceId := ctx.Param("instanceId")
	if instanceId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "instanceId é obrigatório"})
		return
	}

	var dto chatwoot_model.ChatwootDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil || ctx.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := ctx.Request.Host
	if forwardedHost := ctx.GetHeader("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}
	serverUrl := fmt.Sprintf("%s://%s", scheme, host)

	result, err := h.chatwootService.SetChatwootConfig(instanceId, dto, serverUrl)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"chatwoot": result,
	})
}

// FindChatwootConfig godoc
// @Summary Buscar configurações do Chatwoot
// @Description Retorna as configurações ativas do Chatwoot para a instância
// @Tags chatwoot
// @Produce json
// @Param instanceId path string true "ID ou Nome da Instância"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /chatwoot/find/{instanceId} [get]
func (h *chatwootHandler) FindChatwootConfig(ctx *gin.Context) {
	instanceId := ctx.Param("instanceId")
	if instanceId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "instanceId é obrigatório"})
		return
	}

	result, err := h.chatwootService.GetChatwootConfig(instanceId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"chatwoot": result,
	})
}

// DeleteChatwootConfig godoc
// @Summary Desativar integração do Chatwoot
// @Description Remove e desabilita as configurações do Chatwoot para a instância
// @Tags chatwoot
// @Produce json
// @Param instanceId path string true "ID ou Nome da Instância"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /chatwoot/delete/{instanceId} [delete]
func (h *chatwootHandler) DeleteChatwootConfig(ctx *gin.Context) {
	instanceId := ctx.Param("instanceId")
	if instanceId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "instanceId é obrigatório"})
		return
	}

	err := h.chatwootService.DeleteChatwootConfig(instanceId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Configurações do Chatwoot removidas com sucesso",
	})
}

// WebhookHandler godoc
// @Summary Webhook de recebimento do Chatwoot
// @Description Recebe webhooks disparados pelo Chatwoot quando um agente envia mensagens
// @Tags chatwoot
// @Accept json
// @Produce json
// @Param instanceId path string true "ID ou Nome da Instância"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /chatwoot/webhook/{instanceId} [post]
func (h *chatwootHandler) WebhookHandler(ctx *gin.Context) {
	instanceId := ctx.Param("instanceId")
	if instanceId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "instanceId é obrigatório"})
		return
	}

	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "erro ao ler corpo da requisição"})
		return
	}

	err = h.chatwootService.ProcessChatwootWebhook(instanceId, bodyBytes)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "ignored_or_error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "webhook processado com sucesso",
	})
}

// SyncContactsHandler godoc
// @Summary Sincronizar contatos com o Chatwoot
// @Description Exporta os contatos da instância para o Chatwoot
// @Tags chatwoot
// @Produce json
// @Param instanceId path string true "ID ou Nome da Instância"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /chatwoot/syncContacts/{instanceId} [post]
func (h *chatwootHandler) SyncContactsHandler(ctx *gin.Context) {
	instanceId := ctx.Param("instanceId")
	if instanceId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "instanceId é obrigatório"})
		return
	}

	err := h.chatwootService.SyncContacts(instanceId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Sincronização de contatos iniciada com sucesso",
	})
}

// SyncMessagesHandler godoc
// @Summary Sincronizar mensagens com o Chatwoot
// @Description Exporta as mensagens recentes da instância para o Chatwoot
// @Tags chatwoot
// @Produce json
// @Param instanceId path string true "ID ou Nome da Instância"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /chatwoot/syncMessages/{instanceId} [post]
func (h *chatwootHandler) SyncMessagesHandler(ctx *gin.Context) {
	instanceId := ctx.Param("instanceId")
	if instanceId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "instanceId é obrigatório"})
		return
	}

	err := h.chatwootService.SyncMessages(instanceId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Sincronização de mensagens iniciada com sucesso",
	})
}
