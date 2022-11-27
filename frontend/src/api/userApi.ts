import { User } from "../types";

export function getUsers(): Promise<User[]> {
  return fetch("http://localhost:8080/v1/users").then((res) => res.json());
}
