import socket, threading

PACKET_SIZE = 1024
middleman = ("inital.minics.dev",3000)
dest = ("192.168.1.145",25565)

# direct packets from middleman to dest
middle_tcp = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
middle_tcp.connect(middleman)

#direct packets out from dest to middleman
dest_tcp = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
dest_tcp.connect(dest)

def dest_handle(conn: socket.socket, endpoint: socket.socket):
    while True:
        print("waiting for dest to send data")
        data = endpoint.recv(PACKET_SIZE)
        if len(data) == 0:
            print("breaking out")
            break
        print("sending packet back to middle man",len(data))
        conn.sendall(data)
    exit()

# ip = middle_tcp.recv(PACKET_SIZE)
# if len(ip) == 0:
#     print("failed to get ip name")
#     exit()
#
# print(str(ip))

threading.Thread(target=dest_handle, args=(middle_tcp, dest_tcp,)).start()

while True:
    data = middle_tcp.recv(PACKET_SIZE)
    if len(data) == 0:
        break
    print("sending packet to dest",len(data))
    dest_tcp.sendall(data)

middle_tcp.close()
dest_tcp.close()
