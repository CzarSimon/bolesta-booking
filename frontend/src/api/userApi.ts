import { BASE_URL } from "../constants";
import { ChangePasswordRequest, StatusBody, User } from "../types";
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

export async function changePassword(
  userId: string,
  req: ChangePasswordRequest
): Promise<StatusBody> {
  const { body, error, metadata } = await httpclient.put<StatusBody>({
    url: `${BASE_URL}/v1/users/${userId}/password`,
    body: req,
  });

  if (!body) {
    throw wrapAndLogError(
      `failed to create change password of User(id=${userId})`,
      error,
      metadata
    );
  }

  return body;
}
