package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	//dotenv initialization
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//connect server by ssh and download what i need
	downloadBackups()
}

func downloadBackups() {

	user := os.Getenv("SSH_USER")
	//pass := ""
	remote := os.Getenv("REMOTE_SERVER_HOST")
	port := os.Getenv("REMOTE_SERVER_PORT")

	hostKeyCallback, err := kh.New(os.Getenv("FOLDER_KNOWN_HOSTS"))
	if err != nil {
		log.Fatal("could not create hostkeycallback function: ", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(os.Getenv("RSA_KEY")),
		},
		HostKeyCallback: hostKeyCallback,
	}

	// connect
	conn, err := ssh.Dial("tcp", remote+":"+port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// create new SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	/////// start of files downloding ///////////

	//first file
	// create destination file
	dstFile, err := os.Create(os.Getenv("FILE_TO_DOWNLOAD1"))
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	// open source file
	srcFile, err := client.Open(os.Getenv("FILE_TO_DOWNLOAD2"))
	if err != nil {
		log.Fatal(err)
	}

	// copy source file to destination file
	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d bytes copied\n", bytes)

	// flush in-memory copy
	err = dstFile.Sync()
	if err != nil {
		log.Fatal(err)
	}

	//second file
	// create destination file
	dstFile2, err := os.Create(os.Getenv("FILE_TO_DOWNLOAD3"))
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile2.Close()

	// open source file
	srcFile2, err := client.Open(os.Getenv("FILE_TO_DOWNLOAD4"))
	if err != nil {
		log.Fatal(err)
	}

	// copy source file to destination file
	bytes, err = io.Copy(dstFile2, srcFile2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d bytes copied\n", bytes)

	// flush in-memory copy
	err = dstFile2.Sync()
	if err != nil {
		log.Fatal(err)
	}
	/////// end of files downloding ///////////
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}

	return ssh.PublicKeys(key)
}
