import React, { SyntheticEvent, useState } from "react";
import { Button } from "antd";
import { Cabin, User, BookingFilter, Optional } from "../../../../types";

import styles from "./BookingFilterSelector.module.css";

interface Props {
  cabins: Cabin[];
  users: User[];
  updateFilter: (filter: BookingFilter) => void;
}

const ANY_VALUE = "*";

export function BookingFilterSelector({ cabins, users, updateFilter }: Props) {
  const [visible, setVisible] = useState<boolean>(false);
  const [cabinId, setCabinId] = useState<Optional<string>>();
  const [userId, setUserId] = useState<Optional<string>>();

  const updateCabinId = (e: SyntheticEvent<HTMLSelectElement, Event>) => {
    setCabinId(getValue(e.currentTarget.value));
  };

  const updateUserId = (e: SyntheticEvent<HTMLSelectElement, Event>) => {
    setUserId(getValue(e.currentTarget.value));
  };

  const onSubmit = (e: SyntheticEvent) => {
    e.preventDefault();
    updateFilter({
      cabinId: cabinId,
      userId: userId,
    });
  };

  return (
    <div>
      {visible && (
        <form onSubmit={onSubmit} className={styles.Form}>
          <label>
            <p className={styles.LabelText}>Stuga</p>
            <select
              className={styles.Select}
              onChange={updateCabinId}
              defaultValue={cabinId}
            >
              <option value={ANY_VALUE}>Alla</option>
              {cabins.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name}
                </option>
              ))}
            </select>
          </label>
          <label>
            <p className={styles.LabelText}>Bokat av</p>
            <select
              className={styles.Select}
              onChange={updateUserId}
              defaultValue={userId}
            >
              <option value={ANY_VALUE}>Alla</option>
              {users.map((u) => (
                <option key={u.id} value={u.id}>
                  {u.name}
                </option>
              ))}
            </select>
          </label>
          <Button block type="primary" onMouseUp={onSubmit}>
            Filtrera
          </Button>
        </form>
      )}
      {visible ? (
        <Button block type="text" onMouseUp={() => setVisible(false)}>
          Dölj filter
        </Button>
      ) : (
        <Button block onMouseUp={() => setVisible(true)}>
          Visa filter
        </Button>
      )}
    </div>
  );
}

function getValue(val: string): Optional<string> {
  return val === ANY_VALUE ? undefined : val;
}
