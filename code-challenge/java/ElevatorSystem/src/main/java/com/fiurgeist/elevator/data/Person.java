package com.fiurgeist.elevator.data;

import java.util.Random;

public class Person {
	private static final Random RAND = new Random(42);
	
	public final int start;
	public final int destination;
	public final int startStep;
	
	public Person(int start, int destination, int startStep) {
		this.start = start;
		this.destination = destination;
		this.startStep = startStep;
	}
	
	public static Person create(int startStep, int maxFloor) {
		int start = RAND.nextInt(maxFloor);
		int destination;
		do {
			destination = RAND.nextInt(maxFloor);
		} while (start == destination);
		
		return new Person(start, destination, startStep);
	}
}
