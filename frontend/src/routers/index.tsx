import React from "react";
import { createBrowserRouter } from "react-router-dom";
import { BookingListContainer } from "../modules/bookinglist";
import { BookingViewContainer } from "../modules/bookingview";
import { LoginContainer } from "../modules/login";
import { ProfileContainer } from "../modules/profile";
import { StartPageContainer } from "../modules/startpage";
import { ProtectedRoutes } from "./ProtectedRoutes";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <ProtectedRoutes />,
    children: [
      {
        path: "/",
        element: <StartPageContainer />,
      },
      {
        path: "bookings",
        element: <BookingListContainer />,
      },
      {
        path: "bookings/new",
        element: <BookingViewContainer />,
      },
      {
        path: "profile",
        element: <ProfileContainer />,
      },
    ],
  },
  {
    path: "/login",
    element: <LoginContainer />,
  },
]);
