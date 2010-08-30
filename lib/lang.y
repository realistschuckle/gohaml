%package gohaml

%union {
	b bool
	n *node
}

%token<b> tfor trange
%token<n> ident
%type<n> Statement

%%

Statement : tfor ident ',' ident ':' '=' trange ident	{
															$$ = new(node)
															
														}
          ;
