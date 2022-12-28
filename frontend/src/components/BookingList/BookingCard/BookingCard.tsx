import React from "react";
import { Booking } from "../../../types";

import styles from "./BookingCard.module.css";

interface Props {
  booking: Booking;
}

export function BookingCard({ booking }: Props) {
  const { startDate, endDate, cabin, user } = booking;
  return (
    <div className={styles.BookingCard}>
      <h2>{cabin.name}</h2>
      <p>Bokat av: {user.name}</p>
      <p>
        {formatDate(startDate)} - {formatDate(endDate)}
      </p>
    </div>
  );
}

function formatDate(date: string): string {
  return date.split("T")[0];
}
