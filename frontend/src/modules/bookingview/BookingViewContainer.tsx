import React from "react";
import { createBooking } from "../../api";
import { LoadingIndicator } from "../../components/LoadingIndicator";
import { useCabins } from "../../hooks";
import { BookingRequest } from "../../types";
import { BookingView } from "./components/BookingView";

export function BookingViewContainer() {
  const cabins = useCabins();
  const createSharpBooking = (req: BookingRequest) => {
    return createBooking(req, false);
  };

  return cabins ? (
    <BookingView cabins={cabins} handleBookingRequest={createSharpBooking} />
  ) : (
    <LoadingIndicator />
  );
}
