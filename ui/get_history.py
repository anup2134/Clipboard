import socket
import json
socket_path = "/tmp/clipboard.sock"
client = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)

def get_history():
    try:
        client.connect(socket_path)
        client.sendall(b"get_history")
        chunks = []
        while True:
            chunk = client.recv(1024)
            if not chunk:
                break
            chunks.append(chunk)

        data = b''.join(chunks)
        history:list[str] = json.loads(data)
        return history
    except Exception as e:
        print(f"error in connection:{e}")
    finally:
        client.close()