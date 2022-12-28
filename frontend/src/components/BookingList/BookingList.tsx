import React from "react";
import { Booking, Optional } from "../../types";
import { LoadingIndicator } from "../LoadingIndicator";
import { BookingCard } from "./BookingCard";

import styles from "./BookingList.module.css";

interface Props {
  bookings: Optional<Booking[]>;
}

export function BookingList({ bookings }: Props) {
  if (bookings === undefined) {
    return <LoadingIndicator />;
  }

  const futureBookings = filterPastBookings(bookings);
  return (
    <div>
      <ul className={styles.List}>
        {futureBookings.map((booking) => (
          <li key={booking.id}>
            <BookingCard booking={booking} />
          </li>
        ))}
      </ul>
    </div>
  );
}

function filterPastBookings(bookings: Booking[]): Booking[] {
  const now = new Date();
  return bookings.filter((b) => new Date(b.endDate) > now);
}
