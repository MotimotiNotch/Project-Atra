package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 識別子とリテラル
	IDENT  = "IDENT"  // Example: Viscosity, Silent_Melody
	NUMBER = "NUMBER" // Example: 1779719026, +0.05, -10

	// 質量と特異点 (Topology & Mass)
	STAR        = "STAR"        // * (質量の単位)
	SINGULARITY = "SINGULARITY" // ! (特異点)

	// 重力とフロー (Gravity & Flow)
	ATTRACT   = "ATTRACT"   // +++ (引力)
	REPEL     = "REPEL"     // -!- (斥力)
	FLOW      = "FLOW"      // -> (測地線フロー)
	WAVE_FLOW = "WAVE_FLOW" // ~> (揺らぎのフロー)
	WORMHOLE  = "WORMHOLE"  // => (ワームホール)

	// 参照と構造 (Reference)
	REF   = "REF"   // @ (リゾーム的な参照)
	COLON = "COLON" // : (ブロックの開始等)
	DOT   = "DOT"   // . (プロパティアクセス等)
	COMMA = "COMMA" // ,

	// 括弧
	LPAREN = "LPAREN" // (
	RPAREN = "RPAREN" // )
	LBRACE = "LBRACE" // {
	RBRACE = "RBRACE" // }

	// レイアウト (Layout / Context)
	NEWLINE = "NEWLINE" // 改行
	INDENT  = "INDENT"  // インデント（深さを Literal に保持）
)
