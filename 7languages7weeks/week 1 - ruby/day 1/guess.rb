puts 'guess number between 1 and 10' 
number = rand(10)+1
while true
	input = gets.to_i
	if input == number
		puts 'correct'
		break
	elsif input < number
		puts 'too low'
	else
		puts 'too high'
	end
end