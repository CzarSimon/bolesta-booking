import React from "react";
import { BookingList } from "../../../../components/BookingList";
import { NavTitle } from "../../../../components/NavTitle";
import {
  Booking,
  BookingFilter,
  Cabin,
  Optional,
  User,
} from "../../../../types";
import { BookingFilterSelector } from "../BookingFilterSelector";

import styles from "./BookingListView.module.css";

interface Props {
  cabins: Cabin[];
  users: User[];
  bookings: Optional<Booking[]>;
  updateFilter: (filter: BookingFilter) => void;
}

export function BookingListView({
  cabins,
  users,
  bookings,
  updateFilter,
}: Props) {
  return (
    <div className={styles.BookingListView}>
      <NavTitle title="Bokningar" />
      <BookingFilterSelector
        cabins={cabins}
        users={users}
        updateFilter={updateFilter}
      />
      <BookingList bookings={bookings} />
    </div>
  );
}
