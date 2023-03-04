import { v4 as uuid } from "uuid";
import log, { Handlers, ConsoleHandler, level } from "@czarsimon/remotelogger";
import { AUTH_TOKEN_KEY, CLIENT_ID_KEY, CURRENT_USER_KEY } from "../constants";
import { Client, Optional, User } from "../types";
import { initHttpclient } from "../api/httpclient";

export function readUser(): Optional<User> {
  const serializedUser = localStorage.getItem(CURRENT_USER_KEY);
  if (!serializedUser) {
    return;
  }

  const token = readToken();
  if (!token) {
    localStorage.removeItem(CURRENT_USER_KEY);
    return;
  }

  return JSON.parse(serializedUser);
}

export function readToken(): Optional<string> {
  const token = localStorage.getItem(AUTH_TOKEN_KEY);
  if (!token) {
    return;
  }

  const expired = isTokenExpired(token);

  if (expired) {
    localStorage.removeItem(AUTH_TOKEN_KEY);
    return;
  }

  return token;
}

function isTokenExpired(token: string): boolean {
  try {
    const body = parseToken(token);
    const now = new Date();
    const expiry = new Date(body.exp * 1000);
    return expiry < now;
  } catch (error) {
    console.log("parsing token failed", error);
    return true;
  }
}

export function initLoggerAndHttpClient() {
  const client = getClient();
  const handlers = getLogHandlers();
  initHttpclient(client, handlers);

  log.configure(handlers);
  log.debug("initiated remotelogger");
  log.debug("initiated httpclient");
}

function getLogHandlers(): Handlers {
  return {
    console: new ConsoleHandler(level.DEBUG),
  };
}

function getClient(): Client {
  return {
    id: getOrCreateId(CLIENT_ID_KEY),
    sessionId: uuid(),
  };
}

function getOrCreateId(key: string): string {
  const id = localStorage.getItem(key);
  if (id) {
    return id;
  }

  const newId = uuid();
  localStorage.setItem(key, newId);
  return newId;
}

interface TokenBody {
  exp: number;
  iat: number;
  nbf: number;
  role: string;
  sub: string;
}

function parseToken(token: string): TokenBody {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [alg, body, sig] = token.split(".");
  const decodedBody = atob(body);
  return JSON.parse(decodedBody);
}
