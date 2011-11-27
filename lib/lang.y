%{
package gohaml

import "fmt"

var Output inode
%}

%union {
  n inode
  s string
  i interface{}
  c icodenode
}

%type<n> statement
%type<c> rhs
%type<s> complex_ident
%token<s> IDENT
%token<i> ATOM FOR RANGE

%%

statement :  FOR IDENT ',' IDENT ':' '=' RANGE IDENT
            {
              rn := new(rangenode)
              rn._lhs1 = $2
              rn._lhs2 = $4
              rn._rhs = res{$8, true}
              $$ = rn
              Output = $$
            }
          | IDENT ':' '=' rhs
            {
              $4.setLHS($1)
              $$ = $4
              Output = $$
            }
          ;

rhs : ATOM
      {
        dan := new(declassnode)
        dan._rhs = $1
        $$ = dan
      }
    | IDENT complex_ident
      {
        dan := new(vdeclassnode)
        dan._rhs.value = $1 + $2
        dan._rhs.needsResolution = true
        $$ = dan      
      }
    ;

complex_ident : '.' IDENT complex_ident
                {
                  $$ = fmt.Sprintf(".%s%s", $2, $3)
                }
              |
                {
                  $$ = ""
                }
              ;

%%
