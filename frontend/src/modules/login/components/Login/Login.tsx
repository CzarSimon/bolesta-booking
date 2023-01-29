import React from "react";

import styles from "./Login.module.css";


export function Login() {
return (
    <div className={styles.LoginPage}>
        <h1 className={styles.Title}>Bölesta Booking</h1>
         <h3 className={styles.Details}>Logga in till ditt konto</h3>
         <label className={styles.FormElement}>
          <p className={styles.LabelText}>Mailadress</p>
          <input
            type="mail"

            className={styles.InputField}
          /> 
        </label>
         <label className={styles.FormElement}>
          <p className={styles.LabelText}>Lösenord</p>
          <input
            type="password"

            className={styles.InputField}
          /> 
        </label>
        <button className={styles.FormButton} type="submit">
          Logga in
        </button>
        <div className={styles.PoweredBy}>Lindgren & Lundin</div>
    </div>

)
}