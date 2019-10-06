arr := list(1, 2, 3, 4)
arr2 := list(1, "2", 3, 4)

List avg := method(
(self sum) / (self size)
)

List avg_ex := method(
sum := 0
if(self size == 0, 
	return 0,
	self foreach(item, if(item proto == Number, sum = sum + item, Exception raise("list item is NaN"))))
sum / self size
)

arr avg println
list() avg_ex println
arr avg_ex println
arr2 avg_ex println