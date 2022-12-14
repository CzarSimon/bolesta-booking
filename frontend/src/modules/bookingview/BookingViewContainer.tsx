import React from "react";
import { useParams } from "react-router-dom";
import { createBooking } from "../../api";
import { LoadingIndicator } from "../../components/LoadingIndicator";
import { useBookings, useCabin, useUsers } from "../../hooks";
import { BookingView } from "./components/BookingView";

export function BookingViewContainer() {
  const { cabinId } = useParams();
  const cabin = useCabin(cabinId!);
  const users = useUsers();
  // const bookings = useBookings({ cabinId });

  return cabin && users ? (
    <BookingView
      cabin={cabin}
      users={users}
      handleBookingRequest={createBooking}
    />
  ) : (
    <LoadingIndicator />
  );
}
