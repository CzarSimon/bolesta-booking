import React, { ChangeEvent, SyntheticEvent, useState } from "react";
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

  const onSubmit = (e: SyntheticEvent) => {
    e.preventDefault();
    parseLoginRequest(email, password)
      .then((req) => submit(req))
      .catch((e) => setErr(e));
  };

  return (
    <div className={styles.LoginPage}>
      <h1 className={styles.Title}>Bölesta Booking</h1>
      <h3 className={styles.Details}>Logga in till ditt konto</h3>
      <form onSubmit={onSubmit}>
        <label className={styles.FormElement}>
          <p className={styles.LabelText}>Mailadress</p>
          <input
            type="email"
            onChange={updateEmail}
            className={styles.InputField}
          />
        </label>
        <label className={styles.FormElement}>
          <p className={styles.LabelText}>Lösenord</p>
          <input
            type="password"
            onChange={updatePassword}
            className={styles.InputField}
          />
        </label>
        <button className={styles.FormButton} type="submit">
          Logga in
        </button>
        <ErrorText error={err} />
      </form>
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
