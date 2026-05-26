package ast

import (
	"project-atra/token"
	"strings"
)

// Node は AST の全ノードが満たすべきインターフェースです。
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement は「記述」を表現するインターフェースです。
type Statement interface {
	Node
	statementNode()
}

// Expression は「意味の断片」を表現するインターフェースです。
type Expression interface {
	Node
	expressionNode()
}

// Universe は Atra ソースファイルのルートノードです。
type Universe struct {
	Statements []Statement
}

func (u *Universe) TokenLiteral() string {
	if len(u.Statements) > 0 {
		return u.Statements[0].TokenLiteral()
	}
	return ""
}

func (u *Universe) String() string {
	var out string
	for _, s := range u.Statements {
		out += s.String() + "\n"
	}
	return out
}

// ExpressionStatement は式単体で構成される文を表現します。
type ExpressionStatement struct {
	Token      token.Token // 式の最初のトークン
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// CelestialNode は「点（天体）」を表現します。
// 識別子、質量（星の数）、特異点フラグ、および内部の「場（Block）」を持ちます。
type CelestialNode struct {
	Token         token.Token // token.IDENT
	Value         string
	Mass          int             // STAR (*) の数
	IsSingularity bool            // SINGULARITY (!) の有無
	Body          *BlockStatement // 天体の内部にある宇宙（任意）
}

func (cn *CelestialNode) expressionNode()      {}
func (cn *CelestialNode) TokenLiteral() string { return cn.Token.Literal }
func (cn *CelestialNode) String() string {
	out := ""
	for i := 0; i < cn.Mass; i++ {
		out += "*"
	}
	if cn.IsSingularity {
		out += "!"
	}
	out += cn.Value
	if cn.Body != nil {
		out += ":\n" + cn.Body.String()
	}
	return out
}

// BlockStatement はコロン（:）に続く、インデントされた一連の記述を表現します。
type BlockStatement struct {
	Token      token.Token // token.COLON
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out string
	for _, s := range bs.Statements {
		out += "\t" + s.String() + "\n"
	}
	return out
}

// FlowStatement は二点間の「流れ（遷移）」を表現します。
type FlowStatement struct {
	Token    token.Token // FLOW (->), WAVE_FLOW (~>), WORMHOLE (=>)
	Source   Expression  // 出発点
	Target   Expression  // 到着点
	Operator string      // 演算子リテラル
}

func (fs *FlowStatement) statementNode()       {}
func (fs *FlowStatement) expressionNode()      {}
func (fs *FlowStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FlowStatement) String() string {
	return fs.Source.String() + " " + fs.Operator + " " + fs.Target.String()
}

// ForceStatement は二点間の「引力・斥力」を表現します。
type ForceStatement struct {
	Token    token.Token // ATTRACT (+++), REPEL (-!- )
	Left     Expression
	Right    Expression
	Operator string
}

func (fs *ForceStatement) statementNode()       {}
func (fs *ForceStatement) expressionNode()      {}
func (fs *ForceStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForceStatement) String() string {
	return fs.Left.String() + " " + fs.Operator + " " + fs.Right.String()
}

// RhizomeExpression は「@」による参照を表現します。
type RhizomeExpression struct {
	Token token.Token // token.REF (@)
	Target string
}

func (re *RhizomeExpression) expressionNode()      {}
func (re *RhizomeExpression) TokenLiteral() string { return re.Token.Literal }
func (re *RhizomeExpression) String() string       { return "@" + re.Target }

// NumberLiteral は数値を表現します。
type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (nl *NumberLiteral) expressionNode()      {}
func (nl *NumberLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NumberLiteral) String() string       { return nl.Token.Literal }

// CallExpression は関数呼び出しのような記述（例: Phase(Identity)）を表現します。
type CallExpression struct {
	Token     token.Token // '(' トークン
	Function  Expression  // IDENT
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	return ce.Function.String() + "(" + strings.Join(args, ", ") + ")"
}

// KeyValuePair はコロンによる対（例: key: value）を表現します。
type KeyValuePair struct {
	Token token.Token // ':' トークン
	Key   Expression
	Value Expression
}

func (kv *KeyValuePair) expressionNode()      {}
func (kv *KeyValuePair) TokenLiteral() string { return kv.Token.Literal }
func (kv *KeyValuePair) String() string {
	return kv.Key.String() + ": " + kv.Value.String()
}
