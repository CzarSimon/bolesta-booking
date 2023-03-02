import React from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "antd";

import styles from "./NewBookingLink.module.css";

export function NewBookingLink() {
  const navigate = useNavigate();
  const toNewBooking = () => {
    navigate(`bookings/new`);
  };

  return (
    <Button
      size="large"
      type="primary"
      onClick={toNewBooking}
      className={styles.Button}
      block
    >
      Ny bokning
    </Button>
  );
}
