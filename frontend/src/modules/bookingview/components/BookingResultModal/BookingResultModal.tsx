import React from "react";
import { useNavigate } from "react-router-dom";

import styles from "./BookingResultModal.module.css";

interface Props {
  success: boolean;
  onClose: () => void;
}

export function BookingResultModal({ success, onClose }: Props) {
  const navigate = useNavigate();
  const onClick = () => {
    if (!success) {
      onClose();
    }
    navigate("/");
  };

  return (
    <div className={styles.Modal}>
      <div className={styles.ModalContent}>
        {success ? <p>Bokningen lyckades!</p> : <p>Bookningen misslyckades!</p>}
        <button onClick={onClick} className={styles.OkButton}>
          Ok
        </button>
      </div>
    </div>
  );
}
