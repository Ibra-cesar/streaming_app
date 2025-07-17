package helper

import (
	"fmt"
	"strings"

	"github.com/Ibra-cesar/video-streaming/src/internal/query_repo"
)

// Print user from users table
func PrintUser(user []query_repo.User) {
	if len(user) == 0 {
		fmt.Println("No Users found")
		return
	}

	tab_head := []string{"ID", "Name", "Email", "Is Admin", "created At", "Updated At"}

	colW := make(map[string]int)
	for _, h := range tab_head {
		colW[h] = len(h)
	}

	for _, user := range user {
		colW["ID"] = max(colW["ID"], len(user.ID.String()))
		colW["Name"] = max(colW["Name"], len(user.Name))
		colW["Email"] = max(colW["Email"], len(user.Email))
	}

	for _, user := range user {
		row := []string{
			user.ID.String(),
			user.Name,
			user.Email,
		}
		fmt.Println(formatRow(row, colW))
	}
}

func formatRow(data []string, widths map[string]int) string {
	var sb strings.Builder
	sb.WriteString("|")
	sb.WriteString(fmt.Sprintf(" %-*s |", widths["ID"], data[0]))
	sb.WriteString(fmt.Sprintf(" %-*s |", widths["Name"], data[1]))
	sb.WriteString(fmt.Sprintf(" %-*s |", widths["Email"], data[2]))
	return sb.String()
}
