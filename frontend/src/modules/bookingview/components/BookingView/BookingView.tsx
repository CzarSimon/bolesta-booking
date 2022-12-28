import React, { ChangeEvent, SyntheticEvent, useState } from "react";
import { ErrorText } from "../../../../components/ErrorText";
import {
  Cabin,
  User,
  BookingRequest,
  Optional,
  Result,
  Success,
  Failure,
  Booking,
} from "../../../../types";
import { BookingResultModal } from "../BookingResultModal";

import styles from "./BookingView.module.css";

interface Props {
  cabin: Cabin;
  users: User[];
  handleBookingRequest: (req: BookingRequest) => Promise<Booking>;
}

export function BookingView({ cabin, users, handleBookingRequest }: Props) {
  const [from, setFrom] = useState<Optional<Date>>();
  const [to, setTo] = useState<Optional<Date>>();
  const [userId, setUserId] = useState<Optional<string>>();
  const [password, setPassword] = useState<Optional<string>>();
  const [err, setErr] = useState<Optional<string>>();
  const [success, setSuccess] = useState<Optional<boolean>>();

  const updateFrom = (e: ChangeEvent<HTMLInputElement>) => {
    setFrom(new Date(e.target.value));
  };

  const updateTo = (e: ChangeEvent<HTMLInputElement>) => {
    setTo(new Date(e.target.value));
  };

  const updateUserId = (e: SyntheticEvent<HTMLSelectElement, Event>) => {
    setUserId(e.currentTarget.value);
  };

  const updatePassword = (e: ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value);
  };

  const onSubmit = (e: SyntheticEvent) => {
    e.preventDefault();
    parseBookingRequest(cabin, to, from, userId, password)
      .then((r) => {
        setErr(undefined);
        handleBookingRequest(r)
          .then((_) => setSuccess(true))
          .catch((_) => setSuccess(false));
      })
      .catch((e) => setErr(e));
  };

  return (
    <div className={styles.BookingView}>
      <h1 className={styles.CabinName}>{cabin.name}</h1>
      <form className={styles.BookingForm} onSubmit={onSubmit}>
        <h2>Välj datum</h2>
        <label className={styles.FormElement}>
          Från: <input type="date" onChange={updateFrom} />
        </label>
        <label className={styles.FormElement}>
          Till: <input type="date" onChange={updateTo} />
        </label>
        <h2>Personliga detaljer</h2>
        <select className={styles.FormElement} onChange={updateUserId} defaultValue="none-selected">
          <option disabled value="none-selected">
            -- Välj Lundinare --
          </option>
          {users.map((u) => (
            <option key={u.id} value={u.id}>
              {u.name}
            </option>
          ))}
        </select>
        <label className={styles.FormElement}>
          Lösenord: <input type="password" onChange={updatePassword} />
        </label>
        <ErrorText error={err} />
        {success !== undefined ? (
          <BookingResultModal
            success={success}
            onClose={() => setSuccess(undefined)}
          />
        ) : null}
        <button className={styles.FormButton} type="submit">Boka</button>
      </form>
    </div>
  );
}

function parseBookingRequest(
  cabin: Cabin,
  to?: Date,
  from?: Date,
  userId?: string,
  password?: string
): Result<BookingRequest, string> {
  if (!to || !from || !userId || !password) {
    return Failure("Alla fält måste fyllas i");
  }

  if (from < new Date()) {
    return Failure("Startdatumet måste ligga i framtiden!");
  }

  if (to < from) {
    return Failure("Slutdatumet måste vara efter startdatumet");
  }

  return Success({
    startDate: from,
    endDate: to,
    cabinId: cabin.id,
    userId,
    password,
  });
}
