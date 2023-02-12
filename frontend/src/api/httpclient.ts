import { Fetch, HttpClient, Headers } from "@czarsimon/httpclient";
import { Handlers } from "@czarsimon/remotelogger";
import { AUTH_TOKEN_KEY } from "../constants";
import { Client } from "../types";

export let httpclient = new HttpClient({
  baseHeaders: {
    "Content-Type": "application/json",
  },
});

export function initHttpclient(client: Client, handlers: Handlers) {
  const baseHeaders: Headers = {
    ...httpclient.getHeaders(),
    "X-Client-ID": client.id,
    "X-Session-ID": client.sessionId,
  };

  const token = localStorage.getItem(AUTH_TOKEN_KEY);
  if (token) {
    baseHeaders["Authorization"] = `Bearer ${token}`;
  }

  httpclient = new HttpClient({
    logHandlers: handlers,
    transport: new Fetch(),
    baseHeaders,
  });
}

export function setHeader(name: string, value: string) {
  httpclient.setHeaders({
    ...httpclient.getHeaders(),
    [name]: value,
  });
}

export function removeHeader(name: string) {
  const { [name]: _, ...rest } = httpclient.getHeaders();
  httpclient.setHeaders(rest);
}
