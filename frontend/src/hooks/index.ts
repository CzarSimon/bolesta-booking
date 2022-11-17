import { useQuery } from "react-query";
import { getCabins } from "../api";
import { Cabin, Optional } from "../types";

export function useCabins(): Optional<Cabin[]> {
  const { data } = useQuery<Cabin[], Error>("cabins", getCabins);
  return data;
}
