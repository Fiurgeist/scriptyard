arr := list(1,list(2,3),list(4,5,6))
arr_sum := method(array,
sum :=0;
array foreach(item, sum = sum + if(item proto == List, arr_sum(item), item));
sum)
arr_sum(arr) println

arr_sum_sum := method(array,
sum :=0;
array foreach(item, sum = sum + if(item proto == List, item sum, item));
sum)
arr_sum_sum(arr) println

arr_sum_flatten_sum := method(array,
array flatten sum)
arr_sum_flatten_sum(arr) println