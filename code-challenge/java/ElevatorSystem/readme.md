Task: Implement an elevator system to transport persons to there destination in the shortest possible time.
Time: 4h


### Prerequisites
- JDK8 installed
- Maven installed

### Build and run
build with: 
`mvn package`

run with default values: 
`java -jar target/elevator-0.0.1-SNAPSHOT.jar`

you can set the following arguments
- the number of elevators; default 16
- the number of person to simulate; default 1000
- the number of floors; default 13
- the number of persons requesting an elevator per simulation step; default 6

`java -jar target/elevator-0.0.1-SNAPSHOT.jar 16 1000 13 6`

### Scheduler
###### As implemented
The scheduler as implemented works in a way that pickup requests are given to the nearest elevator at that moment.
And when a person enters an elevator with a given goal in mind, the destination is not just added to the end of the queue. 
Instead it's checked if the destination of the new passenger is anywhere between the location of the elevator and any of the other destinations the elevator already has in its queue. Duplicates are also filtered out.
It's not the best solution, but it makes it definitely better than just FIFO. Due to the time constraint given I was not able to fully implement the scheduler I had in mind.

###### As envisioned
For the scheduling of the pickup requests I planned to consider the direction of the elevators to there next target.
As a further step, I thought about using all other destinations an elevator has to really find the elevator with the fewest amount of simulation steps to pickup the given person.
Also due to the time constraint I wasn't able to clean up the code, e.g. there is a lot of logic in the Elevator class which I would move to the ElevatorControlSystem class.
