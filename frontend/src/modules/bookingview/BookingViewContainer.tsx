import React from "react";
import { createBooking } from "../../api";
import { LoadingIndicator } from "../../components/LoadingIndicator";
import { useCabins } from "../../hooks";
import { BookingView } from "./components/BookingView";

export function BookingViewContainer() {
  const cabins = useCabins();

  return cabins ? (
    <BookingView cabins={cabins} handleBookingRequest={createBooking} />
  ) : (
    <LoadingIndicator />
  );
}
