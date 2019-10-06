from decimal import Decimal
import math
from typing import Type, TypeVar, Iterable

NUMBER = TypeVar("NUMBER", int, float, Decimal)

class Vector(object):
	def __init__(self, coordinates: Iterable[NUMBER]):
		if not coordinates:
			raise ValueError("The coordinates must be set.")

		self.coordinates = tuple([Decimal(coordinate) for coordinate in coordinates])
		self.dimension = len(coordinates)
		self._magnitude = None

	def __str__(self) -> str:
		return "Vector: {}".format(self.coordinates)

	def __eq__(self, other: "Vector") -> bool:
		return self.coordinates == other.coordinates

	def __add__(self, other: "Vector") -> "Vector":
		if self.dimension != other.dimension:
			raise ValueError("The other vector must have the same dimension.")

		return Vector(tuple(coordinate + other.coordinates[idx] for idx, coordinate in enumerate(self.coordinates)))

	def __sub__(self, other: "Vector") -> "Vector":
		if self.dimension != other.dimension:
			raise ValueError("The other vector must have the same dimension.")
		
		return Vector(tuple(coordinate - other.coordinates[idx] for idx, coordinate in enumerate(self.coordinates)))

	def __mul__(self, other: "Vector") -> "Vector":
		"""Cross product"""
		if self.dimension != other.dimension:
			raise ValueError("The other vector must have the same dimension.")

		if self.dimension == 3:
			x1, y1, z1 = self.coordinates
			x2, y2, z2 = other.coordinates
		elif self.dimension == 2:
			x1, y1 = self.coordinates
			x2, y2 = other.coordinates
			z1 = z2 = 0
		else:
			raise Exception("Cross product is only defined in two or tree dimensions.")

		return Vector((
			y1*z2 - y2*z1,
			-(x1*z2 - x2*z1),
			x1*y2 - x2*y1,
		))

	def scalar(self, scalar: NUMBER) -> "Vector":
		scalar = Decimal(scalar)
		return Vector(tuple(coordinate * scalar for coordinate in self.coordinates))

	def magnitude(self) -> Decimal:
		if not self._magnitude:
			self._magnitude = Decimal(sum([x*x for x in self.coordinates])).sqrt()
		
		return self._magnitude

	def normalize(self) -> "Vector":
		magnitude = self.magnitude()
		if magnitude == 0:
			raise Exception("The zero vector cannot be normalized.")
		return self.scalar(Decimal(1.0)/magnitude)

	def dot(self, other: "Vector") -> Decimal:
		if self.dimension != other.dimension:
			raise ValueError("The other vector must have the same dimension.")
		
		return sum([coordinate * other.coordinates[idx] for idx, coordinate in enumerate(self.coordinates)])

	def angle_rad(self, other: "Vector") -> float:
		if self.dimension != other.dimension:
			raise ValueError("The other vector must have the same dimension.")
		
		m1 = self.magnitude()
		if m1 == 0:
			raise Exception("Cannot compute angle. Self is a zero vector.")
		
		m2 = other.magnitude()
		if m2 == 0:
			raise Exception("Cannot compute angle. Other vector is a zero vector.")

		return math.acos(self.dot(other) / (m1 * m2))

	def angle_deg(self, other: "Vector") -> float:
		rad = self.angle_rad(other)
		return math.degrees(rad)

	def is_parallel_to(self, other: "Vector") -> bool:
		if self.dimension != other.dimension:
			raise ValueError("The other vector must have the same dimension.")

		if self.magnitude() == 0 or other.magnitude() == 0:
			return True

		angle = self.angle_deg(other)
		return angle == 0 or angle == 180

	def is_orthogonal_to(self, other: "Vector") -> bool:
		if self.dimension != other.dimension:
			raise ValueError("The other vector must have the same dimension.")

		# also viable
		# return math.abs(self.dot(other)) < 1e-10 # a tolerance value
		
		if self.magnitude() == 0 or other.magnitude() == 0:
			return True

		return self.angle_deg(other) == 90

	def projection_onto(self, basis: "Vector") -> "Vector":
		if self.dimension != basis.dimension:
			raise ValueError("The basis vector must have the same dimension.")

		unit_vector = basis.normalize()
		return unit_vector.scalar(self.dot(unit_vector))

	def perpendicular_to(self, basis: "Vector") -> "Vector":
		return self - self.projection_onto(basis)


def area_of_parallelogram(vector1: Vector, vector2: Vector) -> Decimal:
	return (vector1 * vector2).magnitude()