package com.fiurgeist.elevator.data;

public class Status {
	public final int currentFloor;
	public final int nextStop;
	
	public Status(int currentFloor, int nextStop) {
		this.currentFloor = currentFloor;
		this.nextStop = nextStop;
	}
}
