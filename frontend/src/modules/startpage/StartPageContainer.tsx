import React from "react";
import { LoadingIndicator } from "../../components/LoadingIndicator";
import { useMyBookings } from "../../hooks";
import { StartPage } from "./components/StartPage";

export function StartPageContainer() {
  const bookings = useMyBookings();
  return bookings ? <StartPage bookings={bookings} /> : <LoadingIndicator />;
}
