package main

import (
	"fmt"
	"os"

	"github.com/mythicalsystems/nemesis-backend/internal/config"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/internal/services"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		die("Failed to load config: %s", err)
	}
	if err := database.Init(cfg); err != nil {
		die("Failed to connect to database: %s", err)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "user:make-admin":
		cmdMakeAdmin(args)
	case "user:list":
		cmdListUsers()
	case "user:password":
		cmdSetPassword(args)
	case "setting:set":
		cmdSetSetting(args)
	case "setting:get":
		cmdGetSetting(args)
	case "setting:list":
		cmdListSettings()
	default:
		fmt.Printf("Unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Nemesis Panel CLI

Usage:
  nemesis-cli <command> [arguments]

Commands:
  user:make-admin <email>          Grant role_id=1 (super-admin) to a user
  user:list                        List all users
  user:password <email> <password> Reset a user's password

  setting:set <key> <value>        Set a panel setting
  setting:get <key>                Get a panel setting value
  setting:list                     List all panel settings`)
}

// ── user:make-admin ─────────────────────────────────────────────────────────

func cmdMakeAdmin(args []string) {
	if len(args) < 1 {
		die("Usage: user:make-admin <email>")
	}
	email := args[0]
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		die("User not found: %s", email)
	}

	// Ensure role_id 1 (root role) exists
	var role models.Role
	if err := database.DB.First(&role, 1).Error; err != nil {
		role = models.Role{
			Name:        "admin",
			DisplayName: "Administrator",
			Color:       "#ef4444",
		}
		if err := database.DB.Create(&role).Error; err != nil {
			die("Failed to create admin role: %s", err)
		}
	}

	if err := database.DB.Model(&user).Update("role_id", 1).Error; err != nil {
		die("Failed to update user role: %s", err)
	}
	fmt.Printf("✓ User %s (%s) is now a super-admin (role_id=1)\n", user.Username, user.Email)
}

// ── user:list ───────────────────────────────────────────────────────────────

func cmdListUsers() {
	var users []models.User
	database.DB.Find(&users)
	fmt.Printf("%-5s %-20s %-30s %-8s %s\n", "ID", "Username", "Email", "Role ID", "Banned")
	fmt.Println("----------------------------------------------------------------------")
	for _, u := range users {
		fmt.Printf("%-5d %-20s %-30s %-8d %s\n", u.ID, u.Username, u.Email, u.RoleID, u.Banned)
	}
}

// ── user:password ───────────────────────────────────────────────────────────

func cmdSetPassword(args []string) {
	if len(args) < 2 {
		die("Usage: user:password <email> <new-password>")
	}
	email, password := args[0], args[1]
	if len(password) < 8 {
		die("Password must be at least 8 characters")
	}
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		die("User not found: %s", email)
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		die("Failed to hash password: %s", err)
	}
	database.DB.Model(&user).Update("password", hash)
	fmt.Printf("✓ Password updated for %s\n", user.Email)
}

// ── setting:set ─────────────────────────────────────────────────────────────

func cmdSetSetting(args []string) {
	if len(args) < 2 {
		die("Usage: setting:set <key> <value>")
	}
	key, value := args[0], args[1]
	services.SetSetting(key, value)
	fmt.Printf("✓ %s = %s\n", key, value)
}

// ── setting:get ─────────────────────────────────────────────────────────────

func cmdGetSetting(args []string) {
	if len(args) < 1 {
		die("Usage: setting:get <key>")
	}
	key := args[0]
	val := services.GetSetting(key, "<not set>")
	fmt.Printf("%s = %s\n", key, val)
}

// ── setting:list ────────────────────────────────────────────────────────────

func cmdListSettings() {
	var settings []models.Setting
	database.DB.Find(&settings)
	fmt.Printf("%-40s %s\n", "Key", "Value")
	fmt.Println("----------------------------------------------------------------------")
	for _, s := range settings {
		val := s.Value
		if len(val) > 60 {
			val = val[:57] + "..."
		}
		fmt.Printf("%-40s %s\n", s.Key, val)
	}
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}
