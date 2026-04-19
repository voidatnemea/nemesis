package admin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/internal/services"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

const defaultChatbotSystemPrompt = `You are the FeatherPanel AI support agent.
You help users run game servers hosted on this panel.
Be concise, technical, and friendly. Never invent API endpoints or features.
When uncertain, direct the user to an admin or the knowledge base.`

func (s *SettingsController) GetChatbotSystemPrompt(c *gin.Context) {
	prompt := services.GetSetting("chatbot:system_prompt_core", defaultChatbotSystemPrompt)
	utils.Success(c, gin.H{"system_prompt": prompt}, "System prompt retrieved", http.StatusOK)
}

func (s *SettingsController) UpdateChatbotSystemPrompt(c *gin.Context) {
	var req struct {
		SystemPrompt string `json:"system_prompt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	services.SetSetting("chatbot:system_prompt_core", req.SystemPrompt)
	utils.Success(c, gin.H{"system_prompt": req.SystemPrompt}, "System prompt updated", http.StatusOK)
}

type SettingsController struct{}

type settingMeta struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Required    bool     `json:"required"`
	Placeholder string   `json:"placeholder"`
	Validation  string   `json:"validation"`
	Category    string   `json:"category"`
	Options     []string `json:"options"`
	Sensitive   bool     `json:"sensitive,omitempty"`
	Value       any      `json:"value"`
}

type categoryMeta struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	Settings    []string `json:"settings"`
}

var settingsSchema = map[string]settingMeta{
	"app:name":        {Name: "Application Name", Description: "The name of your panel.", Type: "text", Required: true, Placeholder: "Nemesis Panel", Category: "general"},
	"app:url":         {Name: "Application URL", Description: "The public URL of your panel (no trailing slash).", Type: "text", Required: true, Placeholder: "https://panel.example.com", Category: "general"},
	"app:description": {Name: "Description", Description: "Short tagline shown on the login page.", Type: "textarea", Required: false, Placeholder: "The best game server panel.", Category: "general"},
	"app:version":     {Name: "Version", Description: "Displayed version string.", Type: "text", Required: false, Placeholder: "1.0.0", Category: "general"},
	"app:timezone":    {Name: "Timezone", Description: "Default timezone for the panel.", Type: "text", Required: false, Placeholder: "UTC", Category: "general"},
	"app:language":    {Name: "Default Language", Description: "Default UI language.", Type: "select", Required: false, Placeholder: "en", Category: "general", Options: []string{"en"}},
	"app:support_url": {Name: "Support URL", Description: "Link shown in the footer and help menu.", Type: "text", Required: false, Placeholder: "https://discord.gg/example", Category: "general"},
	"app:terms_url":   {Name: "Terms of Service URL", Description: "Link to your Terms of Service page.", Type: "text", Required: false, Placeholder: "https://example.com/tos", Category: "general"},
	"app:privacy_url": {Name: "Privacy Policy URL", Description: "Link to your Privacy Policy page.", Type: "text", Required: false, Placeholder: "https://example.com/privacy", Category: "general"},

	"branding:logo_light":       {Name: "Light Logo URL", Description: "Logo displayed on light-mode pages.", Type: "text", Required: false, Placeholder: "/assets/logo-light.png", Category: "branding"},
	"branding:logo_dark":        {Name: "Dark Logo URL", Description: "Logo displayed on dark-mode pages.", Type: "text", Required: false, Placeholder: "/assets/logo-dark.png", Category: "branding"},
	"branding:favicon":          {Name: "Favicon URL", Description: "Browser tab icon URL.", Type: "text", Required: false, Placeholder: "/favicon.ico", Category: "branding"},
	"branding:login_background": {Name: "Login Background URL", Description: "Image shown behind the login form.", Type: "text", Required: false, Placeholder: "/assets/bg.jpg", Category: "branding"},
	"branding:primary_color":    {Name: "Primary Color", Description: "Primary accent color (hex).", Type: "text", Required: false, Placeholder: "#6366f1", Category: "branding"},
	"branding:accent_color":     {Name: "Accent Color", Description: "Secondary accent color (hex).", Type: "text", Required: false, Placeholder: "#8b5cf6", Category: "branding"},
	"branding:footer_text":      {Name: "Footer Text", Description: "Custom text shown in the footer.", Type: "text", Required: false, Placeholder: "© 2026 Your Company", Category: "branding"},
	"branding:custom_css":       {Name: "Custom CSS", Description: "Custom CSS injected into every page.", Type: "textarea", Required: false, Placeholder: "/* your css */", Category: "branding"},
	"branding:custom_js":        {Name: "Custom JavaScript", Description: "Custom JavaScript injected into every page (before </body>).", Type: "textarea", Required: false, Placeholder: "// your js", Category: "branding"},

	"auth:registration":      {Name: "User Registration", Description: "Allow new users to register accounts.", Type: "toggle", Required: false, Category: "auth"},
	"auth:require_email":     {Name: "Require Email Verification", Description: "Users must verify their email before logging in.", Type: "toggle", Required: false, Category: "auth"},
	"auth:password_min":      {Name: "Minimum Password Length", Description: "Minimum number of characters for a password.", Type: "text", Required: false, Placeholder: "8", Category: "auth"},
	"auth:password_complex":  {Name: "Require Complex Passwords", Description: "Require uppercase, lowercase, numbers and symbols.", Type: "toggle", Required: false, Category: "auth"},
	"auth:oidc_enabled":      {Name: "OIDC Login", Description: "Allow login via OIDC providers.", Type: "toggle", Required: false, Category: "auth"},
	"auth:session_lifetime":  {Name: "Session Lifetime (days)", Description: "How many days a remember-me session lasts.", Type: "text", Required: false, Placeholder: "30", Category: "auth"},
	"auth:max_login_attempts": {Name: "Max Login Attempts", Description: "Lock the account after N failed logins (0 = unlimited).", Type: "text", Required: false, Placeholder: "5", Category: "auth"},
	"auth:lockout_duration":   {Name: "Lockout Duration (minutes)", Description: "How long a locked account stays locked.", Type: "text", Required: false, Placeholder: "15", Category: "auth"},

	"discord:oauth_enabled": {Name: "Discord OAuth", Description: "Allow login via Discord OAuth2.", Type: "toggle", Required: false, Category: "discord"},
	"discord:client_id":     {Name: "Discord Client ID", Description: "OAuth2 client ID from the Discord Developer Portal.", Type: "text", Required: false, Placeholder: "123456789", Category: "discord"},
	"discord:client_secret": {Name: "Discord Client Secret", Description: "OAuth2 client secret from Discord.", Type: "password", Required: false, Placeholder: "••••••••", Category: "discord", Sensitive: true},
	"discord:redirect_uri":  {Name: "Redirect URI", Description: "Callback URL registered with Discord.", Type: "text", Required: false, Placeholder: "https://panel.example.com/auth/discord/callback", Category: "discord"},
	"discord:guild_id":      {Name: "Required Guild ID", Description: "Require users to be in this Discord guild (optional).", Type: "text", Required: false, Placeholder: "987654321", Category: "discord"},
	"discord:bot_token":     {Name: "Bot Token", Description: "Bot token for guild membership checks.", Type: "password", Required: false, Placeholder: "••••••••", Category: "discord", Sensitive: true},
	"discord:webhook_url":   {Name: "Webhook URL", Description: "Discord webhook for panel notifications.", Type: "password", Required: false, Placeholder: "https://discord.com/api/webhooks/...", Category: "discord", Sensitive: true},

	"twofa:enabled":        {Name: "Enable 2FA", Description: "Allow users to enable two-factor authentication.", Type: "toggle", Required: false, Category: "twofa"},
	"twofa:required":       {Name: "Require 2FA", Description: "Force all users to enable 2FA.", Type: "toggle", Required: false, Category: "twofa"},
	"twofa:required_admin": {Name: "Require 2FA for Admins", Description: "Force administrators to enable 2FA.", Type: "toggle", Required: false, Category: "twofa"},
	"twofa:issuer":         {Name: "TOTP Issuer", Description: "The issuer name shown in authenticator apps.", Type: "text", Required: false, Placeholder: "Nemesis Panel", Category: "twofa"},
	"twofa:backup_codes":   {Name: "Backup Codes", Description: "Number of backup codes issued to each user.", Type: "text", Required: false, Placeholder: "8", Category: "twofa"},
	"twofa:window":         {Name: "TOTP Window", Description: "Allowed clock drift window (steps of 30s).", Type: "text", Required: false, Placeholder: "1", Category: "twofa"},

	"mail:driver":     {Name: "Mail Driver", Description: "Transport for outgoing email.", Type: "select", Required: false, Placeholder: "smtp", Category: "mail", Options: []string{"smtp", "log", "null"}},
	"mail:host":       {Name: "SMTP Host", Description: "Hostname of your SMTP server.", Type: "text", Required: false, Placeholder: "smtp.example.com", Category: "mail"},
	"mail:port":       {Name: "SMTP Port", Description: "Port of your SMTP server.", Type: "text", Required: false, Placeholder: "587", Category: "mail"},
	"mail:username":   {Name: "SMTP Username", Description: "Login for your SMTP server.", Type: "text", Required: false, Placeholder: "noreply@example.com", Category: "mail"},
	"mail:password":   {Name: "SMTP Password", Description: "Password for your SMTP server.", Type: "password", Required: false, Placeholder: "••••••••", Category: "mail", Sensitive: true},
	"mail:encryption": {Name: "Encryption", Description: "Encryption method for SMTP.", Type: "select", Required: false, Placeholder: "tls", Category: "mail", Options: []string{"tls", "ssl", "none"}},
	"mail:from":       {Name: "From Address", Description: "The sender address for all panel emails.", Type: "text", Required: false, Placeholder: "noreply@example.com", Category: "mail"},
	"mail:from_name":  {Name: "From Name", Description: "The sender display name.", Type: "text", Required: false, Placeholder: "Nemesis Panel", Category: "mail"},

	"security:recaptcha_enabled":   {Name: "reCAPTCHA", Description: "Enable Google reCAPTCHA on auth forms.", Type: "toggle", Required: false, Category: "security"},
	"security:recaptcha_site_key":  {Name: "reCAPTCHA Site Key", Description: "Your reCAPTCHA v2 site key.", Type: "text", Required: false, Placeholder: "6Le...", Category: "security"},
	"security:recaptcha_secret":    {Name: "reCAPTCHA Secret Key", Description: "Your reCAPTCHA v2 secret key.", Type: "password", Required: false, Placeholder: "••••••••", Category: "security", Sensitive: true},
	"security:cors_origins":        {Name: "Allowed CORS Origins", Description: "Comma-separated list of allowed origins.", Type: "textarea", Required: false, Placeholder: "https://panel.example.com", Category: "security"},
	"security:rate_limit":          {Name: "Global Rate Limiting", Description: "Enable API rate limiting.", Type: "toggle", Required: false, Category: "security"},
	"security:csp_enabled":         {Name: "Content Security Policy", Description: "Send a strict CSP header.", Type: "toggle", Required: false, Category: "security"},
	"security:hsts_enabled":        {Name: "HSTS", Description: "Send the Strict-Transport-Security header.", Type: "toggle", Required: false, Category: "security"},
	"security:ip_allowlist":        {Name: "Admin IP Allowlist", Description: "Comma-separated IPs allowed to access the admin area.", Type: "textarea", Required: false, Placeholder: "1.2.3.4, 5.6.7.0/24", Category: "security"},
	"security:zerotrust_enabled":   {Name: "ZeroTrust Scanner", Description: "Scan uploaded files against known-bad hashes.", Type: "toggle", Required: false, Category: "security"},
	"security:maintenance_mode":    {Name: "Maintenance Mode", Description: "Block all user-facing routes.", Type: "toggle", Required: false, Category: "security"},

	"seo:meta_title":        {Name: "Meta Title", Description: "Default <title> tag.", Type: "text", Required: false, Placeholder: "Nemesis Panel", Category: "seo"},
	"seo:meta_description":  {Name: "Meta Description", Description: "Default meta description.", Type: "textarea", Required: false, Placeholder: "The best game server panel.", Category: "seo"},
	"seo:meta_keywords":     {Name: "Meta Keywords", Description: "Comma-separated keywords.", Type: "text", Required: false, Placeholder: "game, hosting, panel", Category: "seo"},
	"seo:og_image":          {Name: "Open Graph Image", Description: "URL used for social share previews.", Type: "text", Required: false, Placeholder: "/assets/og.png", Category: "seo"},
	"seo:twitter_handle":    {Name: "Twitter Handle", Description: "Twitter @username for card attribution.", Type: "text", Required: false, Placeholder: "@nemesis", Category: "seo"},
	"seo:robots_txt":        {Name: "robots.txt", Description: "Custom robots.txt content.", Type: "textarea", Required: false, Placeholder: "User-agent: *\nAllow: /", Category: "seo"},
	"seo:sitemap_enabled":   {Name: "Enable Sitemap", Description: "Serve /sitemap.xml automatically.", Type: "toggle", Required: false, Category: "seo"},
	"seo:canonical_base":    {Name: "Canonical Base URL", Description: "Base URL for canonical link tags.", Type: "text", Required: false, Placeholder: "https://panel.example.com", Category: "seo"},

	"pwa:enabled":          {Name: "Enable PWA", Description: "Install the panel as a Progressive Web App.", Type: "toggle", Required: false, Category: "pwa"},
	"pwa:app_name":         {Name: "PWA App Name", Description: "Long name in the manifest.", Type: "text", Required: false, Placeholder: "Nemesis Panel", Category: "pwa"},
	"pwa:short_name":       {Name: "PWA Short Name", Description: "Short name on the home screen.", Type: "text", Required: false, Placeholder: "Nemesis", Category: "pwa"},
	"pwa:theme_color":      {Name: "Theme Color", Description: "PWA theme color (hex).", Type: "text", Required: false, Placeholder: "#6366f1", Category: "pwa"},
	"pwa:background_color": {Name: "Background Color", Description: "PWA splash background color (hex).", Type: "text", Required: false, Placeholder: "#0f0f10", Category: "pwa"},
	"pwa:display":          {Name: "Display Mode", Description: "Manifest display mode.", Type: "select", Required: false, Placeholder: "standalone", Category: "pwa", Options: []string{"standalone", "fullscreen", "minimal-ui", "browser"}},
	"pwa:icon_192":         {Name: "Icon 192×192", Description: "URL to the 192×192 icon.", Type: "text", Required: false, Placeholder: "/icons/icon-192.png", Category: "pwa"},
	"pwa:icon_512":         {Name: "Icon 512×512", Description: "URL to the 512×512 icon.", Type: "text", Required: false, Placeholder: "/icons/icon-512.png", Category: "pwa"},

	"chatbot:enabled":       {Name: "Enable Chatbot", Description: "Show the AI chatbot widget.", Type: "toggle", Required: false, Category: "chatbot"},
	"chatbot:provider":      {Name: "Provider", Description: "The chatbot provider.", Type: "select", Required: false, Placeholder: "openai", Category: "chatbot", Options: []string{"openai", "anthropic", "ollama", "custom"}},
	"chatbot:api_key":       {Name: "API Key", Description: "Provider API key.", Type: "password", Required: false, Placeholder: "••••••••", Category: "chatbot", Sensitive: true},
	"chatbot:model":         {Name: "Model", Description: "Model identifier (e.g. gpt-4, claude-sonnet-4-6).", Type: "text", Required: false, Placeholder: "gpt-4", Category: "chatbot"},
	"chatbot:system_prompt": {Name: "System Prompt", Description: "Prompt prepended to every conversation.", Type: "textarea", Required: false, Placeholder: "You are a helpful support agent for our game panel.", Category: "chatbot"},
	"chatbot:greeting":      {Name: "Greeting Message", Description: "Opening message shown to users.", Type: "text", Required: false, Placeholder: "Hi! How can I help?", Category: "chatbot"},
	"chatbot:user_required": {Name: "Require Login", Description: "Only authenticated users can use the chatbot.", Type: "toggle", Required: false, Category: "chatbot"},

	"status:enabled":          {Name: "Status Page", Description: "Enable the public status page.", Type: "toggle", Required: false, Category: "status"},
	"status:public_url":       {Name: "Public URL", Description: "External status page URL (if hosted elsewhere).", Type: "text", Required: false, Placeholder: "https://status.example.com", Category: "status"},
	"status:show_incidents":   {Name: "Show Incidents", Description: "Display active incidents on the status page.", Type: "toggle", Required: false, Category: "status"},
	"status:show_uptime":      {Name: "Show Uptime %", Description: "Display uptime percentages for each node.", Type: "toggle", Required: false, Category: "status"},
	"status:show_metrics":     {Name: "Show Node Metrics", Description: "Display CPU/RAM usage on the status page.", Type: "toggle", Required: false, Category: "status"},
	"status:refresh_interval": {Name: "Refresh Interval (s)", Description: "How often the status page auto-refreshes.", Type: "text", Required: false, Placeholder: "30", Category: "status"},

	"servers:default_cpu":        {Name: "Default CPU Limit (%)", Description: "CPU cap applied when a user doesn't specify.", Type: "text", Required: false, Placeholder: "100", Category: "servers"},
	"servers:default_memory":     {Name: "Default Memory Limit (MB)", Description: "Memory cap applied by default.", Type: "text", Required: false, Placeholder: "1024", Category: "servers"},
	"servers:default_disk":       {Name: "Default Disk Limit (MB)", Description: "Disk cap applied by default.", Type: "text", Required: false, Placeholder: "5120", Category: "servers"},
	"servers:default_swap":       {Name: "Default Swap Limit (MB)", Description: "Swap cap applied by default (-1 = unlimited).", Type: "text", Required: false, Placeholder: "0", Category: "servers"},
	"servers:default_io":         {Name: "Default IO Weight", Description: "Block-IO weight (10–1000).", Type: "text", Required: false, Placeholder: "500", Category: "servers"},
	"servers:max_per_user":       {Name: "Max Servers Per User", Description: "Hard cap on servers per user (0 = unlimited).", Type: "text", Required: false, Placeholder: "10", Category: "servers"},
	"servers:allow_oom":          {Name: "OOM Killer", Description: "Allow the kernel OOM killer to stop servers.", Type: "toggle", Required: false, Category: "servers"},
	"servers:startup_timeout":    {Name: "Startup Timeout (s)", Description: "How long to wait for a server to boot.", Type: "text", Required: false, Placeholder: "30", Category: "servers"},
	"servers:backup_enabled":     {Name: "Backups", Description: "Allow users to create backups.", Type: "toggle", Required: false, Category: "servers"},
	"servers:backup_max_per_srv": {Name: "Max Backups Per Server", Description: "Per-server backup retention cap.", Type: "text", Required: false, Placeholder: "5", Category: "servers"},

	"knowledgebase:enabled":        {Name: "Knowledgebase", Description: "Enable the knowledgebase module.", Type: "toggle", Required: false, Category: "knowledgebase"},
	"knowledgebase:public":         {Name: "Publicly Visible", Description: "Allow non-logged-in users to read articles.", Type: "toggle", Required: false, Category: "knowledgebase"},
	"knowledgebase:search_enabled": {Name: "Full-Text Search", Description: "Enable the search bar on the knowledgebase.", Type: "toggle", Required: false, Category: "knowledgebase"},
	"knowledgebase:show_footer":    {Name: "Show Footer Link", Description: "Display a knowledgebase link in the footer.", Type: "toggle", Required: false, Category: "knowledgebase"},
	"knowledgebase:welcome":        {Name: "Welcome Text", Description: "Text shown on the knowledgebase landing page.", Type: "textarea", Required: false, Placeholder: "Welcome to our help center.", Category: "knowledgebase"},

	"tickets:enabled":         {Name: "Ticket System", Description: "Enable the support ticket module.", Type: "toggle", Required: false, Category: "tickets"},
	"tickets:max_open":        {Name: "Max Open Tickets", Description: "Maximum tickets a user can have open at once (0 = unlimited).", Type: "text", Required: false, Placeholder: "5", Category: "tickets"},
	"tickets:default_dept":    {Name: "Default Department", Description: "Department assigned to new tickets.", Type: "text", Required: false, Placeholder: "Support", Category: "tickets"},
	"tickets:auto_close_days": {Name: "Auto-close Idle (days)", Description: "Close tickets with no reply after N days (0 = never).", Type: "text", Required: false, Placeholder: "7", Category: "tickets"},
	"tickets:allow_attachments": {Name: "Allow Attachments", Description: "Let users upload files to tickets.", Type: "toggle", Required: false, Category: "tickets"},
	"tickets:notify_admins":     {Name: "Notify Admins", Description: "Send admins a notification on new tickets.", Type: "toggle", Required: false, Category: "tickets"},

	"analytics:google_enabled":    {Name: "Google Analytics", Description: "Enable Google Analytics tracking.", Type: "toggle", Required: false, Category: "analytics"},
	"analytics:google_id":         {Name: "GA Measurement ID", Description: "Your GA4 measurement ID (e.g. G-XXXX).", Type: "text", Required: false, Placeholder: "G-XXXXXXXXXX", Category: "analytics"},
	"analytics:plausible_enabled": {Name: "Plausible Analytics", Description: "Enable Plausible tracking.", Type: "toggle", Required: false, Category: "analytics"},
	"analytics:plausible_domain":  {Name: "Plausible Domain", Description: "The domain configured in Plausible.", Type: "text", Required: false, Placeholder: "panel.example.com", Category: "analytics"},
	"analytics:plausible_script":  {Name: "Plausible Script URL", Description: "Script URL if self-hosting Plausible.", Type: "text", Required: false, Placeholder: "https://plausible.io/js/script.js", Category: "analytics"},
	"analytics:matomo_enabled":    {Name: "Matomo Analytics", Description: "Enable Matomo tracking.", Type: "toggle", Required: false, Category: "analytics"},
	"analytics:matomo_url":        {Name: "Matomo URL", Description: "Your Matomo instance URL.", Type: "text", Required: false, Placeholder: "https://matomo.example.com", Category: "analytics"},
	"analytics:matomo_site_id":    {Name: "Matomo Site ID", Description: "The numeric site ID in Matomo.", Type: "text", Required: false, Placeholder: "1", Category: "analytics"},
	"analytics:posthog_key":       {Name: "PostHog API Key", Description: "PostHog project API key.", Type: "password", Required: false, Placeholder: "phc_...", Category: "analytics", Sensitive: true},
	"analytics:posthog_host":      {Name: "PostHog Host", Description: "PostHog host URL.", Type: "text", Required: false, Placeholder: "https://app.posthog.com", Category: "analytics"},

	"notifications:enabled":       {Name: "Notifications", Description: "Enable the panel notification system.", Type: "toggle", Required: false, Category: "features"},
	"notifications:email":         {Name: "Email Notifications", Description: "Also send notifications via email.", Type: "toggle", Required: false, Category: "features"},
	"notifications:discord":       {Name: "Discord Notifications", Description: "Send notifications to the configured Discord webhook.", Type: "toggle", Required: false, Category: "features"},
}

var categoriesSchema = map[string]categoryMeta{
	"general":       {Name: "General", Description: "Core application settings.", Icon: "Settings", Settings: []string{"app:name", "app:url", "app:description", "app:version", "app:timezone", "app:language", "app:support_url", "app:terms_url", "app:privacy_url"}},
	"branding":      {Name: "Branding", Description: "Logos, colors and custom styling.", Icon: "Palette", Settings: []string{"branding:logo_light", "branding:logo_dark", "branding:favicon", "branding:login_background", "branding:primary_color", "branding:accent_color", "branding:footer_text", "branding:custom_css", "branding:custom_js"}},
	"auth":          {Name: "Authentication", Description: "Login, registration, and password policy.", Icon: "Shield", Settings: []string{"auth:registration", "auth:require_email", "auth:password_min", "auth:password_complex", "auth:oidc_enabled", "auth:session_lifetime", "auth:max_login_attempts", "auth:lockout_duration"}},
	"discord":       {Name: "Discord OAuth", Description: "Sign-in and notifications via Discord.", Icon: "MessageCircle", Settings: []string{"discord:oauth_enabled", "discord:client_id", "discord:client_secret", "discord:redirect_uri", "discord:guild_id", "discord:bot_token", "discord:webhook_url"}},
	"twofa":         {Name: "Two-Factor Auth", Description: "TOTP, backup codes, and 2FA policy.", Icon: "KeyRound", Settings: []string{"twofa:enabled", "twofa:required", "twofa:required_admin", "twofa:issuer", "twofa:backup_codes", "twofa:window"}},
	"mail":          {Name: "Mail", Description: "Outgoing email configuration.", Icon: "Mail", Settings: []string{"mail:driver", "mail:host", "mail:port", "mail:username", "mail:password", "mail:encryption", "mail:from", "mail:from_name"}},
	"security":      {Name: "Security", Description: "Security hardening and rate limiting.", Icon: "ShieldCheck", Settings: []string{"security:recaptcha_enabled", "security:recaptcha_site_key", "security:recaptcha_secret", "security:cors_origins", "security:rate_limit", "security:csp_enabled", "security:hsts_enabled", "security:ip_allowlist", "security:zerotrust_enabled", "security:maintenance_mode"}},
	"seo":           {Name: "SEO", Description: "Search engine and social sharing metadata.", Icon: "Globe", Settings: []string{"seo:meta_title", "seo:meta_description", "seo:meta_keywords", "seo:og_image", "seo:twitter_handle", "seo:robots_txt", "seo:sitemap_enabled", "seo:canonical_base"}},
	"pwa":           {Name: "Progressive Web App", Description: "Installability and manifest options.", Icon: "Smartphone", Settings: []string{"pwa:enabled", "pwa:app_name", "pwa:short_name", "pwa:theme_color", "pwa:background_color", "pwa:display", "pwa:icon_192", "pwa:icon_512"}},
	"chatbot":       {Name: "Chatbot", Description: "AI assistant widget configuration.", Icon: "Bot", Settings: []string{"chatbot:enabled", "chatbot:provider", "chatbot:api_key", "chatbot:model", "chatbot:system_prompt", "chatbot:greeting", "chatbot:user_required"}},
	"status":        {Name: "Status Page", Description: "Public service status configuration.", Icon: "Activity", Settings: []string{"status:enabled", "status:public_url", "status:show_incidents", "status:show_uptime", "status:show_metrics", "status:refresh_interval"}},
	"servers":       {Name: "Servers", Description: "Default limits and server behaviour.", Icon: "Server", Settings: []string{"servers:default_cpu", "servers:default_memory", "servers:default_disk", "servers:default_swap", "servers:default_io", "servers:max_per_user", "servers:allow_oom", "servers:startup_timeout", "servers:backup_enabled", "servers:backup_max_per_srv"}},
	"knowledgebase": {Name: "Knowledgebase", Description: "Help center and article visibility.", Icon: "BookOpen", Settings: []string{"knowledgebase:enabled", "knowledgebase:public", "knowledgebase:search_enabled", "knowledgebase:show_footer", "knowledgebase:welcome"}},
	"tickets":       {Name: "Tickets", Description: "Support ticket system configuration.", Icon: "LifeBuoy", Settings: []string{"tickets:enabled", "tickets:max_open", "tickets:default_dept", "tickets:auto_close_days", "tickets:allow_attachments", "tickets:notify_admins"}},
	"analytics":     {Name: "Analytics", Description: "Usage analytics providers.", Icon: "BarChart3", Settings: []string{"analytics:google_enabled", "analytics:google_id", "analytics:plausible_enabled", "analytics:plausible_domain", "analytics:plausible_script", "analytics:matomo_enabled", "analytics:matomo_url", "analytics:matomo_site_id", "analytics:posthog_key", "analytics:posthog_host"}},
	"features":      {Name: "Features", Description: "Notification delivery toggles.", Icon: "Sliders", Settings: []string{"notifications:enabled", "notifications:email", "notifications:discord"}},
}

func (s *SettingsController) Index(c *gin.Context) {
	var rows []models.Setting
	database.DB.Find(&rows)
	dbValues := make(map[string]string, len(rows))
	for _, r := range rows {
		dbValues[r.Key] = r.Value
	}

	allSettings := make(map[string]settingMeta, len(settingsSchema))
	for key, meta := range settingsSchema {
		meta.Options = meta.Options
		if meta.Options == nil {
			meta.Options = []string{}
		}
		rawVal := dbValues[key]
		if meta.Type == "toggle" {
			meta.Value = rawVal == "true"
		} else {
			meta.Value = rawVal
		}
		allSettings[key] = meta
	}

	organizedSettings := make(map[string]gin.H, len(categoriesSchema))
	for catKey, cat := range categoriesSchema {
		catSettings := make(map[string]settingMeta)
		for _, sKey := range cat.Settings {
			if m, ok := allSettings[sKey]; ok {
				catSettings[sKey] = m
			}
		}
		organizedSettings[catKey] = gin.H{
			"category": cat,
			"settings": catSettings,
		}
	}

	utils.Success(c, gin.H{
		"settings":           allSettings,
		"categories":         categoriesSchema,
		"organized_settings": organizedSettings,
	}, "Settings retrieved", http.StatusOK)
}

func (s *SettingsController) Categories(c *gin.Context) {
	cats := make([]string, 0, len(categoriesSchema))
	for k := range categoriesSchema {
		cats = append(cats, k)
	}
	utils.Success(c, cats, "Categories retrieved", http.StatusOK)
}

func (s *SettingsController) GetByCategory(c *gin.Context) {
	category := c.Param("category")
	var settings []models.Setting
	database.DB.Where("`key` LIKE ?", category+":%").Find(&settings)
	utils.Success(c, settings, "Settings retrieved", http.StatusOK)
}

func (s *SettingsController) Show(c *gin.Context) {
	key := c.Param("setting")
	var setting models.Setting
	if err := database.DB.Where("`key` = ?", key).First(&setting).Error; err != nil {
		utils.Error(c, "Setting not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, setting, "Setting retrieved", http.StatusOK)
}

func (s *SettingsController) Update(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	for key, val := range req {
		var strVal string
		switch v := val.(type) {
		case bool:
			if v {
				strVal = "true"
			} else {
				strVal = "false"
			}
		case string:
			strVal = v
		case float64:
			strVal = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", v), "0"), ".")
		default:
			strVal = fmt.Sprintf("%v", v)
		}
		var setting models.Setting
		if database.DB.Where("`key` = ?", key).First(&setting).Error != nil {
			database.DB.Create(&models.Setting{Key: key, Value: strVal})
		} else {
			database.DB.Model(&setting).Update("value", strVal)
		}
	}
	utils.Success(c, nil, "Settings updated", http.StatusOK)
}
