import React, { ChangeEvent, useState } from "react";
import { Button, Input } from "antd";
import { ErrorText } from "../../../../components/ErrorText";
import {
  Failure,
  LoginRequest,
  Optional,
  Result,
  Success,
} from "../../../../types";

import styles from "./Login.module.css";

interface Props {
  submit: (req: LoginRequest) => void;
}

export function Login({ submit }: Props) {
  const [email, setEmail] = useState<Optional<string>>();
  const [password, setPassword] = useState<Optional<string>>();
  const [err, setErr] = useState<Optional<string>>();

  const updateEmail = (e: ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
  };

  const updatePassword = (e: ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value);
  };

  const onSubmit = () => {
    parseLoginRequest(email, password)
      .then((req) => submit(req))
      .catch((e) => setErr(e));
  };

  return (
    <div className={styles.LoginPage}>
      <h1 className={styles.Title}>Bölesta Booking</h1>
      <h2 className={styles.Details}>Logga in</h2>
      <Input
        type="email"
        placeholder="Mailaddress"
        onChange={updateEmail}
        size="large"
        className={styles.FormEmail}
      />
      <Input.Password
        placeholder="Lösenord"
        onChange={updatePassword}
        size="large"
        className={styles.FormPassword}
      />
      <Button block type="primary" onClick={onSubmit} size="large">
        Logga in
      </Button>
      <ErrorText error={err} />
      <div className={styles.PoweredBy}>Lindgren & Lundin</div>
    </div>
  );
}

function parseLoginRequest(
  email?: string,
  password?: string
): Result<LoginRequest, string> {
  if (!email || !password) {
    return Failure("Alla fält måste fyllas i");
  }

  return Success({
    email,
    password,
  });
}
