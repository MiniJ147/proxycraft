# prototype proxy system going to reimplement for actual use

import socket, threading, sys

if len(sys.argv)<3:
    print("please provide python3 minecraft-proxy.py {minecraft-ip} {minecraft-port}")
    exit() 

dest = (sys.argv[1],int(sys.argv[2]))
host = ('127.0.0.1',12345)
print(dest)
PACKET_SIZE = 1024

def endpoint_recv(conn: socket.socket, endpoint: socket.socket):
    print("setup endpoint_recv")
    while True:
        data = endpoint.recv(PACKET_SIZE)
        if len(data)==0:
            break
        print("middle man got packet from dest")
        conn.send(data)

def connection(conn: socket.socket,addr):
    print("setup conn")
    endpoint = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
    endpoint.connect(dest)
    
    threading.Thread(target=endpoint_recv,args=(conn,endpoint,)).start()
    while True:
        data = conn.recv(PACKET_SIZE)
        if len(data) == 0:
            break
        print(addr,"middle man got packet")
        endpoint.sendall(data)
    endpoint.close()
    conn.close()
   
print("middle man")


tcp = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
tcp.bind(host)
tcp.listen()

while True:
    conn,addr = tcp.accept()
    print("new connection",addr)
    threading.Thread(target=connection,args=(conn,addr,)).start()
       

