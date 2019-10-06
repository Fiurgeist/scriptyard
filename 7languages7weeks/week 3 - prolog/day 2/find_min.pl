find_min([Head], Min) :- Min is Head.
find_min([First, Second | Tail]) :-
First > Second,Min is Second,
find_min([Second | Tail],Min).
%find_min([First, Second | Tail]) :-
%First <= Second,Min is First,
%find_min([First | Tail],Min).