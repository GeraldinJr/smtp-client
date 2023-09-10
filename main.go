package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
)

func main() {
	smtpServer := "sandbox.smtp.mailtrap.io"
	smtpPort := 587
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	from := "remetente@email.com"
	to := "destiatario@email.com"
	subject := "Explorando o SMTP"
	body := "Prezad@s, \n\n\n A nuvem não é real, são apenas servidores. 4566"

	addr := fmt.Sprintf("%s:%d", smtpServer, smtpPort)
	host, _, _ := net.SplitHostPort(addr)
	conn, err := net.Dial("tcp", addr)

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor SMTP:", err)
		return
	}
	defer client.Close()

	// Configure a autenticação
	auth := smtp.PlainAuth("", username, password, smtpServer)

	// Connect to the SMTP server.
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer,
	}

	if err := client.StartTLS(tlsconfig); err != nil {
		fmt.Println("Erro ao iniciar a conexão segura TLS:", err)
		return
	}

	// Cumprimente o servidor
	client.Hello("localhost")

	// Autentique-se com o servidor SMTP
	if err := client.Auth(auth); err != nil {
		fmt.Println("Erro na autenticação SMTP:", err)
		return
	}

	// Defina o remetente e destinatário
	err = client.Mail(from)
	if err != nil {
		fmt.Println("Erro ao adicionar remetente:", err)
		return
	}
	err = client.Rcpt(to)
	if err != nil {
		fmt.Println("Erro ao adicionar destinatário:", err)
		return
	}

	// Prepare a mensagem
	data, err := client.Data()
	if err != nil {
		fmt.Println("Erro ao preparar a mensagem:", err)
		return
	}
	defer data.Close()

	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)
	_, err = data.Write([]byte(message))
	if err != nil {
		fmt.Println("Erro ao escrever a mensagem:", err)
		return
	}

	// Conclua a transação de envio
	defer client.Quit()

	fmt.Println("E-mail enviado com sucesso!")
}
