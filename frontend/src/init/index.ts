import { v4 as uuid } from "uuid";
import log, { Handlers, ConsoleHandler, level } from "@czarsimon/remotelogger";
import { CLIENT_ID_KEY, CURRENT_USER_KEY } from "../constants";
import { Client, Optional, User } from "../types";
import { initHttpclient } from "../api/httpclient";

export function readUser(): Optional<User> {
  const serializedUser = localStorage.getItem(CURRENT_USER_KEY);
  if (!serializedUser) {
    return;
  }

  return JSON.parse(serializedUser);
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
