# **📋 Clipboard**

_A super simple, lightweight, and easy-to-use clipboard manager for X11-based Linux desktops. Built with Go and Python._

---

## Prerequisites
- **Python** 3.13 or later
- **Golang** 1.24 or later
- **Linux (amd64 architecture)**
- **X11** graphics manager (not Wayland)

_Check your installed versions:_
```bash
python --version
go version
```

## **Installation**
### 1) Clipboard GUI installation: 
```bash
git clone https://github.com/anup2134/Clipboard
cd clipboard/ui
python -m venv .
pip install -r requirements.txt
```

Create a keyboard shortcut that runs the following command to open the Clipboard GUI:
```bash
<path-to-clipboard-folder>/ui/bin/python <path-to-clipboard-folder>/ui/main_window.py
```
For example, if you cloned the repository into your home directory (~), the shortcut command would be:
```bash
~/Clipboard/ui/bin/python ~/Clipboard/ui/main_window.py
```

### 2) Clipboard background process:
```bash
cd <path-to-clipboard-folder>/backend
go build -o /home/user/.clipboard-daemon .
```

Add the following line to your ~/.xprofile (_create it if it doesn't exist_):
```bash
/home/anup/.clipboard-daemon >> /home/anup/.clipboard-daemon.log 2>&1 &
```

_Note: This will start the clipboard daemon on X11 session startup. Logs will be saved to ~/.clipboard-daemon.log._
