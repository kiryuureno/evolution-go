var m, b, Jt, Fe, Es, it;


function CwIcon(t){
return m.jsx("svg",{
xmlns:"http://www.w3.org/2000/svg",width:"24",height:"24",viewBox:"0 0 24 24",fill:"none",stroke:"currentColor",strokeWidth:"2",strokeLinecap:"round",strokeLinejoin:"round",...t,children:[m.jsx("path",{
d:"M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"}
),m.jsx("path",{
d:"M8 10h.01"}
),m.jsx("path",{
d:"M12 10h.01"}
),m.jsx("path",{
d:"M16 10h.01"}
)].filter(Boolean)}
)}

function cwL({
open:t,onClose:n,instance:r}
){
const[loading,setLoading]=b.useState(!1);
const[fetching,setFetching]=b.useState(!1);
const[syncingContacts,setSyncingContacts]=b.useState(!1);
const[syncingMessages,setSyncingMessages]=b.useState(!1);
const defaultForm={
enabled:!1,url:"",token:"",accountId:"",inboxId:0,autoCreate:!0,signMsg:!1,reopenConversation:!1,conversationPending:!1,importContacts:!1,importMessages:!1,daysLimitImportMessages:0}
;
const[formData,setFormData]=b.useState(defaultForm);
const fetchConfig=b.useCallback(async(showToast=!1)=>{
if(!r)return;
const targetId=r.instanceName||r.id;
setFetching(!0);
try{
const resp=await Jt.get(`/chatwoot/find/${
targetId}
`);
if(resp&&resp.data&&resp.data.chatwoot){
const cw=resp.data.chatwoot;
setFormData({
enabled:!!cw.enabled,url:cw.url||"",token:cw.token||"",accountId:cw.accountId?String(cw.accountId):"",inboxId:cw.inboxId||0,autoCreate:cw.autoCreate!==undefined?!!cw.autoCreate:!0,signMsg:!!cw.signMsg,reopenConversation:!!cw.reopenConversation,conversationPending:!!cw.conversationPending,importContacts:!!cw.importContacts,importMessages:!!cw.importMessages,daysLimitImportMessages:cw.daysLimitImportMessages||0}
);
if(showToast)Fe.success("Configurações do Chatwoot carregadas!")}
}
catch(err){
if(showToast)Fe.info("Nenhuma configuração salva encontrada. Preencha os campos para salvar.")}
finally{
setFetching(!1)}
}
,[r]);
b.useEffect(()=>{
if(t&&r){
setFormData(defaultForm);
fetchConfig(!1)}
}
,[t,r,fetchConfig]);
if(!t||!r)return null;
const targetId=r.instanceName||r.id;
const handleSave=async(e)=>{
if(e)e.preventDefault();
setLoading(!0);
try{
const payload={
enabled:!!formData.enabled,url:(formData.url||"").trim(),token:(formData.token||"").trim(),accountId:(formData.accountId||"").trim(),inboxId:parseInt(formData.inboxId,10)||0,autoCreate:!!formData.autoCreate,signMsg:!!formData.signMsg,reopenConversation:!!formData.reopenConversation,conversationPending:!!formData.conversationPending,importContacts:!!formData.importContacts,importMessages:!!formData.importMessages,daysLimitImportMessages:parseInt(formData.daysLimitImportMessages,10)||0}
;
const resp=await Jt.post(`/chatwoot/set/${
targetId}
`,payload);
if(resp&&resp.data&&resp.data.chatwoot){
const cw=resp.data.chatwoot;
setFormData({
enabled:!!cw.enabled,url:cw.url||"",token:cw.token||"",accountId:cw.accountId?String(cw.accountId):"",inboxId:cw.inboxId||0,autoCreate:!!cw.autoCreate,signMsg:!!cw.signMsg,reopenConversation:!!cw.reopenConversation,conversationPending:!!cw.conversationPending,importContacts:!!cw.importContacts,importMessages:!!cw.importMessages,daysLimitImportMessages:cw.daysLimitImportMessages||0}
)}
Fe.success("Configurações do Chatwoot salvas com sucesso!")}
catch(err){
console.error("Erro ao salvar Chatwoot:",err);
const msg=(err.response&&err.response.data&&err.response.data.error)||err.message||"Erro ao salvar configurações";
Fe.error(msg)}
finally{
setLoading(!1)}
}
;
const handleDelete=async()=>{
if(!window.confirm("Deseja realmente remover as configurações do Chatwoot desta instância?"))return;
setLoading(!0);
try{
await Jt.delete(`/chatwoot/delete/${
targetId}
`);
setFormData(defaultForm);
Fe.success("Configuração do Chatwoot removida!")}
catch(err){
console.error("Erro ao remover Chatwoot:",err);
const msg=(err.response&&err.response.data&&err.response.data.error)||err.message||"Erro ao remover";
Fe.error(msg)}
finally{
setLoading(!1)}
}
;
const handleSyncContacts=async()=>{
setSyncingContacts(!0);
try{
await Jt.post(`/chatwoot/syncContacts/${
targetId}
`);
Fe.success("Sincronização de contatos iniciada!")}
catch(err){
const msg=(err.response&&err.response.data&&err.response.data.error)||err.message||"Erro ao sincronizar contatos";
Fe.error(msg)}
finally{
setSyncingContacts(!1)}
}
;
const handleSyncMessages=async()=>{
setSyncingMessages(!0);
try{
await Jt.post(`/chatwoot/syncMessages/${
targetId}
`);
Fe.success("Sincronização de mensagens iniciada!")}
catch(err){
const msg=(err.response&&err.response.data&&err.response.data.error)||err.message||"Erro ao sincronizar mensagens";
Fe.error(msg)}
finally{
setSyncingMessages(!1)}
}
;
const handleChange=(field,val)=>{
setFormData(prev=>({
...prev,[field]:val}
))}
;
return m.jsx("div",{
className:"fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in duration-200",children:m.jsxs("div",{
className:"w-full max-w-2xl rounded-xl bg-card border border-sidebar-border shadow-2xl max-h-[90vh] flex flex-col overflow-hidden text-foreground",children:[m.jsxs("div",{
className:"flex items-center justify-between p-5 border-b border-sidebar-border bg-sidebar/50",children:[m.jsxs("div",{
className:"flex items-center gap-3",children:[m.jsx("div",{
className:"p-2 rounded-lg bg-teal-500/10 text-teal-500",children:m.jsx(CwIcon,{
className:"h-6 w-6"}
)}
),m.jsxs("div",{
children:[m.jsx("h2",{
className:"text-lg font-bold text-foreground",children:"Configurações do Chatwoot"}
),m.jsxs("p",{
className:"text-xs text-muted-foreground",children:["Instância: ",m.jsx("span",{
className:"font-semibold text-teal-400",children:r.instanceName}
)]}
)]}
)]}
),m.jsxs("div",{
className:"flex items-center gap-2",children:[m.jsxs(it,{
type:"button",variant:"outline",size:"sm",onClick:()=>fetchConfig(!0),disabled:fetching||loading,className:"h-8 px-3 text-xs gap-1.5 text-teal-400 border-teal-500/30 hover:bg-teal-500/10",children:[m.jsx("svg",{
className:`h-3.5 w-3.5 ${
fetching?"animate-spin":""}
`,fill:"none",viewBox:"0 0 24 24",stroke:"currentColor",strokeWidth:"2",children:m.jsx("path",{
strokeLinecap:"round",strokeLinejoin:"round",d:"M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"}
)}
),"Carregar Salvo"]}
),m.jsx("button",{
onClick:n,className:"rounded-lg p-1.5 text-muted-foreground hover:text-foreground hover:bg-sidebar-accent transition-colors",children:m.jsx(Es,{
className:"h-5 w-5"}
)}
)]}
)]}
),m.jsx("form",{
onSubmit:handleSave,className:"flex-1 overflow-y-auto p-6 space-y-6",children:m.jsxs("div",{
className:"space-y-6",children:[m.jsxs("div",{
className:`flex items-center justify-between p-4 rounded-xl border transition-all ${
formData.enabled?"bg-teal-500/10 border-teal-500/30":"bg-sidebar-accent/50 border-sidebar-border"}
`,children:[m.jsxs("div",{
className:"space-y-0.5",children:[m.jsx("label",{
className:"text-sm font-semibold text-foreground cursor-pointer",htmlFor:"cw-enabled",children:"Ativar Integração Chatwoot"}
),m.jsx("p",{
className:"text-xs text-muted-foreground",children:"Sincroniza mensagens e contatos do WhatsApp diretamente com sua plataforma Chatwoot."}
)]}
),m.jsx("input",{
type:"checkbox",id:"cw-enabled",checked:formData.enabled,onChange:e=>handleChange("enabled",e.target.checked),className:"h-5 w-5 rounded border-sidebar-border text-teal-600 focus:ring-teal-500 cursor-pointer accent-teal-500"}
)]}
),m.jsxs("div",{
className:"grid grid-cols-1 md:grid-cols-2 gap-4",children:[m.jsxs("div",{
className:"space-y-1.5 md:col-span-2",children:[m.jsx("label",{
className:"text-xs font-semibold text-foreground/80",children:"URL do Chatwoot *"}
),m.jsx("input",{
type:"url",placeholder:"https://chatwoot.suaempresa.com",value:formData.url,onChange:e=>handleChange("url",e.target.value),className:"w-full rounded-lg border border-sidebar-border bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-teal-500"}
)]}
),m.jsxs("div",{
className:"space-y-1.5 md:col-span-2",children:[m.jsx("label",{
className:"text-xs font-semibold text-foreground/80",children:"Token do Chatwoot (User API Key) *"}
),m.jsx("input",{
type:"password",placeholder:"User API Key do Chatwoot",value:formData.token,onChange:e=>handleChange("token",e.target.value),className:"w-full rounded-lg border border-sidebar-border bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-teal-500 font-mono"}
)]}
),m.jsxs("div",{
className:"space-y-1.5",children:[m.jsx("label",{
className:"text-xs font-semibold text-foreground/80",children:"Account ID (ID da Conta) *"}
),m.jsx("input",{
type:"text",placeholder:"1",value:formData.accountId,onChange:e=>handleChange("accountId",e.target.value),className:"w-full rounded-lg border border-sidebar-border bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-teal-500"}
)]}
),m.jsxs("div",{
className:"space-y-1.5",children:[m.jsx("label",{
className:"text-xs font-semibold text-foreground/80",children:"Inbox ID (ID da Caixa)"}
),m.jsx("input",{
type:"number",placeholder:"0 (Automático se auto-criar)",value:formData.inboxId,onChange:e=>handleChange("inboxId",e.target.value),className:"w-full rounded-lg border border-sidebar-border bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-teal-500"}
)]}
)]}
),m.jsxs("div",{
className:"space-y-3 border-t border-sidebar-border pt-4",children:[m.jsx("h3",{
className:"text-xs font-bold uppercase tracking-wider text-muted-foreground",children:"Opções de Comportamento"}
),m.jsxs("div",{
className:"grid grid-cols-1 md:grid-cols-2 gap-3",children:[m.jsxs("label",{
className:"flex items-center gap-2.5 p-3 rounded-lg border border-sidebar-border bg-sidebar-accent/30 cursor-pointer hover:bg-sidebar-accent/50 transition-colors",children:[m.jsx("input",{
type:"checkbox",checked:formData.autoCreate,onChange:e=>handleChange("autoCreate",e.target.checked),className:"h-4 w-4 rounded accent-teal-500"}
),m.jsxs("div",{
children:[m.jsx("p",{
className:"text-xs font-medium text-foreground",children:"Criar Inbox Automático"}
),m.jsx("p",{
className:"text-[10px] text-muted-foreground",children:"Cria a caixa de entrada no Chatwoot se não existir."}
)]}
)]}
),m.jsxs("label",{
className:"flex items-center gap-2.5 p-3 rounded-lg border border-sidebar-border bg-sidebar-accent/30 cursor-pointer hover:bg-sidebar-accent/50 transition-colors",children:[m.jsx("input",{
type:"checkbox",checked:formData.signMsg,onChange:e=>handleChange("signMsg",e.target.checked),className:"h-4 w-4 rounded accent-teal-500"}
),m.jsxs("div",{
children:[m.jsx("p",{
className:"text-xs font-medium text-foreground",children:"Assinar Mensagens"}
),m.jsx("p",{
className:"text-[10px] text-muted-foreground",children:"Adiciona o nome do atendente ao enviar."}
)]}
)]}
),m.jsxs("label",{
className:"flex items-center gap-2.5 p-3 rounded-lg border border-sidebar-border bg-sidebar-accent/30 cursor-pointer hover:bg-sidebar-accent/50 transition-colors",children:[m.jsx("input",{
type:"checkbox",checked:formData.reopenConversation,onChange:e=>handleChange("reopenConversation",e.target.checked),className:"h-4 w-4 rounded accent-teal-500"}
),m.jsxs("div",{
children:[m.jsx("p",{
className:"text-xs font-medium text-foreground",children:"Reabrir Conversa"}
),m.jsx("p",{
className:"text-[10px] text-muted-foreground",children:"Reabre conversas ao receber novas mensagens."}
)]}
)]}
),m.jsxs("label",{
className:"flex items-center gap-2.5 p-3 rounded-lg border border-sidebar-border bg-sidebar-accent/30 cursor-pointer hover:bg-sidebar-accent/50 transition-colors",children:[m.jsx("input",{
type:"checkbox",checked:formData.conversationPending,onChange:e=>handleChange("conversationPending",e.target.checked),className:"h-4 w-4 rounded accent-teal-500"}
),m.jsxs("div",{
children:[m.jsx("p",{
className:"text-xs font-medium text-foreground",children:"Conversa como Pendente"}
),m.jsx("p",{
className:"text-[10px] text-muted-foreground",children:"Status inicial como pendente no Chatwoot."}
)]}
)]}
),m.jsxs("label",{
className:"flex items-center gap-2.5 p-3 rounded-lg border border-sidebar-border bg-sidebar-accent/30 cursor-pointer hover:bg-sidebar-accent/50 transition-colors",children:[m.jsx("input",{
type:"checkbox",checked:formData.importContacts,onChange:e=>handleChange("importContacts",e.target.checked),className:"h-4 w-4 rounded accent-teal-500"}
),m.jsxs("div",{
children:[m.jsx("p",{
className:"text-xs font-medium text-foreground",children:"Importar Contatos"}
),m.jsx("p",{
className:"text-[10px] text-muted-foreground",children:"Sincroniza os contatos do WhatsApp."}
)]}
)]}
),m.jsxs("label",{
className:"flex items-center gap-2.5 p-3 rounded-lg border border-sidebar-border bg-sidebar-accent/30 cursor-pointer hover:bg-sidebar-accent/50 transition-colors",children:[m.jsx("input",{
type:"checkbox",checked:formData.importMessages,onChange:e=>handleChange("importMessages",e.target.checked),className:"h-4 w-4 rounded accent-teal-500"}
),m.jsxs("div",{
children:[m.jsx("p",{
className:"text-xs font-medium text-foreground",children:"Importar Mensagens"}
),m.jsx("p",{
className:"text-[10px] text-muted-foreground",children:"Sincroniza histórico recente de mensagens."}
)]}
)]}
)]}
),m.jsxs("div",{
className:"space-y-3 border-t border-sidebar-border pt-4",children:[m.jsxs("div",{
className:"flex items-center justify-between gap-4",children:[m.jsxs("div",{
className:"space-y-0.5",children:[m.jsx("label",{
className:"text-xs font-semibold text-foreground/80",children:"Limite de Dias para Importar Mensagens"}
),m.jsx("p",{
className:"text-[10px] text-muted-foreground",children:"0 para importar todo o histórico disponível."}
)]}
),m.jsx("input",{
type:"number",min:"0",value:formData.daysLimitImportMessages,onChange:e=>handleChange("daysLimitImportMessages",e.target.value),className:"w-28 rounded-lg border border-sidebar-border bg-background px-3 py-1.5 text-sm text-foreground text-center focus:outline-none focus:ring-2 focus:ring-teal-500"}
)]}
),formData.enabled&&m.jsxs("div",{
className:"flex flex-wrap gap-2 pt-2",children:[m.jsxs(it,{
type:"button",variant:"outline",size:"sm",onClick:handleSyncContacts,disabled:syncingContacts||loading,className:"text-xs gap-1.5 border-teal-500/30 text-teal-400 hover:bg-teal-500/10",children:[syncingContacts?"Sincronizando Contatos...":"Sincronizar Contatos Agora"]}
),m.jsxs(it,{
type:"button",variant:"outline",size:"sm",onClick:handleSyncMessages,disabled:syncingMessages||loading,className:"text-xs gap-1.5 border-teal-500/30 text-teal-400 hover:bg-teal-500/10",children:[syncingMessages?"Sincronizando Mensagens...":"Sincronizar Mensagens Agora"]}
)]}
)]}
)]}
)}
),m.jsxs("div",{
className:"flex items-center justify-between p-4 border-t border-sidebar-border bg-sidebar/50",children:[m.jsx(it,{
type:"button",variant:"outline",onClick:handleDelete,disabled:loading,className:"text-xs text-red-500 border-red-500/30 hover:bg-red-500/10 hover:text-red-400",children:"Remover / Desativar"}
),m.jsxs("div",{
className:"flex items-center gap-2",children:[m.jsx(it,{
type:"button",variant:"ghost",onClick:n,disabled:loading,className:"text-xs text-muted-foreground hover:text-foreground",children:"Cancelar"}
),m.jsxs(it,{
type:"button",onClick:handleSave,disabled:loading,className:"text-xs bg-teal-600 hover:bg-teal-500 text-white font-medium px-4 py-2 rounded-lg flex items-center gap-1.5 transition-colors",children:[loading&&m.jsx("svg",{
className:"animate-spin h-3.5 w-3.5 text-white",fill:"none",viewBox:"0 0 24 24",children:[m.jsx("circle",{
className:"opacity-25",cx:"12",cy:"12",r:"10",stroke:"currentColor",strokeWidth:"4"}
),m.jsx("path",{
className:"opacity-75",fill:"currentColor",d:"M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"}
)].filter(Boolean)}
),"Salvar Configuração"]}
)]}
)]}
)]}
)}
)}

