import { BASE_URL } from "../constants";
import { AuthenticatedResponse, LoginRequest } from "../types";
import { httpclient } from "./httpclient";
import { wrapAndLogError } from "./util";

export async function requestLogin(
  req: LoginRequest
): Promise<AuthenticatedResponse> {
  const { body, error, metadata } =
    await httpclient.post<AuthenticatedResponse>({
      url: `${BASE_URL}/v1/login`,
      body: req,
    });

  if (error) {
    throw error;
  }

  if (!body) {
    throw wrapAndLogError(
      `failed to login user(email=${req.email})`,
      error,
      metadata
    );
  }

  return body;
}
