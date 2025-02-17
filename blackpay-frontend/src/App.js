import React, { useState } from "react";
import PaymentForm from "./components/PaymentForm";
import PaymentStatus from "./components/PaymentStatus";

function App() {
  const [paymentLink, setPaymentLink] = useState("");
  const [status, setStatus] = useState("");

  const initiatePayment = async (amount, email, currency) => {
    try {
      const response = await fetch("http://localhost:8080/api/payment", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          amount,
          email,
          currency,
        }),
      });

      if (!response.ok) {
        throw new Error("Payment initiation failed");
      }

      const data = await response.json();
      setPaymentLink(data.paymentLink);
      setStatus("Payment initiated. Redirecting...");
    } catch (error) {
      setStatus("Error: " + error.message);
    }
  };

  return (
    <div className="App">
      <h1>BlackPay Payment</h1>
      <PaymentForm initiatePayment={initiatePayment} />
      {paymentLink && (
        <PaymentStatus paymentLink={paymentLink} status={status} />
      )}
    </div>
  );
}

export default App;
