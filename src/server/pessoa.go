package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/diegolikescode/rinha2023q1_sql/src/config"
	"github.com/diegolikescode/rinha2023q1_sql/src/model"
	"github.com/gofiber/fiber/v3"
)

func validaCorpoPessoa(p model.Pessoa) bool {
    if p.Apelido == "" || p.Nome == "" || p.Nascimento == "" {
	return false
    }

    return true
}

func PessoaToCadastro(p model.Pessoa) model.PessoaCadastro {
    saveStack := strings.Join(p.Stack, ";")
    model := model.PessoaCadastro{
	Apelido: p.Apelido,
	Nome: p.Nome,
	Nascimento: p.Nascimento,
	Stack: saveStack,
	CampoBusca: strings.Join([]string{saveStack, p.Apelido, p.Nome}, ";"),
    }

    return model
}

func CriarPessoa (c fiber.Ctx) error {
    var body model.Pessoa
    var err error
    if  err = json.Unmarshal(c.Body(), &body); err != nil {
	fmt.Println("parsing error", err.Error())
	return c.SendStatus(http.StatusBadRequest)
    }

    if !validaCorpoPessoa(body) {
	return c.SendStatus(http.StatusUnprocessableEntity)
    }

    m:= PessoaToCadastro(body)

    conn := config.Connection
    row, err := conn.Query("SELECT * FROM insert_pessoa($1, $2, $3, $4, $5)", m.Apelido, m.Nome, m.Nascimento, m.Stack, m.CampoBusca); if err != nil {
	return c.SendStatus(http.StatusBadRequest)
    }

    var id string
    for row.Next() {
	if err := row.Scan(&id); err != nil {
	    return c.SendStatus(http.StatusBadRequest)
	}
    }

    row.Close()
    c.Set("Location", "/pessoas/"+id)
    return c.SendStatus(http.StatusCreated)
}

func DetalharPessoa (c fiber.Ctx) error {
    var err error

    id := c.Params("id")
    conn := config.Connection
    row, err := conn.Query("SELECT * FROM select_by_id($1)", id); 
    if err != nil {
	return c.SendStatus(http.StatusNotFound)
    }

    var p model.Pessoa
    var stackConvert string
    for row.Next() {
	row.Scan(&p.Id, &p.Apelido, &p.Nome, &p.Nascimento, &stackConvert)
	p.Stack = strings.Split(stackConvert, ";")
    }

    c.SendStatus(http.StatusOK)
    return c.JSON(p)
}

func BuscarTermoPessoa (c fiber.Ctx) error {
    var err error

    t := c.Query("t", "")
    if t == "" {
	return c.SendStatus(http.StatusBadRequest)
    }
    conn := config.Connection
    row, err := conn.Query("SELECT * FROM select_by_t($1)", t); 
    if err != nil {
	c.JSON([]string{})
	return c.SendStatus(http.StatusOK)
    }

    var pessoas []model.Pessoa
    var stackConvert string
    for row.Next() {
	var p model.Pessoa
	row.Scan(&p.Id, &p.Apelido, &p.Nome, &p.Nascimento, &stackConvert)
	p.Stack = strings.Split(stackConvert, ";")
	pessoas = append(pessoas, p)
    }

    c.SendStatus(http.StatusOK)
    return c.JSON(pessoas)
}

func ContagemPessoas (c fiber.Ctx) error {
    var err error

    c.Set("Content-Type", "text/html")
    c.SendStatus(http.StatusOK)

    conn := config.Connection
    var count int
    row, err := conn.Query("SELECT * FROM count_pessoas()"); 
    if err != nil {
	c.SendString(string(count))
    }

    row.Scan(&count)
    return c.SendString(string(count))
}

