package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/api/controllers"
	"github.com/mythicalsystems/nemesis-backend/internal/api/controllers/admin"
	"github.com/mythicalsystems/nemesis-backend/internal/api/controllers/system"
	userctrl "github.com/mythicalsystems/nemesis-backend/internal/api/controllers/user"
	"github.com/mythicalsystems/nemesis-backend/internal/api/controllers/wings"
	"github.com/mythicalsystems/nemesis-backend/internal/api/middleware"
)

func Register(r *gin.Engine) {
	homeController := &controllers.HomeController{}
	webAppController := &controllers.WebAppController{}
	authController := &controllers.AuthController{}
	userController := &controllers.UserController{}

	api := r.Group("/api")

	// ── Public ────────────────────────────────────────────────────────────────
	api.GET("/", homeController.Index)
	api.GET("/manifest.webmanifest", webAppController.Index)
	api.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// Auth (public)
	auth := api.Group("/user/auth")
	{
		auth.PUT("/login", authController.Login)
		auth.PUT("/register", authController.Register)
		auth.DELETE("/logout", authController.Logout)
		// 2FA
		twofa := &userctrl.TwoFactorController{}
		auth.POST("/two-factor", twofa.Verify)
		// Password reset
		forgot := &userctrl.ForgotPasswordController{}
		reset := &userctrl.ResetPasswordController{}
		auth.PUT("/forgot-password", forgot.Send)
		auth.PUT("/reset-password", reset.Reset)
	}

	// System / public
	sysCtrl := &system.SystemController{}
	api.GET("/system/settings", sysCtrl.GetSettings)
	api.GET("/system/oidc/providers", sysCtrl.GetOidcProviders)

	// ── Authenticated user routes ──────────────────────────────────────────────
	user := api.Group("/user")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/profile", userController.Profile)

		// Session
		sess := &userctrl.SessionController{}
		user.GET("/session", sess.Get)
		user.PATCH("/session", sess.Update)

		// API clients
		apiClients := &userctrl.ApiClientController{}
		user.GET("/api-clients", apiClients.GetApiClients)
		user.POST("/api-clients", apiClients.CreateApiClient)
		user.GET("/api-clients/:id", apiClients.GetApiClient)
		user.PUT("/api-clients/:id", apiClients.UpdateApiClient)
		user.DELETE("/api-clients/:id", apiClients.DeleteApiClient)
		user.POST("/api-clients/:id/regenerate", apiClients.RegenerateApiKeys)

		// SSH keys
		sshKeys := &userctrl.UserSshKeyController{}
		user.GET("/ssh-keys", sshKeys.GetUserSshKeys)
		user.POST("/ssh-keys", sshKeys.CreateUserSshKey)
		user.GET("/ssh-keys/:id", sshKeys.GetUserSshKey)
		user.PUT("/ssh-keys/:id", sshKeys.UpdateUserSshKey)
		user.DELETE("/ssh-keys/:id", sshKeys.DeleteUserSshKey)

		// Tickets
		tickets := &userctrl.UserTicketsController{}
		user.GET("/tickets", tickets.Index)
		user.PUT("/tickets", tickets.Create)
		user.GET("/tickets/categories", tickets.GetCategories)
		user.GET("/tickets/priorities", tickets.GetPriorities)
		user.GET("/tickets/statuses", tickets.GetStatuses)
		user.GET("/tickets/:uuid", tickets.Show)
		user.POST("/tickets/:uuid/reply", tickets.Reply)
		user.DELETE("/tickets/:uuid", tickets.Delete)

		// Knowledgebase (user)
		kb := &userctrl.UserKnowledgebaseController{}
		user.GET("/knowledgebase/categories", kb.GetCategories)
		user.GET("/knowledgebase/articles", kb.GetArticles)
		user.GET("/knowledgebase/articles/:id", kb.GetArticle)

		// Notifications
		notifCtrl := &system.NotificationsController{}
		user.GET("/notifications", notifCtrl.GetUserNotifications)

		// Servers
		serverCtrl := &userctrl.ServerUserController{}
		user.GET("/servers", serverCtrl.GetUserServers)
		user.GET("/servers/:uuidShort", serverCtrl.GetServer)
		user.PUT("/servers/:uuidShort", serverCtrl.UpdateServer)
		user.DELETE("/servers/:uuidShort", serverCtrl.DeleteServer)

		// Server backups
		backupCtrl := &userctrl.ServerBackupController{}
		user.GET("/servers/:uuidShort/backups", backupCtrl.GetBackups)
		user.POST("/servers/:uuidShort/backups", backupCtrl.CreateBackup)
		user.GET("/servers/:uuidShort/backups/:backupUuid", backupCtrl.GetBackup)
		user.DELETE("/servers/:uuidShort/backups/:backupUuid", backupCtrl.DeleteBackup)
		user.POST("/servers/:uuidShort/backups/:backupUuid/lock", backupCtrl.LockBackup)
		user.POST("/servers/:uuidShort/backups/:backupUuid/unlock", backupCtrl.UnlockBackup)

		// Server databases
		dbCtrl := &userctrl.ServerDatabaseController{}
		user.GET("/servers/:uuidShort/databases", dbCtrl.GetServerDatabases)
		user.POST("/servers/:uuidShort/databases", dbCtrl.CreateServerDatabase)
		user.GET("/servers/:uuidShort/databases/hosts", dbCtrl.GetAvailableDatabaseHosts)
		user.GET("/servers/:uuidShort/databases/:databaseId", dbCtrl.GetServerDatabase)
		user.DELETE("/servers/:uuidShort/databases/:databaseId", dbCtrl.DeleteServerDatabase)

		// Server allocations
		allocCtrl := &userctrl.ServerAllocationController{}
		user.GET("/servers/:uuidShort/allocations", allocCtrl.GetServerAllocations)
		user.DELETE("/servers/:uuidShort/allocations/:allocationId", allocCtrl.DeleteAllocation)
		user.POST("/servers/:uuidShort/allocations/:allocationId/primary", allocCtrl.SetPrimaryAllocation)

		// Server schedules
		schedCtrl := &userctrl.ServerScheduleController{}
		user.GET("/servers/:uuidShort/schedules", schedCtrl.GetSchedules)
		user.POST("/servers/:uuidShort/schedules", schedCtrl.CreateSchedule)
		user.GET("/servers/:uuidShort/schedules/:scheduleId", schedCtrl.GetSchedule)
		user.PUT("/servers/:uuidShort/schedules/:scheduleId", schedCtrl.UpdateSchedule)
		user.DELETE("/servers/:uuidShort/schedules/:scheduleId", schedCtrl.DeleteSchedule)
		user.POST("/servers/:uuidShort/schedules/:scheduleId/toggle", schedCtrl.ToggleSchedule)
	}

	// ── Admin routes ───────────────────────────────────────────────────────────
	adminGroup := api.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware())
	adminGroup.Use(middleware.AdminMiddleware(""))
	{
		// Dashboard
		dash := &admin.DashboardController{}
		adminGroup.GET("/dashboard", dash.Index)

		// Users
		usersCtrl := &admin.UsersController{}
		adminGroup.GET("/users", usersCtrl.Index)
		adminGroup.PUT("/users", usersCtrl.Create)
		adminGroup.GET("/users/external/:externalId", usersCtrl.ShowByExternalID)
		adminGroup.GET("/users/:uuid", usersCtrl.Show)
		adminGroup.PATCH("/users/:uuid", usersCtrl.Update)
		adminGroup.DELETE("/users/:uuid", usersCtrl.Delete)
		adminGroup.GET("/users/:uuid/servers", usersCtrl.OwnedServers)
		adminGroup.POST("/users/:uuid/ban", usersCtrl.Ban)
		adminGroup.POST("/users/:uuid/unban", usersCtrl.Unban)
		adminGroup.POST("/users/:uuid/sso-token", usersCtrl.CreateSsoToken)

		// Nodes
		nodesCtrl := &admin.NodesController{}
		adminGroup.GET("/nodes", nodesCtrl.Index)
		adminGroup.PUT("/nodes", nodesCtrl.Create)
		adminGroup.GET("/nodes/:id", nodesCtrl.Show)
		adminGroup.PATCH("/nodes/:id", nodesCtrl.Update)
		adminGroup.DELETE("/nodes/:id", nodesCtrl.Delete)
		adminGroup.POST("/nodes/:id/reset-key", nodesCtrl.ResetKey)

		// Servers
		serversCtrl := &admin.ServersController{}
		adminGroup.GET("/servers", serversCtrl.Index)
		adminGroup.PUT("/servers", serversCtrl.Create)
		adminGroup.GET("/servers/external/:externalId", serversCtrl.ShowByExternalID)
		adminGroup.GET("/servers/node/:nodeId", serversCtrl.GetByNode)
		adminGroup.GET("/servers/owner/:ownerId", serversCtrl.GetByOwner)
		adminGroup.GET("/servers/:id", serversCtrl.Show)
		adminGroup.PATCH("/servers/:id", serversCtrl.Update)
		adminGroup.DELETE("/servers/:id", serversCtrl.Delete)
		adminGroup.DELETE("/servers/:id/hard", serversCtrl.HardDelete)
		adminGroup.POST("/servers/:id/suspend", serversCtrl.Suspend)
		adminGroup.POST("/servers/:id/unsuspend", serversCtrl.Unsuspend)

		// Settings
		settingsCtrl := &admin.SettingsController{}
		adminGroup.GET("/settings", settingsCtrl.Index)
		adminGroup.GET("/settings/categories", settingsCtrl.Categories)
		adminGroup.GET("/settings/category/:category", settingsCtrl.GetByCategory)
		adminGroup.GET("/settings/:setting", settingsCtrl.Show)
		adminGroup.PATCH("/settings", settingsCtrl.Update)

		// Roles
		rolesCtrl := &admin.RolesController{}
		adminGroup.GET("/roles", rolesCtrl.Index)
		adminGroup.PUT("/roles", rolesCtrl.Create)
		adminGroup.GET("/roles/:id", rolesCtrl.Show)
		adminGroup.PATCH("/roles/:id", rolesCtrl.Update)
		adminGroup.DELETE("/roles/:id", rolesCtrl.Delete)

		// Permissions
		permsCtrl := &admin.PermissionsController{}
		adminGroup.GET("/permissions", permsCtrl.Index)
		adminGroup.PUT("/permissions", permsCtrl.Create)
		adminGroup.GET("/permissions/:id", permsCtrl.Show)
		adminGroup.PATCH("/permissions/:id", permsCtrl.Update)
		adminGroup.DELETE("/permissions/:id", permsCtrl.Delete)

		// Locations
		locsCtrl := &admin.LocationsController{}
		adminGroup.GET("/locations", locsCtrl.Index)
		adminGroup.PUT("/locations", locsCtrl.Create)
		adminGroup.GET("/locations/:id", locsCtrl.Show)
		adminGroup.PATCH("/locations/:id", locsCtrl.Update)
		adminGroup.DELETE("/locations/:id", locsCtrl.Delete)

		// Allocations
		allocsCtrl := &admin.AllocationsController{}
		adminGroup.GET("/allocations", allocsCtrl.Index)
		adminGroup.PUT("/allocations", allocsCtrl.Create)
		adminGroup.GET("/allocations/available", allocsCtrl.GetAvailable)
		adminGroup.DELETE("/allocations/bulk-delete", allocsCtrl.BulkDelete)
		adminGroup.DELETE("/allocations/delete-unused", allocsCtrl.DeleteUnused)
		adminGroup.GET("/allocations/:id", allocsCtrl.Show)
		adminGroup.PATCH("/allocations/:id", allocsCtrl.Update)
		adminGroup.DELETE("/allocations/:id", allocsCtrl.Delete)

		// Realms
		realmsCtrl := &admin.RealmsController{}
		adminGroup.GET("/realms", realmsCtrl.Index)
		adminGroup.PUT("/realms", realmsCtrl.Create)
		adminGroup.GET("/realms/:id", realmsCtrl.Show)
		adminGroup.PATCH("/realms/:id", realmsCtrl.Update)
		adminGroup.DELETE("/realms/:id", realmsCtrl.Delete)

		// Mounts
		mountsCtrl := &admin.MountsController{}
		adminGroup.GET("/mounts", mountsCtrl.Index)
		adminGroup.PUT("/mounts", mountsCtrl.Create)
		adminGroup.GET("/mounts/:id", mountsCtrl.Show)
		adminGroup.PATCH("/mounts/:id", mountsCtrl.Update)
		adminGroup.DELETE("/mounts/:id", mountsCtrl.Delete)

		// Spells
		spellsCtrl := &admin.SpellsController{}
		adminGroup.GET("/spells", spellsCtrl.Index)
		adminGroup.PUT("/spells", spellsCtrl.Create)
		adminGroup.GET("/spells/realm/:realmId", spellsCtrl.GetByRealm)
		adminGroup.GET("/spells/:id", spellsCtrl.Show)
		adminGroup.PATCH("/spells/:id", spellsCtrl.Update)
		adminGroup.DELETE("/spells/:id", spellsCtrl.Delete)
		adminGroup.GET("/spells/:spellId/variables", spellsCtrl.ListVariables)
		adminGroup.POST("/spells/:spellId/variables", spellsCtrl.CreateVariable)
		adminGroup.PATCH("/spell-variables/:id", spellsCtrl.UpdateVariable)
		adminGroup.DELETE("/spell-variables/:id", spellsCtrl.DeleteVariable)

		// Images
		imagesCtrl := &admin.ImagesController{}
		adminGroup.GET("/images", imagesCtrl.Index)
		adminGroup.POST("/images", imagesCtrl.Create)
		adminGroup.GET("/images/:id", imagesCtrl.Show)
		adminGroup.PATCH("/images/:id", imagesCtrl.Update)
		adminGroup.DELETE("/images/:id", imagesCtrl.Delete)

		// Databases (host instances)
		dbsCtrl := &admin.DatabasesController{}
		adminGroup.GET("/databases", dbsCtrl.Index)
		adminGroup.PUT("/databases", dbsCtrl.Create)
		adminGroup.GET("/databases/node/:nodeId", dbsCtrl.GetByNode)
		adminGroup.GET("/databases/:id", dbsCtrl.Show)
		adminGroup.PATCH("/databases/:id", dbsCtrl.Update)
		adminGroup.DELETE("/databases/:id", dbsCtrl.Delete)

		// Tickets
		ticketsCtrl := &admin.AdminTicketsController{}
		adminGroup.GET("/tickets", ticketsCtrl.Index)
		adminGroup.GET("/tickets/:uuid", ticketsCtrl.Show)
		adminGroup.PATCH("/tickets/:uuid", ticketsCtrl.Update)
		adminGroup.DELETE("/tickets/:uuid", ticketsCtrl.Delete)
		adminGroup.POST("/tickets/:uuid/reply", ticketsCtrl.Reply)
		adminGroup.POST("/tickets/:uuid/close", ticketsCtrl.Close)
		adminGroup.POST("/tickets/:uuid/reopen", ticketsCtrl.Reopen)

		// Ticket categories
		ticketCatsCtrl := &admin.TicketCategoriesController{}
		adminGroup.GET("/tickets/categories", ticketCatsCtrl.Index)
		adminGroup.PUT("/tickets/categories", ticketCatsCtrl.Create)
		adminGroup.GET("/tickets/categories/:id", ticketCatsCtrl.Show)
		adminGroup.PATCH("/tickets/categories/:id", ticketCatsCtrl.Update)
		adminGroup.DELETE("/tickets/categories/:id", ticketCatsCtrl.Delete)

		// Ticket statuses
		ticketStatusCtrl := &admin.TicketStatusesController{}
		adminGroup.GET("/tickets/statuses", ticketStatusCtrl.Index)
		adminGroup.PUT("/tickets/statuses", ticketStatusCtrl.Create)
		adminGroup.GET("/tickets/statuses/:id", ticketStatusCtrl.Show)
		adminGroup.PATCH("/tickets/statuses/:id", ticketStatusCtrl.Update)
		adminGroup.DELETE("/tickets/statuses/:id", ticketStatusCtrl.Delete)

		// Ticket priorities
		ticketPrioCtrl := &admin.TicketPrioritiesController{}
		adminGroup.GET("/tickets/priorities", ticketPrioCtrl.Index)
		adminGroup.PUT("/tickets/priorities", ticketPrioCtrl.Create)
		adminGroup.GET("/tickets/priorities/:id", ticketPrioCtrl.Show)
		adminGroup.PATCH("/tickets/priorities/:id", ticketPrioCtrl.Update)
		adminGroup.DELETE("/tickets/priorities/:id", ticketPrioCtrl.Delete)

		// Ticket messages
		ticketMsgCtrl := &admin.AdminTicketMessagesController{}
		adminGroup.GET("/tickets/:uuid/messages", ticketMsgCtrl.Index)
		adminGroup.POST("/tickets/:uuid/messages", ticketMsgCtrl.Create)
		adminGroup.GET("/tickets/:uuid/messages/:id", ticketMsgCtrl.Show)
		adminGroup.PATCH("/tickets/:uuid/messages/:id", ticketMsgCtrl.Update)
		adminGroup.DELETE("/tickets/:uuid/messages/:id", ticketMsgCtrl.Delete)

		// Knowledgebase
		kbCtrl := &admin.KnowledgebaseController{}
		adminGroup.GET("/knowledgebase/categories", kbCtrl.CategoriesIndex)
		adminGroup.PUT("/knowledgebase/categories", kbCtrl.CategoriesCreate)
		adminGroup.GET("/knowledgebase/categories/:id", kbCtrl.CategoriesShow)
		adminGroup.PATCH("/knowledgebase/categories/:id", kbCtrl.CategoriesUpdate)
		adminGroup.DELETE("/knowledgebase/categories/:id", kbCtrl.CategoriesDelete)
		adminGroup.GET("/knowledgebase/articles", kbCtrl.ArticlesIndex)
		adminGroup.PUT("/knowledgebase/articles", kbCtrl.ArticlesCreate)
		adminGroup.GET("/knowledgebase/articles/:id", kbCtrl.ArticlesShow)
		adminGroup.PATCH("/knowledgebase/articles/:id", kbCtrl.ArticlesUpdate)
		adminGroup.DELETE("/knowledgebase/articles/:id", kbCtrl.ArticlesDelete)
		adminGroup.GET("/knowledgebase/articles/:id/tags", kbCtrl.GetTags)
		adminGroup.POST("/knowledgebase/articles/:id/tags", kbCtrl.CreateTag)
		adminGroup.DELETE("/knowledgebase/articles/:id/tags/:tagId", kbCtrl.DeleteTag)

		// Notifications
		notifsCtrl := &admin.NotificationsController{}
		adminGroup.GET("/notifications", notifsCtrl.Index)
		adminGroup.PUT("/notifications", notifsCtrl.Create)
		adminGroup.GET("/notifications/:id", notifsCtrl.Show)
		adminGroup.PATCH("/notifications/:id", notifsCtrl.Update)
		adminGroup.DELETE("/notifications/:id", notifsCtrl.Delete)

		// Subdomains
		subdomainsCtrl := &admin.AdminSubdomainsController{}
		adminGroup.GET("/subdomains", subdomainsCtrl.Index)
		adminGroup.PUT("/subdomains", subdomainsCtrl.Create)
		adminGroup.GET("/subdomains/:uuid", subdomainsCtrl.Show)
		adminGroup.PATCH("/subdomains/:uuid", subdomainsCtrl.Update)
		adminGroup.DELETE("/subdomains/:uuid", subdomainsCtrl.Delete)
		adminGroup.GET("/subdomains/:uuid/subdomains", subdomainsCtrl.ListSubdomains)
	}

	// ── Wings/remote routes ────────────────────────────────────────────────────
	wingsCtrl := &wings.WingsAdminController{}
	wingsGroup := api.Group("/wings/admin")
	wingsGroup.Use(middleware.AuthMiddleware())
	wingsGroup.Use(middleware.AdminMiddleware(""))
	{
		wingsGroup.GET("", wingsCtrl.Index)
		wingsGroup.GET("/node/:id/system", wingsCtrl.GetNodeInfo)
		wingsGroup.GET("/node/:id/servers", wingsCtrl.GetServers)
	}
}
