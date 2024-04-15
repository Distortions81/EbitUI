package EbitUI

const minSizeX = 1200
const minSizeY = 700

type ITYPE int

const (
	ITYPE_TEXT = iota
	ITYPE_BUTTON
	ITYPE_FLOW
	ITYPE_DIVIDER
	ITYPE_TAB
)

type FLOW_DIR int

const (
	FLOW_NONE = iota
	FLOW_HORIZONTAL
	FLOW_VERTICAL
	FLOW_BOTH
)
