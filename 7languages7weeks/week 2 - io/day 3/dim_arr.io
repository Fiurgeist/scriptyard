TwoDimList := Object clone
TwoDimList list := List clone
TwoDimList dim := method(x, y,
list = List clone
y repeat(
tmp := List clone
x repeat(tmp append(nil))
list append(tmp)
)
)
TwoDimList set := method(x, y, value,
list at(y) atPut(x, value)
)
TwoDimList get := method(x, y,
list at(y) at(x)
)
TwoDimList transpose := method(
tmp := TwoDimList clone
tmp dim(list size, list at(0) size)
list foreach(i, v,
v foreach(j, w, tmp set(i, j, w))
)
list = tmp
)

squareBrackets := method(
r := List clone
call message arguments foreach(arg,
r append(doMessage(arg))
)
r
)

s := File with("matrix.txt") openForReading contents
matrix := doString(s)
matrix println