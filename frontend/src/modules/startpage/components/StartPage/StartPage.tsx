import React from "react";
import { BookingCard } from "../../../../components/BookingList/BookingCard";
import { Booking, Optional, User } from "../../../../types";
import { BookingsLink } from "../BookingsLink";
import { NewBookingLink } from "../NewBookingLink";
import { ProfileLink } from "../ProfileLink";

import styles from "./StartPage.module.css";

interface Props {
  user: User;
  bookings: Booking[];
}

export function StartPage({ user, bookings }: Props) {
  const nextBooking: Optional<Booking> = bookings.length
    ? bookings[0]
    : undefined;
  return (
    <div className={styles.StartPage}>
      <h1 className={styles.Title}>Bölesta booking</h1>
      {nextBooking && (
        <div>
          <h3>Nästa bookning</h3>
          <BookingCard booking={nextBooking} />
        </div>
      )}
      <NewBookingLink />
      <BookingsLink />
      <ProfileLink user={user} />
    </div>
  );
}
