import { BASE_URL } from "../constants";
import {
  BookingRequest,
  Booking,
  BookingFilter,
  StatusBody,
  BookingResult,
} from "../types";
import { httpclient } from "./httpclient";
import { wrapAndLogError } from "./util";

export async function createBooking(
  req: BookingRequest,
  dryRun: boolean
): Promise<BookingResult> {
  const { body, error, metadata } = await httpclient.post<Booking>({
    url: `${BASE_URL}/v1/bookings?dry-run=${dryRun}`,
    body: req,
  });

  const BOOKING_VIOLATION_KEY = "x-booking-rule-violation";

  if (metadata.status === 403 && metadata.headers[BOOKING_VIOLATION_KEY]) {
    const violation = parseInt(metadata.headers[BOOKING_VIOLATION_KEY]);
    return { success: false, errorId: violation };
  }

  if (!body) {
    wrapAndLogError(
      `failed to create booking(cabinId=${req.cabinId})`,
      error,
      metadata
    );
    return { success: false, errorId: 0 };
  }

  return { success: true, bookingId: body?.id };
}

export async function listBookings(filter?: BookingFilter): Promise<Booking[]> {
  const queryString = createBookingFilter(filter);
  const { body, error, metadata } = await httpclient.get<Booking[]>({
    url: `${BASE_URL}/v1/bookings${queryString}`,
  });

  if (!body) {
    throw wrapAndLogError(
      `failed to fetch bookings with filter=${queryString}`,
      error,
      metadata
    );
  }

  return body;
}

export async function deleteBooking(id: string): Promise<StatusBody> {
  const { body, error, metadata } = await httpclient.delete<StatusBody>({
    url: `${BASE_URL}/v1/bookings/${id}`,
  });

  if (!body) {
    throw wrapAndLogError(
      `failed to delelete booking booking(id=${id})`,
      error,
      metadata
    );
  }

  return body;
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
