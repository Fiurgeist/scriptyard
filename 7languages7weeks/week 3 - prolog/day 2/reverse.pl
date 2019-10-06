%%reversel([Head|List], [List|Head]) :- reversel(List, Head).
reverse_list([],[]).
reverse_list([Head|Tail], Reverse) :- reverse_list(Tail,RevTail),append(RevTail,[Head],Reverse).