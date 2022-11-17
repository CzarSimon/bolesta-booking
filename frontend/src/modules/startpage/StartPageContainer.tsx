import React from "react";
import { LoadingIndicator } from "../../components/LoadingIndicator";
import { useCabins } from "../../hooks";
import { StartPage } from "./components/StartPage";

export function StartPageContainer() {
  const cabins = useCabins();
  return cabins ? <StartPage cabins={cabins} /> : <LoadingIndicator />;
}
