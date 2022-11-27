import React from "react";
import { createBrowserRouter } from "react-router-dom";
import { BookingViewContainer } from "../modules/bookingview";
import { StartPageContainer } from "../modules/startpage";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <StartPageContainer />,
  },
  {
    path: "cabins/:cabinId",
    element: <BookingViewContainer />,
  },
]);
