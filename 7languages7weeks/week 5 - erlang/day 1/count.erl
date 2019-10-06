-module(count).
-export([to/1]).
-export([to_ten/1]).

to(1) -> "1";
to(To) -> to(To-1) ++ " " ++ integer_to_list(To).

% part 2, count from N to 10

to_ten(10) -> "10";
to_ten(From) -> if 
				From < 10 -> integer_to_list(From) ++ " " ++ to_ten(From+1);
				From > 10 -> integer_to_list(From) ++ " " ++ to_ten(From-1) 
				end.