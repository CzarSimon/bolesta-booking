import React, { SyntheticEvent, useState } from "react";
import { Button } from "antd";
import { Cabin, User, BookingFilter, Optional } from "../../../../types";

import styles from "./BookingFilterSelector.module.css";
import { ItemSelect } from "../../../../components/ItemSelect";

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

  const updateCabinId = (id: string) => setCabinId(getValue(id));
  const updateUserId = (id: string) => setUserId(getValue(id));

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
        <div className={styles.Form}>
          <p className={styles.LabelText}>Stuga</p>
          <ItemSelect
            items={cabins}
            onChange={updateCabinId}
            anyOption={{ label: "Alla", value: ANY_VALUE }}
          />
          <p className={styles.LabelText}>Bokat av</p>
          <ItemSelect
            items={users}
            onChange={updateUserId}
            anyOption={{ label: "Alla", value: ANY_VALUE }}
          />
          <Button block type="primary" onMouseUp={onSubmit}>
            Filtrera
          </Button>
        </div>
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
