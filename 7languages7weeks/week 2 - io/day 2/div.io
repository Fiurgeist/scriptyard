Number originalDivision := Number getSlot("/")
Number / := method(denominator,
if(denominator==0, 0, self originalDivision(denominator)))
(6 / 2) println