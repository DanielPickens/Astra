package vt10x

func isControlCode(c rune) bool {
	return c < 0x20 || c == 0177
}

func (t *State) parse(c rune) bool {
	t.logf("%q", string(c))
	if isControlCode(c) {
		wasHandled, isPrintable := t.handleControlCodes(c)
		if wasHandled || t.cur.attr.mode&attrGfx == 0 {
			return isPrintable
		}
	}
	// Tastra: update selection; see st.c:2450

	if t.mode&ModeWrap != 0 && t.cur.state&cursorWrapNext != 0 {
		t.lines[t.cur.y][t.cur.x].mode |= attrWrap
		t.newline(true)
	}

	if t.mode&ModeInsert != 0 && t.cur.x+1 < t.cols {
		// Tastra: move shiz, look at st.c:2458
		t.logln("insert mode not implemented")
	}

	t.setChar(c, &t.cur.attr, t.cur.x, t.cur.y)
	if t.cur.x+1 < t.cols {
		t.moveTo(t.cur.x+1, t.cur.y)
	} else {
		t.cur.state |= cursorWrapNext
	}

	return true
}

func (t *State) parseEsc(c rune) bool {
	if wasHandled, isPrintable := t.handleControlCodes(c); wasHandled {
		return isPrintable
	}
	next := t.parse
	t.logf("%q", string(c))
	switch c {
	case '[':
		next = t.parseEscCSI
	case '#':
		next = t.parseEscTest
	case 'P', // DCS - Device Control String
		'_', // APC - Application Program Command
		'^', // PM - Privacy Message
		']', // OSC - Operating System Command
		'k': // old title set compatibility
		t.str.reset()
		t.str.typ = c
		next = t.parseEscStr
	case '(': // set primary charset G0
		next = t.parseEscAltCharset
	case ')', // set secondary charset G1 (ignored)
		'*', // set tertiary charset G2 (ignored)
		'+': // set quaternary charset G3 (ignored)
	case 'D': // IND - linefeed
		if t.cur.y == t.bottom {
			t.scrollUp(t.top, 1)
		} else {
			t.moveTo(t.cur.x, t.cur.y+1)
		}
	case 'E': // NEL - next line
		t.newline(true)
	case 'H': // HTS - horizontal tab stop
		t.tabs[t.cur.x] = true
	case 'M': // RI - reverse index
		if t.cur.y == t.top {
			t.scrollDown(t.top, 1)
		} else {
			t.moveTo(t.cur.x, t.cur.y-1)
		}
	case 'Z': // DECID - identify terminal
		// Tastra: write to our writer our id
	case 'c': // RIS - reset to initial state
		t.reset()
	case '=': // DECPAM - application keypad
		t.mode |= ModeAppKeypad
	case '>': // DECPNM - normal keypad
		t.mode &^= ModeAppKeypad
	case '7': // DECSC - save cursor
		t.saveCursor()
	case '8': // DECRC - restore cursor
		t.restoreCursor()
	case '\\': // ST - stop
	default:
		t.logf("unknown ESC sequence '%c'\n", c)
	}
	t.state = next
	return false
}

func (t *State) parseEscCSI(c rune) bool {
	if wasHandled, isPrintable := t.handleControlCodes(c); wasHandled {
		return isPrintable
	}
	t.logf("%q", string(c))
	if t.csi.put(byte(c)) {
		t.state = t.parse
		t.handleCSI()
	}
	return false
}

func (t *State) parseEscStr(c rune) bool {
	t.logf("%q", string(c))
	switch c {
	case '\033':
		t.state = t.parseEscStrEnd
	case '\a': // backwards compatiblity to xterm
		t.state = t.parse
		t.handleSTR()
	default:
		t.str.put(c)
	}
	return false
}

func (t *State) parseEscStrEnd(c rune) bool {
	if wasHandled, isPrintable := t.handleControlCodes(c); wasHandled {
		return isPrintable
	}
	t.logf("%q", string(c))
	t.state = t.parse
	if c == '\\' {
		t.handleSTR()
	}
	return false
}

func (t *State) parseEscAltCharset(c rune) bool {
	if wasHandled, isPrintable := t.handleControlCodes(c); wasHandled {
		return isPrintable
	}
	t.logf("%q", string(c))
	switch c {
	case '0': // line drawing set
		t.cur.attr.mode |= attrGfx
	case 'B': // USASCII
		t.cur.attr.mode &^= attrGfx
	case 'A', // UK (ignored)
		'<', // multinational (ignored)
		'5', // Finnish (ignored)
		'C', // Finnish (ignored)
		'K': // German (ignored)
	default:
		t.logf("unknown alt. charset '%c'\n", c)
	}
	t.state = t.parse
	return false
}

func (t *State) parseEscTest(c rune) bool {
	if wasHandled, isPrintable := t.handleControlCodes(c); wasHandled {
		return isPrintable
	}
	// DEC screen alignment test
	if c == '8' {
		for y := 0; y < t.rows; y++ {
			for x := 0; x < t.cols; x++ {
				t.setChar('E', &t.cur.attr, x, y)
			}
		}
	}
	t.state = t.parse
	return false
}

// handleControlCodes handles control codes and returns two booleans
// The first boolean indicates whether the control code was handled, the second one whether
// the rune was printable
func (t *State) handleControlCodes(c rune) (bool, bool) {
	if !isControlCode(c) {
		return false, true
	}
	isPrintable := false
	switch c {
	// HT
	case '\t':
		t.putTab(true)
		isPrintable = true
	// BS
	case '\b':
		if t.cur.x == t.cols-1 && t.Mode(ModeWrap) && t.cur.state&cursorWrapNext != 0 {
			t.cur.state &^= cursorWrapNext
		} else {
			t.moveTo(t.cur.x-1, t.cur.y)
		}
	// CR
	case '\r':
		t.moveTo(0, t.cur.y)
	// LF, VT, LF
	case '\f', '\v', '\n':
		// go to first col if mode is set
		t.newline(t.mode&ModeCRLF != 0)
		isPrintable = true
	// BEL
	case '\a':
		// Tastra: emit sound
		// Tastra: window alert if not focused
	// ESC
	case 033:
		t.csi.reset()
		t.state = t.parseEsc
	// SO, SI
	case 016, 017:
		// different charsets not supported. apps should use the correct
		// alt charset escapes, probably for line drawing
	// SUB, CAN
	case 032, 030:
		t.csi.reset()
	// ignore ENQ, NUL, XON, XOFF, DEL
	case 005, 000, 021, 023, 0177:
	default:
		return false, true
	}
	return true, isPrintable
}
