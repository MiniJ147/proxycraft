import socket, time


s = socket.socket(socket.AF_INET, socket.SOCK_STREAM) 

s.connect(("127.0.0.1",25566))
s.send(bytes([123]))
data = s.recv(1024)

print(str(data))


