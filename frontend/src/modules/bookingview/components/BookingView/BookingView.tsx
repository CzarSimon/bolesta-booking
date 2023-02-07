import React, { ChangeEvent, SyntheticEvent, useState } from "react";
import { Button, DatePicker, Input } from "antd";
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
import { ItemSelect } from "../../../../components/ItemSelect";

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

  const updatePassword = (e: ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value);
  };

  const updateDates = (_: any, dates: [string, string]) => {
    const [fromStr, toStr] = dates;
    setFrom(new Date(fromStr));
    setTo(new Date(toStr));
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
      <div className={styles.BookingForm}>
        <h2>Välj datum</h2>
        <DatePicker.RangePicker onChange={updateDates} />
        <h2>Personliga detaljer</h2>
        <ItemSelect
          items={users}
          placeholder="Välj Lundinare"
          onChange={setUserId}
        />
        <Input.Password placeholder="Lösenord" onChange={updatePassword} />
        <ErrorText error={err} />
        {success !== undefined ? (
          <BookingResultModal
            success={success}
            onClose={() => setSuccess(undefined)}
          />
        ) : null}
        <Button type="primary" block onMouseUp={onSubmit}>
          Boka
        </Button>
      </div>
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
