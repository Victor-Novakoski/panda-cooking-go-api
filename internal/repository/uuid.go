package repository

import "regexp"

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// isUUID evita mandar ao Postgres um id que não é UUID: ele responderia com
// erro de sintaxe (500) em vez de "não encontrado" (404).
func isUUID(s string) bool { return uuidPattern.MatchString(s) }
