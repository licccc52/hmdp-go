package middleware

import (
	"hmdp/src/dto"
	"testing"
)

func TestGeratateToken(t *testing.T) {
	j := NewJWT()
	var userDTO dto.UserDTO
	userDTO.Icon = "hello"
	userDTO.Id = 1
	userDTO.NickName = "ni"

	clamis := j.CreateClaims(userDTO)
	token , err := j.CreateToken(clamis)

	if err != nil {
		t.Fatal("expected no err")
	}

	if token == "" {
		t.Fatal("expected a token , but get a empty token")
	}

	t.Log(token)
}



func TestParseToken(t *testing.T) {
	str := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwibmlja05hbWUiOiJuaSIsImljb24iOiJoZWxsbyIsIkJ1ZmZlclRpbWUiOjg2NDAwLCJpc3MiOiJsb3NlciIsImV4cCI6MTc0MTE1MTM5OSwibmJmIjoxNzQwNTQ1NTk5LCJpYXQiOjE3NDA1NDY1OTl9._gWvLZEO415u6zKm4EISXdsmP_ik4trurN1C_qw1ZqA"
	j := NewJWT()	
	clamis , err := j.ParseToken(str)
	if err != nil {
		t.Fatal("err!")
	}
	user := clamis.UserDTO
	t.Log(user)
}
