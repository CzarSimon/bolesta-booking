import React from "react";
import { Booking } from "../../../types";
import { Card } from "antd";

import styles from "./BookingCard.module.css";

interface Props {
  booking: Booking;
}

export function BookingCard({ booking }: Props) {
  const { startDate, endDate, cabin, user } = booking;
  return (
    <>
      <Card title={cabin.name} bordered={false} className={styles.BookingCard}>
        <div className={styles.Content}>
          <p>Bokat av: {user.name}</p>
          <p>
            {formatDate(startDate)} - {formatDate(endDate)}
          </p>
        </div>
      </Card>
    </>
  );
}

function formatDate(date: string): string {
  return date.split("T")[0];
}
