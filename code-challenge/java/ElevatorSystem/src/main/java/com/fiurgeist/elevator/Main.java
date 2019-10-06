package com.fiurgeist.elevator;

import java.util.ArrayList;
import java.util.Deque;
import java.util.Iterator;
import java.util.LinkedList;
import java.util.List;
import java.util.StringJoiner;
import java.util.stream.IntStream;

import com.fiurgeist.elevator.data.Person;
import com.fiurgeist.elevator.system.Elevator;
import com.fiurgeist.elevator.system.ElevatorControlSystem;
import com.fiurgeist.elevator.system.ElevatorControlSystem.Direction;

import gnu.trove.map.TIntIntMap;


/**
 * 
 * TODO convert to unit test
 *
 */
public class Main {
	public static void main(String[] args) {
		int elevatorsCount = args.length > 0 ? Integer.parseInt(args[0]) : 16;
		int maxPersons = args.length > 1 ? Integer.parseInt(args[1]) : 1000;
		int maxFloor = args.length > 2 ? Integer.parseInt(args[2]) : 13;
		int personsPerStep = args.length > 3 ? Integer.parseInt(args[3]) : 6;
		
		List<Elevator> elevators = new ArrayList<>(elevatorsCount);
		for(int i = 0; i < elevatorsCount; ++i) {
			elevators.add(new Elevator(i));
		}

		Deque<Person> persons = new LinkedList<>();
		for(int i = 0; i < maxPersons; ++i) {
			persons.add(Person.create(i / personsPerStep, maxFloor));
		}
		
		List<List<Person>> waitingPersons = new ArrayList<>(maxFloor);
		for(int i = 0; i < maxFloor; ++i) {
			waitingPersons.add(new LinkedList<>());
		}
		
		ElevatorControlSystem system = new ElevatorControlSystem(elevatorsCount);
		
		boolean isDone = false;
		int step;
		for(step = 0; !isDone; ++step) {
			//update system with elevator state
			elevators.forEach(elevator -> system.update(elevator.getId(), elevator.getStatus()));

			//add persons/requests
			for(int i = 0; i < personsPerStep && !persons.isEmpty(); ++i) {
				Person person = persons.pop();
				system.pickup(person.start, person.start < person.destination ? Direction.UP : Direction.DOWN);
				waitingPersons.get(person.start).add(person);
			}
			
			//return elevator pickup requests
			TIntIntMap orders = system.step();
			
			int[] personsInElevator = new int[elevators.size()];
			StringJoiner elevatorLog = new StringJoiner("\t");
			elevators.forEach(elevator -> {
				elevator.addPickupRequest(orders.get(elevator.getId()));
				elevator.move();
				
				List<Person> personsOnFloor = waitingPersons.get(elevator.getCurrentFloor());
				if(!personsOnFloor.isEmpty()) {
					int elevatorDirection = elevator.getCurrentFloor() - elevator.getDestinationFloor();
					Iterator<Person> iter = personsOnFloor.listIterator();
					while(iter.hasNext()){
						Person person = iter.next();
						int personDirection = person.start - person.destination;
					    if(elevatorDirection == 0 
					    		|| elevatorDirection * personDirection > 0) {
					    	elevator.addPerson(person);
					        iter.remove();
					    }
					}
				}
				
				personsInElevator[elevator.getId()] = elevator.getNumberOfPersons();
				elevatorLog.add(String.format("%1$d(%2$d)", elevator.getCurrentFloor(), personsInElevator[elevator.getId()]));
			});

			System.out.println("Elevators [floor(passengers)]: " + elevatorLog);
			isDone = persons.isEmpty() && waitingPersons.stream().mapToInt(list -> list.size()).sum() == 0 && IntStream.of(personsInElevator).sum() == 0;
		}
		System.out.println("All Persons reached there destination after " + step + " simmulation steps");
	}
}
