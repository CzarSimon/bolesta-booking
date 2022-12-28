import React from "react";
import { createBrowserRouter } from "react-router-dom";
import { BookingListContainer } from "../modules/bookinglist";
import { BookingViewContainer } from "../modules/bookingview";
import { StartPageContainer } from "../modules/startpage";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <StartPageContainer />,
  },
  {
    path: "bookings",
    element: <BookingListContainer />,
  },
  {
    path: "cabins/:cabinId",
    element: <BookingViewContainer />,
  },
]);
