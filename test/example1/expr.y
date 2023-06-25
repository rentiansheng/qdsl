
%{

package main

import (
  )


%}

%union {
  obj []map[string]interface{}
  key string
  val interface{}
  list []interface{}
  op string
  errMsg string
}


%token <val> String Number Literal
%token <op> OpType JointType OpTypeBetween
%type <val> array between
%type <list> elements
%type <obj> object input expr
%type <val> value
%type <op> OpType  JointType



%%

input: object{
        yylex.(*lex).result = $$

}

object: expr| object JointType expr {
   	obj := $$
	obj = append(obj,
		map[string]interface{}{"condition": $3},
		 map[string]interface{}{"op": $2})
	$$ = obj
	yylex.(*lex).state = lexStateJoint

}


expr: String OpType value {
      	$$ =  []map[string]interface{}{ {$1.(string): map[string]interface{}{ $2:   $3}}}

      	yylex.(*lex).state = lexStateJoint
      } | String '=' value {
      	$$ =  []map[string]interface{}{{$1.(string): map[string]interface{}{ "=":   $3}}}
      	yylex.(*lex).state = lexStateJoint

      } |  String OpTypeBetween between {
      	$$ =  []map[string]interface{}{{$1.(string): map[string]interface{}{ "between":   $3}}}
      	yylex.(*lex).state = lexStateJoint
      }

between : '(' elements ')' {
	yylex.(*lex).state = LexStateValueElement
	$$ = $2
  	yylex.(*lex).state = lexStateJoint
}

array: '[' elements ']' {
	yylex.(*lex).state = LexStateValueElement
	$$ = $2
	yylex.(*lex).state = lexStateJoint
}

elements: {
	$$ = []interface{}{}
}| value {
	$$ = []interface{}{$1}
}| elements ',' value {
	$$ = append($1, $3)
}

value: String
| Number
| Literal
| object {
        yylex.(*lex).pushState()
	$$ = $1
	yylex.(*lex).popState()
}
| array


%%
