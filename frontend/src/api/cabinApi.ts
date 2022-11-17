import { Cabin } from "../types";

export function getCabins(): Promise<Cabin[]> {
  return fetch("http://localhost:8080/v1/cabins").then((res) => res.json());
}
