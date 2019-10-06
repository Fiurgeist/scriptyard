-module(count_words).
-export([words/1]).

words([]) -> 0;
words(Sentence) -> count(Sentence) + 1.

count([]) -> 0;
count([32|Tail]) -> count(Tail) + 1;
count([_|Tail]) -> count(Tail).