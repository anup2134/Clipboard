from PyQt6.QtCore import QSize, Qt
from PyQt6.QtWidgets import (
    QWidget, QLabel, QPushButton, QHBoxLayout
)
import subprocess

class ClickableLabel(QLabel):
    def __init__(self,text:str,main_self):
        super().__init__(text=text)
        self.main_self = main_self

    def mousePressEvent(self, event):
        subprocess.run(['wl-copy'],input=self.text().encode())
        super().mousePressEvent(event)
        self.main_self.close()

def TopBarWidget(main_self):
    top_bar = QHBoxLayout()
    heading = QLabel("Clipboard")
    heading.setStyleSheet("""
        QLabel {
            padding: 0px;
            color: white;
            font-size: 24px;
        }
    """)
    
    exit_button = QPushButton("✕")
    exit_button.setFixedSize(30, 30)
    exit_button.setStyleSheet("""
        QPushButton {
            color: white;
            background-color: transparent;
            border: none;
            font-size: 18px;
        }
        QPushButton:hover {
            color: red;
        }
    """)
    exit_button.clicked.connect(main_self.close)

    top_bar.addWidget(heading)
    top_bar.addStretch()
    top_bar.addWidget(exit_button)

    top_bar_widget = QWidget()
    top_bar_widget.setLayout(top_bar)

    top_bar_widget.setStyleSheet("""
        border-bottom: rgba(255, 255, 255, 0.2);
    """)

    return top_bar_widget


def get_label(text:str,main_self):
    label = ClickableLabel(text=text,main_self=main_self)
    label.setWordWrap(True)
    label.setFixedSize(QSize(400, 93))
    label.setAlignment(Qt.AlignmentFlag.AlignTop | Qt.AlignmentFlag.AlignLeft)
    label.setMouseTracking(True)
    label.setStyleSheet("""
        QLabel {
            color: white;
            font-size: 18px;
            padding: 0px;
            margin: 0px; 
            border-bottom: 1px solid rgba(255,255,255,0.2);
        }
        QLabel:hover {
            background-color: rgba(93, 110, 122, 0.3);
        }
    """)


    return label
