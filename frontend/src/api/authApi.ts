import { AuthenticatedResponse, LoginRequest } from "../types";

export async function requestLogin(
  req: LoginRequest
): Promise<AuthenticatedResponse> {
  return fetch(`http://localhost:8080/v1/login`, {
    method: "POST",
    body: JSON.stringify(req),
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
  }).then((res) => res.json());
}
