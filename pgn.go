package chess

import (
	"fmt"
	"strings"
	"time"
)

func ExportToPGN(moveLog *MoveLog) string {
	var sb strings.Builder

	// PGN headers
	sb.WriteString(fmt.Sprintf("[Event \"%s\"]\n", "Casual Game"))
	sb.WriteString(fmt.Sprintf("[Site \"%s\"]\n", "Local"))
	sb.WriteString(fmt.Sprintf("[Date \"%s\"]\n", time.Now().Format("2006.01.02")))
	sb.WriteString(fmt.Sprintf("[Round \"%s\"]\n", "1"))
	sb.WriteString(fmt.Sprintf("[White \"%s\"]\n", "Player 1"))
	sb.WriteString(fmt.Sprintf("[Black \"%s\"]\n", "Player 2"))
	sb.WriteString(fmt.Sprintf("[Result \"%s\"]\n", "*")) // TODO: determine result
	sb.WriteString("\n")

	// Move text
	for i, move := range moveLog.moves {
		sb.WriteString(fmt.Sprintf("%d. ", i+1))
		sb.WriteString(move.notation)
		sb.WriteString(" ")
	}

	return sb.String()
}
