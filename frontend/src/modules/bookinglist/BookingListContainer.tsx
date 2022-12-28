import React, { useState } from "react";
import { BookingListView } from "./components/BookingListView";
import { useBookings, useCabins, useUsers } from "../../hooks";
import { BookingFilter } from "../../types";

export function BookingListContainer() {
  const [filter, setFilter] = useState<BookingFilter>({});
  const cabins = useCabins();
  const users = useUsers();
  const bookings = useBookings(filter);

  return (
    <BookingListView
      cabins={cabins || []}
      users={users || []}
      bookings={bookings}
      updateFilter={setFilter}
    />
  );
}
