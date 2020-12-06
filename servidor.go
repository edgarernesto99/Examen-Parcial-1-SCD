package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"net"
	"container/list"
)

var mensajes list.List
var clientes list.List
var idGlobal uint64
var mostrar bool = false

type Archivo struct {
	BS []byte
	Name string
	UserName string
}

func servidor() {
	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleCliente(c)
	}
}

func handleCliente(c net.Conn) {
	var op uint64
	var err error
	nuevoCliente(c)
	for {
		err = gob.NewDecoder(c).Decode(&op)
		if err != nil {
			//fmt.Println(err)
			continue
		}
		if op == 1 {
			recibirMensaje(c)
		} else if op == 2 {
			recibirArchivo(c)
		} else if op == 3 {
			desconectarCliente(c)
			return
		} else if op == 0 {
			nuevoCliente(c)
		}
	}
}

func nuevoCliente(c net.Conn) {
	var msg string
	err := gob.NewDecoder(c).Decode(&msg)
	if err != nil {
		fmt.Println(err)
	}
	clientes.PushBack(c)
	msg = "Se conecto "+msg
	if (mostrar) {
		fmt.Println(msg)
	}
	mensajes.PushBack(msg)
}

func recibirMensaje(c net.Conn) {
	var msg string
	err := gob.NewDecoder(c).Decode(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	if mostrar {
		fmt.Println(msg)
	}
	mensajes.PushBack(msg)
	enviarMensajeATodos(msg, c)
}

func enviarMensajeATodos(msg string, c net.Conn) {
	var op uint64 = 1
	for e:=clientes.Front(); e != nil; e = e.Next() {
		if e.Value.(net.Conn) != c {
			err := gob.NewEncoder(e.Value.(net.Conn)).Encode(op)
			if err != nil {
				//fmt.Println(err)
				continue
			}
			err = gob.NewEncoder(e.Value.(net.Conn)).Encode(msg)
			if err != nil {
				//fmt.Println(err)
				continue
			}
		}
	}
}

func recibirArchivo(c net.Conn) {
	var archivo Archivo
	err := gob.NewDecoder(c).Decode(&archivo)
	if err != nil {
		fmt.Println(err)
		return
	}
	//Guardar archivo
	file, err := os.Create(archivo.Name)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	file.Write(archivo.BS)

	msg := archivo.UserName+": "+archivo.Name
	mensajes.PushBack(msg)
	if (mostrar) {
		fmt.Println(msg)
	}
	//Mandar archivo a todos
	enviarArchivoATodos(c, archivo)
}

func enviarArchivoATodos(c net.Conn, archivo Archivo) {
	var op uint64 = 2
	for e:=clientes.Front(); e != nil; e = e.Next() {
		if e.Value.(net.Conn) != c {
			err := gob.NewEncoder(e.Value.(net.Conn)).Encode(op)
			if err != nil {
				//fmt.Println(err)
				continue
			}
			err = gob.NewEncoder(e.Value.(net.Conn)).Encode(archivo)
			if err != nil {
				//fmt.Println(err)
				continue
			}
		}
	}
}

func mostrarMensajes() {
	for e:=mensajes.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func desconectarCliente(c net.Conn) {
	for e:=clientes.Front(); e != nil; e = e.Next() {
		if e.Value.(net.Conn) == c {
			clientes.Remove(e)
			return
		}
	}
}

func respaldarMensajes() {
	//Respaldar en archivo txt
	file, err := os.Create("Respaldo de mensajes.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	for e:=mensajes.Front(); e != nil; e = e.Next() {
		file.WriteString(e.Value.(string)+ "\n")
	}
}

func main() {
	go servidor()
	var op uint64
	for {
		fmt.Println(" ============================ ")
		fmt.Println("|            MENU            |")
		fmt.Println("|============================|")
		fmt.Println("| 1. Mostrar mensajes        |")
		fmt.Println("| 2. Respaldar mensajes      |")
		fmt.Println("| 0. Terminar servidor       |")
		fmt.Println(" ============================ ")
		fmt.Print("Opción: ")
		fmt.Scanln(&op)
		if (op == 1) {
			mostrar = true
			mostrarMensajes()
			fmt.Scanln(&op)
			mostrar = false
		} else if (op == 2) {
			respaldarMensajes();
		} else if (op == 0) {
			break
		} else {
			fmt.Println("Opcion no valida")
		}
	}
}