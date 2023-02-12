import { BASE_URL } from "../constants";
import { User } from "../types";
import { httpclient } from "./httpclient";
import { wrapAndLogError } from "./util";

export async function getUsers(): Promise<User[]> {
  const { body, error, metadata } = await httpclient.get<User[]>({
    url: `${BASE_URL}/v1/users`,
  });

  if (!body) {
    throw wrapAndLogError(`failed to fetch users`, error, metadata);
  }

  return body;
}
