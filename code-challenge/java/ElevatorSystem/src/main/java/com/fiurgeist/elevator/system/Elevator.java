package com.fiurgeist.elevator.system;

import java.util.LinkedList;
import java.util.stream.IntStream;

import com.fiurgeist.elevator.data.Person;
import com.fiurgeist.elevator.data.Status;

import gnu.trove.map.TIntIntMap;
import gnu.trove.map.hash.TIntIntHashMap;

/**
 * 
 * TODO move most of the logic into {@link ElevatorControlSystem}
 *
 */
public class Elevator {
	private final int id;
	private int currentFloor;
	private int nextStop;
	private LinkedList<Integer> destinations;
	private LinkedList<Integer> pickupRequest;
	private TIntIntMap personsPerDestination;
	
	public Elevator(int id) {
		this.id = id;
		this.currentFloor = 0;
		this.nextStop = 0;
		destinations = new LinkedList<>();
		pickupRequest = new LinkedList<>();
		personsPerDestination = new TIntIntHashMap();
	}

	public int getId() {
		return id;
	}

	public int getCurrentFloor() {
		return currentFloor;
	}

	public int getDestinationFloor() {
		return nextStop;
	}

	public void move() {
		if(currentFloor < nextStop) {
			currentFloor++;
		} else if(currentFloor > nextStop) {
			currentFloor--;
		}
		
		if(currentFloor == nextStop) {
			if(!pickupRequest.isEmpty() && pickupRequest.peek().intValue() == currentFloor) {
				pickupRequest.removeFirst();
				if(!pickupRequest.isEmpty()) {
					nextStop = pickupRequest.peek();
				}
			}
			
			if(!destinations.isEmpty() && destinations.peek().intValue() == currentFloor) {
				destinations.removeFirst();
				personsPerDestination.remove(currentFloor);
			}
			if(!destinations.isEmpty() && (currentFloor == nextStop || 
					Math.abs(currentFloor - destinations.peek()) < Math.abs(currentFloor - nextStop))) {
				nextStop = destinations.peek();
			}
		}
	}
	
	public Status getStatus() {
		return new Status(currentFloor, nextStop);
	}

	public void addPerson(Person person) {
		personsPerDestination.adjustOrPutValue(person.destination, 1, 1);
		
		if(destinations.isEmpty()) {
			destinations.add(person.destination);
		} else {// TODO consider current direction of elevator
			boolean wasAdded = false;
			int elevatorPosition = currentFloor;
			for(int i = 0; i < destinations.size(); ++i) {
				int floor = destinations.get(i);
				if(floor == person.destination) {
					wasAdded = true;
					break;
				}
				if(isBetween(person.destination, elevatorPosition, floor)) {
					destinations.add(i, person.destination);
					wasAdded = true;
					break;
				}
				elevatorPosition = floor;
			}
			if(!wasAdded) {
				destinations.add(person.destination);
			}
		}
		if(currentFloor == nextStop || 
				isBetween(person.destination, currentFloor, nextStop)) {
			nextStop = person.destination;
		}
	}

	public int getNumberOfPersons() {
		return IntStream.of(personsPerDestination.values()).sum();
	}

	public void addPickupRequest(int newFloor) {
		if(pickupRequest.isEmpty()) {
			pickupRequest.add(newFloor);
		} else {// TODO consider current direction of elevator
			boolean wasAdded = false;
			int elevatorPosition = currentFloor;
			for(int i = 0; i < pickupRequest.size(); ++i) {
				int floor = pickupRequest.get(i);
				if(floor == newFloor) {
					wasAdded = true;
					break;
				}
				if(isBetween(newFloor, elevatorPosition, floor)) {
					pickupRequest.add(i, newFloor);
					wasAdded = true;
					break;
				}
				elevatorPosition = floor;
			}
			if(!wasAdded) {
				pickupRequest.add(newFloor);
			}
		}
		if(currentFloor == nextStop || 
				isBetween(newFloor,currentFloor, nextStop)) {
			nextStop = newFloor;
		}
	}
	
	public static boolean isBetween(int value, int start, int end) {
		if(start < end
				&& start < value
				&& value < end) {
			return true;
		} else if(start > end
				&& start > value
				&& value > end) {
			return true;
		}
		return false;
	}
}
