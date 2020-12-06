package main

import (
	"encoding/gob"
	"fmt"
	"bufio"
	"os"
	"net"
	"container/list"
)

var nombre string
var mensajes list.List
var mostrar bool = false

type Archivo struct {
	BS []byte
	Name string
	UserName string
}

func cliente() {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(c).Encode(nombre)
	if err != nil {
		fmt.Println(err)
	}
	go recibirMensajes(c)
	accionesCliente(c)
	c.Close()
}

func accionesCliente(c net.Conn) {
	var op string
	scaner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println(" ============================ ")
		fmt.Println("|            MENU            |")
		fmt.Println("|============================|")
		fmt.Println("| 1. Enviar mensaje          |")
		fmt.Println("| 2. Enviar archivo          |")
		fmt.Println("| 3. Mostrar mensajes        |")
		fmt.Println("| 0. Salir                   |")
		fmt.Println(" ============================ ")
		fmt.Print("Opción: ")
		scaner.Scan()
		op = scaner.Text()
		if op == "1" {
			var msg string
			fmt.Print("Mensaje: ")
			scaner.Scan()
			msg = scaner.Text()
			var submenu uint64 = 1
			err := gob.NewEncoder(c).Encode(submenu)
			if err != nil {
				fmt.Println(err)
			} else {
				mensajes.PushBack("Tú: "+msg)
				msg := nombre+": "+msg
				gob.NewEncoder(c).Encode(msg)
			}
		} else if op == "2" {
			var msg string
			fmt.Print("Ruta: ")
			scaner.Scan()
			msg = scaner.Text()
			enviarArchivo(c, msg)
		} else if op == "3" {
			mostrar = true
			mostrarMensajes()
			scaner.Scan()
			mostrar = false
		} else if op == "0" {
			var submenu uint64 = 3
			err := gob.NewEncoder(c).Encode(submenu)
			if err != nil {
				fmt.Println(err)
			}
			break
		} else {
			fmt.Println("Opcion no valida")
		}
	}
}

func mostrarMensajes() {
	fmt.Println()
	for e:=mensajes.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func recibirMensajes(c net.Conn) {
	var op uint64
	var msg string
	for {
		err := gob.NewDecoder(c).Decode(&op)
		if err != nil {
			//fmt.Println(err)
			continue
		}
		if (op == 1) {
			err = gob.NewDecoder(c).Decode(&msg)
			if err != nil {
				//fmt.Println(err)
				continue
			}
			mensajes.PushBack(msg)
			if (mostrar) {
				fmt.Println(msg)
			} else {
				mostrarMensajes()
			}
		} else if (op == 2) {
			recibirArchivo(c)
		}
	}
}

func enviarArchivo(c net.Conn, route string) {
	file, err := os.Open(route)
	if (err != nil) {
		fmt.Println(err)
		return
	}
	defer file.Close()

	//Para ver el estado del archivo (tamaño , etc)
	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	total := stat.Size() //tamaño del contenido
	bs := make([]byte, total) //slice de bytes
	count, err := file.Read(bs)
	if err != nil{
		fmt.Println(err,count)
		return
	}
	//Enviar slice de bytes
	nameFile := file.Name()

	var submenu uint64 = 2
	err = gob.NewEncoder(c).Encode(submenu)
	if err != nil {
		fmt.Println(err)
		return
	} 
	archivo := Archivo{BS: bs, Name: nameFile, UserName: nombre}
	terminarEnvioArchivo(c, archivo)
}

func terminarEnvioArchivo(c net.Conn, archivo Archivo) {
	err := gob.NewEncoder(c).Encode(&archivo)
	if err != nil {
		fmt.Println(err)
	} else {
		mensajes.PushBack("Tú: "+archivo.Name)
	}
}

//Aqui se redirecciona cuando labandera indica archivo
func recibirArchivo(c net.Conn) {
	var archivo Archivo
	err := gob.NewDecoder(c).Decode(&archivo)
	if err != nil {
		//fmt.Println(err)
		return
	}
	//Guardar archivo
	file, err := os.Create(archivo.Name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	file.Write(archivo.BS)

	//Guardar mensaje
	msg := archivo.UserName+": "+archivo.Name
	mensajes.PushBack(msg)
	if (mostrar) {
		fmt.Println(msg)
	} else {
		mostrarMensajes()
	}
}

func main() {
	scaner := bufio.NewScanner(os.Stdin)
	fmt.Print("Nombre: ")
	scaner.Scan()
	nombre = scaner.Text()
	cliente()
}
