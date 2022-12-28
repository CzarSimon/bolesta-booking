import React, { SyntheticEvent, useState } from "react";
import { Cabin, User, BookingFilter, Optional } from "../../../../types";

interface Props {
  cabins: Cabin[];
  users: User[];
  updateFilter: (filter: BookingFilter) => void;
}

const ANY_VALUE = "*";

export function BookingFilterSelector({ cabins, users, updateFilter }: Props) {
  const [visible, setVisible] = useState<boolean>(true);
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
        <form onSubmit={onSubmit}>
          <select onChange={updateCabinId} defaultValue={cabinId}>
            <option value={ANY_VALUE}>Alla</option>
            {cabins.map((c) => (
              <option key={c.id} value={c.id}>
                {c.name}
              </option>
            ))}
          </select>
          <select onChange={updateUserId} defaultValue={userId}>
            <option value={ANY_VALUE}>Alla</option>
            {users.map((u) => (
              <option key={u.id} value={u.id}>
                {u.name}
              </option>
            ))}
          </select>
          <button type="submit">Filtrera</button>
        </form>
      )}
      <button onClick={() => setVisible(!visible)}>
        {visible ? "Dölj filter" : "Visa filter"}
      </button>
    </div>
  );
}

function getValue(val: string): Optional<string> {
  return val === ANY_VALUE ? undefined : val;
}
