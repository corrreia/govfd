package escpos

// Character code table page constants for SetCharacterCodeTable().
// These represent the most commonly used character encoding pages for ESC/POS displays.
const (
	// Western/Latin character sets
	CharsetPC437     = 0  // PC437: USA, Standard Europe (default)
	CharsetKatakana  = 1  // Katakana
	CharsetPC850     = 2  // PC850: Multilingual
	CharsetPC860     = 3  // PC860: Portuguese
	CharsetPC863     = 4  // PC863: Canadian-French
	CharsetPC865     = 5  // PC865: Nordic
	CharsetPC851     = 11 // PC851: Greek
	CharsetPC853     = 12 // PC853: Turkish
	CharsetPC857     = 13 // PC857: Turkish
	CharsetPC737     = 14 // PC737: Greek
	CharsetISO8859_7 = 15 // ISO8859-7: Greek
	CharsetWPC1252   = 16 // WPC1252
	CharsetPC866     = 17 // PC866: Cyrillic 2
	CharsetPC852     = 18 // PC852: Latin 2
	CharsetPC858     = 19 // PC858: Euro

	// Asian/Unicode character sets
	CharsetTCVN3_1 = 30 // TCVN-3: Vietnamese
	CharsetTCVN3_2 = 31 // TCVN-3: Vietnamese

	// Middle Eastern/Arabic character sets
	CharsetPC720   = 32 // PC720: Arabic
	CharsetPC862   = 36 // PC862: Hebrew
	CharsetPC864   = 37 // PC864: Arabic
	CharsetPC1098  = 41 // PC1098: Farsi
	CharsetWPC1256 = 50 // WPC1256: Arabic

	// Eastern European character sets
	CharsetWPC775     = 33 // WPC775: Baltic Rim
	CharsetPC855      = 34 // PC855: Cyrillic
	CharsetPC861      = 35 // PC861: Icelandic
	CharsetPC869      = 38 // PC869: Greek
	CharsetISO8859_2  = 39 // ISO8859-2: Latin2
	CharsetISO8859_15 = 40 // ISO8859-15: Latin9
	CharsetPC1118     = 42 // PC1118: Lithuanian
	CharsetPC1119     = 43 // PC1119: Lithuanian
	CharsetPC1125     = 44 // PC1125: Ukrainian
	CharsetWPC1250    = 45 // WPC1250: Latin 2
	CharsetWPC1251    = 46 // WPC1251: Cyrillic
	CharsetWPC1253    = 47 // WPC1253: Greek
	CharsetWPC1254    = 48 // WPC1254: Turkish
	CharsetWPC1255    = 49 // WPC1255: Hebrew
	CharsetWPC1257    = 51 // WPC1257: Baltic Rim
	CharsetWPC1258    = 52 // WPC1258: Vietnamese
	CharsetKZ1048     = 53 // KZ1048: Kazakhstan

	// Special character sets
	CharsetSpace254 = 254 // Space characters
	CharsetSpace255 = 255 // Space characters
)
