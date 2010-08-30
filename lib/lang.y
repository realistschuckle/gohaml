%package gohaml

%union {
	n inode
	s string
}

%type<n> Statement
%type<s> ident str

%%

Statement : tfor ident ',' ident ':' '=' trange ident	{
															rn := new(rangenode)
															rn._first = $2
															rn._second = $4
															rn._third = $8
															$$ = rn
														}
		  | ident ':' '=' str							{
															dan := new(declassnode)
															dan._lhs = $1
															dan._rhs = res{$4[1:len($4) - 1], false}
															$$ = dan
														}
          ;
