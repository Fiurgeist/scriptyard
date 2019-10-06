-module(shopping).
-export([total_price/1]).

total_price(Books) -> [{Item,Price*Amount}||{Item,Price,Amount}<-Books,Price>30].