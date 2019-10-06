package com.fiurgeist.elevator.system;

import java.util.Map.Entry;

import com.fiurgeist.elevator.data.Status;

import java.util.SortedMap;
import java.util.TreeMap;

import gnu.trove.iterator.TIntObjectIterator;
import gnu.trove.map.TIntIntMap;
import gnu.trove.map.TIntObjectMap;
import gnu.trove.map.hash.TIntIntHashMap;
import gnu.trove.map.hash.TIntObjectHashMap;

public class ElevatorControlSystem {
	private TIntObjectMap<Status> elevatorStatuses;
	private TIntObjectMap<FloorPickupRequest> pickupRequests;

	public enum Direction {
		UP,
		DOWN
	}
	
	public final class PickupOrder {
		public final int elevatorId;
		public final int floor;
		
		public PickupOrder(int elevatorId, int floor) {
			this.elevatorId = elevatorId;
			this.floor = floor;
		}
	}
	
	private class FloorPickupRequest {
		private int[] requests = {0, 0};
		
		public void addRequest(Direction direction) {
			requests[direction.ordinal()]++;
		}
		
		public int[] getRequests() {
			return requests;
		}
	}
	
	public ElevatorControlSystem(int elevatorsCount) {
		elevatorStatuses = new TIntObjectHashMap<>(elevatorsCount);
		pickupRequests = new TIntObjectHashMap<>(16);
	}
	
	public TIntObjectMap<Status> getElevatorStatuses() {
		return elevatorStatuses;
	}
	
	public void update(int elevatorId, Status status) {
		elevatorStatuses.put(elevatorId, status);
	}

	public void pickup(int floor, Direction direction) {
		FloorPickupRequest floorReq = pickupRequests.get(floor);
		if(floorReq == null) {
			floorReq = new FloorPickupRequest();
			pickupRequests.put(floor, floorReq);
		}
		floorReq.addRequest(direction);
	}
	
	public TIntIntMap step() {
		TIntIntMap orders = new TIntIntHashMap(pickupRequests.size());
		pickupRequests.forEachEntry((floor, requests) -> {
			if(requests.getRequests()[0]==0 && requests.getRequests()[1]==0) {
				return true;
			}
			// create sorted map of distances between the elevators current position and the pickup request position
			SortedMap<Integer, TIntObjectMap<Status>> distanceMap = new TreeMap<>();
			elevatorStatuses.forEachEntry((id, status) -> {
				int distance = Math.abs(status.currentFloor - floor);
				TIntObjectMap<Status> elevators = distanceMap.get(distance);
				if(elevators == null) {
					elevators = new TIntObjectHashMap<>();
					distanceMap.put(distance, elevators);
				}
				elevators.put(id, status);
				return true;
			});
			
			// find the elevator with the shortest distance where the pickup request is between the elevators current position and next stop
			int elevatorId = -1;
			for(Entry<Integer, TIntObjectMap<Status>> entry : distanceMap.entrySet()) {
				TIntObjectIterator<Status> it = entry.getValue().iterator();
				while(it.hasNext()) {
					it.advance();
					if(Elevator.isBetween(floor, it.value().currentFloor, it.value().nextStop)) {
						elevatorId = it.key();
						break;
					}
				}
				if(elevatorId != -1) {
					break;
				}
			}
			if(elevatorId == -1) {
				elevatorId = distanceMap.get(distanceMap.firstKey()).keys()[0];
			}
			// TODO consider the pickup request direction when choosing the nearest elevator
			// TODO not just checking if it's between the next stop, but checking all the destinations an elevator has 
			// to calculate the number of simulation steps needed to reach the floor of the pickup request 
			orders.put(elevatorId, floor);
			return true;
		});
		pickupRequests.clear();
		return orders;
	}
}
