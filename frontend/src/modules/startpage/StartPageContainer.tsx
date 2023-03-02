import React from "react";
import { LoadingIndicator } from "../../components/LoadingIndicator";
import { useMyBookings } from "../../hooks";
import { useAuth } from "../../state/auth/hooks";
import { StartPage } from "./components/StartPage";

export function StartPageContainer() {
  const { user } = useAuth();
  const bookings = useMyBookings();
  return user && bookings ? (
    <StartPage user={user} bookings={bookings} />
  ) : (
    <LoadingIndicator />
  );
}
