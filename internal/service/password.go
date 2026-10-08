package service

import (
	"strings"
	"unicode/utf8"
)

// MaxPasswordBytes é o limite do bcrypt: o que passa de 72 bytes seria
// ignorado (o Go recusa com erro). Letra com acento ocupa 2 bytes.
const MaxPasswordBytes = 72

// commonPasswords são senhas de 10+ caracteres que aparecem no topo das
// listas de vazamentos (as menores já caem no mínimo de 10). A senha dos
// usuários de demonstração entra porque está no README.
var commonPasswords = map[string]bool{
	"0000000000": true, "0123456789": true, "1111111111": true, "1122334455": true,
	"1212121212": true, "1231231231": true, "123123123123": true, "1234567890": true,
	"12345678910": true, "1234567890a": true, "123456789a": true, "123456789abc": true,
	"123mudar123": true, "1q2w3e4r5t": true, "1q2w3e4r5t6y": true, "1qaz2wsx3edc": true,
	"2222222222": true, "9876543210": true, "a123456789": true, "abc1234567": true,
	"abcd123456": true, "abcdefghij": true, "admin12345": true, "admin123456": true,
	"administrador": true, "administrator": true, "asdfghjkl1": true, "asdfghjkl12": true,
	"batman1234": true, "botafogo123": true, "brasil1234": true, "brasil12345": true,
	"changeme123": true, "corinthians": true, "corinthians1": true, "cruzeiro123": true,
	"deus123456": true, "deusefiel": true, "dragon12345": true, "flamengo123": true,
	"flamengo1234": true, "fluminense": true, "football123": true, "gremio1234": true,
	"iloveyou12": true, "iloveyou123": true, "internacional": true, "jesus12345": true,
	"letmein123": true, "minhasenha": true, "minhasenha1": true, "minhasenha123": true,
	"monkey12345": true, "mudar12345": true, "mudar123456": true, "palmeiras1": true,
	"palmeiras123": true, "panda-cooking": true, "panda-cooking-demo": true, "panda123456": true,
	"pandacooking": true, "passw0rd123": true, "password12": true, "password123": true,
	"password1234": true, "princesa123": true, "q1w2e3r4t5": true, "qazwsxedc1": true,
	"qwerty1234": true, "qwerty12345": true, "qwerty123456": true, "qwertyuiop": true,
	"receitas123": true, "santos1234": true, "saopaulo123": true, "senha12345": true,
	"senha123456": true, "senha1234567": true, "senhaforte": true, "senhaforte123": true,
	"senhasenha": true, "starwars12": true, "sunshine12": true, "superman12": true,
	"teamo12345": true, "trocar1234": true, "vasco12345": true, "welcome123": true,
	"zxcvbnm123": true,
}

// checkPassword completa o que a validação do corpo já conferiu (mínimo de
// 10 caracteres): limite do bcrypt, senha comum e senha igual ao e-mail.
func checkPassword(password, email string) error {
	if len(password) > MaxPasswordBytes {
		return fieldError(KindInvalid, "password", "use no máximo 72 bytes (letra com acento conta como 2)")
	}

	lower := strings.ToLower(password)
	if commonPasswords[lower] || sameChar(lower) {
		return fieldError(KindInvalid, "password", "senha muito comum, escolha outra")
	}

	if email != "" {
		local, _, _ := strings.Cut(email, "@")
		if lower == strings.ToLower(email) || lower == strings.ToLower(local) {
			return fieldError(KindInvalid, "password", "a senha não pode ser o seu e-mail")
		}
	}
	return nil
}

// sameChar diz se o texto é um caractere só repetido ("aaaaaaaaaa").
func sameChar(s string) bool {
	first, _ := utf8.DecodeRuneInString(s)
	return strings.Trim(s, string(first)) == ""
}
