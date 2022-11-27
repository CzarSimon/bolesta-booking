import { useQuery } from "react-query";
import { getCabin, getCabins, getUsers, listBookings } from "../api";
import { Booking, BookingFilter, Cabin, Optional, User } from "../types";

export function useCabin(id: string): Optional<Cabin> {
  const { data } = useQuery<Cabin, Error>(["cabin", id], () => getCabin(id));
  return data;
}

export function useCabins(): Optional<Cabin[]> {
  const { data } = useQuery<Cabin[], Error>("cabins", getCabins);
  return data;
}

export function useUsers(): Optional<User[]> {
  const { data } = useQuery<User[], Error>("users", getUsers);
  return data;
}

export function useBookings(filter?: BookingFilter): Optional<Booking[]> {
  const { data } = useQuery<Booking[], Error>(
    ["bookings", filter?.userId, filter?.cabinId],
    () => listBookings(filter)
  );
  return data;
}
