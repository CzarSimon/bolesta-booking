import React from "react";
import { useNavigate } from "react-router-dom";
import { BookingResult, Optional } from "../../../../types";

import styles from "./BookingResultModal.module.css";

interface Props {
  result?: BookingResult;
  onClose: () => void;
}

const VIOLATIONS: { [key: number]: string } = {
  0: "Okänt fel",
  1: "Felaktig bokning",
  2: "Bokningen är för lång, den får max vara 1 vecka",
  3: "Bokningen ligger för långt i framtiden, du kan endast boka 3 månader framåt",
  4: "Du har för många aktiva bokningar, man får max ha två",
};

export function BookingResultModal({ result, onClose }: Props) {
  const navigate = useNavigate();
  const onClick = () => {
    if (result && !result.success) {
      onClose();
      return;
    }
    navigate("/");
  };

  const violationText: Optional<string> = result?.errorId
    ? VIOLATIONS[result.errorId]
    : undefined;

  return (
    <div className={styles.Modal}>
      <div className={styles.ModalContent}>
        {result?.success ? (
          <p>Bokningen lyckades!</p>
        ) : (
          <div>
            <h3>Bookningen misslyckades!</h3>
            <p>{violationText}</p>
          </div>
        )}
        <button onClick={onClick} className={styles.OkButton}>
          Ok
        </button>
      </div>
    </div>
  );
}
