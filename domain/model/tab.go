package model

type Tab struct {
	ID           int
	Windows      []*Window
	ActiveWindow *Window
}

func NewTab(id int) *Tab {
	return &Tab{
		ID:      id,
		Windows: make([]*Window, 0),
	}
}

func (t *Tab) AddWindow(w *Window) {
	t.Windows = append(t.Windows, w)
	if t.ActiveWindow == nil {
		t.ActiveWindow = w
	}
}

func (t *Tab) RemoveWindow(windowID int) {
	for i, w := range t.Windows {
		if w.ID == windowID {
			t.Windows = append(t.Windows[:i], t.Windows[i+1:]...)
			if t.ActiveWindow.ID == windowID && len(t.Windows) > 0 {
				t.ActiveWindow = t.Windows[0]
			}
			break
		}
	}
}

func (t *Tab) SetActiveWindow(windowID int) {
	for _, w := range t.Windows {
		if w.ID == windowID {
			t.ActiveWindow = w
			break
		}
	}
}
