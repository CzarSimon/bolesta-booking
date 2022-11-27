import { Cabin } from "../types";

export function getCabin(id: string): Promise<Cabin> {
  return fetch(`http://localhost:8080/v1/cabins/${id}`).then((res) =>
    res.json()
  );
}

export function getCabins(): Promise<Cabin[]> {
  return fetch("http://localhost:8080/v1/cabins").then((res) => res.json());
}
