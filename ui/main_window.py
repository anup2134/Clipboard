from PyQt6.QtCore import Qt, QSize
from PyQt6.QtWidgets import (
    QApplication, QWidget, QVBoxLayout,
    QScrollArea, QMainWindow
)
from widgets import get_label, TopBarWidget
from get_history import get_history

app = QApplication([])

class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        content_widget = QWidget()
        
        layout = QVBoxLayout()
        layout.setSpacing(15) 

        top_bar_widget = TopBarWidget(self)

        history = get_history()
        if not history:
            layout.addWidget(get_label("Empty history",self))
        else:
            print(len(history))
            for h in reversed(history):
                label = get_label(h,self)
                layout.addWidget(label)


        content_widget.setLayout(layout)

        scroll_area = QScrollArea()
        scroll_area.setWidgetResizable(True)
        scroll_area.setWidget(content_widget)
        scroll_area.setVerticalScrollBarPolicy(Qt.ScrollBarPolicy.ScrollBarAlwaysOff)
        scroll_area.setHorizontalScrollBarPolicy(Qt.ScrollBarPolicy.ScrollBarAlwaysOff)
        

        main_layout = QVBoxLayout()
        main_layout.addWidget(top_bar_widget)
        main_layout.addWidget(scroll_area)

        main_widget = QWidget()
        main_widget.setLayout(main_layout)

        main_widget.setStyleSheet("""
            background-color: rgb(0, 0, 0);
        """)

        self.setMinimumSize(QSize(400,400))

        self.setCentralWidget(main_widget)
        self.setWindowTitle("Scrollable Layout Example")

window = MainWindow()
window.setAttribute(Qt.WidgetAttribute.WA_TranslucentBackground)
window.setWindowFlags(Qt.WindowType.FramelessWindowHint)
window.show()  

app.exec()