-module(print_error).
-export([error/1]).

error(success) -> "success";
error({error, Message}) -> "error: " ++ Message.