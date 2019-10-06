book(anhalter, adams).
book(x, clamp).
book(rg_veda, clamp).
book(nausicaa,miyazaki).
book(hdr,tolkin).

%% not really needed 
book_by_author(Which,Author) :- book(Which,Author). 