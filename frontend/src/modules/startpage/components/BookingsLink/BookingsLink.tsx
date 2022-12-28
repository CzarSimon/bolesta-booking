import React from "react";
import { useNavigate } from "react-router-dom";

import styles from "./BookingsLink.module.css";

export function BookingsLink() {
  const navigate = useNavigate();
  const toBookings = () => {
    navigate(`bookings`);
  };

  return (
    <div className={styles.BookingsLink} onClick={toBookings}>
      <h2>Se bokningar</h2>
    </div>
  );
}
