-module(map).
-export([get/2]).
-export([get2/2]).

get([],Key) -> null;
get([Head|Tail], Key) -> 
if element(1,Head) == Key -> element(2,Head);
true -> get(Tail, Key)
end.

get2(List,Key)-> 
Tupel = lists:keytake(Key, 1, List),
if Tupel ==false -> null;
true -> element(2,Tupel)
end.