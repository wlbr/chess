package chess

import (
	"fmt"
	"strings"
	"time"
)

func ExportToPGN(moveLog *MoveLog, whitePlayer, blackPlayer string) string {
	var sb strings.Builder

	// PGN headers
	sb.WriteString(fmt.Sprintf("[Event \"%s\"]\n", "Casual Game"))
	sb.WriteString(fmt.Sprintf("[Site \"%s\"]\n", "Local"))
	sb.WriteString(fmt.Sprintf("[Date \"%s\"]\n", time.Now().Format("2006.01.02")))
	sb.WriteString(fmt.Sprintf("[Round \"%s\"]\n", "1"))
	sb.WriteString(fmt.Sprintf("[White \"%s\"]\n", whitePlayer))
	sb.WriteString(fmt.Sprintf("[Black \"%s\"]\n", blackPlayer))
	sb.WriteString(fmt.Sprintf("[Result \"%s\"]\n", "*")) // TODO: determine result
	sb.WriteString("\n")

	// Move text
	for i, move := range moveLog.moves {
		if i%2 == 0 {
			sb.WriteString(fmt.Sprintf("%d. ", i/2+1))
		}
		sb.WriteString(move.notation)
		sb.WriteString(" ")
	}

	return sb.String()
}
