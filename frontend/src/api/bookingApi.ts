import { BookingRequest, Booking, BookingFilter } from "../types";

export async function createBooking(req: BookingRequest): Promise<Booking> {
  return fetch(`http://localhost:8080/v1/bookings`, {
    method: "POST",
    body: JSON.stringify(req),
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
  }).then((res) => res.json());
}

export async function listBookings(filter?: BookingFilter): Promise<Booking[]> {
  const queryString = createBookingFilter(filter);
  return fetch(`http://localhost:8080/v1/bookings${queryString}`).then((res) =>
    res.json()
  );
}

function createBookingFilter(filter?: BookingFilter): string {
  if (!filter) {
    return "";
  }

  const { userId, cabinId } = filter;
  let queryString = "";
  if (userId) {
    queryString += `userId=${userId}`;
  }

  if (cabinId) {
    queryString = queryString
      ? `${queryString}&cabinId=${cabinId}`
      : `cabinId=${cabinId}`;
  }

  return queryString ? `?${queryString}` : "";
}
