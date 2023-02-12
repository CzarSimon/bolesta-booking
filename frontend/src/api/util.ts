import { ResponseMetadata } from "@czarsimon/httpclient";
import log from "@czarsimon/remotelogger";
import { Optional } from "../types";

export function wrapAndLogError(
  info: string,
  error: Optional<Error>,
  metadata: ResponseMetadata
): Error {
  const { method, requestId, status, url } = metadata;
  const message = `${info}, error=${error?.message} endpoint=${method} ${url} status=${status} requestId=${requestId}`;
  log.error(message);
  return new Error(message);
}
