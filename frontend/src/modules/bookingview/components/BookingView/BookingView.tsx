import React, { SyntheticEvent, useState } from "react";
import { Button, DatePicker } from "antd";
import { ErrorText } from "../../../../components/ErrorText";
import {
  Cabin,
  BookingRequest,
  Optional,
  Result,
  Success,
  Failure,
  Booking,
} from "../../../../types";
import { BookingResultModal } from "../BookingResultModal";
import { ItemSelect } from "../../../../components/ItemSelect";

import styles from "./BookingView.module.css";
import { NavTitle } from "../../../../components/NavTitle";

interface Props {
  cabins: Cabin[];
  handleBookingRequest: (req: BookingRequest) => Promise<Booking>;
}

export function BookingView({ cabins, handleBookingRequest }: Props) {
  const [cabinId, setCabinId] = useState<Optional<string>>();
  const [from, setFrom] = useState<Optional<Date>>();
  const [to, setTo] = useState<Optional<Date>>();
  const [err, setErr] = useState<Optional<string>>();
  const [success, setSuccess] = useState<Optional<boolean>>();

  const updateDates = (_: any, dates: [string, string]) => {
    const [fromStr, toStr] = dates;
    setFrom(new Date(fromStr));
    setTo(new Date(toStr));
  };

  const onSubmit = (e: SyntheticEvent) => {
    e.preventDefault();
    parseBookingRequest(cabinId, to, from)
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
      <NavTitle title="Ny bokning" />
      <div className={styles.BookingForm}>
        <ItemSelect
          items={cabins}
          onChange={setCabinId}
          large
          placeholder="Välj stuga"
        />
        <DatePicker.RangePicker
          onChange={updateDates}
          size="large"
          className={styles.FormElement}
          placeholder={["Från datum", "Till datum"]}
        />
        <ErrorText error={err} />
        {success !== undefined ? (
          <BookingResultModal
            success={success}
            onClose={() => setSuccess(undefined)}
          />
        ) : null}
        <Button
          type="primary"
          block
          onMouseUp={onSubmit}
          size="large"
          className={styles.FormButton}
        >
          Boka
        </Button>
      </div>
    </div>
  );
}

function parseBookingRequest(
  cabinId?: string,
  to?: Date,
  from?: Date
): Result<BookingRequest, string> {
  if (!cabinId || !to || !from) {
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
    cabinId,
  });
}
