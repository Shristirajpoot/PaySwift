import React, { useState } from "react";

function PaymentForm({ initiatePayment }) {
  const [amount, setAmount] = useState("");
  const [email, setEmail] = useState("");
  const [currency, setCurrency] = useState("NGN");

  const handleSubmit = (event) => {
    event.preventDefault();
    initiatePayment(amount, email, currency);
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>Amount: </label>
        <input
          type="number"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
        />
      </div>
      <div>
        <label>Email: </label>
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
      </div>
      <div>
        <label>Currency: </label>
        <select
          value={currency}
          onChange={(e) => setCurrency(e.target.value)}
        >
          <option value="NGN">NGN</option>
          <option value="USD">USD</option>
        </select>
      </div>
      <button type="submit">Pay</button>
    </form>
  );
}

export default PaymentForm;
