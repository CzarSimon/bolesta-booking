import { BASE_URL } from "../constants";
import { Cabin } from "../types";
import { httpclient } from "./httpclient";
import { wrapAndLogError } from "./util";

export async function getCabin(id: string): Promise<Cabin> {
  const { body, error, metadata } = await httpclient.get<Cabin>({
    url: `${BASE_URL}/v1/cabins/${id}`,
  });

  if (!body) {
    throw wrapAndLogError(`failed to fetch cabin(id=${id})`, error, metadata);
  }

  return body;
}

export async function getCabins(): Promise<Cabin[]> {
  const { body, error, metadata } = await httpclient.get<Cabin[]>({
    url: `${BASE_URL}/v1/cabins`,
  });

  if (!body) {
    throw wrapAndLogError(`failed to fetch cabins`, error, metadata);
  }

  return body;
}
