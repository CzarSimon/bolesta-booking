import React from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "antd";

export function BookingsLink() {
  const navigate = useNavigate();
  const toBookings = () => {
    navigate(`bookings`);
  };

  return (
    <Button size="large" block onClick={toBookings}>
      Se bokningar
    </Button>
  );
}
