package service

import (
	"net/url"
	"regexp"
	"strings"
)

// MaxTimeMinutes é o maior tempo de preparo aceito: dois dias, para caber
// massa que descansa e carne que marina de um dia para o outro.
const MaxTimeMinutes = 2880

var youtubeID = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

// YouTubeID tira o id do vídeo de um link do YouTube, nos formatos que o
// site e o aplicativo entregam ao compartilhar: watch?v=, youtu.be/,
// /shorts/, /embed/ e /live/. Qualquer outro endereço é recusado: o front
// monta o player com esse id, e só o YouTube pode aparecer embutido.
func YouTubeID(raw string) (string, bool) {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Port() != "" {
		return "", false
	}

	var id string
	switch strings.ToLower(u.Hostname()) {
	case "youtu.be":
		id = strings.TrimPrefix(u.Path, "/")
	case "youtube.com", "www.youtube.com", "m.youtube.com", "www.youtube-nocookie.com":
		if u.Path == "/watch" {
			id = u.Query().Get("v")
			break
		}
		for _, prefix := range []string{"/shorts/", "/embed/", "/live/"} {
			if rest, ok := strings.CutPrefix(u.Path, prefix); ok {
				id = rest
				break
			}
		}
	}
	id = strings.TrimSuffix(id, "/")
	if !youtubeID.MatchString(id) {
		return "", false
	}
	return id, true
}

// canonicalVideoURL guarda sempre o mesmo formato, sem os parâmetros de
// rastreio que vêm no link compartilhado (?si=, &t=, &list=). Vazio fica
// vazio: a receita não tem vídeo.
func canonicalVideoURL(raw string) string {
	id, ok := YouTubeID(raw)
	if !ok {
		return ""
	}
	return "https://www.youtube.com/watch?v=" + id
}
