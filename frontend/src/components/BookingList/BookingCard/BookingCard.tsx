import React from "react";
import { Card } from "antd";
import { useQueryClient } from "react-query";
import { Booking } from "../../../types";
import { useAuth } from "../../../state/auth/hooks";
import { DeleteButton } from "./DeleteButton";
import { deleteBooking } from "../../../api";

import styles from "./BookingCard.module.css";

interface Props {
  booking: Booking;
}

export function BookingCard({ booking }: Props) {
  const queryClient = useQueryClient();
  const { user: loggedInUser } = useAuth();

  const { startDate, endDate, cabin, user } = booking;
  const showDelete = loggedInUser && loggedInUser.id === booking.user.id;

  const onDelete = () => {
    deleteBooking(booking.id).then(() => {
      queryClient.invalidateQueries("bookings");
    });
  };

  return (
    <>
      <Card
        title={cabin.name}
        extra={showDelete ? <DeleteButton onConfirm={onDelete} /> : null}
        bordered={false}
        className={styles.BookingCard}
      >
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
